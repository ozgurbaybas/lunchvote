package application

import (
	"context"
	"strings"

	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

type inMemoryUserRepository struct {
	byID    map[string]domain.User
	byEmail map[string]domain.User
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		byID:    make(map[string]domain.User),
		byEmail: make(map[string]domain.User),
	}
}

func (r *inMemoryUserRepository) Save(_ context.Context, user domain.User) error {
	r.byID[user.ID] = user
	r.byEmail[strings.ToLower(user.Email)] = user
	return nil
}

func (r *inMemoryUserRepository) GetByID(_ context.Context, id string) (domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepository) GetByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

type inMemoryTeamRepository struct {
	byID map[string]domain.Team
}

func newInMemoryTeamRepository() *inMemoryTeamRepository {
	return &inMemoryTeamRepository{
		byID: make(map[string]domain.Team),
	}
}

func (r *inMemoryTeamRepository) Save(_ context.Context, team domain.Team) error {
	r.byID[team.ID] = team
	return nil
}

func (r *inMemoryTeamRepository) GetByID(_ context.Context, id string) (domain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return domain.Team{}, domain.ErrTeamNotFound
	}
	return team, nil
}
