package graph_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/ericls/imgdd/graph/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func tAllUsersUnauthorizedAccess(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}

	// Test with no authentication
	tc.clearAuthenticationInfo()
	err := tc.client.Post(`
	query {
		viewer {
			allUsers {
				id
				email
				name
			}
		}
	}`, &resp)
	require.Error(t, err)
	require.Nil(t, resp.Viewer)

	// Test with regular user (not site owner)
	tc.forceAuthenticate()
	err = tc.client.Post(`
	query {
		viewer {
			allUsers {
				id
				email
				name
			}
		}
	}`, &resp)
	require.Error(t, err)
}

func tAllUsersBasicFunctionality(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}

	// Authenticate as site owner
	orgUser := tc.forceAuthenticate(asSiteOwner)

	// Create additional test users
	testEmails := []string{
		"alice@example.com",
		"bob@example.com", 
		"charlie@example.com",
	}
	
	for _, email := range testEmails {
		_, err := tc.identityRepo.CreateUser(email, "password123")
		require.NoError(t, err)
	}

	// Test basic query
	err := tc.client.Post(`
	query {
		viewer {
			allUsers {
				id
				email
				name
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.AllUsers)
	
	// Should have at least 4 users (site owner + 3 test users)
	require.GreaterOrEqual(t, len(resp.Viewer.AllUsers), 4)
	
	// Verify the site owner is in the results
	var foundSiteOwner bool
	for _, user := range resp.Viewer.AllUsers {
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.Email)
		require.NotEmpty(t, user.Name)
		if user.Email == orgUser.User.Email {
			foundSiteOwner = true
		}
	}
	require.True(t, foundSiteOwner, "Site owner should be in the results")
}

func tAllUsersPagination(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}

	// Authenticate as site owner
	tc.forceAuthenticate(asSiteOwner)

	// Create multiple test users
	for i := 0; i < 10; i++ {
		email := uuid.NewString() + "@example.com"
		_, err := tc.identityRepo.CreateUser(email, "password123")
		require.NoError(t, err)
	}

	// Test pagination with limit
	err := tc.client.Post(`
	query {
		viewer {
			allUsers(limit: 5) {
				id
				email
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.Len(t, resp.Viewer.AllUsers, 5)

	// Test pagination with offset
	var resp2 struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}
	err = tc.client.Post(`
	query {
		viewer {
			allUsers(limit: 5, offset: 5) {
				id
				email
			}
		}
	}`, &resp2)
	require.NoError(t, err)
	require.NotNil(t, resp2.Viewer)
	require.LessOrEqual(t, len(resp2.Viewer.AllUsers), 5)

	// Verify no overlap between pages
	firstPageEmails := make(map[string]bool)
	for _, user := range resp.Viewer.AllUsers {
		firstPageEmails[user.Email] = true
	}
	
	for _, user := range resp2.Viewer.AllUsers {
		require.False(t, firstPageEmails[user.Email], "User should not appear in both pages")
	}
}

func tAllUsersSearch(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}

	// Authenticate as site owner
	tc.forceAuthenticate(asSiteOwner)

	// Create test users with specific emails
	testUsers := []string{
		"john.doe@company.com",
		"jane.smith@company.com",
		"alice.jones@different.org",
		"bob.wilson@another.net",
	}
	
	for _, email := range testUsers {
		_, err := tc.identityRepo.CreateUser(email, "password123")
		require.NoError(t, err)
	}

	// Test search by domain
	err := tc.client.Post(`
	query {
		viewer {
			allUsers(search: "company.com") {
				id
				email
				name
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.Len(t, resp.Viewer.AllUsers, 2)
	
	// Verify all results contain the search term
	for _, user := range resp.Viewer.AllUsers {
		require.Contains(t, user.Email, "company.com")
	}

	// Test search by name prefix
	var resp2 struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}
	err = tc.client.Post(`
	query {
		viewer {
			allUsers(search: "john") {
				id
				email
				name
			}
		}
	}`, &resp2)
	require.NoError(t, err)
	require.NotNil(t, resp2.Viewer)
	require.Len(t, resp2.Viewer.AllUsers, 1)
	require.Contains(t, resp2.Viewer.AllUsers[0].Email, "john.doe")

	// Test case-insensitive search
	var resp3 struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}
	err = tc.client.Post(`
	query {
		viewer {
			allUsers(search: "ALICE") {
				id
				email
				name
			}
		}
	}`, &resp3)
	require.NoError(t, err)
	require.NotNil(t, resp3.Viewer)
	require.Len(t, resp3.Viewer.AllUsers, 1)
	require.Contains(t, resp3.Viewer.AllUsers[0].Email, "alice.jones")
}

func tAllUsersParameterValidation(t *testing.T, tc *TestContext) {
	// Authenticate as site owner
	tc.forceAuthenticate(asSiteOwner)

	// Test negative limit
	var resp1 struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}
	err := tc.client.Post(`
	query {
		viewer {
			allUsers(limit: -1) {
				id
				email
			}
		}
	}`, &resp1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit must be non-negative")

	// Test negative offset
	var resp2 struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}
	err = tc.client.Post(`
	query {
		viewer {
			allUsers(offset: -1) {
				id
				email
			}
		}
	}`, &resp2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset must be non-negative")

	// Test limit exceeding maximum
	var resp3 struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}
	err = tc.client.Post(`
	query {
		viewer {
			allUsers(limit: 1001) {
				id
				email
			}
		}
	}`, &resp3)
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit cannot exceed 1000")
}

func tAllUsersWithVariables(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}

	// Authenticate as site owner
	tc.forceAuthenticate(asSiteOwner)

	// Create test users
	for i := 0; i < 3; i++ {
		email := uuid.NewString() + "@test.com"
		_, err := tc.identityRepo.CreateUser(email, "password123")
		require.NoError(t, err)
	}

	// Test with variables
	limit := 2
	offset := 0
	search := "test.com"
	
	err := tc.client.Post(`
	query GetUsers($limit: Int, $offset: Int, $search: String) {
		viewer {
			allUsers(limit: $limit, offset: $offset, search: $search) {
				id
				email
				name
			}
		}
	}`, &resp, 
		client.Var("limit", limit),
		client.Var("offset", offset),
		client.Var("search", search),
	)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.LessOrEqual(t, len(resp.Viewer.AllUsers), 2)
	
	// Verify all results contain the search term
	for _, user := range resp.Viewer.AllUsers {
		require.Contains(t, user.Email, "test.com")
	}
}

func tAllUsersEmptyResults(t *testing.T, tc *TestContext) {
	var resp struct {
		Viewer *struct {
			AllUsers []*model.User
		}
	}

	// Authenticate as site owner
	tc.forceAuthenticate(asSiteOwner)

	// Search for something that doesn't exist
	err := tc.client.Post(`
	query {
		viewer {
			allUsers(search: "nonexistent-domain-12345.com") {
				id
				email
				name
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer)
	require.NotNil(t, resp.Viewer.AllUsers)
	require.Len(t, resp.Viewer.AllUsers, 0)
}

func TestUserManagementResolvers(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tAllUsersUnauthorizedAccess,
		tAllUsersBasicFunctionality,
		tAllUsersPagination,
		tAllUsersSearch,
		tAllUsersParameterValidation,
		tAllUsersWithVariables,
		tAllUsersEmptyResults,
	)
}