package application

import (
	"context"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type inMemoryRepository struct {
	restaurants []domain.Restaurant
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		restaurants: make([]domain.Restaurant, 0),
	}
}

func (r *inMemoryRepository) Save(_ context.Context, restaurant domain.Restaurant) error {
	r.restaurants = append(r.restaurants, restaurant)
	return nil
}

func (r *inMemoryRepository) List(_ context.Context) ([]domain.Restaurant, error) {
	result := make([]domain.Restaurant, len(r.restaurants))
	copy(result, r.restaurants)
	return result, nil
}
