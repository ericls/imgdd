package httpserver

import (
	"imgdd/graph"
	"imgdd/image"
	"imgdd/storage"
)

type ContextKey string

func NewGqlResolver(identityManager *IdentityManager, storageRepo storage.StorageRepo, imageRepo image.ImageRepo, imageDomain string) *graph.Resolver {
	return &graph.Resolver{
		IdentityRepo:       identityManager.IdentityRepo,
		StorageRepo:        storageRepo,
		ImageRepo:          imageRepo,
		ContextUserManager: identityManager.ContextUserManager,
		LoginFn:            identityManager.AuthenticateContext,
		LogoutFn:           identityManager.LogoutContext,
		ImageDomain:        imageDomain,
		IsHttps:            IsHttps,
	}
}

func NewGraphConfig(resolver *graph.Resolver) graph.Config {
	config := graph.Config{
		Resolvers: resolver,
	}
	config.Directives.IsSiteOwner = graph.IsSiteOwner(resolver)
	config.Directives.IsAuthenticated = graph.IsAuthenticated(resolver)
	return config
}
