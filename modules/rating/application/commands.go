package application

type CreateRatingCommand struct {
	ID           string
	RestaurantID string
	UserID       string
	Score        int
	Comment      string
}
