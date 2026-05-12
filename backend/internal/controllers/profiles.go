package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// ProfileController обрабатывает HTTP-запросы профиля текущего пользователя.
type ProfileController struct {
	profiles services.ProfileService
}

// NewProfileController создаёт контроллер профиля текущего пользователя.
func NewProfileController(profiles services.ProfileService) *ProfileController {
	return &ProfileController{profiles: profiles}
}

// RegisterRoutes регистрирует маршруты профиля текущего пользователя.
func (controller *ProfileController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/me", controller.GetMe)
}

// GetMe обрабатывает получение профиля текущего пользователя.
// @Summary Получить профиль текущего пользователя
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Profile
// @Failure 401 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /profile/me [get]
func (controller *ProfileController) GetMe(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.profiles.GetMe(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
