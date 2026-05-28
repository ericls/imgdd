package graph_test

import (
	"crypto/md5"
	"fmt"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/require"
)

func tAvatarURLIsGravatarURL(t *testing.T, tc *TestContext) {
	orgUser := tc.forceAuthenticate()

	var resp struct {
		Viewer *struct {
			OrganizationUser *struct {
				User struct {
					AvatarURL string
				}
			}
		}
	}
	err := tc.client.Post(`
	query {
		viewer {
			organizationUser {
				user {
					avatarUrl
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer.OrganizationUser)

	email := strings.ToLower(strings.TrimSpace(orgUser.User.Email))
	hash := md5.Sum([]byte(email))
	expectedURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=identicon&s=80", hash)

	require.Equal(t, expectedURL, resp.Viewer.OrganizationUser.User.AvatarURL)
}

func tOrganizationUserByIDSiteOwnerAccess(t *testing.T, tc *TestContext) {
	target := tc.forceAuthenticate()
	tc.forceAuthenticate(asSiteOwner)

	var resp struct {
		Viewer *struct {
			OrganizationUserByID *struct {
				ID   string
				User struct {
					ID    string
					Email string
				}
			}
		}
	}
	err := tc.client.Post(`
	query($id: ID!) {
		viewer {
			organizationUserById(id: $id) {
				id
				user {
					id
					email
				}
			}
		}
	}`, &resp, client.Var("id", target.Id))
	require.NoError(t, err)
	require.NotNil(t, resp.Viewer.OrganizationUserByID)
	require.Equal(t, target.Id, resp.Viewer.OrganizationUserByID.ID)
	require.Equal(t, target.User.Email, resp.Viewer.OrganizationUserByID.User.Email)
}

func tOrganizationUserByIDNonSiteOwnerDenied(t *testing.T, tc *TestContext) {
	target := tc.forceAuthenticate()
	tc.forceAuthenticate() // regular user

	var resp struct {
		Viewer *struct {
			OrganizationUserByID *struct{ ID string }
		}
	}
	err := tc.client.Post(`
	query($id: ID!) {
		viewer {
			organizationUserById(id: $id) {
				id
			}
		}
	}`, &resp, client.Var("id", target.Id))
	require.Error(t, err)
}

func TestUserResolvers(t *testing.T) {
	tc := newTestContext(t)
	tc.runTestCases(
		tAvatarURLIsGravatarURL,
		tOrganizationUserByIDSiteOwnerAccess,
		tOrganizationUserByIDNonSiteOwnerDenied,
	)
}
