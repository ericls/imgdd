package storage

import (
	"fmt"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/logging"
	"github.com/ericls/imgdd/utils"
)

var logger = logging.GetLogger("cleanup")

func CleanupStoredImage(storedImageRepo StoredImageRepo, storageDefRepo StorageDefRepo) (int, error) {
	storedImagesToDelete, err := storedImageRepo.GetStoredImagesToDelete()
	if err != nil {
		return 0, err
	}
	if len(storedImagesToDelete) == 0 {
		return 0, nil
	}
	storageDefIds := map[string]struct{}{}
	for _, si := range storedImagesToDelete {
		storageDefIds[si.StorageDefinitionId] = struct{}{}
	}
	storageDefIdsList := make([]string, 0, len(storageDefIds))
	for id := range storageDefIds {
		storageDefIdsList = append(storageDefIdsList, id)
	}
	storageDefs, err := storageDefRepo.GetStorageDefinitionsByIds(storageDefIdsList)
	if err != nil {
		return 0, err
	}
	storageDefMap := map[string]*dm.StorageDefinition{}
	for _, sd := range storageDefs {
		storageDefMap[sd.Id] = sd
	}
	storageByStorageDefId := map[string]Storage{}
	for _, si := range storedImagesToDelete {
		storageDef, ok := storageDefMap[si.StorageDefinitionId]
		if !ok {
			return 0, fmt.Errorf("storage definition not found for stored image %s", si.Id)
		}
		storage := storageByStorageDefId[si.StorageDefinitionId]
		if storage == nil {
			storage, err = GetBackend(storageDef.StorageType).FromJSONConfig([]byte(storageDef.Config))
			if err != nil {
				return 0, err
			}
			storageByStorageDefId[si.StorageDefinitionId] = storage
		}
	}

	var tasks []utils.WorkerTask[*dm.StoredImage]
	for _, storedImage := range storedImagesToDelete {
		storage, ok := storageByStorageDefId[storedImage.StorageDefinitionId]
		if !ok {
			logger.Warn().
				Str("stored_image_id", storedImage.Id).
				Str("storage_definition_id", storedImage.StorageDefinitionId).
				Msg("Storage not found for stored image")
			continue
		}

		tasks = append(tasks, utils.WorkerTask[*dm.StoredImage]{
			Job: storedImage,
			Execute: func(img *dm.StoredImage) error {
				return deleteStoredImage(img, storage)
			},
		})
	}

	maxWorkers := 10
	results := utils.RunWorkerPool(tasks, maxWorkers)
	deletedStoredImageIds := []string{}
	for result := range results {
		if result.Err != nil {
			logger.Error().Str("stored_image_id", result.Job.Id).Err(result.Err).Msg("Error deleting stored image")
		} else {
			deletedStoredImageIds = append(deletedStoredImageIds, result.Job.Id)
		}
	}
	err = storedImageRepo.MarkStoredImagesAsDeleted(deletedStoredImageIds)
	return len(deletedStoredImageIds), err
}

func deleteStoredImage(storedImage *dm.StoredImage, s Storage) error {
	meta := s.GetMeta(storedImage.FileIdentifier)
	if meta.ByteSize != 0 {
		err := s.Delete(storedImage.FileIdentifier)
		if err != nil {
			return err
		}
	}
	return nil
}
