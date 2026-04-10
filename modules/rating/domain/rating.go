package domain

import (
	"strings"
	"time"
)

type Rating struct {
	ID           string
	RestaurantID string
	UserID       string
	Score        int
	Comment      string
	CreatedAt    time.Time
}

func NewRating(
	id string,
	restaurantID string,
	userID string,
	score int,
	comment string,
	now time.Time,
) (Rating, error) {
	if strings.TrimSpace(id) == "" {
		return Rating{}, ErrInvalidRatingID
	}

	if strings.TrimSpace(restaurantID) == "" {
		return Rating{}, ErrInvalidRestaurantID
	}

	if strings.TrimSpace(userID) == "" {
		return Rating{}, ErrInvalidUserID
	}

	if score < 1 || score > 5 {
		return Rating{}, ErrInvalidRatingScore
	}

	return Rating{
		ID:           strings.TrimSpace(id),
		RestaurantID: strings.TrimSpace(restaurantID),
		UserID:       strings.TrimSpace(userID),
		Score:        score,
		Comment:      strings.TrimSpace(comment),
		CreatedAt:    now,
	}, nil
}
