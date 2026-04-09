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

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}
