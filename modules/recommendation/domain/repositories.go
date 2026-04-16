package domain

import (
	"context"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type TeamReader interface {
	GetByID(ctx context.Context, id string) (identitydomain.Team, error)
}

type RestaurantReader interface {
	List(ctx context.Context) ([]restaurantdomain.Restaurant, error)
}

type RatingReader interface {
	ListByRestaurantID(ctx context.Context, restaurantID string) ([]ratingdomain.Rating, error)
}

type PollReader interface {
	ListByTeamID(ctx context.Context, teamID string) ([]polldomain.Poll, error)
}
