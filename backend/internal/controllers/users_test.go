package controllers

import (
	"net/http"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func TestUserControllerCreatePassesActorAndRequest(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	users := &fakeUsersService{output: &dto.User{
		ID:    testUUID(),
		Email: "student@example.com",
		Role:  repositories.UserRoleStudent,
	}}
	router := authenticatedRouter(actor)
	NewUserController(users).RegisterRoutes(router.Group("/users"))

	recorder := requestWithToken(router, http.MethodPost, "/users", dto.CreateUser{
		Email:    "student@example.com",
		Password: "password123",
		Role:     repositories.UserRoleStudent,
	}, "token")

	requireStatus(t, recorder, http.StatusCreated)
	output := decodeResponse[dto.User](t, recorder)
	requireEqual(t, users.output.ID, output.ID)
	requireEqual(t, actor.ID, users.actor.ID)
	requireEqual(t, repositories.UserRoleStudent, users.createInput.Role)
}

func TestUserControllerCreateRejectsInvalidJSON(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	users := &fakeUsersService{}
	router := authenticatedRouter(actor)
	NewUserController(users).RegisterRoutes(router.Group("/users"))

	recorder := requestWithToken(router, http.MethodPost, "/users", map[string]any{"email": 10}, "token")

	requireStatus(t, recorder, http.StatusBadRequest)
	if users.createCalled {
		t.Fatal("service must not be called for invalid JSON")
	}
}

func TestUserControllerGetByIDPassesPathID(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	userID := testUUID()
	users := &fakeUsersService{output: &dto.User{ID: userID, Email: "admin@example.com", Role: repositories.UserRoleAdministrator}}
	router := authenticatedRouter(actor)
	NewUserController(users).RegisterRoutes(router.Group("/users"))

	recorder := requestWithToken(router, http.MethodGet, "/users/"+userID, nil, "token")

	requireStatus(t, recorder, http.StatusOK)
	requireEqual(t, userID, users.id)
	requireEqual(t, actor.ID, users.actor.ID)
}

func TestUserControllerUpdateMapsDuplicateEmail(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	userID := testUUID()
	users := &fakeUsersService{err: services.ErrDuplicateUserEmail}
	router := authenticatedRouter(actor)
	NewUserController(users).RegisterRoutes(router.Group("/users"))

	recorder := requestWithToken(router, http.MethodPatch, "/users/"+userID, dto.UpdateUser{
		Email: ptr("taken@example.com"),
	}, "token")

	requireStatus(t, recorder, http.StatusConflict)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "duplicate_email", output.Error)
}

func TestUserControllerDeleteReturnsNoContent(t *testing.T) {
	actor := testActor(repositories.UserRoleAdministrator)
	userID := testUUID()
	users := &fakeUsersService{}
	router := authenticatedRouter(actor)
	NewUserController(users).RegisterRoutes(router.Group("/users"))

	recorder := requestWithToken(router, http.MethodDelete, "/users/"+userID, nil, "token")

	requireStatus(t, recorder, http.StatusNoContent)
	requireEqual(t, userID, users.id)
	requireEqual(t, actor.ID, users.actor.ID)
}

func authenticatedRouter(actor dto.Actor) *gin.Engine {
	router := testRouter()
	router.Use(AuthMiddleware(&fakeTokenParser{actor: actor}))
	return router
}
