package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Переменные окружения
var (
	JWTSecret     string
	AdminUsername string
	AdminPassword string
	DBName        string
)

// InitConfig загружает переменные окружения
func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, использую переменные окружения")
	}

	JWTSecret = getEnv("JWT_SECRET", "default_secret")
	AdminUsername = getEnv("ADMIN_USERNAME", "admin")
	AdminPassword = getEnv("ADMIN_PASSWORD", "password")
	DBName = getEnv("DB_NAME", "database.db")
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
