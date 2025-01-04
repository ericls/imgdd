package graph

import (
	"context"
	"imgdd/identity"
	"imgdd/image"
	"imgdd/storage"
)

//go:generate go run github.com/99designs/gqlgen generate
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	IdentityRepo       identity.IdentityRepo
	StorageRepo        storage.StorageRepo
	ImageRepo          image.ImageRepo
	ContextUserManager identity.ContextUserManager
	LoginFn            func(c context.Context, userId string, organizationUserId string)
	LogoutFn           func(c context.Context)
	ImageDomain        string
	IsHttps            func(c context.Context) bool
}
