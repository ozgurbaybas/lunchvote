package httpserver

import (
	"net/http"

	"github.com/ozgurbaybas/lunchvote/platform/config"
	"github.com/ozgurbaybas/lunchvote/platform/logger"
)

type healthResponse struct {
	Status string `json:"status"`
}

func New(
	cfg config.Config,
	logg *logger.Logger,
	registerRoutes ...func(mux *http.ServeMux),
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, healthResponse{Status: "ok"})
	})

	for _, register := range registerRoutes {
		register(mux)
	}

	handler := Chain(
		mux,
		WithRecovery(logg),
		WithRequestID,
		WithRequestLogging(logg),
	)

	return &http.Server{
		Addr:    cfg.HTTPAddress(),
		Handler: handler,
	}
}
