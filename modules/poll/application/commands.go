package application

type CreatePollCommand struct {
	ID            string
	TeamID        string
	Title         string
	RestaurantIDs []string
	CreatorUserID string
}

type VoteCommand struct {
	PollID       string
	UserID       string
	RestaurantID string
}
