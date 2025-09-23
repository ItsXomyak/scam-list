package dto

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// HTTPError represent http error response
type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrResourceNotFoundResponse = &HTTPError{
		Code:    http.StatusNotFound,
		Message: "the requested resource could not be found",
	}

	ErrUnprocessableEntityResponse = &HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: "unprocessable entity",
	}

	ErrDuplicateKeyResponse = &HTTPError{
		Code:    http.StatusConflict,
		Message: "requested resource already exists",
	}
)

func FromError(err error) *HTTPError {
	if err == nil {
		return &HTTPError{
			Code:    http.StatusOK,
			Message: "ok",
		}
	}

	switch {
	case errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows):
		return ErrResourceNotFoundResponse
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrDuplicateKeyResponse
			// можно добавить и другие SQLSTATE коды при необходимости
		}
	}

	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}
