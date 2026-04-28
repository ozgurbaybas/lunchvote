package http

import polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"

type pollOptionResponse struct {
	RestaurantID string `json:"restaurant_id"`
}

type voteResponse struct {
	UserID       string `json:"user_id"`
	RestaurantID string `json:"restaurant_id"`
	VotedAt      string `json:"voted_at"`
}

type pollResponse struct {
	ID        string               `json:"id"`
	TeamID    string               `json:"team_id"`
	Title     string               `json:"title"`
	Status    string               `json:"status"`
	Options   []pollOptionResponse `json:"options"`
	Votes     []voteResponse       `json:"votes"`
	CreatedAt string               `json:"created_at"`
	ClosedAt  *string              `json:"closed_at"`
}

type pollResultsResponse struct {
	PollID  string         `json:"poll_id"`
	Results map[string]int `json:"results"`
}

func toPollResponse(poll polldomain.Poll) pollResponse {
	options := make([]pollOptionResponse, 0, len(poll.Options))
	for _, option := range poll.Options {
		options = append(options, pollOptionResponse{
			RestaurantID: option.RestaurantID,
		})
	}

	votes := make([]voteResponse, 0, len(poll.Votes))
	for _, vote := range poll.Votes {
		votes = append(votes, voteResponse{
			UserID:       vote.UserID,
			RestaurantID: vote.RestaurantID,
			VotedAt:      vote.VotedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}

	var closedAt *string
	if poll.ClosedAt != nil {
		value := poll.ClosedAt.UTC().Format("2006-01-02T15:04:05Z")
		closedAt = &value
	}

	return pollResponse{
		ID:        poll.ID,
		TeamID:    poll.TeamID,
		Title:     poll.Title,
		Status:    string(poll.Status),
		Options:   options,
		Votes:     votes,
		CreatedAt: poll.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		ClosedAt:  closedAt,
	}
}
