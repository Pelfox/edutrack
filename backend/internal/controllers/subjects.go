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
