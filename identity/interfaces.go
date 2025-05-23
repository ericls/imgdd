package identity

import (
	"context"

	dm "github.com/ericls/imgdd/domainmodels"
)

type IdentityRepo interface {
	CreateUserWithOrganization(email string, organizationName string, password string) (*dm.OrganizationUser, error)
	CreateUser(email string, password string) (*dm.User, error)
	AddRoleToOrganizationUser(organizationUserId string, roleKey string) error

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
	UpdateUserPassword(id string, password string) error

	// GetAllUsers returns a paginated list of users with optional search criteria.
	GetAllUsers(limit int, offset int, search *string) ([]*dm.User, int)
}

type ContextUserManager interface {
	// GetAuthenticationInfo returns the authentication info from the given context.
	GetAuthenticationInfo(c context.Context) *AuthenticationInfo
	// WithAuthenticationInfo returns a new context with the given authentication info.
	WithAuthenticationInfo(c context.Context, authenticationInfo *AuthenticationInfo) context.Context
	// SetAuthenticationInfo sets the given authentication info to the given context.
	SetAuthenticationInfo(c context.Context, authenticationInfo *AuthenticationInfo)
}
