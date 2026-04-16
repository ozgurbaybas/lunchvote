package domain

import "context"

type Repository interface {
	Save(ctx context.Context, poll Poll) error
	GetByID(ctx context.Context, id string) (Poll, error)
}
