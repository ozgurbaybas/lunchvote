package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	identityapp "github.com/ozgurbaybas/lunchvote/modules/identity/application"
	identitypostgres "github.com/ozgurbaybas/lunchvote/modules/identity/infrastructure/postgres"
	identityhttp "github.com/ozgurbaybas/lunchvote/modules/identity/interfaces/http"
	"github.com/ozgurbaybas/lunchvote/platform/config"
	"github.com/ozgurbaybas/lunchvote/platform/httpserver"
	"github.com/ozgurbaybas/lunchvote/platform/logger"
	"github.com/ozgurbaybas/lunchvote/platform/postgres"
)

type App struct {
	cfg    config.Config
	logger *logger.Logger
	db     *postgres.DB
	server *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg := config.Load()

	logg := logger.New(cfg.AppEnv)

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres connection: %w", err)
	}

	userRepo := identitypostgres.NewUserRepository(db.Pool)
	teamRepo := identitypostgres.NewTeamRepository(db.Pool)
	identityService := identityapp.NewService(userRepo, teamRepo, nil)
	identityHandler := identityhttp.NewHandler(identityService)

	server := httpserver.New(
		cfg,
		logg,
		func(mux *http.ServeMux) {
			identityhttp.RegisterRoutes(mux, identityHandler)
		},
	)

	return &App{
		cfg:    cfg,
		logger: logg,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("http server starting", "addr", a.server.Addr)

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server failed: %w", err)
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")

		if err := a.server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		if err := a.db.Close(context.Background()); err != nil {
			return fmt.Errorf("close postgres connection: %w", err)
		}

		a.logger.Info("application stopped")
		return nil
	}
}
