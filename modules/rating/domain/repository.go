package domain

import "context"

type Repository interface {
	Save(ctx context.Context, rating Rating) error
	GetByRestaurantAndUser(ctx context.Context, restaurantID string, userID string) (Rating, error)
	ListByRestaurantID(ctx context.Context, restaurantID string) ([]Rating, error)
}
