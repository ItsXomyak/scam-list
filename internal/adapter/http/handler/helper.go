package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func getCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // unique_violation
			return http.StatusConflict
		}
	}

	// default
	return http.StatusInternalServerError
}
