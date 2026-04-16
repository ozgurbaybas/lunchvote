package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

func TestRepository_SaveAndGetByID(t *testing.T) {
	pool := newTestPool(t)
	resetPollTables(t, pool)
	seedPollDependencies(t, pool)

	repo := NewRepository(pool)
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	poll, err := polldomain.NewPoll(
		"poll-1",
		testPollTeamID,
		"Friday Lunch",
		[]string{testPollRestaurantOne, testPollRestaurantTwo},
		now,
	)
	if err != nil {
		t.Fatalf("new poll: %v", err)
	}

	if err := poll.Vote(testPollOwnerID, testPollRestaurantOne, now); err != nil {
		t.Fatalf("vote owner: %v", err)
	}

	if err := poll.Vote(testPollMemberID, testPollRestaurantTwo, now.Add(time.Minute)); err != nil {
		t.Fatalf("vote member: %v", err)
	}

	if err := repo.Save(context.Background(), poll); err != nil {
		t.Fatalf("save poll: %v", err)
	}

	got, err := repo.GetByID(context.Background(), "poll-1")
	if err != nil {
		t.Fatalf("get poll by id: %v", err)
	}

	if got.ID != "poll-1" {
		t.Fatalf("expected poll id poll-1, got %s", got.ID)
	}

	if got.TeamID != testPollTeamID {
		t.Fatalf("expected team id %s, got %s", testPollTeamID, got.TeamID)
	}

	if len(got.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(got.Options))
	}

	if len(got.Votes) != 2 {
		t.Fatalf("expected 2 votes, got %d", len(got.Votes))
	}
}

func TestRepository_GetByID_ReturnsErrPollNotFound(t *testing.T) {
	pool := newTestPool(t)
	resetPollTables(t, pool)

	repo := NewRepository(pool)

	_, err := repo.GetByID(context.Background(), "missing-poll")
	if !errors.Is(err, polldomain.ErrPollNotFound) {
		t.Fatalf("expected error %v, got %v", polldomain.ErrPollNotFound, err)
	}
}

func TestRepository_ListByTeamID(t *testing.T) {
	pool := newTestPool(t)
	resetPollTables(t, pool)
	seedPollDependencies(t, pool)

	repo := NewRepository(pool)
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	firstPoll, err := polldomain.NewPoll(
		"poll-1",
		testPollTeamID,
		"Friday Lunch",
		[]string{testPollRestaurantOne, testPollRestaurantTwo},
		now,
	)
	if err != nil {
		t.Fatalf("new first poll: %v", err)
	}

	secondPoll, err := polldomain.NewPoll(
		"poll-2",
		testPollTeamID,
		"Monday Lunch",
		[]string{testPollRestaurantOne, testPollRestaurantTwo},
		now.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("new second poll: %v", err)
	}

	if err := repo.Save(context.Background(), firstPoll); err != nil {
		t.Fatalf("save first poll: %v", err)
	}

	if err := repo.Save(context.Background(), secondPoll); err != nil {
		t.Fatalf("save second poll: %v", err)
	}

	polls, err := repo.ListByTeamID(context.Background(), testPollTeamID)
	if err != nil {
		t.Fatalf("list polls by team: %v", err)
	}

	if len(polls) != 2 {
		t.Fatalf("expected 2 polls, got %d", len(polls))
	}

	if polls[0].ID != "poll-1" {
		t.Fatalf("expected first poll poll-1, got %s", polls[0].ID)
	}
}
