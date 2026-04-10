package application

import (
	"context"
	"errors"
	"testing"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

func TestService_CreateRating(t *testing.T) {
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	ratingRepo := newInMemoryRatingRepository()
	userRepo := newInMemoryUserRepository()
	restaurantRepo := newInMemoryRestaurantRepository()
	service := NewService(ratingRepo, userRepo, restaurantRepo, func() time.Time { return now })

	user, err := identitydomain.NewUser("user-1", "Ozgur", "ozgur@example.com", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	if err := userRepo.Save(context.Background(), user); err != nil {
		t.Fatalf("save user: %v", err)
	}

	restaurant, err := restaurantdomain.NewRestaurant(
		"restaurant-1",
		"Kebap House",
		"Ataturk Caddesi No 10",
		"Istanbul",
		"Bakirkoy",
		[]restaurantdomain.MealCard{restaurantdomain.MealCardTicket},
		now,
	)
	if err != nil {
		t.Fatalf("new restaurant: %v", err)
	}
	if err := restaurantRepo.Save(context.Background(), restaurant); err != nil {
		t.Fatalf("save restaurant: %v", err)
	}

	rating, err := service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-1",
		RestaurantID: "restaurant-1",
		UserID:       "user-1",
		Score:        5,
		Comment:      " very good ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rating.Score != 5 {
		t.Fatalf("expected score 5, got %d", rating.Score)
	}

	if rating.Comment != "very good" {
		t.Fatalf("expected trimmed comment, got %q", rating.Comment)
	}
}

func TestService_CreateRating_ReturnsErrorWhenUserMissing(t *testing.T) {
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	service := NewService(
		newInMemoryRatingRepository(),
		newInMemoryUserRepository(),
		newInMemoryRestaurantRepository(),
		func() time.Time { return now },
	)

	_, err := service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-1",
		RestaurantID: "restaurant-1",
		UserID:       "missing-user",
		Score:        5,
	})
	if !errors.Is(err, identitydomain.ErrUserNotFound) {
		t.Fatalf("expected error %v, got %v", identitydomain.ErrUserNotFound, err)
	}
}

func TestService_CreateRating_ReturnsErrorWhenRestaurantMissing(t *testing.T) {
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	ratingRepo := newInMemoryRatingRepository()
	userRepo := newInMemoryUserRepository()
	restaurantRepo := newInMemoryRestaurantRepository()
	service := NewService(ratingRepo, userRepo, restaurantRepo, func() time.Time { return now })

	user, err := identitydomain.NewUser("user-1", "Ozgur", "ozgur@example.com", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	if err := userRepo.Save(context.Background(), user); err != nil {
		t.Fatalf("save user: %v", err)
	}

	_, err = service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-1",
		RestaurantID: "missing-restaurant",
		UserID:       "user-1",
		Score:        5,
	})
	if !errors.Is(err, restaurantdomain.ErrRestaurantNotFound) {
		t.Fatalf("expected error %v, got %v", restaurantdomain.ErrRestaurantNotFound, err)
	}
}

func TestService_CreateRating_ReturnsErrorWhenRatingAlreadyExists(t *testing.T) {
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	ratingRepo := newInMemoryRatingRepository()
	userRepo := newInMemoryUserRepository()
	restaurantRepo := newInMemoryRestaurantRepository()
	service := NewService(ratingRepo, userRepo, restaurantRepo, func() time.Time { return now })

	user, _ := identitydomain.NewUser("user-1", "Ozgur", "ozgur@example.com", now)
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

	_, err := service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-1",
		RestaurantID: "restaurant-1",
		UserID:       "user-1",
		Score:        5,
	})
	if err != nil {
		t.Fatalf("first create rating: %v", err)
	}

	_, err = service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-2",
		RestaurantID: "restaurant-1",
		UserID:       "user-1",
		Score:        4,
	})
	if !errors.Is(err, ratingdomain.ErrRatingAlreadyExists) {
		t.Fatalf("expected error %v, got %v", ratingdomain.ErrRatingAlreadyExists, err)
	}
}

func TestService_ListRatingsByRestaurant(t *testing.T) {
	now := time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC)

	ratingRepo := newInMemoryRatingRepository()
	userRepo := newInMemoryUserRepository()
	restaurantRepo := newInMemoryRestaurantRepository()
	service := NewService(ratingRepo, userRepo, restaurantRepo, func() time.Time { return now })

	userOne, _ := identitydomain.NewUser("user-1", "Ozgur", "ozgur@example.com", now)
	userTwo, _ := identitydomain.NewUser("user-2", "Yasmin", "yasmin@example.com", now)
	_ = userRepo.Save(context.Background(), userOne)
	_ = userRepo.Save(context.Background(), userTwo)

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

	_, _ = service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-1",
		RestaurantID: "restaurant-1",
		UserID:       "user-1",
		Score:        5,
	})
	_, _ = service.CreateRating(context.Background(), CreateRatingCommand{
		ID:           "rating-2",
		RestaurantID: "restaurant-1",
		UserID:       "user-2",
		Score:        4,
	})

	ratings, err := service.ListRatingsByRestaurant(context.Background(), "restaurant-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(ratings) != 2 {
		t.Fatalf("expected 2 ratings, got %d", len(ratings))
	}
}
