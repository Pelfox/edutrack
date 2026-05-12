package app

import (
	"context"
	"fmt"

	"github.com/Pelfox/edutrack/backend/internal/config"
	"github.com/Pelfox/edutrack/backend/internal/controllers"
	"github.com/Pelfox/edutrack/backend/internal/database"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// StartApp запускает приложение, настраивая корректный роутинг между компонентами приложения.
func StartApp(logger zerolog.Logger, appConfig *config.AppConfig) error {
	db, err := database.ConnectPostgres(context.Background(), appConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	router := gin.Default()
	userRepository := repositories.NewUserRepository(db)
	specialtyRepository := repositories.NewSpecialtyRepository(db)
	groupRepository := repositories.NewGroupRepository(db)
	studentRepository := repositories.NewStudentRepository(db)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userRepository, appConfig.JWTSecret)
	specialtyService := services.NewSpecialtyService(specialtyRepository)
	groupService := services.NewGroupService(groupRepository)
	studentService := services.NewStudentService(studentRepository)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	specialtyController := controllers.NewSpecialtyController(specialtyService)
	groupController := controllers.NewGroupController(groupService)
	studentController := controllers.NewStudentController(studentService)

	api := router.Group("/api")
	authController.RegisterRoutes(api.Group("/auth"))
	authorized := api.Group("", controllers.AuthMiddleware(authService))
	userController.RegisterRoutes(authorized.Group("/users"))
	specialtyController.RegisterRoutes(authorized.Group("/specialties"))
	groupController.RegisterRoutes(authorized.Group("/groups"))
	studentController.RegisterRoutes(authorized.Group("/students"))

	// Запускаем HTTP-сервер и ожидаем новые подключения.
	if err := router.Run(appConfig.ListenAddr); err != nil {
		return fmt.Errorf("failed to start up the server: %w", err)
	}

	return nil
}
