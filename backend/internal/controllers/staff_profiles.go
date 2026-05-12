package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// StaffProfileController обрабатывает HTTP-запросы профилей сотрудников.
type StaffProfileController struct {
	profiles services.StaffProfilesService
	role     repositories.UserRole
}

// NewTeacherController создаёт контроллер профилей преподавателей.
func NewTeacherController(profiles services.StaffProfilesService) *StaffProfileController {
	return &StaffProfileController{profiles: profiles, role: repositories.UserRoleTeacher}
}

// NewAdministratorController создаёт контроллер профилей администраторов.
func NewAdministratorController(profiles services.StaffProfilesService) *StaffProfileController {
	return &StaffProfileController{profiles: profiles, role: repositories.UserRoleAdministrator}
}

// RegisterRoutes регистрирует маршруты профилей сотрудников.
func (controller *StaffProfileController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("", controller.List)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание профиля сотрудника.
// @Summary Создать профиль сотрудника
// @Tags staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProfile true "Данные профиля"
// @Success 201 {object} dto.Profile
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /teachers [post]
// @Router /administrators [post]
func (controller *StaffProfileController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateProfile
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.profiles.Create(ctx.Request.Context(), actor, controller.role, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// List обрабатывает получение списка профилей сотрудников.
// @Summary Получить список профилей сотрудников
// @Tags staff
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Profile
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /teachers [get]
// @Router /administrators [get]
func (controller *StaffProfileController) List(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.profiles.List(ctx.Request.Context(), actor, controller.role)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// GetByID обрабатывает получение профиля сотрудника по идентификатору.
// @Summary Получить профиль сотрудника по идентификатору
// @Tags staff
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор профиля"
// @Success 200 {object} dto.Profile
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /teachers/{id} [get]
// @Router /administrators/{id} [get]
func (controller *StaffProfileController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.profiles.GetByID(ctx.Request.Context(), actor, controller.role, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает частичное обновление профиля сотрудника.
// @Summary Частично обновить профиль сотрудника
// @Tags staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор профиля"
// @Param request body dto.UpdateProfile true "Данные профиля"
// @Success 200 {object} dto.Profile
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /teachers/{id} [patch]
// @Router /administrators/{id} [patch]
func (controller *StaffProfileController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateProfile
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.profiles.Update(ctx.Request.Context(), actor, controller.role, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление профиля сотрудника.
// @Summary Удалить профиль сотрудника
// @Tags staff
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор профиля"
// @Success 204
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /teachers/{id} [delete]
// @Router /administrators/{id} [delete]
func (controller *StaffProfileController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.profiles.Delete(ctx.Request.Context(), actor, controller.role, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
