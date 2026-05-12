package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// GradeController обрабатывает HTTP-запросы модуля оценок.
type GradeController struct {
	grades services.GradesService
}

// NewGradeController создаёт контроллер оценок.
func NewGradeController(grades services.GradesService) *GradeController {
	return &GradeController{grades: grades}
}

// RegisterRoutes регистрирует маршруты модуля оценок.
func (controller *GradeController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("", controller.List)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание оценки.
// @Summary Создать оценку
// @Tags grades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateGrade true "Данные оценки"
// @Success 201 {object} dto.Grade
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /grades [post]
func (controller *GradeController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateGrade
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.grades.Create(ctx.Request.Context(), actor, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// List обрабатывает получение списка оценок.
// @Summary Получить список оценок
// @Tags grades
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Grade
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /grades [get]
func (controller *GradeController) List(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.grades.List(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetByID обрабатывает получение оценки по идентификатору.
// @Summary Получить оценку по идентификатору
// @Tags grades
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор оценки"
// @Success 200 {object} dto.Grade
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /grades/{id} [get]
func (controller *GradeController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.grades.GetByID(ctx.Request.Context(), actor, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает частичное обновление оценки.
// @Summary Частично обновить оценку
// @Tags grades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор оценки"
// @Param request body dto.UpdateGrade true "Данные оценки"
// @Success 200 {object} dto.Grade
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /grades/{id} [patch]
func (controller *GradeController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateGrade
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.grades.Update(ctx.Request.Context(), actor, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление оценки.
// @Summary Удалить оценку
// @Tags grades
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор оценки"
// @Success 204
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /grades/{id} [delete]
func (controller *GradeController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.grades.Delete(ctx.Request.Context(), actor, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
