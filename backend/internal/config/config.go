package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port          string
	DBName        string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	JWTSecret     string
	JWTExpireTime int
}

var Config EnvConfig = loadConfig()

func loadConfig() EnvConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load .env file: %s", err)
	}

	return EnvConfig{
		Port:          getEnv("PORT", "8080"),
		DBName:        getEnv("DB_NAME", "subguard_db"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		JWTSecret:     getEnv("JWT_SECRET", "very-secret"),
		JWTExpireTime: getEnvAsInt("JWT_EXPIRE_TIME", 3600),
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	valueAsInt, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return valueAsInt
}
