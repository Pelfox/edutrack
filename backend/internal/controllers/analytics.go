package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// AnalyticsController обрабатывает HTTP-запросы модуля аналитики.
type AnalyticsController struct {
	analytics services.AnalyticsService
}

// NewAnalyticsController создаёт контроллер аналитики.
func NewAnalyticsController(analytics services.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{analytics: analytics}
}

// RegisterRoutes регистрирует маршруты модуля аналитики.
func (controller *AnalyticsController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/overview", controller.GetOverview)
}

// GetOverview обрабатывает получение сводных показателей админ-панели.
// @Summary Получить сводные показатели админ-панели
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AnalyticsOverview
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /analytics/overview [get]
func (controller *AnalyticsController) GetOverview(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.analytics.GetOverview(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
