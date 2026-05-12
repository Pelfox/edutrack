package controllers

import (
	"net/http"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthController обрабатывает HTTP-запросы авторизации.
type AuthController struct {
	auth services.AuthService
}

// NewAuthController создаёт контроллер авторизации.
func NewAuthController(auth services.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

// RegisterRoutes регистрирует маршруты авторизации.
func (controller *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", controller.Login)
}

// Login обрабатывает вход пользователя в аккаунт.
func (controller *AuthController) Login(ctx *gin.Context) {
	var request dto.Login
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondError(ctx, services.ErrInvalidInput)
		return
	}

	output, err := controller.auth.Login(ctx.Request.Context(), request)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
