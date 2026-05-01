package app

import (
	"fmt"

	"github.com/Pelfox/edutrack/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// StartApp запускает приложение, настраивая корректный роутинг между компонентами приложения.
func StartApp(logger zerolog.Logger, appConfig *config.AppConfig) error {
	router := gin.Default()

	// Запускаем HTTP-сервер и ожидаем новые подключения.
	if err := router.Run(appConfig.ListenAddr); err != nil {
		return fmt.Errorf("failed to start up the server: %w", err)
	}

	return nil
}
