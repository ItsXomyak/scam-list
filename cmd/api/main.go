package main

import (
	"context"

	"github.com/ItsXomyak/scam-list/config"
	"github.com/ItsXomyak/scam-list/internal/app"
	"github.com/ItsXomyak/scam-list/pkg/logger"
)

const (
	serviceName = "scam-list"
	configPath  = ".env"
)

func main() {
	ctx := context.Background()

	// Initialize logger
	logger := logger.InitLogger(serviceName, logger.LevelDebug)

	// Load configuration
	cfg, err := config.New(configPath)
	if err != nil {
		logger.Error(ctx, "failed to load config", err)
		return
	}

	// Print loaded configuration
	config.PrintConfig(cfg)

	// Create and run the app
	app, err := app.NewApp(ctx, cfg, logger)
	if err != nil {
		logger.Error(ctx, "failed to create app", err)
		return
	}

	// Run the application
	if err := app.Run(ctx); err != nil {
		logger.Error(ctx, "app run failed", err)
		return
	}
}
