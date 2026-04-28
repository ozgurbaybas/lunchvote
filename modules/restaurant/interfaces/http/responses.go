package http

import "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"

type restaurantResponse struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Address            string   `json:"address"`
	City               string   `json:"city"`
	District           string   `json:"district"`
	SupportedMealCards []string `json:"supported_meal_cards"`
	IsActive           bool     `json:"is_active"`
	CreatedAt          string   `json:"created_at"`
}

func toRestaurantResponse(restaurant domain.Restaurant) restaurantResponse {
	mealCards := make([]string, 0, len(restaurant.SupportedMealCards))
	for _, card := range restaurant.SupportedMealCards {
		mealCards = append(mealCards, string(card))
	}

	return restaurantResponse{
		ID:                 restaurant.ID,
		Name:               restaurant.Name,
		Address:            restaurant.Address,
		City:               restaurant.City,
		District:           restaurant.District,
		SupportedMealCards: mealCards,
		IsActive:           restaurant.IsActive,
		CreatedAt:          restaurant.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
