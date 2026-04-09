package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/ozgurbaybas/lunchvote/platform/config"
	"github.com/ozgurbaybas/lunchvote/platform/logger"
)

type healthResponse struct {
	Status string `json:"status"`
}

func New(cfg config.Config, logg *logger.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
	})

	return &http.Server{
		Addr:    cfg.HTTPAddress(),
		Handler: withLogging(logg, mux),
	}
}

func withLogging(logg *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg.Info("http request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
