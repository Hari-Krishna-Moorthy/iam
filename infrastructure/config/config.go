package config

import (
	"os"
)

type Config struct {
	DBURL    string
	RedisURL string
	Port     string
}

func Load() *Config {
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
