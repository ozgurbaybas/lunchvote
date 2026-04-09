package application

import (
	"context"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type Clock func() time.Time

type Service struct {
	repository domain.Repository
	now        Clock
}

func NewService(repository domain.Repository, now Clock) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		repository: repository,
		now:        now,
	}
}

func (s *Service) CreateRestaurant(ctx context.Context, cmd CreateRestaurantCommand) (domain.Restaurant, error) {
	mealCards := make([]domain.MealCard, 0, len(cmd.SupportedMealCards))
	for _, card := range cmd.SupportedMealCards {
		mealCards = append(mealCards, domain.MealCard(card))
	}

	restaurant, err := domain.NewRestaurant(
		cmd.ID,
		cmd.Name,
		cmd.Address,
		cmd.City,
		cmd.District,
		mealCards,
		s.now(),
	)
	if err != nil {
		return domain.Restaurant{}, err
	}

	if err := s.repository.Save(ctx, restaurant); err != nil {
		return domain.Restaurant{}, err
	}

	return restaurant, nil
}

func (s *Service) ListRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	return s.repository.List(ctx)
}
