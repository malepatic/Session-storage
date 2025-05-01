package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	RedisURL        string
	PostgresURL     string
	JWTSecret       string
	TokenExpiration time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set default values
	config := &Config{
		Port:            "8080",
		RedisURL:        "redis://localhost:6379/0",
		PostgresURL:     "postgres://postgres:postgres@localhost:5432/session_db?sslmode=disable",
		JWTSecret:       "your-secret-key",
		TokenExpiration: 24 * time.Hour,
	}

	// Override with environment variables if they exist
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		config.RedisURL = redisURL
	}

	if postgresURL := os.Getenv("DB_URL"); postgresURL != "" {
		config.PostgresURL = postgresURL
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	}

	if expStr := os.Getenv("TOKEN_EXPIRATION"); expStr != "" {
		exp, err := strconv.Atoi(expStr)
		if err == nil {
			config.TokenExpiration = time.Duration(exp) * time.Hour
		}
	}

	return config, nil
}
