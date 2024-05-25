package httpserver

import (
	"context"
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

func (cu *HttpContextUserManager) ValidateUserPassword(userId string, suppliedPassword string) bool {
	hashedPassword := cu.identityRepo.GetUserPassword(userId)
	return identity.CheckPasswordHash(suppliedPassword, hashedPassword)
}

type IdentityManager struct {
	IdentityRepo       identity.IdentityRepo
	ContextUserManager identity.ContextUserManager
}

func NewIdentityManager(identityRepo identity.IdentityRepo) *IdentityManager {
	return &IdentityManager{
		IdentityRepo:       identityRepo,
		ContextUserManager: NewContextUserManager("authInfo", identityRepo),
	}
}

func (i *IdentityManager) makeAuthenticationInfoFromRequest(r *http.Request) *identity.AuthenticationInfo {
	authenticatedUserId := GetSessionValue(r, authenticated_user_id_session_key)
	authorizedUserId := GetSessionValue(r, authroized_user_id_session_key)
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

func Authenticate(w http.ResponseWriter, r *http.Request, userId string, organizationUserId string) {
	SetSessionValue(w, r, authenticated_user_id_session_key, userId)
	SetSessionValue(w, r, authroized_user_id_session_key, organizationUserId)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	SetSessionValue(w, r, authenticated_user_id_session_key, "")
	SetSessionValue(w, r, authroized_user_id_session_key, "")
	ClearSession(w, r)
}

func (i *IdentityManager) AuthenticateContext(c context.Context, userId string, organzationUserId string) {
	w := GetResponseWriter(c)
	r := GetRequest(c)
	Authenticate(w, r, userId, organzationUserId)
	user := i.IdentityRepo.GetUserById(userId)
	orgUser := i.IdentityRepo.GetOrganizationUserById(organzationUserId)
	authContext := i.ContextUserManager.GetAuthenticationInfo(c)
	if authContext == nil {
		authContext = &identity.AuthenticationInfo{}
	}
	authContext.AuthenticatedUser.User = user
	authContext.AuthorizedUser.OrganizationUser = orgUser
}

func (i *IdentityManager) LogoutContext(c context.Context) {
	authContext := i.ContextUserManager.GetAuthenticationInfo(c)
	if authContext != nil {
		authContext.AuthenticatedUser.User = nil
		authContext.AuthorizedUser.OrganizationUser = nil
	}
	w := GetResponseWriter(c)
	r := GetRequest(c)
	Logout(w, r)
}
