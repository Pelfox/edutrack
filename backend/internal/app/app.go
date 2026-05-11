package app

import (
	"context"
	"fmt"

	"github.com/Pelfox/edutrack/backend/internal/config"
	"github.com/Pelfox/edutrack/backend/internal/database"
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

	// Запускаем HTTP-сервер и ожидаем новые подключения.
	if err := router.Run(appConfig.ListenAddr); err != nil {
		return fmt.Errorf("failed to start up the server: %w", err)
	}

	return nil
}
