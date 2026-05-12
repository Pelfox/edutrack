package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// GroupController обрабатывает HTTP-запросы модуля групп.
type GroupController struct {
	groups services.GroupsService
}

// NewGroupController создаёт контроллер групп.
func NewGroupController(groups services.GroupsService) *GroupController {
	return &GroupController{groups: groups}
}

// RegisterRoutes регистрирует маршруты модуля групп.
func (controller *GroupController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("", controller.List)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание группы.
func (controller *GroupController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateGroup
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.groups.Create(ctx.Request.Context(), actor, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// List обрабатывает получение списка групп.
func (controller *GroupController) List(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.groups.List(ctx.Request.Context(), actor)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetByID обрабатывает получение группы по идентификатору.
func (controller *GroupController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.groups.GetByID(ctx.Request.Context(), actor, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает частичное обновление группы.
func (controller *GroupController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateGroup
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.groups.Update(ctx.Request.Context(), actor, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление группы.
func (controller *GroupController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.groups.Delete(ctx.Request.Context(), actor, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
