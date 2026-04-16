package application

import (
	"context"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type inMemoryTeamReader struct {
	byID map[string]identitydomain.Team
}

func newInMemoryTeamReader() *inMemoryTeamReader {
	return &inMemoryTeamReader{byID: make(map[string]identitydomain.Team)}
}

func (r *inMemoryTeamReader) GetByID(_ context.Context, id string) (identitydomain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return identitydomain.Team{}, identitydomain.ErrTeamNotFound
	}
	return team, nil
}

type inMemoryRestaurantReader struct {
	items []restaurantdomain.Restaurant
}

func newInMemoryRestaurantReader() *inMemoryRestaurantReader {
	return &inMemoryRestaurantReader{items: make([]restaurantdomain.Restaurant, 0)}
}

func (r *inMemoryRestaurantReader) List(_ context.Context) ([]restaurantdomain.Restaurant, error) {
	result := make([]restaurantdomain.Restaurant, len(r.items))
	copy(result, r.items)
	return result, nil
}

type inMemoryRatingReader struct {
	byRestaurantID map[string][]ratingdomain.Rating
}

func newInMemoryRatingReader() *inMemoryRatingReader {
	return &inMemoryRatingReader{byRestaurantID: make(map[string][]ratingdomain.Rating)}
}

func (r *inMemoryRatingReader) ListByRestaurantID(_ context.Context, restaurantID string) ([]ratingdomain.Rating, error) {
	result := make([]ratingdomain.Rating, len(r.byRestaurantID[restaurantID]))
	copy(result, r.byRestaurantID[restaurantID])
	return result, nil
}

type inMemoryPollReader struct {
	byTeamID map[string][]polldomain.Poll
}

func newInMemoryPollReader() *inMemoryPollReader {
	return &inMemoryPollReader{byTeamID: make(map[string][]polldomain.Poll)}
}

func (r *inMemoryPollReader) ListByTeamID(_ context.Context, teamID string) ([]polldomain.Poll, error) {
	result := make([]polldomain.Poll, len(r.byTeamID[teamID]))
	copy(result, r.byTeamID[teamID])
	return result, nil
}
