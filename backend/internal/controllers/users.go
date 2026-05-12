package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// UserController обрабатывает HTTP-запросы пользовательского модуля.
type UserController struct {
	users services.UsersService
}

// NewUserController создаёт контроллер пользовательского модуля.
func NewUserController(users services.UsersService) *UserController {
	return &UserController{users: users}
}

// RegisterRoutes регистрирует маршруты пользовательского модуля.
func (controller *UserController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", controller.Create)
	router.GET("/:id", controller.GetByID)
	router.PATCH("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}

// Create обрабатывает создание пользователя.
// @Summary Создать пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUser true "Данные пользователя"
// @Success 201 {object} dto.User
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 409 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users [post]
func (controller *UserController) Create(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.CreateUser
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.users.Create(ctx.Request.Context(), actor, request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

// GetByID обрабатывает получение пользователя по идентификатору.
// @Summary Получить пользователя по идентификатору
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор пользователя"
// @Success 200 {object} dto.User
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users/{id} [get]
func (controller *UserController) GetByID(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	output, err := controller.users.GetByID(ctx.Request.Context(), actor, ctx.Param("id"))
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Update обрабатывает частичное обновление пользователя.
// @Summary Частично обновить пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор пользователя"
// @Param request body dto.UpdateUser true "Данные пользователя"
// @Success 200 {object} dto.User
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 409 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users/{id} [patch]
func (controller *UserController) Update(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	var request dto.UpdateUser
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.users.Update(ctx.Request.Context(), actor, ctx.Param("id"), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Delete обрабатывает удаление пользователя.
// @Summary Удалить пользователя
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "Идентификатор пользователя"
// @Success 204
// @Failure 400 {object} dto.Error
// @Failure 401 {object} dto.Error
// @Failure 403 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users/{id} [delete]
func (controller *UserController) Delete(ctx *gin.Context) {
	actor, ok := actorFromContext(ctx)
	if !ok {
		respondError(ctx, services.ErrUnauthenticatedUser)
		return
	}

	if err := controller.users.Delete(ctx.Request.Context(), actor, ctx.Param("id")); err != nil {
		respondError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
