package domain

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, restaurant Restaurant) error
	List(ctx context.Context) ([]Restaurant, error)
}
