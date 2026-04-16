package postgres

import (
	"context"
	"testing"
	"time"

	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
)

func TestRepository_SaveAndList(t *testing.T) {
	pool := newTestPool(t)
	resetRatingTables(t, pool)
	seedRatingDependencies(t, pool)

	repo := NewRepository(pool)

	now := time.Date(2026, time.April, 11, 12, 0, 0, 0, time.UTC)

	rating, err := ratingdomain.NewRating(
		"rating-1",
		testRatingRestaurantID,
		testRatingUserID,
		5,
		"Excellent!",
		now,
	)
	if err != nil {
		t.Fatalf("new rating: %v", err)
	}

	if err := repo.Save(context.Background(), rating); err != nil {
		t.Fatalf("save rating: %v", err)
	}

	list, err := repo.ListByRestaurantID(context.Background(), testRatingRestaurantID)
	if err != nil {
		t.Fatalf("list ratings: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 rating, got %d", len(list))
	}

	if list[0].UserID != testRatingUserID {
		t.Fatalf("expected user id %s, got %s", testRatingUserID, list[0].UserID)
	}
}
