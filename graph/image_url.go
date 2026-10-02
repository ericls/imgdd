package graph

import (
	"context"

	"github.com/ericls/imgdd/domainmodels"
)

// externalImageIdentifiers returns the stored copies of an image, used to build direct image URLs.
func (r *Resolver) externalImageIdentifiers(ctx context.Context, imageId string) ([]*domainmodels.ExternalImageIdentifier, error) {
	storedImages, err := LoadersFor(ctx).StoredImagesByImageIdsLoader.Load(ctx, imageId)
	if err != nil {
		return nil, err
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
	return externalIdentifiers, nil
}
