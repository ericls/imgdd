package identity

import (
	dm "imgdd/domainmodels"
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
