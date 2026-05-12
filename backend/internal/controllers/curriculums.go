package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// CurriculumController обрабатывает HTTP-запросы модуля учебных планов.
type CurriculumController struct {
	curriculums services.CurriculumsService
}

// NewCurriculumController создаёт контроллер учебных планов.
func NewCurriculumController(curriculums services.CurriculumsService) *CurriculumController {
	return &CurriculumController{curriculums: curriculums}
}

// RegisterRoutes регистрирует маршруты модуля учебных планов.
func (controller *CurriculumController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("", controller.List)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание учебного плана.
// @Summary Создать учебный план
// @Tags curriculums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCurriculum true "Данные учебного плана"
// @Success 201 {object} dto.Curriculum
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /curriculums [post]
func (controller *CurriculumController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateCurriculum
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.curriculums.Create(ctx.Request.Context(), actor, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// List обрабатывает получение списка учебных планов.
// @Summary Получить список учебных планов
// @Tags curriculums
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Curriculum
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /curriculums [get]
func (controller *CurriculumController) List(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.curriculums.List(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetByID обрабатывает получение учебного плана по идентификатору.
// @Summary Получить учебный план по идентификатору
// @Tags curriculums
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор учебного плана"
// @Success 200 {object} dto.Curriculum
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /curriculums/{id} [get]
func (controller *CurriculumController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.curriculums.GetByID(ctx.Request.Context(), actor, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает частичное обновление учебного плана.
// @Summary Частично обновить учебный план
// @Tags curriculums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор учебного плана"
// @Param request body dto.UpdateCurriculum true "Данные учебного плана"
// @Success 200 {object} dto.Curriculum
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /curriculums/{id} [patch]
func (controller *CurriculumController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateCurriculum
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.curriculums.Update(ctx.Request.Context(), actor, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление учебного плана.
// @Summary Удалить учебный план
// @Tags curriculums
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор учебного плана"
// @Success 204
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /curriculums/{id} [delete]
func (controller *CurriculumController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.curriculums.Delete(ctx.Request.Context(), actor, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
