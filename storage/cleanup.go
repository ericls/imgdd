package storage

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/logging"
	"github.com/ericls/imgdd/utils"
)

var logger = logging.GetLogger("cleanup")
var lockKey = utils.GenerateLockKey("cleanup_stored_image")

type CleanupConfig struct {
	Enabled  bool
	Interval time.Duration
}

func ReadCleanupConfigFromEnv() *CleanupConfig {
	if os.Getenv("CLEANUP_ENABLED") == "" {
		return nil
	}
	enabled := utils.IsStrTruthy(os.Getenv("CLEANUP_ENABLED"))
	if !enabled {
		return &CleanupConfig{Enabled: false}
	}
	intervalStr := os.Getenv("CLEANUP_INTERVAL")
	if intervalStr == "" {
		logger.Warn().Msg("CLEANUP_INTERVAL not set. Using default value of 600 seconds")
		intervalStr = "600"
	}
	intervalInt, err := strconv.Atoi(intervalStr)
	if err != nil {
		logger.Warn().Err(err).Msg("Error parsing CLEANUP_INTERVAL. Using default value of 600 seconds")
		intervalInt = 600
	}
	return &CleanupConfig{
		Enabled:  true,
		Interval: time.Duration(intervalInt) * time.Second,
	}
}

func CleanupStoredImageTask(lock utils.MutexLock, storedImageRepo StoredImageRepo, storageDefRepo StorageDefRepo) error {
	return utils.RunWithLock(lock, func() error {
		deletedCount, err := CleanupStoredImage(storedImageRepo, storageDefRepo)
		if err != nil {
			logger.Error().Err(err).Msg("Error cleaning up stored images")
		} else {
			logger.Info().Int("deleted_count", deletedCount).Msg("Cleaned up stored images")
		}
		return nil
	})
}

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

func RunCleanupTask(lock utils.MutexLock, storedImageRepo StoredImageRepo, storageDefRepo StorageDefRepo, interval time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		logger.Info().
			Str("interval", interval.String()).
			Msg("Cleanup task started")

		for {
			select {
			case <-ticker.C:
				CleanupStoredImageTask(lock, storedImageRepo, storageDefRepo)
			case <-ctx.Done():
				logger.Info().Msg("Cleanup task shutting down gracefully...")
				return
			}
		}
	}()

	// Wait for termination signal
	<-stop
	logger.Info().Msg("Received termination signal. Initiating shutdown...")
	cancel()
	time.Sleep(1 * time.Second) // Give tasks some time to finish
	logger.Info().Msg("Cleanup task shutdown complete.")
}
