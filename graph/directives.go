package graph

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func IsSiteOwner(r *Resolver) func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		currentUser := r.ContextUserManager.GetCurrentOrganizationUser(ctx)
		for _, role := range currentUser.Roles {
			if role.Key == "site_owner" {
				return next(ctx)
			}
		}
		return nil, fmt.Errorf("not site owner")
	}
}
