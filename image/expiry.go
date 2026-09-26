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

var expiryLogger = logging.GetLogger("expiry")

const defaultExpiryIntervalSeconds = 60

// ExpiryConfig controls the periodic task that marks expired images as deleted.
// It is separate from the stored image cleanup task: the expiry sweep is a
// single cheap UPDATE, while cleanup talks to external storage, so the sweep
// can run much more often.
type ExpiryConfig struct {
	Enabled  bool
	Interval time.Duration
}

func ReadExpiryConfigFromEnv() *ExpiryConfig {
	if os.Getenv("EXPIRY_ENABLED") == "" {
		return nil
	}
	enabled := utils.IsStrTruthy(os.Getenv("EXPIRY_ENABLED"))
	if !enabled {
		return &ExpiryConfig{Enabled: false}
	}
	intervalStr := os.Getenv("EXPIRY_INTERVAL")
	if intervalStr == "" {
		expiryLogger.Warn().Msg("EXPIRY_INTERVAL not set. Using default value of 60 seconds")
		intervalStr = strconv.Itoa(defaultExpiryIntervalSeconds)
	}
	intervalInt, err := strconv.Atoi(intervalStr)
	if err != nil {
		expiryLogger.Warn().Err(err).Msg("Error parsing EXPIRY_INTERVAL. Using default value of 60 seconds")
		intervalInt = defaultExpiryIntervalSeconds
	}
	return &ExpiryConfig{
		Enabled:  true,
		Interval: time.Duration(intervalInt) * time.Second,
	}
}

// DeleteExpiredImagesTask runs one expiry sweep under lock. Expired images are
// auto marked as deleted here; their files are then removed by the regular
// stored image cleanup task like any other deleted image.
func DeleteExpiredImagesTask(lock utils.MutexLock, imageRepo ImageRepo) error {
	return utils.RunWithLock(lock, func() error {
		count, err := imageRepo.DeleteExpiredImages()
		if err != nil {
			expiryLogger.Error().Err(err).Msg("Error deleting expired images")
		} else if count > 0 {
			expiryLogger.Info().Int("expired_count", count).Msg("Deleted expired images")
		}
		return nil
	})
}

func RunExpiryTask(lock utils.MutexLock, imageRepo ImageRepo, interval time.Duration) {
	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		expiryLogger.Info().
			Str("interval", interval.String()).
			Msg("Expiry task started")

		for {
			select {
			case <-ticker.C:
				DeleteExpiredImagesTask(lock, imageRepo)
			case <-ctx.Done():
				expiryLogger.Info().Msg("Expiry task shutting down gracefully...")
				return
			}
		}
	}()

	<-stop
	expiryLogger.Info().Msg("Received termination signal. Initiating shutdown...")
	cancel()
	time.Sleep(1 * time.Second)
	expiryLogger.Info().Msg("Expiry task shutdown complete.")
}
