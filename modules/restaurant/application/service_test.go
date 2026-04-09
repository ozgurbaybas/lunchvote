package application

import (
	"context"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

func TestService_CreateRestaurant(t *testing.T) {
	now := time.Date(2026, time.April, 9, 12, 0, 0, 0, time.UTC)

	repo := newInMemoryRepository()
	service := NewService(repo, func() time.Time { return now })

	restaurant, err := service.CreateRestaurant(context.Background(), CreateRestaurantCommand{
		ID:                 "restaurant-1",
		Name:               "Kebap House",
		Address:            "Ataturk Caddesi No 10",
		City:               "Istanbul",
		District:           "Bakirkoy",
		SupportedMealCards: []string{"ticket", "multinet"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if restaurant.ID != "restaurant-1" {
		t.Fatalf("expected restaurant id restaurant-1, got %s", restaurant.ID)
	}

	if len(repo.restaurants) != 1 {
		t.Fatalf("expected repository to contain 1 restaurant, got %d", len(repo.restaurants))
	}
}

func TestService_CreateRestaurant_ReturnsDomainError(t *testing.T) {
	now := time.Date(2026, time.April, 9, 12, 0, 0, 0, time.UTC)

	repo := newInMemoryRepository()
	service := NewService(repo, func() time.Time { return now })

	_, err := service.CreateRestaurant(context.Background(), CreateRestaurantCommand{
		ID:                 "restaurant-1",
		Name:               "",
		Address:            "Ataturk Caddesi No 10",
		City:               "Istanbul",
		District:           "Bakirkoy",
		SupportedMealCards: []string{"ticket"},
	})
	if err != domain.ErrInvalidRestaurantName {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidRestaurantName, err)
	}
}

func TestService_ListRestaurants(t *testing.T) {
	now := time.Date(2026, time.April, 9, 12, 0, 0, 0, time.UTC)

	repo := newInMemoryRepository()
	service := NewService(repo, func() time.Time { return now })

	_, err := service.CreateRestaurant(context.Background(), CreateRestaurantCommand{
		ID:                 "restaurant-1",
		Name:               "Kebap House",
		Address:            "Ataturk Caddesi No 10",
		City:               "Istanbul",
		District:           "Bakirkoy",
		SupportedMealCards: []string{"ticket"},
	})
	if err != nil {
		t.Fatalf("create first restaurant: %v", err)
	}

	_, err = service.CreateRestaurant(context.Background(), CreateRestaurantCommand{
		ID:                 "restaurant-2",
		Name:               "Burger Spot",
		Address:            "Inonu Caddesi No 20",
		City:               "Istanbul",
		District:           "Kadikoy",
		SupportedMealCards: []string{"multinet"},
	})
	if err != nil {
		t.Fatalf("create second restaurant: %v", err)
	}

	restaurants, err := service.ListRestaurants(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(restaurants) != 2 {
		t.Fatalf("expected 2 restaurants, got %d", len(restaurants))
	}
}
