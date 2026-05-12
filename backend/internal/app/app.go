package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Pelfox/edutrack/backend/docs"
	"github.com/Pelfox/edutrack/backend/internal/config"
	"github.com/Pelfox/edutrack/backend/internal/controllers"
	"github.com/Pelfox/edutrack/backend/internal/database"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/Pelfox/edutrack/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// StartApp запускает приложение, настраивая корректный роутинг между компонентами приложения.
func StartApp(logger zerolog.Logger, appConfig *config.AppConfig) error {
	db, err := database.ConnectPostgres(context.Background(), appConfig.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(corsMiddleware())

	userRepository := repositories.NewUserRepository(db)
	specialtyRepository := repositories.NewSpecialtyRepository(db)
	groupRepository := repositories.NewGroupRepository(db)
	studentRepository := repositories.NewStudentRepository(db)
	profileRepository := repositories.NewProfileRepository(db)
	subjectRepository := repositories.NewSubjectRepository(db)
	curriculumRepository := repositories.NewCurriculumRepository(db)
	gradeRepository := repositories.NewGradeRepository(db)
	analyticsRepository := repositories.NewAnalyticsRepository(db)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(userRepository, appConfig.JWTSecret)
	specialtyService := services.NewSpecialtyService(specialtyRepository)
	groupService := services.NewGroupService(groupRepository)
	studentService := services.NewStudentService(studentRepository)
	profileService := services.NewProfileService(profileRepository)
	staffProfileService := services.NewStaffProfilesService(profileRepository, userRepository)
	subjectService := services.NewSubjectService(subjectRepository)
	curriculumService := services.NewCurriculumService(curriculumRepository)
	gradeService := services.NewGradeService(gradeRepository, studentRepository)
	analyticsService := services.NewAnalyticsService(analyticsRepository)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	specialtyController := controllers.NewSpecialtyController(specialtyService)
	groupController := controllers.NewGroupController(groupService)
	studentController := controllers.NewStudentController(studentService)
	profileController := controllers.NewProfileController(profileService)
	teacherController := controllers.NewTeacherController(staffProfileService)
	administratorController := controllers.NewAdministratorController(staffProfileService)
	subjectController := controllers.NewSubjectController(subjectService)
	curriculumController := controllers.NewCurriculumController(curriculumService)
	gradeController := controllers.NewGradeController(gradeService)
	analyticsController := controllers.NewAnalyticsController(analyticsService)

	api := router.Group("/api")
	authController.RegisterRoutes(api.Group("/auth"))
	authorized := api.Group("", controllers.AuthMiddleware(authService))
	userController.RegisterRoutes(authorized.Group("/users"))
	specialtyController.RegisterRoutes(authorized.Group("/specialties"))
	groupController.RegisterRoutes(authorized.Group("/groups"))
	studentController.RegisterRoutes(authorized.Group("/students"))
	profileController.RegisterRoutes(authorized.Group("/profile"))
	teacherController.RegisterRoutes(authorized.Group("/teachers"))
	administratorController.RegisterRoutes(authorized.Group("/administrators"))
	subjectController.RegisterRoutes(authorized.Group("/subjects"))
	curriculumController.RegisterRoutes(authorized.Group("/curriculums"))
	gradeController.RegisterRoutes(authorized.Group("/grades"))
	analyticsController.RegisterRoutes(authorized.Group("/analytics"))
	router.GET("/openapi.json", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(docs.SwaggerInfo.ReadDoc()))
	})
	router.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/openapi.json"),
		httpSwagger.PersistAuthorization(true),
	)))

	// Запускаем HTTP-сервер и ожидаем новые подключения.
	if err := router.Run(appConfig.ListenAddr); err != nil {
		return fmt.Errorf("failed to start up the server: %w", err)
	}

	return nil
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]struct{}{
		"http://localhost:3000": {},
		"http://127.0.0.1:3000": {},
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			ctx.Header("Vary", "Origin")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
