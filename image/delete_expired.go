package image

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ericls/imgdd/logging"
	"github.com/ericls/imgdd/utils"
)

var deleteExpiredLogger = logging.GetLogger("delete-expired-images")

const DefaultDeleteExpiredImagesInterval = 60 * time.Second

// DeleteExpiredImagesConfig controls the periodic task that marks expired
// images as deleted. It is enabled by default.
//
// It is separate from the stored image cleanup task: marking images as deleted
// is a single cheap UPDATE, while cleanup talks to external storage, so this
// task can run much more often.
type DeleteExpiredImagesConfig struct {
	Enabled  bool
	Interval time.Duration
}

func DefaultDeleteExpiredImagesConfig() *DeleteExpiredImagesConfig {
	return &DeleteExpiredImagesConfig{
		Enabled:  true,
		Interval: DefaultDeleteExpiredImagesInterval,
	}
}

// ReadDeleteExpiredImagesConfigFromEnv reads DELETE_EXPIRED_IMAGES_ENABLED and
// DELETE_EXPIRED_IMAGES_INTERVAL (seconds). It returns nil when neither is set,
// so that other config sources or the default apply.
func ReadDeleteExpiredImagesConfigFromEnv() *DeleteExpiredImagesConfig {
	enabledStr := os.Getenv("DELETE_EXPIRED_IMAGES_ENABLED")
	intervalStr := os.Getenv("DELETE_EXPIRED_IMAGES_INTERVAL")
	if enabledStr == "" && intervalStr == "" {
		return nil
	}
	conf := DefaultDeleteExpiredImagesConfig()
	if enabledStr != "" {
		conf.Enabled = utils.IsStrTruthy(enabledStr)
	}
	if intervalStr != "" {
		intervalInt, err := strconv.Atoi(intervalStr)
		if err != nil || intervalInt <= 0 {
			deleteExpiredLogger.Warn().Err(err).Str("value", intervalStr).
				Msg("Invalid DELETE_EXPIRED_IMAGES_INTERVAL. Using default value of 60 seconds")
		} else {
			conf.Interval = time.Duration(intervalInt) * time.Second
		}
	}
	return conf
}

// DeleteExpiredImagesTask runs one pass under lock. Expired images are auto
// marked as deleted here; their files are then removed by the regular stored
// image cleanup task like any other deleted image.
func DeleteExpiredImagesTask(lock utils.MutexLock, imageRepo ImageRepo) error {
	return utils.RunWithLock(lock, func() error {
		count, err := imageRepo.DeleteExpiredImages()
		if err != nil {
			deleteExpiredLogger.Error().Err(err).Msg("Error deleting expired images")
		} else if count > 0 {
			deleteExpiredLogger.Info().Int("deleted_count", count).Msg("Deleted expired images")
		}
		return nil
	})
}

func RunDeleteExpiredImagesTask(lock utils.MutexLock, imageRepo ImageRepo, interval time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		deleteExpiredLogger.Info().
			Str("interval", interval.String()).
			Msg("Delete expired images task started")

		for {
			select {
			case <-ticker.C:
				DeleteExpiredImagesTask(lock, imageRepo)
			case <-ctx.Done():
				deleteExpiredLogger.Info().Msg("Delete expired images task shutting down gracefully...")
				return
			}
		}
	}()

	<-stop
	deleteExpiredLogger.Info().Msg("Received termination signal. Initiating shutdown...")
	cancel()
	time.Sleep(1 * time.Second)
	deleteExpiredLogger.Info().Msg("Delete expired images task shutdown complete.")
}
