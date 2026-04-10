package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/platform/config"
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

func resetRestaurantTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queries := []string{
		`DELETE FROM restaurant_meal_cards`,
		`DELETE FROM restaurants`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			t.Fatalf("reset table with query %q: %v", query, err)
		}
	}
}
