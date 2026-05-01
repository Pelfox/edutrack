package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

// parser содержит в себе готовый к использованию глобальный парсер для конфигурации.
var parser = koanf.New(".")

// AppConfig содержит все конфигурационные поля для приложения.
type AppConfig struct {
	// ListenAddr - адрес, на котором необходимо запустить веб-сервер.
	// Значение по умолчанию: ":8000".
	ListenAddr string `koanf:"listen_addr"`
}

// LoadConfig загружает все значения для конфигурации из .env-файла или переменных окружения.
func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()
	config := &AppConfig{
		ListenAddr: ":8000",
	}

	// Загружаем конфигурационные значения из переменных окружения (или .env-файла).
	parser.Load(env.Provider(".", env.Opt{
		Prefix: "BACKEND_",
		TransformFunc: func(key string, value string) (string, any) {
			key = strings.ToLower(strings.TrimPrefix(key, "BACKEND_"))
			if strings.Contains(value, " ") {
				return key, strings.Split(value, " ")
			}
			return key, value
		},
	}), nil)

	// Загружаем значения в структуру конфигурации.
	if err := parser.Unmarshal("", config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return config, nil
}
