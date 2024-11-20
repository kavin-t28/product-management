package config

import (
	"os"
)

var (
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	CacheHost  string
	CachePort  string
)

func LoadConfig() {
	DBHost = getEnv("DB_HOST", "localhost")
	DBPort = getEnv("DB_PORT", "5432")
	DBUser = getEnv("DB_USER", "pm_user")
	DBPassword = getEnv("DB_PASSWORD", "password")
	DBName = getEnv("DB_NAME", "product_management")
	CacheHost = getEnv("CACHE_HOST", "localhost")
	CachePort = getEnv("CACHE_PORT", "6379")
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
