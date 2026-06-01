package editing

import (
	"fmt"
	"io"
	"sort"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/storage"
)

// NewFetchImageFunc creates a FetchImageFunc that reads image bytes from storage.
func NewFetchImageFunc(
	storedImageRepo storage.StoredImageRepo,
	storageDefRepo storage.StorageDefRepo,
) FetchImageFunc {
	return func(imageId string) ([]byte, error) {
		storedImages, err := storedImageRepo.GetStoredImagesByImageId(imageId)
		if err != nil {
			return nil, fmt.Errorf("failed to get stored images for %s: %w", imageId, err)
		}
		if len(storedImages) == 0 {
			return nil, fmt.Errorf("no stored images found for %s", imageId)
		}

		// Collect storage definition IDs
		defIds := make([]string, 0, len(storedImages))
		for _, si := range storedImages {
			defIds = append(defIds, si.StorageDefinitionId)
		}
		defs, err := storageDefRepo.GetStorageDefinitionsByIds(defIds)
		if err != nil {
			return nil, fmt.Errorf("failed to get storage definitions for %s: %w", imageId, err)
		}

		// Build lookup and filter enabled
		defMap := make(map[string]*dm.StorageDefinition)
		for _, d := range defs {
			if d != nil && d.IsEnabled {
				defMap[d.Id] = d
			}
		}

		type candidate struct {
			si  *dm.StoredImage
			def *dm.StorageDefinition
		}
		var candidates []candidate
		for _, si := range storedImages {
			if d, ok := defMap[si.StorageDefinitionId]; ok {
				candidates = append(candidates, candidate{si, d})
			}
		}
		if len(candidates) == 0 {
			return nil, fmt.Errorf("no enabled storage backends for image %s", imageId)
		}

		sort.SliceStable(candidates, func(i, j int) bool {
			return candidates[i].def.Priority < candidates[j].def.Priority
		})

		best := candidates[0]
		storageInstance, err := storage.GetStorage(best.def)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage instance: %w", err)
		}

		reader := storageInstance.GetReader(best.si.FileIdentifier)
		if reader == nil {
			return nil, fmt.Errorf("failed to get reader for file %s", best.si.FileIdentifier)
		}
		defer reader.Close()

		const maxImageBytes = 10 * 1024 * 1024 // 10 MB
		limitedReader := io.LimitReader(reader, maxImageBytes+1)
		data, err := io.ReadAll(limitedReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read image data: %w", err)
		}
		if len(data) > maxImageBytes {
			return nil, fmt.Errorf("image exceeds maximum size of %d bytes", maxImageBytes)
		}
		return data, nil
	}
}
