package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

func TestService_CreateUser(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	user, err := service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-1",
		Name:  "Ozgur",
		Email: "ozgur@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != "user-1" {
		t.Fatalf("expected user id user-1, got %s", user.ID)
	}

	if user.Email != "ozgur@example.com" {
		t.Fatalf("expected email ozgur@example.com, got %s", user.Email)
	}
}

func TestService_CreateUser_ReturnsErrorWhenEmailAlreadyExists(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-1",
		Name:  "Ozgur",
		Email: "ozgur@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-2",
		Name:  "Another",
		Email: "ozgur@example.com",
	})
	if !errors.Is(err, ErrUserEmailAlreadyExists) {
		t.Fatalf("expected error %v, got %v", ErrUserEmailAlreadyExists, err)
	}
}

func TestService_CreateTeam(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-1",
		Name:  "Ozgur",
		Email: "ozgur@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	team, err := service.CreateTeam(context.Background(), CreateTeamCommand{
		ID:      "team-1",
		Name:    "Backend Team",
		OwnerID: "user-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if team.OwnerID != "user-1" {
		t.Fatalf("expected owner id user-1, got %s", team.OwnerID)
	}
}

func TestService_CreateTeam_ReturnsErrorWhenOwnerDoesNotExist(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.CreateTeam(context.Background(), CreateTeamCommand{
		ID:      "team-1",
		Name:    "Backend Team",
		OwnerID: "missing-user",
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrUserNotFound, err)
	}
}

func TestService_AddTeamMember(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-1",
		Name:  "Owner",
		Email: "owner@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-2",
		Name:  "Member",
		Email: "member@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateTeam(context.Background(), CreateTeamCommand{
		ID:      "team-1",
		Name:    "Backend Team",
		OwnerID: "user-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	team, err := service.AddTeamMember(context.Background(), AddTeamMemberCommand{
		TeamID: "team-1",
		UserID: "user-2",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !team.HasMember("user-2") {
		t.Fatalf("expected team to contain user-2")
	}
}

func TestService_AddTeamMember_ReturnsErrorWhenTeamDoesNotExist(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.AddTeamMember(context.Background(), AddTeamMemberCommand{
		TeamID: "missing-team",
		UserID: "user-1",
	})
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrTeamNotFound, err)
	}
}

func TestService_AddTeamMember_ReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-1",
		Name:  "Owner",
		Email: "owner@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateTeam(context.Background(), CreateTeamCommand{
		ID:      "team-1",
		Name:    "Backend Team",
		OwnerID: "user-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.AddTeamMember(context.Background(), AddTeamMemberCommand{
		TeamID: "team-1",
		UserID: "missing-user",
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrUserNotFound, err)
	}
}

func TestService_AddTeamMember_ReturnsErrorWhenMemberAlreadyExists(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	service := NewService(users, teams, func() time.Time { return now })

	_, err := service.CreateUser(context.Background(), CreateUserCommand{
		ID:    "user-1",
		Name:  "Owner",
		Email: "owner@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateTeam(context.Background(), CreateTeamCommand{
		ID:      "team-1",
		Name:    "Backend Team",
		OwnerID: "user-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.AddTeamMember(context.Background(), AddTeamMemberCommand{
		TeamID: "team-1",
		UserID: "user-1",
	})
	if !errors.Is(err, domain.ErrMemberAlreadyExists) {
		t.Fatalf("expected error %v, got %v", domain.ErrMemberAlreadyExists, err)
	}
}
