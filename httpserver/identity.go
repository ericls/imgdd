package httpserver

import (
	"context"
	"imgdd/httpserver/persister"
	"imgdd/identity"
	"net/http"
)

const authenticated_user_id_session_key = "authenticated_user_id"
const authroized_user_id_session_key = "authorized_user_id"

type authContextKey string

type HttpContextUserManager struct {
	contextKey   authContextKey
	identityRepo identity.IdentityRepo
}

func NewContextUserManager(contextKey string, identityRepo identity.IdentityRepo) *HttpContextUserManager {
	return &HttpContextUserManager{
		contextKey:   authContextKey(contextKey),
		identityRepo: identityRepo,
	}
}

func (cu *HttpContextUserManager) GetAuthenticationInfo(c context.Context) *identity.AuthenticationInfo {
	v, ok := c.Value(cu.contextKey).(*identity.AuthenticationInfo)
	if v == nil || !ok {
		return nil
	}
	return v
}

func (cu *HttpContextUserManager) WithAuthenticationInfo(c context.Context, authenticationInfo *identity.AuthenticationInfo) context.Context {
	return context.WithValue(c, cu.contextKey, authenticationInfo)
}

func (cu *HttpContextUserManager) SetAuthenticationInfo(c context.Context, authenticationInfo *identity.AuthenticationInfo) {
	existing := cu.GetAuthenticationInfo(c)
	if existing == nil {
		return
	}
	existing.AuthenticatedUser = authenticationInfo.AuthenticatedUser
	existing.AuthorizedUser = authenticationInfo.AuthorizedUser
}

type IdentityManager struct {
	IdentityRepo       identity.IdentityRepo
	ContextUserManager identity.ContextUserManager
	Persister          persister.Persister
}

func NewIdentityManager(identityRepo identity.IdentityRepo, persister persister.Persister) *IdentityManager {
	return &IdentityManager{
		IdentityRepo:       identityRepo,
		ContextUserManager: NewContextUserManager("authInfo", identityRepo),
		Persister:          persister,
	}
}

func (i *IdentityManager) makeAuthenticationInfoFromRequest(r *http.Request) *identity.AuthenticationInfo {
	authenticatedUserId, err := i.Persister.Get(r, authenticated_user_id_session_key)
	if err != nil {
		return nil
	}
	authorizedUserId, err := i.Persister.Get(r, authroized_user_id_session_key)
	if err != nil {
		return nil
	}
	authenticatedUser := identity.AuthenticatedUser{}
	authorizedUser := identity.AuthorizedUser{}
	if authenticatedUserId != "" {
		authenticatedUser.User = i.IdentityRepo.GetUserById(authenticatedUserId)
	}
	if authorizedUserId != "" {
		authorizedUser.OrganizationUser = i.IdentityRepo.GetOrganizationUserById(authorizedUserId)
	}
	authenticationInfo := identity.AuthenticationInfo{
		AuthenticatedUser: &authenticatedUser,
		AuthorizedUser:    &authorizedUser,
	}
	return &authenticationInfo
}

func (i *IdentityManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authentication_info := i.makeAuthenticationInfoFromRequest(r)
		newContext := i.ContextUserManager.WithAuthenticationInfo(r.Context(), authentication_info)
		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}

func (i *IdentityManager) Authenticate(w http.ResponseWriter, r *http.Request, userId string, organizationUserId string) {
	i.Persister.Set(w, r, authenticated_user_id_session_key, userId)
	i.Persister.Set(w, r, authroized_user_id_session_key, organizationUserId)
}

func (i *IdentityManager) Logout(w http.ResponseWriter, r *http.Request) {
	i.Persister.Set(w, r, authenticated_user_id_session_key, "")
	i.Persister.Set(w, r, authroized_user_id_session_key, "")
	i.Persister.Clear(w, r)
}

func (i *IdentityManager) AuthenticateContext(c context.Context, userId string, organzationUserId string) {
	w := GetResponseWriter(c)
	r := GetRequest(c)
	// Set authentication info on the session for next requests
	i.Authenticate(w, r, userId, organzationUserId)
	user := i.IdentityRepo.GetUserById(userId)
	orgUser := i.IdentityRepo.GetOrganizationUserById(organzationUserId)
	// Set authentication info on the context for this request
	i.ContextUserManager.SetAuthenticationInfo(c, &identity.AuthenticationInfo{
		AuthenticatedUser: &identity.AuthenticatedUser{User: user},
		AuthorizedUser:    &identity.AuthorizedUser{OrganizationUser: orgUser},
	})
}

func (i *IdentityManager) LogoutContext(c context.Context) {
	authContext := i.ContextUserManager.GetAuthenticationInfo(c)
	if authContext != nil {
		authContext.AuthenticatedUser.User = nil
		authContext.AuthorizedUser.OrganizationUser = nil
	}
	w := GetResponseWriter(c)
	r := GetRequest(c)
	i.Logout(w, r)
}
