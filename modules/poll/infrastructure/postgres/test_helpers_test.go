package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/platform/config"
)

const (
	testPollTeamID        = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	testPollOwnerID       = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	testPollMemberID      = "dddddddd-dddd-dddd-dddd-dddddddddddd"
	testPollRestaurantOne = "poll-restaurant-1"
	testPollRestaurantTwo = "poll-restaurant-2"
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

func resetPollTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queries := []string{
		`DELETE FROM poll_votes`,
		`DELETE FROM poll_options`,
		`DELETE FROM polls`,
		`DELETE FROM team_memberships WHERE team_id = 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'`,
		`DELETE FROM teams WHERE id = 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'`,
		`DELETE FROM restaurants WHERE id IN ('poll-restaurant-1', 'poll-restaurant-2')`,
		`DELETE FROM users WHERE id IN ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'dddddddd-dddd-dddd-dddd-dddddddddddd')`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			t.Fatalf("reset table with query %q: %v", query, err)
		}
	}
}

func seedPollDependencies(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2026, time.April, 12, 11, 0, 0, 0, time.UTC)

	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, name, email, created_at)
		VALUES ($1, $2, $3, $4)
	`, testPollOwnerID, "Poll Owner", "poll-owner@example.com", now)
	if err != nil {
		t.Fatalf("insert poll owner: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, name, email, created_at)
		VALUES ($1, $2, $3, $4)
	`, testPollMemberID, "Poll Member", "poll-member@example.com", now)
	if err != nil {
		t.Fatalf("insert poll member: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, owner_id, created_at)
		VALUES ($1, $2, $3, $4)
	`, testPollTeamID, "Backend Team", testPollOwnerID, now)
	if err != nil {
		t.Fatalf("insert team: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO team_memberships (team_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4), ($1, $5, $6, $4)
	`,
		testPollTeamID,
		testPollOwnerID, "owner",
		now,
		testPollMemberID, "member",
	)
	if err != nil {
		t.Fatalf("insert team memberships: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO restaurants (id, name, address, city, district, is_active, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7),
			($8, $9, $10, $11, $12, $13, $14)
	`,
		testPollRestaurantOne, "Poll Restaurant One", "Address 1", "Istanbul", "Bakirkoy", true, now,
		testPollRestaurantTwo, "Poll Restaurant Two", "Address 2", "Istanbul", "Kadikoy", true, now,
	)
	if err != nil {
		t.Fatalf("insert restaurants: %v", err)
	}
}
