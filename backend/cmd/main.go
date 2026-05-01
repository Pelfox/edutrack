package main

import (
	"os"

	"github.com/Pelfox/edutrack/backend/internal/app"
	"github.com/Pelfox/edutrack/backend/internal/config"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Загружаем конфигурационные значения для приложения.
	appConfig, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load application configuration")
	}

	// Запускаем приложение.
	if err := app.StartApp(logger, appConfig); err != nil {
		logger.Fatal().Err(err).Msg("failed to start up the application")
	}
}
