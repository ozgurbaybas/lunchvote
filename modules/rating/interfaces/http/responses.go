package http

import ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"

type errorResponse struct {
	Error string `json:"error"`
}

type ratingResponse struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	UserID       string `json:"user_id"`
	Score        int    `json:"score"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"created_at"`
}

func toRatingResponse(rating ratingdomain.Rating) ratingResponse {
	return ratingResponse{
		ID:           rating.ID,
		RestaurantID: rating.RestaurantID,
		UserID:       rating.UserID,
		Score:        rating.Score,
		Comment:      rating.Comment,
		CreatedAt:    rating.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
