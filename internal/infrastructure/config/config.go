package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string
}

var Cfg Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: No .env file found. Using system environment variables.")
	}

	Cfg = Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "merch_shop"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	log.Println("✅ Env successfully loaded")
}

func GetDBUrl() string {
	return "postgres://" + Cfg.DBUser + ":" + Cfg.DBPassword + "@" + Cfg.DBHost + ":" + Cfg.DBPort + "/" + Cfg.DBName + "?sslmode=" + Cfg.DBSSLMode
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
