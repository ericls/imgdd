package test_support

import (
	dm "imgdd/domainmodels"

	"github.com/google/uuid"
)

type TestIdentityRepo struct {
	users             []*dm.User
	organizations     []*dm.Organization
	organizationUsers []*dm.OrganizationUser
	userPasswords     map[string]string
}

func (repo *TestIdentityRepo) Reset() {
	repo.users = make([]*dm.User, 0)
	repo.organizations = make([]*dm.Organization, 0)
	repo.userPasswords = make(map[string]string)
}

func (repo *TestIdentityRepo) GetUserById(id string) *dm.User {
	for _, user := range repo.users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

func (repo *TestIdentityRepo) GetUserByEmail(email string) *dm.User {
	for _, user := range repo.users {
		if user.Email == email {
			return user
		}
	}
	return nil
}

func (repo *TestIdentityRepo) GetOrganizationUserById(id string) *dm.OrganizationUser {
	for _, orgUser := range repo.organizationUsers {
		if orgUser.Id == id {
			return orgUser
		}
	}
	return nil
}

func (repo *TestIdentityRepo) GetOrganizationUsersByIds(ids []string) []*dm.OrganizationUser {
	var orgUsers []*dm.OrganizationUser
	for _, id := range ids {
		orgUsers = append(orgUsers, repo.GetOrganizationUserById(id))
	}
	return orgUsers
}

func (repo *TestIdentityRepo) CreateUser(email string, orangizationId string, password string) (*dm.User, error) {
	id := uuid.New().String()
	user := &dm.User{
		Id:             id,
		Email:          email,
		OrganizationId: orangizationId,
	}
	repo.users = append(repo.users, user)
	repo.setPassword(id, password)
	return user, nil
}

func (repo *TestIdentityRepo) CreateUserWithOrganization(email string, organizationName string, password string) (*dm.OrganizationUser, error) {
	org := &dm.Organization{
		Id:          uuid.New().String(),
		DisplayName: organizationName,
		Slug:        uuid.New().String(),
	}
	user, err := repo.CreateUser(email, org.Id, password)
	if err != nil {
		return nil, err
	}
	organizationUser := &dm.OrganizationUser{
		Id:           uuid.New().String(),
		Organization: org,
		User:         user,
		Roles: []*dm.Role{
			{
				Id:   uuid.New().String(),
				Name: "admin",
			},
		},
	}
	repo.organizations = append(repo.organizations, org)
	repo.organizationUsers = append(repo.organizationUsers, organizationUser)
	return organizationUser, nil
}

func (repo *TestIdentityRepo) GetOrganizationById(id string) *dm.Organization {
	for _, org := range repo.organizations {
		if org.Id == id {
			return org
		}
	}
	return nil
}

func (repo *TestIdentityRepo) GetUsersByIds(ids []string) []*dm.User {
	var users []*dm.User
	for _, id := range ids {
		users = append(users, repo.GetUserById(id))
	}
	return users
}

func (repo *TestIdentityRepo) GetUserPassword(id string) string {
	return repo.userPasswords[id]
}

func (reop *TestIdentityRepo) GetOrganizationsByIds(ids []string) []*dm.Organization {
	var orgs []*dm.Organization
	for _, id := range ids {
		orgs = append(orgs, reop.GetOrganizationById(id))
	}
	return orgs
}

func (repo *TestIdentityRepo) setPassword(userId string, password string) {
	if repo.userPasswords == nil {
		repo.userPasswords = make(map[string]string)
	}
	repo.userPasswords[userId] = password
}

type orgIdOrgUserIdPair struct {
	orgId     string
	orgUserId string
}

func (repo *TestIdentityRepo) getAllOrganizationsForUser(userId string) []*orgIdOrgUserIdPair {
	var orgs []*orgIdOrgUserIdPair
	for _, orgUser := range repo.organizationUsers {
		if orgUser.User.Id == userId {
			orgs = append(orgs, &orgIdOrgUserIdPair{
				orgId:     orgUser.Organization.Id,
				orgUserId: orgUser.Id,
			})
		}
	}
	return orgs
}

func (repo *TestIdentityRepo) GetOrganizationForUser(userId string, maybeOrganizationId string) (*dm.Organization, *dm.OrganizationUser) {
	possibleOrgs := repo.getAllOrganizationsForUser(userId)
	if maybeOrganizationId == "" {
		if len(possibleOrgs) == 1 {
			orgId := possibleOrgs[0].orgId
			orgUserId := possibleOrgs[0].orgUserId
			return repo.GetOrganizationById(orgId), repo.GetOrganizationUserById(orgUserId)
		}
		return nil, nil
	}
	for _, orgOrgUser := range possibleOrgs {
		orgId := orgOrgUser.orgId
		orgUserId := orgOrgUser.orgUserId
		if orgId == maybeOrganizationId {
			return repo.GetOrganizationById(orgId), repo.GetOrganizationUserById(orgUserId)
		}
	}
	return nil, nil
}
