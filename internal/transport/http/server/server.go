package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ItsXomyak/scam-list/config"
	"github.com/ItsXomyak/scam-list/internal/transport/http/handler"
	"github.com/ItsXomyak/scam-list/pkg/logger"
	"github.com/gin-gonic/gin"
)

const serverIPAddress = "%s:%d"

type API struct {
	router *gin.Engine
	server *http.Server
	routes *handlers

	addr string
	log  logger.Logger
}

type handlers struct {
	verify *handler.Verify
}

func New(cfg config.Config, verifier Verifier, logger logger.Logger) *API {
	addr := fmt.Sprintf(serverIPAddress, "0.0.0.0", cfg.HTTPServer.Port)

	// Set Gin mode based on environment
	switch cfg.HTTPServer.GinEnviroment {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
		gin.SetMode(cfg.HTTPServer.GinEnviroment)
	default:
		logger.Warn(context.Background(), "invalid gin environment, setting to debug mode", "env", cfg.HTTPServer.GinEnviroment)
		gin.SetMode(gin.DebugMode)
	}

	// Initialize handlers
	handlers := &handlers{
		verify: handler.NewVerify(verifier, logger),
	}

	router := gin.New()

	api := &API{
		router: router,
		routes: handlers,
		addr:   addr,
		log:    logger,
	}

	api.server = &http.Server{
		Addr:    api.addr,
		Handler: api.withMiddleware(),
	}

	api.setupRoutes()

	return api
}

func (a *API) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx = logger.WithAction(ctx, "http_server_shutdown")

	a.log.Debug(ctx, "shutting down HTTP server...", "address", a.addr)
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}
	a.log.Debug(ctx, "shutting down HTTP server completed", "address", a.addr)

	return nil
}

func (a *API) Start(ctx context.Context, errCh chan<- error) {
	go func() {
		a.log.Info(ctx, "started http server", "address", a.addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}
