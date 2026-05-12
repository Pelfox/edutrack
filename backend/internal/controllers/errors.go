package controllers

import (
	"errors"
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func respondError(ctx *gin.Context, err error) {
	status, code, message := mapError(err)
	ctx.JSON(status, dto.Error{
		Error:   code,
		Message: message,
	})
}

func mapError(err error) (int, string, string) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		return http.StatusBadRequest, "invalid_request", err.Error()
	case errors.Is(err, services.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid_credentials", "invalid email or password"
	case errors.Is(err, services.ErrUnauthenticatedUser):
		return http.StatusUnauthorized, "unauthenticated", "authentication is required"
	case errors.Is(err, services.ErrForbidden):
		return http.StatusForbidden, "forbidden", "access denied"
	case errors.Is(err, services.ErrUserNotFound):
		return http.StatusNotFound, "not_found", "user not found"
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, "not_found", "resource not found"
	case errors.Is(err, services.ErrDuplicateUserEmail):
		return http.StatusConflict, "duplicate_email", "user with this email already exists"
	default:
		return http.StatusInternalServerError, "internal_error", "internal server error"
	}
}
