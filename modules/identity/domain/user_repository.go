package domain

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Save(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
