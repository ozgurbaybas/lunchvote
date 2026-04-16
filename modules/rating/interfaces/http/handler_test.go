package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	ratingapp "github.com/ozgurbaybas/lunchvote/modules/rating/application"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type inMemoryRatingRepository struct {
	ratings []ratingdomain.Rating
}

func newInMemoryRatingRepository() *inMemoryRatingRepository {
	return &inMemoryRatingRepository{
		ratings: make([]ratingdomain.Rating, 0),
	}
}

func (r *inMemoryRatingRepository) Save(_ context.Context, rating ratingdomain.Rating) error {
	r.ratings = append(r.ratings, rating)
	return nil
}

func (r *inMemoryRatingRepository) GetByRestaurantAndUser(_ context.Context, restaurantID string, userID string) (ratingdomain.Rating, error) {
	for _, rating := range r.ratings {
		if rating.RestaurantID == strings.TrimSpace(restaurantID) && rating.UserID == strings.TrimSpace(userID) {
			return rating, nil
		}
	}

	return ratingdomain.Rating{}, ratingdomain.ErrRatingNotFound
}

func (r *inMemoryRatingRepository) ListByRestaurantID(_ context.Context, restaurantID string) ([]ratingdomain.Rating, error) {
	result := make([]ratingdomain.Rating, 0)
	for _, rating := range r.ratings {
		if rating.RestaurantID == restaurantID {
			result = append(result, rating)
		}
	}
	return result, nil
}

type inMemoryUserRepository struct {
	byID    map[string]identitydomain.User
	byEmail map[string]identitydomain.User
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		byID:    make(map[string]identitydomain.User),
		byEmail: make(map[string]identitydomain.User),
	}
}

func (r *inMemoryUserRepository) Save(_ context.Context, user identitydomain.User) error {
	r.byID[user.ID] = user
	r.byEmail[strings.ToLower(user.Email)] = user
	return nil
}

func (r *inMemoryUserRepository) GetByID(_ context.Context, id string) (identitydomain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepository) GetByEmail(_ context.Context, email string) (identitydomain.User, error) {
	user, ok := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}
	return user, nil
}

type inMemoryRestaurantRepository struct {
	restaurants []restaurantdomain.Restaurant
}

func newInMemoryRestaurantRepository() *inMemoryRestaurantRepository {
	return &inMemoryRestaurantRepository{
		restaurants: make([]restaurantdomain.Restaurant, 0),
	}
}

func (r *inMemoryRestaurantRepository) Save(_ context.Context, restaurant restaurantdomain.Restaurant) error {
	r.restaurants = append(r.restaurants, restaurant)
	return nil
}

func (r *inMemoryRestaurantRepository) List(_ context.Context) ([]restaurantdomain.Restaurant, error) {
	result := make([]restaurantdomain.Restaurant, len(r.restaurants))
	copy(result, r.restaurants)
	return result, nil
}

func newTestMux() *nethttp.ServeMux {
	ratingRepo := newInMemoryRatingRepository()
	userRepo := newInMemoryUserRepository()
	restaurantRepo := newInMemoryRestaurantRepository()

	now := time.Date(2026, time.April, 11, 12, 0, 0, 0, time.UTC)

	user, _ := identitydomain.NewUser(
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"Rating User",
		"rating-user@example.com",
		now,
	)
	_ = userRepo.Save(context.Background(), user)

	restaurant, _ := restaurantdomain.NewRestaurant(
		"restaurant-1",
		"Kebap House",
		"Ataturk Caddesi No 10",
		"Istanbul",
		"Bakirkoy",
		[]restaurantdomain.MealCard{restaurantdomain.MealCardTicket},
		now,
	)
	_ = restaurantRepo.Save(context.Background(), restaurant)

	service := ratingapp.NewService(
		ratingRepo,
		userRepo,
		restaurantRepo,
		func() time.Time { return now },
	)
	handler := NewHandler(service)

	mux := nethttp.NewServeMux()
	RegisterRoutes(mux, handler)
	return mux
}

func TestCreateRating_ReturnsCreated(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/ratings",
		bytes.NewReader([]byte(`{
			"id":"rating-1",
			"restaurant_id":"restaurant-1",
			"user_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"score":5,
			"comment":"great food"
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestCreateRating_ReturnsBadRequestWhenBodyInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/ratings", bytes.NewReader([]byte(`{`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateRating_ReturnsBadRequestWhenRequiredFieldsMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/ratings",
		bytes.NewReader([]byte(`{
			"id":"",
			"restaurant_id":"restaurant-1",
			"user_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"score":5
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateRating_ReturnsBadRequestWhenScoreInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/ratings",
		bytes.NewReader([]byte(`{
			"id":"rating-1",
			"restaurant_id":"restaurant-1",
			"user_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"score":6
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateRating_ReturnsNotFoundWhenUserMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/ratings",
		bytes.NewReader([]byte(`{
			"id":"rating-1",
			"restaurant_id":"restaurant-1",
			"user_id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			"score":5
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestCreateRating_ReturnsNotFoundWhenRestaurantMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/ratings",
		bytes.NewReader([]byte(`{
			"id":"rating-1",
			"restaurant_id":"missing-restaurant",
			"user_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"score":5
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestCreateRating_ReturnsConflictWhenDuplicate(t *testing.T) {
	mux := newTestMux()

	body := []byte(`{
		"id":"rating-1",
		"restaurant_id":"restaurant-1",
		"user_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"score":5
	}`)

	firstReq := httptest.NewRequest(nethttp.MethodPost, "/v1/ratings", bytes.NewReader(body))
	firstRec := httptest.NewRecorder()
	mux.ServeHTTP(firstRec, firstReq)

	secondReq := httptest.NewRequest(nethttp.MethodPost, "/v1/ratings", bytes.NewReader(body))
	secondRec := httptest.NewRecorder()
	mux.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != nethttp.StatusConflict {
		t.Fatalf("expected status 409, got %d", secondRec.Code)
	}
}

func TestListRatingsByRestaurant_ReturnsEmptyArray(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/restaurants/restaurant-1/ratings", nil)
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
		t.Fatalf("expected empty array, got %d", len(response))
	}
}

func TestListRatingsByRestaurant_ReturnsRatings(t *testing.T) {
	mux := newTestMux()

	createReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/ratings",
		bytes.NewReader([]byte(`{
			"id":"rating-1",
			"restaurant_id":"restaurant-1",
			"user_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"score":5,
			"comment":"great food"
		}`)),
	)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/restaurants/restaurant-1/ratings", nil)
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
		t.Fatalf("expected 1 rating, got %d", len(response))
	}
}
