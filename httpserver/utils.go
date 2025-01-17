package httpserver

import (
	"github.com/ericls/imgdd/graph"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
)

type ContextKey string

func NewGqlResolver(identityManager *IdentityManager, storageDefRepo storage.StorageDefRepo, imageRepo image.ImageRepo, imageDomain string) *graph.Resolver {
	return &graph.Resolver{
		IdentityRepo:       identityManager.IdentityRepo,
		StorageDefRepo:     storageDefRepo,
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
