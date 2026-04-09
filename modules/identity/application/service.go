package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

var (
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
)

type Clock func() time.Time

type Service struct {
	users domain.UserRepository
	teams domain.TeamRepository
	now   Clock
}

func NewService(
	users domain.UserRepository,
	teams domain.TeamRepository,
	now Clock,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		users: users,
		teams: teams,
		now:   now,
	}
}

func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (domain.User, error) {
	_, err := s.users.GetByEmail(ctx, cmd.Email)
	if err == nil {
		return domain.User{}, ErrUserEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}

	user, err := domain.NewUser(cmd.ID, cmd.Name, cmd.Email, s.now())
	if err != nil {
		return domain.User{}, err
	}

	if err := s.users.Save(ctx, user); err != nil {
		return domain.User{}, fmt.Errorf("save user: %w", err)
	}

	return user, nil
}

func (s *Service) CreateTeam(ctx context.Context, cmd CreateTeamCommand) (domain.Team, error) {
	_, err := s.users.GetByID(ctx, cmd.OwnerID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.Team{}, domain.ErrUserNotFound
		}
		return domain.Team{}, fmt.Errorf("get owner by id: %w", err)
	}

	team, err := domain.NewTeam(cmd.ID, cmd.Name, cmd.OwnerID, s.now())
	if err != nil {
		return domain.Team{}, err
	}

	if err := s.teams.Save(ctx, team); err != nil {
		return domain.Team{}, fmt.Errorf("save team: %w", err)
	}

	return team, nil
}

func (s *Service) AddTeamMember(ctx context.Context, cmd AddTeamMemberCommand) (domain.Team, error) {
	team, err := s.teams.GetByID(ctx, cmd.TeamID)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			return domain.Team{}, domain.ErrTeamNotFound
		}
		return domain.Team{}, fmt.Errorf("get team by id: %w", err)
	}

	_, err = s.users.GetByID(ctx, cmd.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.Team{}, domain.ErrUserNotFound
		}
		return domain.Team{}, fmt.Errorf("get user by id: %w", err)
	}

	if err := team.AddMember(cmd.UserID, s.now()); err != nil {
		return domain.Team{}, err
	}

	if err := s.teams.Save(ctx, team); err != nil {
		return domain.Team{}, fmt.Errorf("save team: %w", err)
	}

	return team, nil
}
