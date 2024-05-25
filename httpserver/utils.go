package httpserver

import (
	"imgdd/graph"
)

type ContextKey string

func NewGqlResolver(identityManager *IdentityManager) *graph.Resolver {
	return &graph.Resolver{
		IdentityRepo:       identityManager.IdentityRepo,
		ContextUserManager: identityManager.ContextUserManager,
		LoginFn:            identityManager.AuthenticateContext,
		LogoutFn:           identityManager.LogoutContext,
	}
}
