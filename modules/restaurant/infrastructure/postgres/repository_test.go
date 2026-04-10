package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

func TestRepository_SaveAndList(t *testing.T) {
	pool := newTestPool(t)
	resetRestaurantTables(t, pool)

	repo := NewRepository(pool)
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	restaurant, err := domain.NewRestaurant(
		"restaurant-1",
		"Kebap House",
		"Ataturk Caddesi No 10",
		"Istanbul",
		"Bakirkoy",
		[]domain.MealCard{domain.MealCardTicket, domain.MealCardMultinet},
		now,
	)
	if err != nil {
		t.Fatalf("new restaurant: %v", err)
	}

	if err := repo.Save(context.Background(), restaurant); err != nil {
		t.Fatalf("save restaurant: %v", err)
	}

	restaurants, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list restaurants: %v", err)
	}

	if len(restaurants) != 1 {
		t.Fatalf("expected 1 restaurant, got %d", len(restaurants))
	}

	got := restaurants[0]

	if got.ID != restaurant.ID {
		t.Fatalf("expected id %s, got %s", restaurant.ID, got.ID)
	}

	if got.Name != restaurant.Name {
		t.Fatalf("expected name %s, got %s", restaurant.Name, got.Name)
	}

	if len(got.SupportedMealCards) != 2 {
		t.Fatalf("expected 2 meal cards, got %d", len(got.SupportedMealCards))
	}
}

func TestRepository_List_ReturnsEmptySliceWhenNoRestaurants(t *testing.T) {
	pool := newTestPool(t)
	resetRestaurantTables(t, pool)

	repo := NewRepository(pool)

	restaurants, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list restaurants: %v", err)
	}

	if len(restaurants) != 0 {
		t.Fatalf("expected 0 restaurants, got %d", len(restaurants))
	}
}

func TestRepository_List_ReturnsRestaurantsOrderedByCreatedAt(t *testing.T) {
	pool := newTestPool(t)
	resetRestaurantTables(t, pool)

	repo := NewRepository(pool)

	firstTime := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)
	secondTime := time.Date(2026, time.April, 10, 13, 0, 0, 0, time.UTC)

	firstRestaurant, err := domain.NewRestaurant(
		"restaurant-1",
		"Kebap House",
		"Ataturk Caddesi No 10",
		"Istanbul",
		"Bakirkoy",
		[]domain.MealCard{domain.MealCardTicket},
		firstTime,
	)
	if err != nil {
		t.Fatalf("new first restaurant: %v", err)
	}

	secondRestaurant, err := domain.NewRestaurant(
		"restaurant-2",
		"Burger Spot",
		"Inonu Caddesi No 20",
		"Istanbul",
		"Kadikoy",
		[]domain.MealCard{domain.MealCardMultinet},
		secondTime,
	)
	if err != nil {
		t.Fatalf("new second restaurant: %v", err)
	}

	if err := repo.Save(context.Background(), secondRestaurant); err != nil {
		t.Fatalf("save second restaurant: %v", err)
	}

	if err := repo.Save(context.Background(), firstRestaurant); err != nil {
		t.Fatalf("save first restaurant: %v", err)
	}

	restaurants, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list restaurants: %v", err)
	}

	if len(restaurants) != 2 {
		t.Fatalf("expected 2 restaurants, got %d", len(restaurants))
	}

	if restaurants[0].ID != "restaurant-1" {
		t.Fatalf("expected first restaurant id restaurant-1, got %s", restaurants[0].ID)
	}

	if restaurants[1].ID != "restaurant-2" {
		t.Fatalf("expected second restaurant id restaurant-2, got %s", restaurants[1].ID)
	}
}
