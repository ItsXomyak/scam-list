package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// setupRoutes - setups http routes
func (a *API) setupRoutes() {
	a.setupDefaultRoutes()

	// API routes
	api := a.router.Group("/api")
	{
		api.GET("/verify/:domain", a.routes.verify.VerifyDomain)
	}
}

// setupDefaultRoutes - setups default http routes
func (a *API) setupDefaultRoutes() {
	// System Health
	a.router.GET("/health", a.HealthCheck)
}

// HealthCheck - returns system information.
func (a *API) HealthCheck(c *gin.Context) {
	response := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"address": a.addr,
			"mode":    gin.Mode(),
		},
	}

	c.JSON(http.StatusOK, response)
}
