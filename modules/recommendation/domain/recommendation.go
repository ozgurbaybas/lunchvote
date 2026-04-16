package domain

import (
	"sort"
	"strings"
)

type Recommendation struct {
	RestaurantID string
	Score        float64
	Reasons      []string
}

func NewRecommendation(restaurantID string, score float64, reasons []string) (Recommendation, error) {
	if strings.TrimSpace(restaurantID) == "" {
		return Recommendation{}, ErrInvalidRestaurantID
	}

	cleanReasons := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		reason = strings.TrimSpace(reason)
		if reason == "" {
			continue
		}
		cleanReasons = append(cleanReasons, reason)
	}

	return Recommendation{
		RestaurantID: strings.TrimSpace(restaurantID),
		Score:        score,
		Reasons:      cleanReasons,
	}, nil
}

func SortByScore(items []Recommendation) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].RestaurantID < items[j].RestaurantID
		}
		return items[i].Score > items[j].Score
	})
}
