package http

import recommendationdomain "github.com/ozgurbaybas/lunchvote/modules/recommendation/domain"

type errorResponse struct {
	Error string `json:"error"`
}

type recommendationResponse struct {
	RestaurantID string   `json:"restaurant_id"`
	Score        float64  `json:"score"`
	Reasons      []string `json:"reasons"`
}

func toRecommendationResponse(item recommendationdomain.Recommendation) recommendationResponse {
	return recommendationResponse{
		RestaurantID: item.RestaurantID,
		Score:        item.Score,
		Reasons:      item.Reasons,
	}
}
