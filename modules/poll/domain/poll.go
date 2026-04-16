package domain

import (
	"strings"
	"time"
)

type PollStatus string

const (
	PollStatusOpen   PollStatus = "open"
	PollStatusClosed PollStatus = "closed"
)

type PollOption struct {
	RestaurantID string
}

type Vote struct {
	UserID       string
	RestaurantID string
	VotedAt      time.Time
}

type Poll struct {
	ID        string
	TeamID    string
	Title     string
	Options   []PollOption
	Votes     []Vote
	Status    PollStatus
	CreatedAt time.Time
	ClosedAt  *time.Time
}

func NewPoll(id, teamID, title string, restaurantIDs []string, now time.Time) (Poll, error) {
	if strings.TrimSpace(id) == "" {
		return Poll{}, ErrInvalidPollID
	}

	if strings.TrimSpace(teamID) == "" {
		return Poll{}, ErrInvalidTeamID
	}

	if strings.TrimSpace(title) == "" {
		return Poll{}, ErrInvalidPollTitle
	}

	options, err := buildOptions(restaurantIDs)
	if err != nil {
		return Poll{}, err
	}

	return Poll{
		ID:        strings.TrimSpace(id),
		TeamID:    strings.TrimSpace(teamID),
		Title:     strings.TrimSpace(title),
		Options:   options,
		Votes:     make([]Vote, 0),
		Status:    PollStatusOpen,
		CreatedAt: now,
	}, nil
}

func (p *Poll) Vote(userID, restaurantID string, now time.Time) error {
	if p.Status != PollStatusOpen {
		return ErrPollClosed
	}

	userID = strings.TrimSpace(userID)
	restaurantID = strings.TrimSpace(restaurantID)

	if userID == "" {
		return ErrUserNotTeamMember
	}

	if !p.hasOption(restaurantID) {
		return ErrPollOptionNotFound
	}

	for _, vote := range p.Votes {
		if vote.UserID == userID {
			return ErrVoteAlreadyExists
		}
	}

	p.Votes = append(p.Votes, Vote{
		UserID:       userID,
		RestaurantID: restaurantID,
		VotedAt:      now,
	})

	return nil
}

func (p *Poll) Close(now time.Time) {
	p.Status = PollStatusClosed
	p.ClosedAt = &now
}

func (p Poll) Results() map[string]int {
	results := make(map[string]int, len(p.Options))
	for _, option := range p.Options {
		results[option.RestaurantID] = 0
	}

	for _, vote := range p.Votes {
		results[vote.RestaurantID]++
	}

	return results
}

func buildOptions(restaurantIDs []string) ([]PollOption, error) {
	if len(restaurantIDs) < 2 {
		return nil, ErrNotEnoughPollOptions
	}

	seen := make(map[string]struct{}, len(restaurantIDs))
	options := make([]PollOption, 0, len(restaurantIDs))

	for _, restaurantID := range restaurantIDs {
		restaurantID = strings.TrimSpace(restaurantID)
		if restaurantID == "" {
			return nil, ErrInvalidRestaurantID
		}

		if _, ok := seen[restaurantID]; ok {
			return nil, ErrDuplicatePollOption
		}

		seen[restaurantID] = struct{}{}
		options = append(options, PollOption{RestaurantID: restaurantID})
	}

	return options, nil
}

func (p Poll) hasOption(restaurantID string) bool {
	for _, option := range p.Options {
		if option.RestaurantID == restaurantID {
			return true
		}
	}
	return false
}
