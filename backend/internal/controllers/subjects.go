package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SubjectController обрабатывает HTTP-запросы модуля предметов.
type SubjectController struct {
	subjects services.SubjectsService
}

// NewSubjectController создаёт контроллер предметов.
func NewSubjectController(subjects services.SubjectsService) *SubjectController {
	return &SubjectController{subjects: subjects}
}

// RegisterRoutes регистрирует маршруты модуля предметов.
func (controller *SubjectController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("", controller.List)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание предмета.
// @Summary Создать предмет
// @Tags subjects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateSubject true "Данные предмета"
// @Success 201 {object} dto.Subject
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /subjects [post]
func (controller *SubjectController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateSubject
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.subjects.Create(ctx.Request.Context(), actor, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// List обрабатывает получение списка предметов.
// @Summary Получить список предметов
// @Tags subjects
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Subject
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /subjects [get]
func (controller *SubjectController) List(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.subjects.List(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetByID обрабатывает получение предмета по идентификатору.
// @Summary Получить предмет по идентификатору
// @Tags subjects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор предмета"
// @Success 200 {object} dto.Subject
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /subjects/{id} [get]
func (controller *SubjectController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.subjects.GetByID(ctx.Request.Context(), actor, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает обновление предмета.
// @Summary Обновить предмет
// @Tags subjects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор предмета"
// @Param request body dto.UpdateSubject true "Данные предмета"
// @Success 200 {object} dto.Subject
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /subjects/{id} [patch]
func (controller *SubjectController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateSubject
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.subjects.Update(ctx.Request.Context(), actor, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление предмета.
// @Summary Удалить предмет
// @Tags subjects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор предмета"
// @Success 204
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /subjects/{id} [delete]
func (controller *SubjectController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.subjects.Delete(ctx.Request.Context(), actor, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
