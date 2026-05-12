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
