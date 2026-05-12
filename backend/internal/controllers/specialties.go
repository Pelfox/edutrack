package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SpecialtyController обрабатывает HTTP-запросы модуля специальностей.
type SpecialtyController struct {
	specialties services.SpecialtiesService
}

// NewSpecialtyController создаёт контроллер специальностей.
func NewSpecialtyController(specialties services.SpecialtiesService) *SpecialtyController {
	return &SpecialtyController{specialties: specialties}
}

// RegisterRoutes регистрирует маршруты модуля специальностей.
func (controller *SpecialtyController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("", controller.List)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание специальности.
// @Summary Создать специальность
// @Tags specialties
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateSpecialty true "Данные специальности"
// @Success 201 {object} dto.Specialty
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /specialties [post]
func (controller *SpecialtyController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateSpecialty
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.specialties.Create(ctx.Request.Context(), actor, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// List обрабатывает получение списка специальностей.
// @Summary Получить список специальностей
// @Tags specialties
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Specialty
// @Failure 401 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /specialties [get]
func (controller *SpecialtyController) List(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.specialties.List(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetByID обрабатывает получение специальности по идентификатору.
// @Summary Получить специальность по идентификатору
// @Tags specialties
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор специальности"
// @Success 200 {object} dto.Specialty
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /specialties/{id} [get]
func (controller *SpecialtyController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.specialties.GetByID(ctx.Request.Context(), actor, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает обновление специальности.
// @Summary Обновить специальность
// @Tags specialties
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор специальности"
// @Param request body dto.UpdateSpecialty true "Данные специальности"
// @Success 200 {object} dto.Specialty
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /specialties/{id} [patch]
func (controller *SpecialtyController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateSpecialty
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.specialties.Update(ctx.Request.Context(), actor, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление специальности.
// @Summary Удалить специальность
// @Tags specialties
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор специальности"
// @Success 204
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /specialties/{id} [delete]
func (controller *SpecialtyController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.specialties.Delete(ctx.Request.Context(), actor, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
