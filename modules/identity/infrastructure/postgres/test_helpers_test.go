package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/platform/config"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := config.Load()
	dsn := cfg.PostgresDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
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

func resetIdentityTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queries := []string{
		`DELETE FROM team_memberships`,
		`DELETE FROM teams`,
		`DELETE FROM users`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			t.Fatalf("reset table with query %q: %v", query, err)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
