package domain

import (
	"testing"
	"time"
)

func TestNewRating(t *testing.T) {
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		id           string
		restaurantID string
		userID       string
		score        int
		comment      string
		wantErr      error
	}{
		{
			name:         "creates valid rating",
			id:           "rating-1",
			restaurantID: "restaurant-1",
			userID:       "user-1",
			score:        5,
			comment:      " great food ",
		},
		{
			name:         "returns error when rating id empty",
			id:           "",
			restaurantID: "restaurant-1",
			userID:       "user-1",
			score:        5,
			wantErr:      ErrInvalidRatingID,
		},
		{
			name:         "returns error when restaurant id empty",
			id:           "rating-1",
			restaurantID: "",
			userID:       "user-1",
			score:        5,
			wantErr:      ErrInvalidRestaurantID,
		},
		{
			name:         "returns error when user id empty",
			id:           "rating-1",
			restaurantID: "restaurant-1",
			userID:       "",
			score:        5,
			wantErr:      ErrInvalidUserID,
		},
		{
			name:         "returns error when score too low",
			id:           "rating-1",
			restaurantID: "restaurant-1",
			userID:       "user-1",
			score:        0,
			wantErr:      ErrInvalidRatingScore,
		},
		{
			name:         "returns error when score too high",
			id:           "rating-1",
			restaurantID: "restaurant-1",
			userID:       "user-1",
			score:        6,
			wantErr:      ErrInvalidRatingScore,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := NewRating(
				tt.id,
				tt.restaurantID,
				tt.userID,
				tt.score,
				tt.comment,
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

			if rating.Comment != "great food" {
				t.Fatalf("expected trimmed comment, got %q", rating.Comment)
			}
		})
	}
}
