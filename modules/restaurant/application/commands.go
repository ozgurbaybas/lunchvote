package application

type CreateRestaurantCommand struct {
	ID                 string
	Name               string
	Address            string
	City               string
	District           string
	SupportedMealCards []string
}
