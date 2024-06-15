package identity

import (
	"context"
	dm "imgdd/domainmodels"
)

type IdentityRepo interface {
	CreateUserWithOrganization(email string, organizationName string, password string) (*dm.OrganizationUser, error)
	CreateUser(email string, password string) (*dm.User, error)

	GetUserById(id string) *dm.User
	GetUsersByIds(ids []string) []*dm.User
	GetUserByEmail(email string) *dm.User

	// If maybeOrganizationId is not empty, it will be used as the organization ID.
	// Otherwise, the implementation should guess the best organization ID for the user.
	// E.g. if the user is a member of only one organization, that organization ID should be used.
	// If the user is a member of multiple organizations, the implementation should return nil.
	GetOrganizationForUser(userId string, maybeOrganizationId string) (*dm.Organization, *dm.OrganizationUser)

	GetOrganizationUsersByIds(ids []string) []*dm.OrganizationUser
	GetOrganizationUserById(id string) *dm.OrganizationUser

	GetUserPassword(id string) string
}

type ContextUserManager interface {
	GetAuthenticationInfo(c context.Context) *AuthenticationInfo
	WithAuthenticationInfo(c context.Context, authenticationInfo *AuthenticationInfo) context.Context
	ValidateUserPassword(userId string, suppliedPassword string) bool
	GetCurrentOrganizationUser(c context.Context) *dm.OrganizationUser
}
