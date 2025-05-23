package identity

import (
	"testing"

	"github.com/ericls/imgdd/db"
	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/test_support"
)

func TestIdentityRepo(t *testing.T) {
	conn := db.GetConnection(TestServiceMan.GetDBConfig())
	identityRepo := NewDBIdentityRepo(conn)
	emailAddr := "test@home.arpa"
	var assertUser = func(ou *dm.OrganizationUser, err error) {
		if err != nil {
			t.Errorf("Error creating user: %s", err)
			t.FailNow()
		}
		if ou.User.Email != emailAddr {
			t.Errorf("Expected user.Email to be %s, got %s", emailAddr, ou.User.Email)
		}
		if len(ou.Roles) != 1 {
			t.Errorf("Expected user to have 1 role, got %d", len(ou.Roles))
		}
		if ou.Roles[0].Key != "owner" {
			t.Errorf("Expected user to have role owner, got %s", ou.Roles[0].Key)
		}
	}
	orgUser, err := identityRepo.CreateUserWithOrganization(emailAddr, "test", "123")
	assertUser(orgUser, err)
	orgUser = identityRepo.GetOrganizationUserById(orgUser.Id)
	assertUser(orgUser, nil)
	orgUsers := identityRepo.GetOrganizationUsersByIds([]string{orgUser.Id})
	if len(orgUsers) != 1 {
		t.Errorf("Failed to fetch user in bulk")
	} else {
		assertUser(orgUsers[0], nil)
	}
	assertUser(orgUser, nil)
	identityRepo.AddRoleToOrganizationUser(orgUser.Id, "member")
	orgUser = identityRepo.GetOrganizationUserById(orgUser.Id)
	if len(orgUser.Roles) != 2 {
		t.Errorf("Expected user to have 2 roles, got %d", len(orgUser.Roles))
	}
}

func TestGetAllUsers(t *testing.T) {
	dbConfig := TestServiceMan.GetDBConfig()
	test_support.ResetDatabase(dbConfig)

	conn := db.GetConnection(dbConfig)
	identityRepo := NewDBIdentityRepo(conn)

	// Create test users
	emails := []string{"alice@example.com", "bob@example.com", "test1@example.com", "test2@example.com"}
	for _, email := range emails {
		_, err := identityRepo.CreateUser(email, "password123")
		if err != nil {
			t.Errorf("Failed to create user with email %s: %s", email, err)
		}
	}

	// Test fetching all users without filters
	allUsers, totalCount := identityRepo.GetAllUsers(10, 0, nil)
	if len(allUsers) != len(emails) {
		t.Errorf("Expected %d users, got %d", len(emails), len(allUsers))
	}
	if totalCount != len(emails) {
		t.Errorf("Expected total count %d, got %d", len(emails), totalCount)
	}

	// Test fetching users with a search filter
	searchTerm := "test"
	filteredUsers, filteredCount := identityRepo.GetAllUsers(10, 0, &searchTerm)
	if len(filteredUsers) != 2 {
		t.Errorf("Expected 2 users matching '%s', got %d", searchTerm, len(filteredUsers))
	}
	if filteredCount != 2 {
		t.Errorf("Expected filtered count 2, got %d", filteredCount)
	}

	// Test pagination
	paginatedUsers, _ := identityRepo.GetAllUsers(2, 0, nil)
	if len(paginatedUsers) != 2 {
		t.Errorf("Expected 2 users in paginated result, got %d", len(paginatedUsers))
	}
	
	// Test pagination offset
	paginatedUsersPage2, _ := identityRepo.GetAllUsers(2, 2, nil)
	if len(paginatedUsersPage2) != 2 {
		t.Errorf("Expected 2 users in second page, got %d", len(paginatedUsersPage2))
	}
	// Ensure we got different results on different pages
	if paginatedUsers[0].Id == paginatedUsersPage2[0].Id {
		t.Errorf("Expected different users on different pages")
	}
}