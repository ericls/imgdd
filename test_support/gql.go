package test_support

import (
	"context"
	"imgdd/identity"
	"net/http"
)

type TestIdentityManager struct {
	IdentityRepo       identity.IdentityRepo
	ContextUserManager *TestContextUserManager
}

func NewTestIdentityManager(identityRepo identity.IdentityRepo) *TestIdentityManager {
	return &TestIdentityManager{
		IdentityRepo:       identityRepo,
		ContextUserManager: NewTestContextUserManager("authInfo", identityRepo),
	}
}

func (i *TestIdentityManager) AuthenticateContext(c context.Context, userId string, organizationUserId string) {
	user := i.IdentityRepo.GetUserById(userId)
	orgUser := i.IdentityRepo.GetOrganizationUserById(organizationUserId)
	authContext := i.ContextUserManager.GetAuthenticationInfo(c)
	if authContext == nil {
		authContext = &identity.AuthenticationInfo{
			AuthenticatedUser: &identity.AuthenticatedUser{
				User: user,
			},
			AuthorizedUser: &identity.AuthorizedUser{
				OrganizationUser: orgUser,
			},
		}
	}
	authContext.AuthenticatedUser.User = user
	authContext.AuthorizedUser.OrganizationUser = orgUser
}

func (i *TestIdentityManager) LogoutContext(c context.Context) {
	authContext := i.ContextUserManager.GetAuthenticationInfo(c)
	if authContext != nil {
		authContext.AuthenticatedUser.User = nil
		authContext.AuthorizedUser.OrganizationUser = nil
	}
}

func (i *TestIdentityManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authentication_info := identity.AuthenticationInfo{
			AuthenticatedUser: &identity.AuthenticatedUser{
				User: nil,
			},
			AuthorizedUser: &identity.AuthorizedUser{
				OrganizationUser: nil,
			},
		}
		newContext := i.ContextUserManager.WithAuthenticationInfo(r.Context(), &authentication_info)
		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}
