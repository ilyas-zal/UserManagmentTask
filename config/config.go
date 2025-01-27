package config

import (
	"os"
)

// DBConfig содержит параметры подключения к базе данных
type DBConfig struct {
	DSN string
}

// LoadConfig загружает конфигурацию из окружения
func LoadConfig() DBConfig {
	return DBConfig{
		DSN: os.Getenv("DATABASE_URL"), // Убедитесь, что переменная окружения DATABASE_URL установлена
	}
}
