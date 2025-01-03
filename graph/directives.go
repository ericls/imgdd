package graph

import (
	"context"
	"fmt"
	"imgdd/identity"

	"github.com/99designs/gqlgen/graphql"
)

func IsSiteOwner(r *Resolver) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		currentUser := identity.GetCurrentOrganizationUser(r.ContextUserManager, ctx)
		if currentUser == nil {
			return nil, fmt.Errorf("not authenticated")
		}
		if currentUser.IsSiteOwner() {
			return next(ctx)
		}
		return nil, fmt.Errorf("not site owner")
	}
}
