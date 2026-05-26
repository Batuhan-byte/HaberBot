package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"haberbot/internal/infrastructure/config"
	"haberbot/internal/infrastructure/container"
	"haberbot/internal/infrastructure/server"
)

func main() {
	// Configure logging format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting HaberBot application server...")

	// 1. Load configuration
	cfg := config.Load()

	// 2. Build dependency container
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	container, err := container.NewContainer(ctx, cfg)
	if err != nil {
		slog.Error("failed to initialize container", "error", err)
		os.Exit(1)
	}
	defer container.DBPool.Close()

	// 3. Start cron scheduler
	err = container.Scheduler.Start(cfg.CronFetchSchedule, cfg.CronProcessSchedule)
	if err != nil {
		slog.Error("failed to start scheduler", "error", err)
		os.Exit(1)
	}
	defer container.Scheduler.Stop()

	// 4. Setup web server
	app := server.SetupFiberServer(container)

	// 5. Graceful shutdown orchestration
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("http server listening", "port", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server failed", "error", err)
			shutdownChan <- syscall.SIGTERM
		}
	}()

	<-shutdownChan
	slog.Info("shutting down HTTP server gracefully...")

	// Wait up to 10 seconds for current requests to finish
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
	} else {
		slog.Info("server shutdown complete")
	}
}
