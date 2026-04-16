package domain

import (
	"testing"
	"time"
)

func TestNewPoll(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		id            string
		teamID        string
		title         string
		restaurantIDs []string
		wantErr       error
	}{
		{
			name:          "creates valid poll",
			id:            "poll-1",
			teamID:        "team-1",
			title:         "Friday Lunch",
			restaurantIDs: []string{"restaurant-1", "restaurant-2"},
		},
		{
			name:          "returns error when not enough options",
			id:            "poll-1",
			teamID:        "team-1",
			title:         "Friday Lunch",
			restaurantIDs: []string{"restaurant-1"},
			wantErr:       ErrNotEnoughPollOptions,
		},
		{
			name:          "returns error when duplicate option",
			id:            "poll-1",
			teamID:        "team-1",
			title:         "Friday Lunch",
			restaurantIDs: []string{"restaurant-1", "restaurant-1"},
			wantErr:       ErrDuplicatePollOption,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poll, err := NewPoll(tt.id, tt.teamID, tt.title, tt.restaurantIDs, now)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if poll.Status != PollStatusOpen {
				t.Fatalf("expected poll status open, got %s", poll.Status)
			}
		})
	}
}

func TestPoll_Vote(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	poll, err := NewPoll("poll-1", "team-1", "Friday Lunch", []string{"restaurant-1", "restaurant-2"}, now)
	if err != nil {
		t.Fatalf("new poll: %v", err)
	}

	if err := poll.Vote("user-1", "restaurant-1", now); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(poll.Votes) != 1 {
		t.Fatalf("expected 1 vote, got %d", len(poll.Votes))
	}
}

func TestPoll_Vote_ReturnsErrorWhenUserVotesTwice(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	poll, _ := NewPoll("poll-1", "team-1", "Friday Lunch", []string{"restaurant-1", "restaurant-2"}, now)
	_ = poll.Vote("user-1", "restaurant-1", now)

	err := poll.Vote("user-1", "restaurant-2", now)
	if err != ErrVoteAlreadyExists {
		t.Fatalf("expected error %v, got %v", ErrVoteAlreadyExists, err)
	}
}

func TestPoll_Vote_ReturnsErrorWhenOptionNotFound(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	poll, _ := NewPoll("poll-1", "team-1", "Friday Lunch", []string{"restaurant-1", "restaurant-2"}, now)

	err := poll.Vote("user-1", "restaurant-3", now)
	if err != ErrPollOptionNotFound {
		t.Fatalf("expected error %v, got %v", ErrPollOptionNotFound, err)
	}
}

func TestPoll_Vote_ReturnsErrorWhenClosed(t *testing.T) {
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	poll, _ := NewPoll("poll-1", "team-1", "Friday Lunch", []string{"restaurant-1", "restaurant-2"}, now)
	poll.Close(now)

	err := poll.Vote("user-1", "restaurant-1", now)
	if err != ErrPollClosed {
		t.Fatalf("expected error %v, got %v", ErrPollClosed, err)
	}
}
