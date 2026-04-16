package domain

import "testing"

func TestNewRecommendation(t *testing.T) {
	item, err := NewRecommendation("restaurant-1", 14.5, []string{"high rating", "team votes"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.RestaurantID != "restaurant-1" {
		t.Fatalf("expected restaurant id restaurant-1, got %s", item.RestaurantID)
	}

	if len(item.Reasons) != 2 {
		t.Fatalf("expected 2 reasons, got %d", len(item.Reasons))
	}
}

func TestSortByScore(t *testing.T) {
	items := []Recommendation{
		{RestaurantID: "restaurant-2", Score: 10},
		{RestaurantID: "restaurant-1", Score: 20},
		{RestaurantID: "restaurant-3", Score: 10},
	}

	SortByScore(items)

	if items[0].RestaurantID != "restaurant-1" {
		t.Fatalf("expected restaurant-1 first, got %s", items[0].RestaurantID)
	}

	if items[1].RestaurantID != "restaurant-2" {
		t.Fatalf("expected restaurant-2 second, got %s", items[1].RestaurantID)
	}
}
