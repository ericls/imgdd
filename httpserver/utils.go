package httpserver

import (
	"context"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/graph"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
)

type ContextKey string

func NewGqlResolver(
	identityManager *IdentityManager,
	storageDefRepo storage.StorageDefRepo,
	storedImageRepo storage.StoredImageRepo,
	imageRepo image.ImageRepo,
	imageRelRepo image.ImageRelationshipRepo,
	imageDomain string,
	defaultURLFormat domainmodels.ImageURLFormat,
	getEmailBackend func(c context.Context) email.EmailBackend,
	secretKey string,
	captchaClient captcha.CaptchaClient,
	allowNewUser bool,
) *graph.Resolver {
	return &graph.Resolver{
		IdentityRepo:       identityManager.IdentityRepo,
		StorageDefRepo:     storageDefRepo,
		StoredImageRepo:    storedImageRepo,
		ImageRepo:          imageRepo,
		ImageRelRepo:       imageRelRepo,
		ContextUserManager: identityManager.ContextUserManager,
		LoginFn:            identityManager.AuthenticateContext,
		LogoutFn:           identityManager.LogoutContext,
		ImageDomain:        imageDomain,
		DefaultURLFormat:   defaultURLFormat,
		IsHttps:            IsHttps,
		GetBaseURL:         GetBaseURL,
		GetEmailBackend:    getEmailBackend,
		SecretKey:          secretKey,
		CaptchaClient:      captchaClient,
		AllowNewUser:       allowNewUser,
	}
}

func NewGraphConfig(resolver *graph.Resolver) graph.Config {
	config := graph.Config{
		Resolvers: resolver,
	}
	config.Directives.IsSiteOwner = graph.IsSiteOwner(resolver)
	config.Directives.IsAuthenticated = graph.IsAuthenticated(resolver)
	config.Directives.CaptchaProtected = graph.CaptchaProtected(resolver)
	return config
}
