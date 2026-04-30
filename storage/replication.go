package storage

import (
	"bytes"
	"fmt"
	"io"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/logging"
	"github.com/ericls/imgdd/utils"

	"github.com/google/uuid"
)

var replicationLogger = logging.GetLogger("replication")

// replicateStoredImage copies a file from sourceStorage to targetStorage and creates a
// StoredImage record linked to the source via CopiedFromId.
func replicateStoredImage(
	source *dm.StoredImage,
	sourceStorage Storage,
	targetStorageDef *dm.StorageDefinition,
	targetStorage Storage,
	imageRepo image.ImageRepo,
) (*dm.StoredImage, error) {
	if source.Image == nil {
		return nil, fmt.Errorf("source stored image %s has no associated image", source.Id)
	}
	mimeType := source.Image.MIMEType
	imageId := source.Image.Id

	reader := sourceStorage.GetReader(source.FileIdentifier)
	if reader == nil {
		return nil, fmt.Errorf("could not open source file %s", source.FileIdentifier)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading source file %s: %w", source.FileIdentifier, err)
	}

	fileIdentifier := uuid.New().String() + utils.GetExtFromMIMEType(mimeType)
	if err := targetStorage.Save(bytes.NewReader(data), fileIdentifier, mimeType); err != nil {
		return nil, fmt.Errorf("saving to target storage: %w", err)
	}

	copiedFromId := source.Id
	return imageRepo.CreateStoredImage(imageId, targetStorageDef.Id, fileIdentifier, &copiedFromId)
}

// ReplicateImageToStorageDefinition replicates a single image from the source backend to the
// target backend. Returns an error if the image is not found in source or already exists in target.
func ReplicateImageToStorageDefinition(
	imageId string,
	sourceStorageDefId string,
	targetStorageDefId string,
	storedImageRepo StoredImageRepo,
	imageRepo image.ImageRepo,
	storageDefRepo StorageDefRepo,
) (*dm.StoredImage, error) {
	sourceDef, err := storageDefRepo.GetStorageDefinitionById(sourceStorageDefId)
	if err != nil {
		return nil, fmt.Errorf("getting source storage definition: %w", err)
	}
	targetDef, err := storageDefRepo.GetStorageDefinitionById(targetStorageDefId)
	if err != nil {
		return nil, fmt.Errorf("getting target storage definition: %w", err)
	}

	storedImages, err := storedImageRepo.GetStoredImagesByImageId(imageId)
	if err != nil {
		return nil, fmt.Errorf("getting stored images for image %s: %w", imageId, err)
	}

	var source *dm.StoredImage
	for _, si := range storedImages {
		if si.StorageDefinitionId == sourceStorageDefId {
			source = si
		}
		if si.StorageDefinitionId == targetStorageDefId {
			return nil, fmt.Errorf("image %s already exists in target storage definition %s", imageId, targetStorageDefId)
		}
	}
	if source == nil {
		return nil, fmt.Errorf("image %s not found in source storage definition %s", imageId, sourceStorageDefId)
	}

	sourceStorage, err := GetStorage(sourceDef)
	if err != nil {
		return nil, fmt.Errorf("instantiating source storage: %w", err)
	}
	targetStorage, err := GetStorage(targetDef)
	if err != nil {
		return nil, fmt.Errorf("instantiating target storage: %w", err)
	}

	return replicateStoredImage(source, sourceStorage, targetDef, targetStorage, imageRepo)
}

// BulkReplicateToStorageDefinition replicates all images present in source but missing from
// target. workers controls the number of concurrent copy goroutines; the internal job buffer
// is workers*2. Returns the number of successfully replicated images.
func BulkReplicateToStorageDefinition(
	sourceStorageDefId string,
	targetStorageDefId string,
	storedImageRepo StoredImageRepo,
	imageRepo image.ImageRepo,
	storageDefRepo StorageDefRepo,
	workers int,
) (int, error) {
	sourceDef, err := storageDefRepo.GetStorageDefinitionById(sourceStorageDefId)
	if err != nil {
		return 0, fmt.Errorf("getting source storage definition: %w", err)
	}
	targetDef, err := storageDefRepo.GetStorageDefinitionById(targetStorageDefId)
	if err != nil {
		return 0, fmt.Errorf("getting target storage definition: %w", err)
	}

	sourceStorage, err := GetStorage(sourceDef)
	if err != nil {
		return 0, fmt.Errorf("instantiating source storage: %w", err)
	}
	targetStorage, err := GetStorage(targetDef)
	if err != nil {
		return 0, fmt.Errorf("instantiating target storage: %w", err)
	}

	total, seq, err := storedImageRepo.GetStoredImagesForReplication(sourceStorageDefId, targetStorageDefId)
	if err != nil {
		return 0, fmt.Errorf("preparing replication: %w", err)
	}
	replicationLogger.Info().
		Str("source", sourceDef.Identifier).
		Str("target", targetDef.Identifier).
		Int("total", total).
		Msg("Starting bulk replication")

	if total == 0 {
		return 0, nil
	}
	execute := func(source *dm.StoredImage) error {
		_, err := replicateStoredImage(source, sourceStorage, targetDef, targetStorage, imageRepo)
		return err
	}

	const progressInterval = 100
	succeeded := 0
	failed := 0
	for result := range utils.RunWorkerPoolIter(seq, execute, workers) {
		if result.Err != nil {
			failed++
			replicationLogger.Error().
				Str("stored_image_id", result.Job.Id).
				Err(result.Err).
				Msg("Failed to replicate stored image")
		} else {
			succeeded++
		}
		done := succeeded + failed
		if done%progressInterval == 0 {
			replicationLogger.Info().
				Int("done", done).
				Int("total", total).
				Int("succeeded", succeeded).
				Int("failed", failed).
				Msg("Replication progress")
		}
	}
	replicationLogger.Info().
		Int("total", total).
		Int("succeeded", succeeded).
		Int("failed", failed).
		Msg("Bulk replication complete")
	return succeeded, nil
}
