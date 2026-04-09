package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

func TestTeamRepository_SaveAndGetByID(t *testing.T) {
	pool := newTestPool(t)
	resetIdentityTables(t, pool)

	userRepo := NewUserRepository(pool)
	teamRepo := NewTeamRepository(pool)
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	owner, err := domain.NewUser("66666666-6666-6666-6666-666666666666", "Owner", "owner@example.com", now)
	if err != nil {
		t.Fatalf("new owner: %v", err)
	}

	member, err := domain.NewUser("77777777-7777-7777-7777-777777777777", "Member", "member@example.com", now)
	if err != nil {
		t.Fatalf("new member: %v", err)
	}

	if err := userRepo.Save(context.Background(), owner); err != nil {
		t.Fatalf("save owner: %v", err)
	}

	if err := userRepo.Save(context.Background(), member); err != nil {
		t.Fatalf("save member: %v", err)
	}

	team, err := domain.NewTeam("88888888-8888-8888-8888-888888888888", "Backend Team", owner.ID, now)
	if err != nil {
		t.Fatalf("new team: %v", err)
	}

	if err := team.AddMember(member.ID, now.Add(time.Minute)); err != nil {
		t.Fatalf("add member: %v", err)
	}

	if err := teamRepo.Save(context.Background(), team); err != nil {
		t.Fatalf("save team: %v", err)
	}

	got, err := teamRepo.GetByID(context.Background(), team.ID)
	if err != nil {
		t.Fatalf("get team by id: %v", err)
	}

	if got.ID != team.ID {
		t.Fatalf("expected team id %s, got %s", team.ID, got.ID)
	}

	if got.OwnerID != owner.ID {
		t.Fatalf("expected owner id %s, got %s", owner.ID, got.OwnerID)
	}

	if len(got.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(got.Members))
	}
}

func TestTeamRepository_GetByID_ReturnsErrTeamNotFound(t *testing.T) {
	pool := newTestPool(t)
	resetIdentityTables(t, pool)

	repo := NewTeamRepository(pool)

	_, err := repo.GetByID(context.Background(), "99999999-9999-9999-9999-999999999999")
	if !errors.Is(err, domain.ErrTeamNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrTeamNotFound, err)
	}
}
