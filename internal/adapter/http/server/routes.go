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

	// для бауки
	admin := a.router.Group("/admin")
	{
		admin.POST("/domain/create", a.routes.admin.CreateDomain)
		admin.GET("/domain", a.routes.admin.GetAllDomains)
		admin.GET("/domain/:domain", a.routes.admin.GetDomain)
		admin.PATCH("/domain/:domain", a.routes.admin.PatchDomain)
		admin.DELETE("/domain/:domain", a.routes.admin.DeleteDomain)
	}
}

// setupDefaultRoutes - setups default http routes
func (a *API) setupDefaultRoutes() {
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
