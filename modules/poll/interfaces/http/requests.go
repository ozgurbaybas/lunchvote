package http

type createPollRequest struct {
	ID            string   `json:"id"`
	TeamID        string   `json:"team_id"`
	Title         string   `json:"title"`
	RestaurantIDs []string `json:"restaurant_ids"`
	CreatorUserID string   `json:"creator_user_id"`
}

type voteRequest struct {
	UserID       string `json:"user_id"`
	RestaurantID string `json:"restaurant_id"`
}
