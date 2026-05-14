package controllers

import (
	"net/http"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
)

func TestAuthControllerLoginReturnsResult(t *testing.T) {
	userID := testUUID()
	auth := &fakeAuthService{loginResult: &dto.LoginResult{
		User: dto.User{
			ID:    userID,
			Email: "admin@example.com",
			Role:  repositories.UserRoleAdministrator,
		},
		Token: "jwt-token",
	}}
	router := testRouter()
	NewAuthController(auth).RegisterRoutes(router.Group("/auth"))

	recorder := request(router, http.MethodPost, "/auth/login", dto.Login{
		Email:    "admin@example.com",
		Password: "password123",
	})

	requireStatus(t, recorder, http.StatusOK)
	output := decodeResponse[dto.LoginResult](t, recorder)
	requireEqual(t, userID, output.User.ID)
	requireEqual(t, "jwt-token", output.Token)
	requireEqual(t, "admin@example.com", auth.loginInput.Email)
}

func TestAuthControllerLoginRejectsInvalidJSON(t *testing.T) {
	auth := &fakeAuthService{}
	router := testRouter()
	NewAuthController(auth).RegisterRoutes(router.Group("/auth"))

	recorder := request(router, http.MethodPost, "/auth/login", map[string]any{"email": 10})

	requireStatus(t, recorder, http.StatusBadRequest)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "invalid_request", output.Error)
}

func TestAuthControllerLoginMapsServiceError(t *testing.T) {
	auth := &fakeAuthService{loginErr: services.ErrInvalidCredentials}
	router := testRouter()
	NewAuthController(auth).RegisterRoutes(router.Group("/auth"))

	recorder := request(router, http.MethodPost, "/auth/login", dto.Login{
		Email:    "admin@example.com",
		Password: "password123",
	})

	requireStatus(t, recorder, http.StatusUnauthorized)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "invalid_credentials", output.Error)
}
