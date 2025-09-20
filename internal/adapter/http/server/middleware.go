package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/ItsXomyak/scam-list/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (a *API) withMiddleware() http.Handler {
	// Apply Gin's built-in middleware
	a.router.Use(gin.Recovery())

	// Custom middlewares
	a.router.Use(a.RequestIDMiddleware())

	return a.router
}

// RequestIDMiddleware injects request_id to the request ctx
func (a *API) RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate request_id
		reqID := newRequestID()

		// Echo to clients for debugging / tracing
		c.Writer.Header().Set("X-Request-ID", reqID)

		ctx := logger.WithRequestID(c.Request.Context(), reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RequestLoggingMiddleware injects a request ID into the context and logs the request details.
func (a *API) RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ctx := c.Request.Context()

		// 1. Log request start
		a.log.Debug(
			ctx,
			"started",
			"method", c.Request.Method,
			"URL", c.Request.URL.Path,
			"request-host", c.Request.Host,
		)

		// Process request
		c.Next()

		// 3. Log request end
		duration := time.Since(start)
		a.log.Debug(
			ctx,
			"completed",
			"method", c.Request.Method,
			"URL", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
		)
	}
}

// newRequestID returns a 16-byte random hex string, e.g. “9f86d081884c7d65…”
func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// fallback to timestamp if crypto/rand fails
		return hex.EncodeToString(fmt.Appendf(nil, "%d", time.Now().UnixNano()))
	}
	return hex.EncodeToString(b)
}
