package identity

import (
	"context"
	"imgdd/db"
	dm "imgdd/domainmodels"
	"testing"
)

func TestUserManager(t *testing.T) {
	conn := db.GetConnection(&TEST_DB_CONF)
	identityRepo := DBIdentityRepo{conn}
	emailAddr := "test@home.arpa"
	var assertUser = func(ou *dm.OrganizationUser, err error) {
		if err != nil {
			t.Errorf("Error creating user: %s", err)
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
	orgUser, err := identityRepo.CreateUserWithOrganization(context.Background(), emailAddr, "test", "123")
	assertUser(orgUser, err)
	orgUser = identityRepo.GetOrganizationUserById(context.Background(), orgUser.Id)
	assertUser(orgUser, nil)
	orgUsers := identityRepo.GetOrganizationUsersByIds(context.Background(), []string{orgUser.Id})
	if len(orgUsers) != 1 {
		t.Errorf("Failed to fetch user in bulk")
	} else {
		assertUser(orgUsers[0], nil)
	}
	assertUser(orgUser, nil)
	identityRepo.AddRoleToOrganizationUser(context.Background(), orgUser.Id, "member")
	orgUser = identityRepo.GetOrganizationUserById(context.Background(), orgUser.Id)
	if len(orgUser.Roles) != 2 {
		t.Errorf("Expected user to have 2 roles, got %d", len(orgUser.Roles))
	}
}
