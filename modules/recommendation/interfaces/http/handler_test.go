package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	recommendationapp "github.com/ozgurbaybas/lunchvote/modules/recommendation/application"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type inMemoryTeamReader struct {
	byID map[string]identitydomain.Team
}

func newInMemoryTeamReader() *inMemoryTeamReader {
	return &inMemoryTeamReader{byID: make(map[string]identitydomain.Team)}
}

func (r *inMemoryTeamReader) GetByID(_ context.Context, id string) (identitydomain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return identitydomain.Team{}, identitydomain.ErrTeamNotFound
	}
	return team, nil
}

type inMemoryRestaurantReader struct {
	items []restaurantdomain.Restaurant
}

func newInMemoryRestaurantReader() *inMemoryRestaurantReader {
	return &inMemoryRestaurantReader{items: make([]restaurantdomain.Restaurant, 0)}
}

func (r *inMemoryRestaurantReader) List(_ context.Context) ([]restaurantdomain.Restaurant, error) {
	result := make([]restaurantdomain.Restaurant, len(r.items))
	copy(result, r.items)
	return result, nil
}

type inMemoryRatingReader struct {
	byRestaurantID map[string][]ratingdomain.Rating
}

func newInMemoryRatingReader() *inMemoryRatingReader {
	return &inMemoryRatingReader{byRestaurantID: make(map[string][]ratingdomain.Rating)}
}

func (r *inMemoryRatingReader) ListByRestaurantID(_ context.Context, restaurantID string) ([]ratingdomain.Rating, error) {
	result := make([]ratingdomain.Rating, len(r.byRestaurantID[restaurantID]))
	copy(result, r.byRestaurantID[restaurantID])
	return result, nil
}

type inMemoryPollReader struct {
	byTeamID map[string][]polldomain.Poll
}

func newInMemoryPollReader() *inMemoryPollReader {
	return &inMemoryPollReader{byTeamID: make(map[string][]polldomain.Poll)}
}

func (r *inMemoryPollReader) ListByTeamID(_ context.Context, teamID string) ([]polldomain.Poll, error) {
	result := make([]polldomain.Poll, len(r.byTeamID[teamID]))
	copy(result, r.byTeamID[teamID])
	return result, nil
}

func newTestMux() *nethttp.ServeMux {
	now := time.Date(2026, time.April, 13, 12, 0, 0, 0, time.UTC)

	teams := newInMemoryTeamReader()
	restaurants := newInMemoryRestaurantReader()
	ratings := newInMemoryRatingReader()
	polls := newInMemoryPollReader()

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "user-1", now)
	teams.byID["team-1"] = team

	restaurantOne, _ := restaurantdomain.NewRestaurant(
		"restaurant-1", "Kebap House", "Addr1", "Istanbul", "Bakirkoy",
		[]restaurantdomain.MealCard{restaurantdomain.MealCardTicket}, now,
	)
	restaurantTwo, _ := restaurantdomain.NewRestaurant(
		"restaurant-2", "Burger Spot", "Addr2", "Istanbul", "Kadikoy",
		[]restaurantdomain.MealCard{restaurantdomain.MealCardMultinet}, now,
	)
	restaurants.items = []restaurantdomain.Restaurant{restaurantOne, restaurantTwo}

	ratings.byRestaurantID["restaurant-1"] = []ratingdomain.Rating{
		{RestaurantID: "restaurant-1", UserID: "user-1", Score: 5},
	}

	polls.byTeamID["team-1"] = []polldomain.Poll{
		{
			ID:     "poll-1",
			TeamID: "team-1",
			Votes: []polldomain.Vote{
				{UserID: "user-1", RestaurantID: "restaurant-1"},
			},
		},
	}

	service := recommendationapp.NewService(teams, restaurants, ratings, polls)
	handler := NewHandler(service)

	mux := nethttp.NewServeMux()
	RegisterRoutes(mux, handler)
	return mux
}

func TestRecommendRestaurants_ReturnsOK(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/teams/team-1/recommendations", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(response))
	}

	if response[0]["restaurant_id"] != "restaurant-1" {
		t.Fatalf("expected restaurant-1 first, got %v", response[0]["restaurant_id"])
	}
}

func TestRecommendRestaurants_ReturnsNotFoundWhenTeamMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/teams/missing-team/recommendations", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestRecommendRestaurants_ReturnsBadRequestWhenLimitInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/teams/team-1/recommendations?limit=abc", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}
