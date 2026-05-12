package controllers

import (
	"strings"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

const actorContextKey = "actor"

// TokenParser описывает сервис, способный разобрать JWT-токен.
type TokenParser interface {
	// ParseToken проверяет JWT-токен и возвращает данные авторизованного пользователя.
	ParseToken(tokenValue string) (dto.Actor, error)
}

// AuthMiddleware проверяет Bearer-токен и сохраняет данные пользователя в контексте запроса.
func AuthMiddleware(parser TokenParser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := bearerToken(ctx.GetHeader("Authorization"))
		if tokenValue == "" {
			respondError(ctx, services.ErrUnauthenticatedUser)
			ctx.Abort()
			return
		}

		actor, err := parser.ParseToken(tokenValue)
		if err != nil {
			respondError(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set(actorContextKey, actor)
		ctx.Next()
	}
}

func actorFromContext(ctx *gin.Context) (dto.Actor, bool) {
	value, ok := ctx.Get(actorContextKey)
	if !ok {
		return dto.Actor{}, false
	}

	actor, ok := value.(dto.Actor)
	return actor, ok
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
