package domain

import (
	"testing"
	"time"
)

func TestNewRestaurant(t *testing.T) {
	now := time.Date(2026, time.April, 9, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        string
		rname     string
		address   string
		city      string
		district  string
		mealCards []MealCard
		wantErr   error
	}{
		{
			name:      "creates valid restaurant",
			id:        "restaurant-1",
			rname:     "Kebap House",
			address:   "Ataturk Caddesi No 10",
			city:      "Istanbul",
			district:  "Bakirkoy",
			mealCards: []MealCard{MealCardTicket, MealCardMultinet},
		},
		{
			name:      "returns error when name is empty",
			id:        "restaurant-1",
			rname:     "",
			address:   "Ataturk Caddesi No 10",
			city:      "Istanbul",
			district:  "Bakirkoy",
			mealCards: []MealCard{MealCardTicket},
			wantErr:   ErrInvalidRestaurantName,
		},
		{
			name:      "returns error when city is empty",
			id:        "restaurant-1",
			rname:     "Kebap House",
			address:   "Ataturk Caddesi No 10",
			city:      "",
			district:  "Bakirkoy",
			mealCards: []MealCard{MealCardTicket},
			wantErr:   ErrInvalidRestaurantCity,
		},
		{
			name:      "returns error when district is empty",
			id:        "restaurant-1",
			rname:     "Kebap House",
			address:   "Ataturk Caddesi No 10",
			city:      "Istanbul",
			district:  "",
			mealCards: []MealCard{MealCardTicket},
			wantErr:   ErrInvalidRestaurantDistrict,
		},
		{
			name:      "returns error when meal cards are empty",
			id:        "restaurant-1",
			rname:     "Kebap House",
			address:   "Ataturk Caddesi No 10",
			city:      "Istanbul",
			district:  "Bakirkoy",
			mealCards: nil,
			wantErr:   ErrInvalidMealCard,
		},
		{
			name:      "returns error when meal card is invalid",
			id:        "restaurant-1",
			rname:     "Kebap House",
			address:   "Ataturk Caddesi No 10",
			city:      "Istanbul",
			district:  "Bakirkoy",
			mealCards: []MealCard{"unknown"},
			wantErr:   ErrInvalidMealCard,
		},
		{
			name:      "returns error when meal card is duplicate",
			id:        "restaurant-1",
			rname:     "Kebap House",
			address:   "Ataturk Caddesi No 10",
			city:      "Istanbul",
			district:  "Bakirkoy",
			mealCards: []MealCard{MealCardTicket, MealCardTicket},
			wantErr:   ErrDuplicateMealCard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRestaurant(
				tt.id,
				tt.rname,
				tt.address,
				tt.city,
				tt.district,
				tt.mealCards,
				now,
			)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if got.ID != tt.id {
				t.Fatalf("expected id %s, got %s", tt.id, got.ID)
			}

			if !got.IsActive {
				t.Fatalf("expected restaurant to be active")
			}

			if got.CreatedAt != now {
				t.Fatalf("expected created at %v, got %v", now, got.CreatedAt)
			}
		})
	}
}
