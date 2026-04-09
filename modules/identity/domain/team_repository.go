package domain

import (
	"context"
	"errors"
)

var ErrTeamNotFound = errors.New("team not found")

type TeamRepository interface {
	Save(ctx context.Context, team Team) error
	GetByID(ctx context.Context, id string) (Team, error)
}
