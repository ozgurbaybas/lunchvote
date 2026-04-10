package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/application"
	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type inMemoryRepository struct {
	restaurants []domain.Restaurant
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		restaurants: make([]domain.Restaurant, 0),
	}
}

func (r *inMemoryRepository) Save(_ context.Context, restaurant domain.Restaurant) error {
	r.restaurants = append(r.restaurants, restaurant)
	return nil
}

func (r *inMemoryRepository) List(_ context.Context) ([]domain.Restaurant, error) {
	result := make([]domain.Restaurant, len(r.restaurants))
	copy(result, r.restaurants)
	return result, nil
}

func newTestMux() *nethttp.ServeMux {
	repo := newInMemoryRepository()
	service := application.NewService(repo, func() time.Time {
		return time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)
	})
	handler := NewHandler(service)

	mux := nethttp.NewServeMux()
	RegisterRoutes(mux, handler)
	return mux
}

func TestCreateRestaurant_ReturnsCreated(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/restaurants",
		bytes.NewReader([]byte(`{
			"id":"restaurant-1",
			"name":"Kebap House",
			"address":"Ataturk Caddesi No 10",
			"city":"Istanbul",
			"district":"Bakirkoy",
			"supported_meal_cards":["ticket","multinet"]
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response["id"] != "restaurant-1" {
		t.Fatalf("expected id restaurant-1, got %v", response["id"])
	}
}

func TestCreateRestaurant_ReturnsBadRequestWhenBodyInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/restaurants", bytes.NewReader([]byte(`{`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateRestaurant_ReturnsBadRequestWhenRequiredFieldsMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/restaurants",
		bytes.NewReader([]byte(`{
			"id":"restaurant-1",
			"name":"",
			"city":"Istanbul",
			"district":"Bakirkoy",
			"supported_meal_cards":["ticket"]
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateRestaurant_ReturnsBadRequestWhenMealCardInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/restaurants",
		bytes.NewReader([]byte(`{
			"id":"restaurant-1",
			"name":"Kebap House",
			"address":"Ataturk Caddesi No 10",
			"city":"Istanbul",
			"district":"Bakirkoy",
			"supported_meal_cards":["unknown"]
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestListRestaurants_ReturnsEmptyArray(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/restaurants", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response []any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(response) != 0 {
		t.Fatalf("expected empty array, got %d items", len(response))
	}
}

func TestListRestaurants_ReturnsRestaurants(t *testing.T) {
	mux := newTestMux()

	createReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/restaurants",
		bytes.NewReader([]byte(`{
			"id":"restaurant-1",
			"name":"Kebap House",
			"address":"Ataturk Caddesi No 10",
			"city":"Istanbul",
			"district":"Bakirkoy",
			"supported_meal_cards":["ticket"]
		}`)),
	)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/restaurants", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 restaurant, got %d", len(response))
	}
}
