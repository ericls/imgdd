package identity

import (
	"testing"

	"github.com/ericls/imgdd/db"
	dm "github.com/ericls/imgdd/domainmodels"
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
