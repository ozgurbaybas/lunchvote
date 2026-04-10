package http

type createRestaurantRequest struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Address            string   `json:"address"`
	City               string   `json:"city"`
	District           string   `json:"district"`
	SupportedMealCards []string `json:"supported_meal_cards"`
}
