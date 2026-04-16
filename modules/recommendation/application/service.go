package application

import (
	"context"
	"errors"
	"fmt"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	recommendationdomain "github.com/ozgurbaybas/lunchvote/modules/recommendation/domain"
)

type Service struct {
	teams       recommendationdomain.TeamReader
	restaurants recommendationdomain.RestaurantReader
	ratings     recommendationdomain.RatingReader
	polls       recommendationdomain.PollReader
}

func NewService(
	teams recommendationdomain.TeamReader,
	restaurants recommendationdomain.RestaurantReader,
	ratings recommendationdomain.RatingReader,
	polls recommendationdomain.PollReader,
) *Service {
	return &Service{
		teams:       teams,
		restaurants: restaurants,
		ratings:     ratings,
		polls:       polls,
	}
}

func (s *Service) RecommendRestaurants(
	ctx context.Context,
	query RecommendRestaurantsQuery,
) ([]recommendationdomain.Recommendation, error) {
	_, err := s.teams.GetByID(ctx, query.TeamID)
	if err != nil {
		if errors.Is(err, identitydomain.ErrTeamNotFound) {
			return nil, identitydomain.ErrTeamNotFound
		}
		return nil, fmt.Errorf("get team by id: %w", err)
	}

	restaurants, err := s.restaurants.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list restaurants: %w", err)
	}

	polls, err := s.polls.ListByTeamID(ctx, query.TeamID)
	if err != nil {
		return nil, fmt.Errorf("list polls by team: %w", err)
	}

	voteCounts := countVotesByRestaurant(polls)

	items := make([]recommendationdomain.Recommendation, 0)
	for _, restaurant := range restaurants {
		if !restaurant.IsActive {
			continue
		}

		score := 0.0
		reasons := make([]string, 0)

		voteCount := voteCounts[restaurant.ID]
		if voteCount > 0 {
			score += float64(voteCount * 2)
			reasons = append(reasons, "team poll history")
		}

		ratings, err := s.ratings.ListByRestaurantID(ctx, restaurant.ID)
		if err != nil {
			return nil, fmt.Errorf("list ratings by restaurant: %w", err)
		}

		if avg, ok := averageRating(ratings); ok {
			score += avg * 10
			reasons = append(reasons, "strong ratings")
		}

		item, err := recommendationdomain.NewRecommendation(restaurant.ID, score, reasons)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	recommendationdomain.SortByScore(items)

	if query.Limit > 0 && len(items) > query.Limit {
		items = items[:query.Limit]
	}

	return items, nil
}

func countVotesByRestaurant(polls []polldomain.Poll) map[string]int {
	result := make(map[string]int)
	for _, poll := range polls {
		for _, vote := range poll.Votes {
			result[vote.RestaurantID]++
		}
	}
	return result
}

func averageRating(ratings []ratingdomain.Rating) (float64, bool) {
	if len(ratings) == 0 {
		return 0, false
	}

	total := 0
	for _, rating := range ratings {
		total += rating.Score
	}

	return float64(total) / float64(len(ratings)), true
}
