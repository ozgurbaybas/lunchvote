package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_UsesDefaultsWhenEnvMissing(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_PORT", "")
	t.Setenv("POSTGRES_HOST", "")

	cfg := Load()

	if cfg.AppEnv != "local" {
		t.Fatalf("expected default app env local, got %s", cfg.AppEnv)
	}

	if cfg.AppPort != "8080" {
		t.Fatalf("expected default app port 8080, got %s", cfg.AppPort)
	}

	if cfg.PostgresHost != "localhost" {
		t.Fatalf("expected default postgres host localhost, got %s", cfg.PostgresHost)
	}
}

func TestLoad_UsesEnvValuesWhenProvided(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "15s")
	t.Setenv("POSTGRES_HOST", "db")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_DB", "lunchvote_test")
	t.Setenv("POSTGRES_USER", "tester")
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("POSTGRES_SSLMODE", "disable")
	t.Setenv("POSTGRES_MAX_CONNS", "20")
	t.Setenv("POSTGRES_MIN_CONNS", "2")

	cfg := Load()

	if cfg.AppEnv != "test" {
		t.Fatalf("expected app env test, got %s", cfg.AppEnv)
	}

	if cfg.AppPort != "9090" {
		t.Fatalf("expected app port 9090, got %s", cfg.AppPort)
	}

	if cfg.AppShutdownTimeout != 15*time.Second {
		t.Fatalf("expected shutdown timeout 15s, got %s", cfg.AppShutdownTimeout)
	}

	if cfg.PostgresHost != "db" {
		t.Fatalf("expected postgres host db, got %s", cfg.PostgresHost)
	}
}

func TestLoad_InvalidEnvFallsBackToDefaults(t *testing.T) {
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "invalid")
	t.Setenv("POSTGRES_MAX_CONNS", "invalid")
	t.Setenv("POSTGRES_MIN_CONNS", "invalid")

	cfg := Load()

	if cfg.AppShutdownTimeout != 10*time.Second {
		t.Fatalf("expected default shutdown timeout 10s, got %s", cfg.AppShutdownTimeout)
	}

	if cfg.PostgresMaxConns != 10 {
		t.Fatalf("expected default postgres max conns 10, got %d", cfg.PostgresMaxConns)
	}

	if cfg.PostgresMinConns != 1 {
		t.Fatalf("expected default postgres min conns 1, got %d", cfg.PostgresMinConns)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
