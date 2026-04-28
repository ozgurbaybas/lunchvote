package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ozgurbaybas/lunchvote/platform/config"
	"github.com/ozgurbaybas/lunchvote/platform/logger"
)

func TestServer_HealthEndpointReturnsOK(t *testing.T) {
	cfg := config.Config{
		AppEnv:  "test",
		AppPort: "8080",
	}
	logg := logger.New("test")

	server := New(cfg, logg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get(RequestIDHeader) == "" {
		t.Fatalf("expected request id header")
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}

func TestServer_PreservesRequestID(t *testing.T) {
	cfg := config.Config{
		AppEnv:  "test",
		AppPort: "8080",
	}
	logg := logger.New("test")

	server := New(cfg, logg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(RequestIDHeader, "test-request-id")

	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Header().Get(RequestIDHeader) != "test-request-id" {
		t.Fatalf("expected request id to be preserved")
	}
}

func TestWithRecovery_ReturnsInternalServerErrorOnPanic(t *testing.T) {
	logg := logger.New("test")

	handler := Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("boom")
		}),
		WithRecovery(logg),
		WithRequestID,
	)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}
