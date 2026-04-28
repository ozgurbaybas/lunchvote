package httpserver

import (
	"net/http"
	"time"

	"github.com/ozgurbaybas/lunchvote/platform/logger"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
		}

		w.Header().Set(RequestIDHeader, requestID)
		next.ServeHTTP(w, r)
	})
}

func WithRecovery(logg *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logg.Error(
						"panic recovered",
						"method", r.Method,
						"path", r.URL.Path,
						"request_id", r.Header.Get(RequestIDHeader),
						"error", recovered,
					)

					WriteError(w, http.StatusInternalServerError, "internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func WithRequestLogging(logg *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()

			next.ServeHTTP(w, r)

			logg.Info(
				"http request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(startedAt).Milliseconds(),
				"request_id", r.Header.Get(RequestIDHeader),
			)
		})
	}
}
