package http

type createRatingRequest struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	UserID       string `json:"user_id"`
	Score        int    `json:"score"`
	Comment      string `json:"comment"`
}
