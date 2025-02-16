package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser          string
	DBPassword      string
	DBHost          string
	DBPort          string
	DBContainerPort string
	DBName          string
	DBTestName      string // для тестов
	DBSSLMode       string
	SecretJWTKey    string
}

var Cfg Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
	}

	Cfg = Config{
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBContainerPort: getEnv("DB_CONTAINER_PORT", "5435"),
		DBName:          getEnv("DB_NAME", "merch_shop"),
		DBTestName:      getEnv("DB_TEST_NAME", "merch_shop_test"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		SecretJWTKey:    getEnv("SECRET_JWT_KEY", "secret-jwt-key"),
	}

	log.Println("Env successfully loaded")
}

func GetDBUrl(isTest bool) string {
	dbName := Cfg.DBName
	if isTest {
		dbName = Cfg.DBTestName
	}
	return "postgres://" + Cfg.DBUser + ":" + Cfg.DBPassword + "@" + Cfg.DBHost + ":" + Cfg.DBPort + "/" + dbName + "?sslmode=" + Cfg.DBSSLMode
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
