package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

func TestUserRepository_SaveAndGetByID(t *testing.T) {
	pool := newTestPool(t)
	resetIdentityTables(t, pool)

	repo := NewUserRepository(pool)
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	user, err := domain.NewUser("11111111-1111-1111-1111-111111111111", "Ozgur", "ozgur@example.com", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}

	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("save user: %v", err)
	}

	got, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}

	if got.ID != user.ID {
		t.Fatalf("expected id %s, got %s", user.ID, got.ID)
	}

	if got.Email != user.Email {
		t.Fatalf("expected email %s, got %s", user.Email, got.Email)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	pool := newTestPool(t)
	resetIdentityTables(t, pool)

	repo := NewUserRepository(pool)
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	user, err := domain.NewUser("22222222-2222-2222-2222-222222222222", "Yasmin", "yasmin@example.com", now)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}

	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("save user: %v", err)
	}

	got, err := repo.GetByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("get user by email: %v", err)
	}

	if got.ID != user.ID {
		t.Fatalf("expected id %s, got %s", user.ID, got.ID)
	}
}

func TestUserRepository_GetByID_ReturnsErrUserNotFound(t *testing.T) {
	pool := newTestPool(t)
	resetIdentityTables(t, pool)

	repo := NewUserRepository(pool)

	_, err := repo.GetByID(context.Background(), "33333333-3333-3333-3333-333333333333")
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrUserNotFound, err)
	}
}

func TestUserRepository_Save_ReturnsEmailConflict(t *testing.T) {
	pool := newTestPool(t)
	resetIdentityTables(t, pool)

	repo := NewUserRepository(pool)
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	firstUser, err := domain.NewUser("44444444-4444-4444-4444-444444444444", "User One", "dup@example.com", now)
	if err != nil {
		t.Fatalf("new first user: %v", err)
	}

	secondUser, err := domain.NewUser("55555555-5555-5555-5555-555555555555", "User Two", "dup@example.com", now)
	if err != nil {
		t.Fatalf("new second user: %v", err)
	}

	if err := repo.Save(context.Background(), firstUser); err != nil {
		t.Fatalf("save first user: %v", err)
	}

	err = repo.Save(context.Background(), secondUser)
	if !errors.Is(err, ErrUserEmailConflict) {
		t.Fatalf("expected error %v, got %v", ErrUserEmailConflict, err)
	}
}
