package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ItsXomyak/scam-list/config"
	httpserver "github.com/ItsXomyak/scam-list/internal/adapter/http/server"
	"github.com/ItsXomyak/scam-list/internal/adapter/postgres"
	"github.com/ItsXomyak/scam-list/internal/services/domain"
	"github.com/ItsXomyak/scam-list/internal/services/pipeline"
	"github.com/ItsXomyak/scam-list/pkg/logger"
	postgresclient "github.com/ItsXomyak/scam-list/pkg/postgres"
)

// App struct represents the application
type App struct {
	postgresDB *postgresclient.Postgres
	httpServer *httpserver.API

	cfg config.Config
	log logger.Logger
}

// NewApp creates a new instance of the application
func NewApp(ctx context.Context, cfg config.Config, log logger.Logger) (*App, error) {
	// Initialize Postgres client
	postgresDB, err := postgresclient.New(ctx, cfg.Postgres.GetDsn(), &postgresclient.Config{
		MaxPoolSize:  cfg.Postgres.MaxPoolSize,
		ConnAttempts: cfg.Postgres.ConnAttempts,
		ConnTimeout:  cfg.Postgres.ConnTimeout,
	})
	if err != nil {
		return nil, err
	}

	// repositories
	domainRepo := postgres.NewDomain(postgresDB.Pool)

	// services
	domainSvc := domain.NewDomainService(domainRepo)

	// core pipeline
	domainPipeline := pipeline.NewDomainPipeline(nil, domainSvc)

	// Initialize HTTP server
	server := httpserver.New(cfg, domainPipeline, domainRepo, log)

	return &App{
		postgresDB: postgresDB,
		httpServer: server,
		cfg:        cfg,
		log:        log,
	}, nil
}

// Run starts the application
func (app *App) Run(ctx context.Context) error {
	// Graceful shutdown
	defer func() {
		ctx = logger.WithAction(ctx, "app_shutdown")

		if err := app.Shutdown(ctx); err != nil {
			app.log.Error(logger.ErrorCtx(ctx, err), "error during app shutdown", err)
			return
		}
		app.log.Info(ctx, "graceful shutdown completed")
	}()

	ctx = logger.WithAction(ctx, "app_run")

	errCh := make(chan error, 1)
	app.httpServer.Start(ctx, errCh)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	app.log.Info(ctx, "application started")
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdownCh:
		app.log.Info(ctx, "shuting down application", "signal", sig.String())
		return nil
	}

	return nil
}

// Shutdown gracefully shuts down the application
func (app *App) Shutdown(ctx context.Context) error {
	t, err := time.ParseDuration(fmt.Sprintf("%ds", app.cfg.HTTPServer.ShutdownTimeoutSeconds))
	if err != nil {
		t = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()

	// Shutdown HTTP server
	if app.httpServer != nil {
		if err := app.httpServer.Stop(ctx); err != nil {
			return logger.WrapError(ctx, fmt.Errorf("failed to stop http server: %w", err))
		}
	}

	// Close Postgres connection
	if app.postgresDB != nil {
		app.postgresDB.Close()
	}

	return nil
}
