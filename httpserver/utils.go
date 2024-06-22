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

func NewGraphConfig(resolver *graph.Resolver) graph.Config {
	config := graph.Config{
		Resolvers: resolver,
	}
	config.Directives.IsSiteOwner = graph.IsSiteOwner(resolver)
	return config
}
