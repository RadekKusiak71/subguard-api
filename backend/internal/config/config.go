package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port       string
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
}

var Config EnvConfig = loadConfig()

func loadConfig() EnvConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load .env file: %s", err)
	}

	return EnvConfig{
		Port:       getEnv("PORT", "8080"),
		DBName:     getEnv("DB_NAME", "subguard_db"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return value
}
