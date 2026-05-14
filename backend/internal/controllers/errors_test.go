package controllers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func TestRespondErrorMapsServiceErrors(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		status    int
		errorCode string
		message   string
	}{
		{
			name:      "invalid input",
			err:       services.ErrInvalidInput,
			status:    http.StatusBadRequest,
			errorCode: "invalid_request",
			message:   "invalid input",
		},
		{
			name:      "invalid credentials",
			err:       services.ErrInvalidCredentials,
			status:    http.StatusUnauthorized,
			errorCode: "invalid_credentials",
			message:   "invalid email or password",
		},
		{
			name:      "unauthenticated",
			err:       services.ErrUnauthenticatedUser,
			status:    http.StatusUnauthorized,
			errorCode: "unauthenticated",
			message:   "authentication is required",
		},
		{
			name:      "forbidden",
			err:       services.ErrForbidden,
			status:    http.StatusForbidden,
			errorCode: "forbidden",
			message:   "access denied",
		},
		{
			name:      "not found",
			err:       services.ErrNotFound,
			status:    http.StatusNotFound,
			errorCode: "not_found",
			message:   "resource not found",
		},
		{
			name:      "duplicate email",
			err:       services.ErrDuplicateUserEmail,
			status:    http.StatusConflict,
			errorCode: "duplicate_email",
			message:   "user with this email already exists",
		},
		{
			name:      "internal",
			err:       errors.New("boom"),
			status:    http.StatusInternalServerError,
			errorCode: "internal_error",
			message:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := testRouter()
			router.GET("/test", func(ctx *gin.Context) {
				respondError(ctx, tt.err)
			})

			recorder := request(router, http.MethodGet, "/test", nil)
			requireStatus(t, recorder, tt.status)
			output := decodeResponse[map[string]string](t, recorder)
			requireEqual(t, tt.errorCode, output["error"])
			requireEqual(t, tt.message, output["message"])
		})
	}
}
