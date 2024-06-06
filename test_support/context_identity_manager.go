package test_support

import (
	"context"
	"imgdd/identity"
)

type authContextKey string

type TestContextUserManager struct {
	contextKey   authContextKey
	identityRepo *TestIdentityRepo
}

func NewTestContextUserManager(contextKey string, identityRepo *TestIdentityRepo) *TestContextUserManager {
	return &TestContextUserManager{
		contextKey:   authContextKey(contextKey),
		identityRepo: identityRepo,
	}
}

func NewContextUserManager(contextKey string, identityRepo *TestIdentityRepo) *TestContextUserManager {
	return &TestContextUserManager{
		contextKey:   authContextKey(contextKey),
		identityRepo: identityRepo,
	}
}

func (cu *TestContextUserManager) GetAuthenticationInfo(c context.Context) *identity.AuthenticationInfo {
	v, ok := c.Value(cu.contextKey).(*identity.AuthenticationInfo)
	if v == nil || !ok {
		return nil
	}
	return v
}

func (cu *TestContextUserManager) WithAuthenticationInfo(c context.Context, authenticationInfo *identity.AuthenticationInfo) context.Context {
	return context.WithValue(c, cu.contextKey, authenticationInfo)
}

func (cu *TestContextUserManager) ValidateUserPassword(userId string, suppliedPassword string) bool {
	hashedPassword := cu.identityRepo.GetUserPassword(context.Background(), userId) // not really hashsed in tests
	return hashedPassword == suppliedPassword
}
