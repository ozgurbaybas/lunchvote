package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/platform/config"
)

const (
	testRatingUserID       = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	testRatingRestaurantID = "rating-restaurant-1"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := config.Load()

	poolConfig, err := pgxpool.ParseConfig(cfg.PostgresDSN())
	if err != nil {
		t.Fatalf("parse postgres config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		t.Fatalf("create postgres pool: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping postgres: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func resetRatingTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queries := []string{
		`DELETE FROM ratings`,
		`DELETE FROM restaurants WHERE id = 'rating-restaurant-1'`,
		`DELETE FROM users WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			t.Fatalf("reset table with query %q: %v", query, err)
		}
	}
}

func seedRatingDependencies(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, name, email, created_at)
		VALUES ($1, $2, $3, $4)
	`,
		testRatingUserID,
		"Rating Test User",
		"rating-test-user@example.com",
		time.Date(2026, time.April, 11, 11, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("insert test user: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO restaurants (id, name, address, city, district, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		testRatingRestaurantID,
		"Rating Test Restaurant",
		"Test Address No 1",
		"Istanbul",
		"Bakirkoy",
		true,
		time.Date(2026, time.April, 11, 11, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("insert test restaurant: %v", err)
	}
}
