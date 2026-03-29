package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName               string
	AppEnv                string
	AppPort               string
	MySQLDSN              string
	RedisAddr             string
	RedisPassword         string
	RedisDB               int
	JWTSecret             string
	JWTAccessExpireMinute int
	JWTRefreshExpireHour  int
	OpenAPIKey            string
	OpenAPISecret         string
	OpenAPITimeSkewSec    int
}

func Load() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	envPath := filepath.Join(wd, ".env")
	if err := godotenv.Overload(envPath); err != nil {
		return nil, fmt.Errorf("load env file %s: %w", envPath, err)
	}

	return &Config{
		AppName:               getEnv("APP_NAME", "go-admin"),
		AppEnv:                getEnv("APP_ENV", "dev"),
		AppPort:               getEnv("APP_PORT", "8080"),
		MySQLDSN:              getEnv("MYSQL_DSN", "root:123456@tcp(127.0.0.1:3306)/go_admin?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:             getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               getEnvInt("REDIS_DB", 0),
		JWTSecret:             getEnv("JWT_SECRET", "replace-with-long-secret"),
		JWTAccessExpireMinute: getEnvInt("JWT_ACCESS_EXPIRE_MINUTE", 120),
		JWTRefreshExpireHour:  getEnvInt("JWT_REFRESH_EXPIRE_HOUR", 168),
		OpenAPIKey:            getEnv("OPEN_API_KEY", "replace-with-open-api-key"),
		OpenAPISecret:         getEnv("OPEN_API_SECRET", "replace-with-open-api-secret"),
		OpenAPITimeSkewSec:    getEnvInt("OPEN_API_TIME_SKEW_SEC", 300),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return result
}
