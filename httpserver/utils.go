package httpserver

import (
	"imgdd/graph"
	"imgdd/storage"
)

type ContextKey string

func NewGqlResolver(identityManager *IdentityManager, storageRepo storage.StorageRepo) *graph.Resolver {
	return &graph.Resolver{
		IdentityRepo:       identityManager.IdentityRepo,
		StorageRepo:        storageRepo,
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
