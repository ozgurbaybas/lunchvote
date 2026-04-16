package application

import (
	"context"
	"errors"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

type Clock func() time.Time

type Service struct {
	polls polldomain.Repository
	teams identitydomain.TeamRepository
	now   Clock
}

func NewService(
	polls polldomain.Repository,
	teams identitydomain.TeamRepository,
	now Clock,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		polls: polls,
		teams: teams,
		now:   now,
	}
}

func (s *Service) CreatePoll(ctx context.Context, cmd CreatePollCommand) (polldomain.Poll, error) {
	team, err := s.teams.GetByID(ctx, cmd.TeamID)
	if err != nil {
		if errors.Is(err, identitydomain.ErrTeamNotFound) {
			return polldomain.Poll{}, identitydomain.ErrTeamNotFound
		}
		return polldomain.Poll{}, err
	}

	if !team.HasMember(cmd.CreatorUserID) {
		return polldomain.Poll{}, polldomain.ErrUserNotTeamMember
	}

	poll, err := polldomain.NewPoll(cmd.ID, cmd.TeamID, cmd.Title, cmd.RestaurantIDs, s.now())
	if err != nil {
		return polldomain.Poll{}, err
	}

	if err := s.polls.Save(ctx, poll); err != nil {
		return polldomain.Poll{}, err
	}

	return poll, nil
}

func (s *Service) Vote(ctx context.Context, cmd VoteCommand) (polldomain.Poll, error) {
	poll, err := s.polls.GetByID(ctx, cmd.PollID)
	if err != nil {
		if errors.Is(err, polldomain.ErrPollNotFound) {
			return polldomain.Poll{}, polldomain.ErrPollNotFound
		}
		return polldomain.Poll{}, err
	}

	team, err := s.teams.GetByID(ctx, poll.TeamID)
	if err != nil {
		if errors.Is(err, identitydomain.ErrTeamNotFound) {
			return polldomain.Poll{}, identitydomain.ErrTeamNotFound
		}
		return polldomain.Poll{}, err
	}

	if !team.HasMember(cmd.UserID) {
		return polldomain.Poll{}, polldomain.ErrUserNotTeamMember
	}

	if err := poll.Vote(cmd.UserID, cmd.RestaurantID, s.now()); err != nil {
		return polldomain.Poll{}, err
	}

	if err := s.polls.Save(ctx, poll); err != nil {
		return polldomain.Poll{}, err
	}

	return poll, nil
}

func (s *Service) Results(ctx context.Context, pollID string) (map[string]int, error) {
	poll, err := s.polls.GetByID(ctx, pollID)
	if err != nil {
		if errors.Is(err, polldomain.ErrPollNotFound) {
			return nil, polldomain.ErrPollNotFound
		}
		return nil, err
	}

	return poll.Results(), nil
}
