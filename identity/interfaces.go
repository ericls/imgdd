package identity

import (
	"context"
	dm "imgdd/domainmodels"
)

type IdentityRepo interface {
	CreateUserWithOrganization(ctx context.Context, email string, organizationName string, password string) (*dm.OrganizationUser, error)

	GetUserById(ctx context.Context, id string) *dm.User
	GetUsersByIds(ctx context.Context, ids []string) []*dm.User
	GetUserByEmail(ctx context.Context, email string) *dm.User

	// If maybeOrganizationId is not empty, it will be used as the organization ID.
	// Otherwise, the implementation should guess the best organization ID for the user.
	// E.g. if the user is a member of only one organization, that organization ID should be used.
	// If the user is a member of multiple organizations, the implementation should return nil.
	GetOrganizationForUser(ctx context.Context, userId string, maybeOrganizationId string) (*dm.Organization, *dm.OrganizationUser)

	GetOrganizationUsersByIds(ctx context.Context, ids []string) []*dm.OrganizationUser
	GetOrganizationUserById(ctx context.Context, id string) *dm.OrganizationUser

	GetUserPassword(ctx context.Context, id string) string
}

type ContextUserManager interface {
	GetAuthenticationInfo(c context.Context) *AuthenticationInfo
	WithAuthenticationInfo(c context.Context, authenticationInfo *AuthenticationInfo) context.Context
	ValidateUserPassword(userId string, suppliedPassword string) bool
}
