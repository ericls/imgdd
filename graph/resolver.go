package graph

import (
	"context"
	"imgdd/identity"
	"imgdd/storage"
)

//go:generate go run github.com/99designs/gqlgen generate
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	IdentityRepo       identity.IdentityRepo
	StorageRepo        storage.StorageRepo
	ContextUserManager identity.ContextUserManager
	LoginFn            func(c context.Context, userId string, organizationUserId string)
	LogoutFn           func(c context.Context)
}
