package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv             string
	AppPort            string
	AppShutdownTimeout time.Duration

	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresSSLMode  string
	PostgresMaxConns int32
	PostgresMinConns int32
}

func Load() Config {
	return Config{
		AppEnv:             getEnv("APP_ENV", "local"),
		AppPort:            getEnv("APP_PORT", "8080"),
		AppShutdownTimeout: getDurationEnv("APP_SHUTDOWN_TIMEOUT", 10*time.Second),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5433"),
		PostgresDB:       getEnv("POSTGRES_DB", "lunchvote"),
		PostgresUser:     getEnv("POSTGRES_USER", "lunchvote"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "lunchvote"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		PostgresMaxConns: int32(getIntEnv("POSTGRES_MAX_CONNS", 10)),
		PostgresMinConns: int32(getIntEnv("POSTGRES_MIN_CONNS", 1)),
	}
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
}

func (c Config) HTTPAddress() string {
	return ":" + c.AppPort
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
