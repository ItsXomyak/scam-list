package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func internalErrorResponse(c *gin.Context, message any) {
	errorResponse(c, http.StatusInternalServerError, message)
}

func errorResponse(c *gin.Context, status int, message any) {
	env := gin.H{"error": message}

	c.JSON(status, env)
}

func badRequestResponse(c *gin.Context, message any) {
	errorResponse(c, http.StatusBadRequest, message)
}
