package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nasshu2916/dmx_viewer/internal/app"
	"github.com/nasshu2916/dmx_viewer/internal/config"
	"github.com/nasshu2916/dmx_viewer/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logger.NewLogger("error").Fatal("Failed to load configuration: ", err)
		return
	}

	appLogger := logger.NewLogger(cfg.App.LogLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		appLogger.Info("Shutting down...")
		cancel()
	}()

	app.Run(ctx, cfg, appLogger)
}
