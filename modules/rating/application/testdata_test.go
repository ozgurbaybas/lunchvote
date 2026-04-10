package application

import (
	"context"
	"strings"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type inMemoryRatingRepository struct {
	ratings []ratingdomain.Rating
}

func newInMemoryRatingRepository() *inMemoryRatingRepository {
	return &inMemoryRatingRepository{
		ratings: make([]ratingdomain.Rating, 0),
	}
}

func (r *inMemoryRatingRepository) Save(_ context.Context, rating ratingdomain.Rating) error {
	r.ratings = append(r.ratings, rating)
	return nil
}

func (r *inMemoryRatingRepository) GetByRestaurantAndUser(_ context.Context, restaurantID string, userID string) (ratingdomain.Rating, error) {
	for _, rating := range r.ratings {
		if rating.RestaurantID == strings.TrimSpace(restaurantID) && rating.UserID == strings.TrimSpace(userID) {
			return rating, nil
		}
	}

	return ratingdomain.Rating{}, ratingdomain.ErrRatingNotFound
}

func (r *inMemoryRatingRepository) ListByRestaurantID(_ context.Context, restaurantID string) ([]ratingdomain.Rating, error) {
	result := make([]ratingdomain.Rating, 0)
	for _, rating := range r.ratings {
		if rating.RestaurantID == restaurantID {
			result = append(result, rating)
		}
	}
	return result, nil
}

type inMemoryUserRepository struct {
	byID    map[string]identitydomain.User
	byEmail map[string]identitydomain.User
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		byID:    make(map[string]identitydomain.User),
		byEmail: make(map[string]identitydomain.User),
	}
}

func (r *inMemoryUserRepository) Save(_ context.Context, user identitydomain.User) error {
	r.byID[user.ID] = user
	r.byEmail[strings.ToLower(user.Email)] = user
	return nil
}

func (r *inMemoryUserRepository) GetByID(_ context.Context, id string) (identitydomain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepository) GetByEmail(_ context.Context, email string) (identitydomain.User, error) {
	user, ok := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}
	return user, nil
}

type inMemoryRestaurantRepository struct {
	restaurants []restaurantdomain.Restaurant
}

func newInMemoryRestaurantRepository() *inMemoryRestaurantRepository {
	return &inMemoryRestaurantRepository{
		restaurants: make([]restaurantdomain.Restaurant, 0),
	}
}

func (r *inMemoryRestaurantRepository) Save(_ context.Context, restaurant restaurantdomain.Restaurant) error {
	r.restaurants = append(r.restaurants, restaurant)
	return nil
}

func (r *inMemoryRestaurantRepository) List(_ context.Context) ([]restaurantdomain.Restaurant, error) {
	result := make([]restaurantdomain.Restaurant, len(r.restaurants))
	copy(result, r.restaurants)
	return result, nil
}
