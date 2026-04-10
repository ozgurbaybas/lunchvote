package application

import (
	"context"
	"errors"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type Clock func() time.Time

type Service struct {
	ratings     ratingdomain.Repository
	users       identitydomain.UserRepository
	restaurants restaurantdomain.Repository
	now         Clock
}

func NewService(
	ratings ratingdomain.Repository,
	users identitydomain.UserRepository,
	restaurants restaurantdomain.Repository,
	now Clock,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		ratings:     ratings,
		users:       users,
		restaurants: restaurants,
		now:         now,
	}
}

func (s *Service) CreateRating(ctx context.Context, cmd CreateRatingCommand) (ratingdomain.Rating, error) {
	if _, err := s.users.GetByID(ctx, cmd.UserID); err != nil {
		if errors.Is(err, identitydomain.ErrUserNotFound) {
			return ratingdomain.Rating{}, identitydomain.ErrUserNotFound
		}
		return ratingdomain.Rating{}, err
	}

	restaurants, err := s.restaurants.List(ctx)
	if err != nil {
		return ratingdomain.Rating{}, err
	}

	foundRestaurant := false
	for _, restaurant := range restaurants {
		if restaurant.ID == cmd.RestaurantID {
			foundRestaurant = true
			break
		}
	}
	if !foundRestaurant {
		return ratingdomain.Rating{}, restaurantdomain.ErrRestaurantNotFound
	}

	_, err = s.ratings.GetByRestaurantAndUser(ctx, cmd.RestaurantID, cmd.UserID)
	if err == nil {
		return ratingdomain.Rating{}, ratingdomain.ErrRatingAlreadyExists
	}
	if !errors.Is(err, ratingdomain.ErrRatingNotFound) {
		return ratingdomain.Rating{}, err
	}

	rating, err := ratingdomain.NewRating(
		cmd.ID,
		cmd.RestaurantID,
		cmd.UserID,
		cmd.Score,
		cmd.Comment,
		s.now(),
	)
	if err != nil {
		return ratingdomain.Rating{}, err
	}

	if err := s.ratings.Save(ctx, rating); err != nil {
		return ratingdomain.Rating{}, err
	}

	return rating, nil
}

func (s *Service) ListRatingsByRestaurant(ctx context.Context, restaurantID string) ([]ratingdomain.Rating, error) {
	return s.ratings.ListByRestaurantID(ctx, restaurantID)
}
