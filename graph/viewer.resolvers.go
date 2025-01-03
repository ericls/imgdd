package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.62

import (
	"context"
	"imgdd/graph/model"
)

// ID is the resolver for the id field.
func (r *viewerResolver) ID(ctx context.Context, obj *model.Viewer) (string, error) {
	return "viewer", nil
}

// OrganizationUser is the resolver for the organizationUser field.
func (r *viewerResolver) OrganizationUser(ctx context.Context, obj *model.Viewer) (*model.OrganizationUser, error) {
	authInfo := r.ContextUserManager.GetAuthenticationInfo(ctx)
	var orgUser *model.OrganizationUser
	if authInfo != nil && authInfo.AuthorizedUser != nil {
		orgUser = model.FromIdentityOrganizationUser(authInfo.AuthorizedUser.OrganizationUser)
	}
	return orgUser, nil
}

// Viewer returns ViewerResolver implementation.
func (r *Resolver) Viewer() ViewerResolver { return &viewerResolver{r} }

type viewerResolver struct{ *Resolver }
