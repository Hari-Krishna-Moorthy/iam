package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL    string
	RedisURL string
	Port     string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		DBURL:    getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/iam?sslmode=disable"),
		RedisURL: getEnv("REDIS_URL", "localhost:6379"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
