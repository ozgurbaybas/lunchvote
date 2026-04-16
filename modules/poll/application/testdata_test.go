package application

import (
	"context"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

type inMemoryPollRepository struct {
	byID map[string]polldomain.Poll
}

func newInMemoryPollRepository() *inMemoryPollRepository {
	return &inMemoryPollRepository{
		byID: make(map[string]polldomain.Poll),
	}
}

func (r *inMemoryPollRepository) Save(_ context.Context, poll polldomain.Poll) error {
	r.byID[poll.ID] = poll
	return nil
}

func (r *inMemoryPollRepository) GetByID(_ context.Context, id string) (polldomain.Poll, error) {
	poll, ok := r.byID[id]
	if !ok {
		return polldomain.Poll{}, polldomain.ErrPollNotFound
	}
	return poll, nil
}

type inMemoryTeamRepository struct {
	byID map[string]identitydomain.Team
}

func newInMemoryTeamRepository() *inMemoryTeamRepository {
	return &inMemoryTeamRepository{
		byID: make(map[string]identitydomain.Team),
	}
}

func (r *inMemoryTeamRepository) Save(_ context.Context, team identitydomain.Team) error {
	r.byID[team.ID] = team
	return nil
}

func (r *inMemoryTeamRepository) GetByID(_ context.Context, id string) (identitydomain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return identitydomain.Team{}, identitydomain.ErrTeamNotFound
	}
	return team, nil
}
