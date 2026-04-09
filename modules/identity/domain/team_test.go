package domain

import (
	"testing"
	"time"
)

func TestNewTeam_AddsOwnerAsMember(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	team, err := NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if team.OwnerID != "user-1" {
		t.Fatalf("expected owner id user-1, got %s", team.OwnerID)
	}

	if len(team.Members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(team.Members))
	}

	if team.Members[0].UserID != "user-1" {
		t.Fatalf("expected owner member user-1, got %s", team.Members[0].UserID)
	}

	if team.Members[0].Role != MembershipRoleOwner {
		t.Fatalf("expected owner role, got %s", team.Members[0].Role)
	}
}

func TestTeam_AddMember(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	team, err := NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := team.AddMember("user-2", now); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !team.HasMember("user-2") {
		t.Fatalf("expected team to contain user-2")
	}
}

func TestTeam_AddMember_ReturnsErrorWhenMemberAlreadyExists(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	team, err := NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = team.AddMember("user-1", now)
	if err != ErrMemberAlreadyExists {
		t.Fatalf("expected error %v, got %v", ErrMemberAlreadyExists, err)
	}
}

func TestTeam_RemoveMember(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	team, err := NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := team.AddMember("user-2", now); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := team.RemoveMember("user-2"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if team.HasMember("user-2") {
		t.Fatalf("expected user-2 to be removed")
	}
}

func TestTeam_RemoveMember_ReturnsErrorWhenMemberNotFound(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	team, err := NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = team.RemoveMember("user-2")
	if err != ErrMemberNotFound {
		t.Fatalf("expected error %v, got %v", ErrMemberNotFound, err)
	}
}

func TestTeam_RemoveMember_ReturnsErrorWhenRemovingOwner(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	team, err := NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = team.RemoveMember("user-1")
	if err != ErrOwnerCannotBeRemoved {
		t.Fatalf("expected error %v, got %v", ErrOwnerCannotBeRemoved, err)
	}
}
