package application

import (
	"context"
	"errors"
	"testing"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

func TestService_RecommendRestaurants(t *testing.T) {
	now := time.Date(2026, time.April, 13, 12, 0, 0, 0, time.UTC)

	teams := newInMemoryTeamReader()
	restaurants := newInMemoryRestaurantReader()
	ratings := newInMemoryRatingReader()
	polls := newInMemoryPollReader()

	service := NewService(teams, restaurants, ratings, polls)

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
		{RestaurantID: "restaurant-1", UserID: "user-2", Score: 4},
	}
	ratings.byRestaurantID["restaurant-2"] = []ratingdomain.Rating{
		{RestaurantID: "restaurant-2", UserID: "user-3", Score: 3},
	}

	polls.byTeamID["team-1"] = []polldomain.Poll{
		{
			ID:     "poll-1",
			TeamID: "team-1",
			Votes: []polldomain.Vote{
				{UserID: "user-1", RestaurantID: "restaurant-1"},
				{UserID: "user-2", RestaurantID: "restaurant-1"},
			},
		},
	}

	items, err := service.RecommendRestaurants(context.Background(), RecommendRestaurantsQuery{
		TeamID: "team-1",
		Limit:  5,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(items))
	}

	if items[0].RestaurantID != "restaurant-1" {
		t.Fatalf("expected restaurant-1 first, got %s", items[0].RestaurantID)
	}
}

func TestService_RecommendRestaurants_ReturnsErrorWhenTeamMissing(t *testing.T) {
	service := NewService(
		newInMemoryTeamReader(),
		newInMemoryRestaurantReader(),
		newInMemoryRatingReader(),
		newInMemoryPollReader(),
	)

	_, err := service.RecommendRestaurants(context.Background(), RecommendRestaurantsQuery{
		TeamID: "missing-team",
	})
	if !errors.Is(err, identitydomain.ErrTeamNotFound) {
		t.Fatalf("expected error %v, got %v", identitydomain.ErrTeamNotFound, err)
	}
}

func TestService_RecommendRestaurants_FiltersInactiveRestaurants(t *testing.T) {
	now := time.Date(2026, time.April, 13, 12, 0, 0, 0, time.UTC)

	teams := newInMemoryTeamReader()
	restaurants := newInMemoryRestaurantReader()
	ratings := newInMemoryRatingReader()
	polls := newInMemoryPollReader()

	service := NewService(teams, restaurants, ratings, polls)

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "user-1", now)
	teams.byID["team-1"] = team

	activeRestaurant, _ := restaurantdomain.NewRestaurant(
		"restaurant-1", "Kebap House", "Addr1", "Istanbul", "Bakirkoy",
		[]restaurantdomain.MealCard{restaurantdomain.MealCardTicket}, now,
	)

	inactiveRestaurant, _ := restaurantdomain.NewRestaurant(
		"restaurant-2", "Burger Spot", "Addr2", "Istanbul", "Kadikoy",
		[]restaurantdomain.MealCard{restaurantdomain.MealCardMultinet}, now,
	)
	inactiveRestaurant.IsActive = false

	restaurants.items = []restaurantdomain.Restaurant{activeRestaurant, inactiveRestaurant}

	items, err := service.RecommendRestaurants(context.Background(), RecommendRestaurantsQuery{
		TeamID: "team-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(items))
	}

	if items[0].RestaurantID != "restaurant-1" {
		t.Fatalf("expected restaurant-1, got %s", items[0].RestaurantID)
	}
}
