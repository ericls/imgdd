package identity

import (
	"context"

	dm "github.com/ericls/imgdd/domainmodels"
)

type AuthenticatedUser struct {
	User *dm.User
}

type AuthorizedUser struct {
	OrganizationUser *dm.OrganizationUser
}

type AuthenticationInfo struct {
	AuthenticatedUser *AuthenticatedUser
	AuthorizedUser    *AuthorizedUser
}

func GetCurrentOrganizationUser(cu ContextUserManager, c context.Context) *dm.OrganizationUser {
	authInfo := cu.GetAuthenticationInfo(c)
	if authInfo == nil {
		return nil
	}
	authUser := authInfo.AuthorizedUser
	if authUser == nil {
		return nil
	}
	return authUser.OrganizationUser
}
