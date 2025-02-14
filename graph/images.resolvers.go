package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"
	"fmt"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/graph/model"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
)

// URL is the resolver for the url field.
func (r *imageResolver) URL(ctx context.Context, obj *model.Image) (string, error) {
	image := domainmodels.Image{
		Id:         obj.ID,
		Identifier: obj.Identifier,
		MIMEType:   obj.MIMEType,
	}
	loader := LoadersFor(ctx).StoredImagesByImageIdsLoader
	storedImages, err := loader.Load(ctx, obj.ID)
	if err != nil {
		return "", err
	}
	var externalIdentifiers []*domainmodels.ExternalImageIdentifier
	for _, storedImage := range storedImages {
		if storedImage == nil || storedImage.StorageDefinition == nil {
			continue
		}
		externalIdentifiers = append(externalIdentifiers, &domainmodels.ExternalImageIdentifier{
			StorageDefinitionIdentifier: storedImage.StorageDefinition.Identifier,
			FileIdentifier:              storedImage.FileIdentifier,
		})
	}
	return image.GetURL(r.ImageDomain, r.IsHttps(ctx), externalIdentifiers, r.DefaultURLFormat), nil
}

// Root is the resolver for the root field.
func (r *imageResolver) Root(ctx context.Context, obj *model.Image) (*model.Image, error) {
	// TODO: Implement this
	return nil, nil
}

// Revisions is the resolver for the revisions field.
func (r *imageResolver) Revisions(ctx context.Context, obj *model.Image) ([]*model.Image, error) {
	// TODO: Implement this method after we support revisions.
	return []*model.Image{}, nil
}

// StoredImages is the resolver for the storedImages field.
func (r *imageResolver) StoredImages(ctx context.Context, obj *model.Image) ([]*model.StoredImage, error) {
	loader := LoadersFor(ctx).StoredImagesByImageIdsLoader
	return loader.Load(ctx, obj.ID)
}

// DeleteImage is the resolver for the deleteImage field.
func (r *mutationResolver) DeleteImage(ctx context.Context, input model.DeleteImageInput) (*model.DeleteImageResult, error) {
	currentUser := identity.GetCurrentOrganizationUser(r.ContextUserManager, ctx)
	if img, err := r.ImageRepo.GetImageById(input.ID); err != nil {
		return nil, err
	} else if img == nil {
		return nil, fmt.Errorf("image not found")
	} else {
		createdBy := r.IdentityRepo.GetOrganizationUserById(img.CreatedById)
		if !currentUser.CanManage(createdBy) {
			return nil, fmt.Errorf("unauthorized")
		}
		if err := r.ImageRepo.DeleteImageById(input.ID); err != nil {
			return nil, err
		} else {
			return &model.DeleteImageResult{ID: &input.ID}, nil
		}
	}
}

// Images is the resolver for the images field.
func (r *viewerResolver) Images(ctx context.Context, obj *model.Viewer, orderBy *model.ImageOrderByInput, filters *model.ImageFilterInput, after *string, before *string) (*model.ImagesResult, error) {
	currentUser := identity.GetCurrentOrganizationUser(r.ContextUserManager, ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	if after != nil && before != nil {
		return nil, fmt.Errorf("only one of after or before can be specified")
	}
	if filters == nil {
		filters = &model.ImageFilterInput{}
	}
	if !currentUser.IsSiteOwner() {
		if filters.CreatedBy == nil {
			filters.CreatedBy = &currentUser.Id
		}
	}

	if filters.CreatedBy != nil {
		filterByUser := r.IdentityRepo.GetOrganizationUserById(*filters.CreatedBy)
		if filterByUser == nil {
			return nil, fmt.Errorf("unathorized")
		}
		if !currentUser.CanManage(filterByUser) {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	paginator := model.MakeImagePaginator(orderBy, filters)
	listImagesFilters := image.FromPaginationFilter(paginator.Filter)
	count, err := r.ImageRepo.CountImages(listImagesFilters)
	if err != nil {
		return nil, err
	}

	if err = paginator.ContributeCursorStringToFilter(after, true); err != nil {
		return nil, err
	}
	if err = paginator.ContributeCursorStringToFilter(before, false); err != nil {
		return nil, err
	}

	listImagesFiltersWithCursor := image.FromPaginationFilter(paginator.Filter)
	var listImagesOrdering *image.ListImagesOrdering
	if before != nil {
		listImagesOrdering = image.FromPaginationOrder(paginator.Order.Reverse())
	} else {
		listImagesOrdering = image.FromPaginationOrder(paginator.Order)
	}

	listImageResult, err := r.ImageRepo.ListImages(listImagesFilters, listImagesFiltersWithCursor, listImagesOrdering, before != nil)
	if err != nil {
		return nil, err
	}
	cursorEncoder := model.CursorEncoder(func(i *domainmodels.Image) string {
		return paginator.Order.EncodeCursor(i)
	})
	result := model.FromListImageResult(&listImageResult, count, cursorEncoder)
	return result, nil
}

// Image returns ImageResolver implementation.
func (r *Resolver) Image() ImageResolver { return &imageResolver{r} }

type imageResolver struct{ *Resolver }
