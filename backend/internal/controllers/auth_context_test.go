package controllers

import (
	"net/http"
	"testing"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareRequiresBearerToken(t *testing.T) {
	router := testRouter()
	router.GET("/protected", AuthMiddleware(&fakeTokenParser{}), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	recorder := request(router, http.MethodGet, "/protected", nil)
	requireStatus(t, recorder, http.StatusUnauthorized)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "unauthenticated", output.Error)
}

func TestAuthMiddlewareStoresActor(t *testing.T) {
	actor := dto.Actor{ID: testUUID(), Role: repositories.UserRoleTeacher}
	parser := &fakeTokenParser{actor: actor}
	router := testRouter()
	router.GET("/protected", AuthMiddleware(parser), func(ctx *gin.Context) {
		contextActor, ok := actorFromContext(ctx)
		if !ok {
			t.Fatal("expected actor in context")
		}
		requireEqual(t, actor.ID, contextActor.ID)
		requireEqual(t, actor.Role, contextActor.Role)
		ctx.Status(http.StatusNoContent)
	})

	recorder := requestWithToken(router, http.MethodGet, "/protected", nil, "token-value")
	requireStatus(t, recorder, http.StatusNoContent)
	requireEqual(t, "token-value", parser.token)
}

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) {
	router := testRouter()
	router.GET("/protected", AuthMiddleware(&fakeTokenParser{err: services.ErrUnauthenticatedUser}), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	recorder := requestWithToken(router, http.MethodGet, "/protected", nil, "bad-token")
	requireStatus(t, recorder, http.StatusUnauthorized)
	output := decodeResponse[dto.Error](t, recorder)
	requireEqual(t, "unauthenticated", output.Error)
}
