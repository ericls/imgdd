package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulServe(srv *http.Server, shutdownTimeout time.Duration) error {
	go func() {
		httpLogger.Info().Str("addr", srv.Addr).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			httpLogger.Fatal().Err(err).Msg("Server failed to start")
			os.Exit(1)
		}
	}()

	// Create a channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a stop signal
	<-stop
	httpLogger.Info().Msg("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		httpLogger.Error().Err(err).Msg("Server failed to shut down gracefully")
		// Fallback: force shutdown
		if closeErr := srv.Close(); closeErr != nil {
			httpLogger.Error().Err(closeErr).Msg("Server failed to shut down forcefully")
		}
		return err
	}

	httpLogger.Info().Msg("Server shut down.")
	return nil
}
