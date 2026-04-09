package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/platform/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = cfg.PostgresMaxConns
	poolConfig.MinConns = cfg.PostgresMinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		db.Pool.Close()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
