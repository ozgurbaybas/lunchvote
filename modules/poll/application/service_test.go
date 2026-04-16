package application

import (
	"context"
	"errors"
	"testing"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

func TestService_CreatePoll(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	pollRepo := newInMemoryPollRepository()
	teamRepo := newInMemoryTeamRepository()
	service := NewService(pollRepo, teamRepo, func() time.Time { return now })

	team, err := identitydomain.NewTeam("team-1", "Backend Team", "user-1", now)
	if err != nil {
		t.Fatalf("new team: %v", err)
	}
	if err := teamRepo.Save(context.Background(), team); err != nil {
		t.Fatalf("save team: %v", err)
	}

	poll, err := service.CreatePoll(context.Background(), CreatePollCommand{
		ID:            "poll-1",
		TeamID:        "team-1",
		Title:         "Friday Lunch",
		RestaurantIDs: []string{"restaurant-1", "restaurant-2"},
		CreatorUserID: "user-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if poll.ID != "poll-1" {
		t.Fatalf("expected poll id poll-1, got %s", poll.ID)
	}
}

func TestService_CreatePoll_ReturnsErrorWhenTeamMissing(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	service := NewService(newInMemoryPollRepository(), newInMemoryTeamRepository(), func() time.Time { return now })

	_, err := service.CreatePoll(context.Background(), CreatePollCommand{
		ID:            "poll-1",
		TeamID:        "missing-team",
		Title:         "Friday Lunch",
		RestaurantIDs: []string{"restaurant-1", "restaurant-2"},
		CreatorUserID: "user-1",
	})
	if !errors.Is(err, identitydomain.ErrTeamNotFound) {
		t.Fatalf("expected error %v, got %v", identitydomain.ErrTeamNotFound, err)
	}
}

func TestService_CreatePoll_ReturnsErrorWhenUserNotTeamMember(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	pollRepo := newInMemoryPollRepository()
	teamRepo := newInMemoryTeamRepository()
	service := NewService(pollRepo, teamRepo, func() time.Time { return now })

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "owner-1", now)
	_ = teamRepo.Save(context.Background(), team)

	_, err := service.CreatePoll(context.Background(), CreatePollCommand{
		ID:            "poll-1",
		TeamID:        "team-1",
		Title:         "Friday Lunch",
		RestaurantIDs: []string{"restaurant-1", "restaurant-2"},
		CreatorUserID: "user-1",
	})
	if !errors.Is(err, polldomain.ErrUserNotTeamMember) {
		t.Fatalf("expected error %v, got %v", polldomain.ErrUserNotTeamMember, err)
	}
}

func TestService_Vote(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	pollRepo := newInMemoryPollRepository()
	teamRepo := newInMemoryTeamRepository()
	service := NewService(pollRepo, teamRepo, func() time.Time { return now })

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "user-1", now)
	_ = team.AddMember("user-2", now)
	_ = teamRepo.Save(context.Background(), team)

	poll, _ := service.CreatePoll(context.Background(), CreatePollCommand{
		ID:            "poll-1",
		TeamID:        "team-1",
		Title:         "Friday Lunch",
		RestaurantIDs: []string{"restaurant-1", "restaurant-2"},
		CreatorUserID: "user-1",
	})

	got, err := service.Vote(context.Background(), VoteCommand{
		PollID:       poll.ID,
		UserID:       "user-2",
		RestaurantID: "restaurant-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got.Votes) != 1 {
		t.Fatalf("expected 1 vote, got %d", len(got.Votes))
	}
}

func TestService_Vote_ReturnsErrorWhenPollMissing(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	service := NewService(newInMemoryPollRepository(), newInMemoryTeamRepository(), func() time.Time { return now })

	_, err := service.Vote(context.Background(), VoteCommand{
		PollID:       "missing-poll",
		UserID:       "user-1",
		RestaurantID: "restaurant-1",
	})
	if !errors.Is(err, polldomain.ErrPollNotFound) {
		t.Fatalf("expected error %v, got %v", polldomain.ErrPollNotFound, err)
	}
}

func TestService_Vote_ReturnsErrorWhenUserNotTeamMember(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	pollRepo := newInMemoryPollRepository()
	teamRepo := newInMemoryTeamRepository()
	service := NewService(pollRepo, teamRepo, func() time.Time { return now })

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "owner-1", now)
	_ = teamRepo.Save(context.Background(), team)

	poll, _ := service.CreatePoll(context.Background(), CreatePollCommand{
		ID:            "poll-1",
		TeamID:        "team-1",
		Title:         "Friday Lunch",
		RestaurantIDs: []string{"restaurant-1", "restaurant-2"},
		CreatorUserID: "owner-1",
	})

	_, err := service.Vote(context.Background(), VoteCommand{
		PollID:       poll.ID,
		UserID:       "user-2",
		RestaurantID: "restaurant-1",
	})
	if !errors.Is(err, polldomain.ErrUserNotTeamMember) {
		t.Fatalf("expected error %v, got %v", polldomain.ErrUserNotTeamMember, err)
	}
}

func TestService_Results(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	pollRepo := newInMemoryPollRepository()
	teamRepo := newInMemoryTeamRepository()
	service := NewService(pollRepo, teamRepo, func() time.Time { return now })

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "user-1", now)
	_ = team.AddMember("user-2", now)
	_ = teamRepo.Save(context.Background(), team)

	poll, _ := service.CreatePoll(context.Background(), CreatePollCommand{
		ID:            "poll-1",
		TeamID:        "team-1",
		Title:         "Friday Lunch",
		RestaurantIDs: []string{"restaurant-1", "restaurant-2"},
		CreatorUserID: "user-1",
	})

	_, _ = service.Vote(context.Background(), VoteCommand{
		PollID:       poll.ID,
		UserID:       "user-1",
		RestaurantID: "restaurant-1",
	})
	_, _ = service.Vote(context.Background(), VoteCommand{
		PollID:       poll.ID,
		UserID:       "user-2",
		RestaurantID: "restaurant-2",
	})

	results, err := service.Results(context.Background(), poll.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if results["restaurant-1"] != 1 || results["restaurant-2"] != 1 {
		t.Fatalf("unexpected results: %#v", results)
	}
}
