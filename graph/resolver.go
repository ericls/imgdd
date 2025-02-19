package graph

import (
	"context"
	"net/url"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
)

//go:generate go run github.com/99designs/gqlgen generate
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	IdentityRepo       identity.IdentityRepo
	StorageDefRepo     storage.StorageDefRepo
	ImageRepo          image.ImageRepo
	ContextUserManager identity.ContextUserManager
	LoginFn            func(c context.Context, userId string, organizationUserId string)
	LogoutFn           func(c context.Context)
	ImageDomain        string
	DefaultURLFormat   domainmodels.ImageURLFormat
	IsHttps            func(c context.Context) bool
	GetBaseURL         func(c context.Context) *url.URL
	GetEmailBackend    func(c context.Context) email.EmailBackend
	SecretKey          string
	CaptchaClient      captcha.CaptchaClient
	AllowNewUser       bool
}
