package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, status int, message any) {
	env := gin.H{"error": message}

	c.JSON(status, env)
}

func internalErrorResponse(c *gin.Context, message any) {
	errorResponse(c, http.StatusInternalServerError, message)
}

func notFoundResponse(c *gin.Context, message any) {
	errorResponse(c, http.StatusNotFound, message)
}

func badRequestResponse(c *gin.Context, message any) {
	errorResponse(c, http.StatusBadRequest, message)
}
