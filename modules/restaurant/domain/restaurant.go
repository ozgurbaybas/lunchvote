package domain

import (
	"strings"
	"time"
)

type Restaurant struct {
	ID                 string
	Name               string
	Address            string
	City               string
	District           string
	SupportedMealCards []MealCard
	IsActive           bool
	CreatedAt          time.Time
}

func NewRestaurant(
	id string,
	name string,
	address string,
	city string,
	district string,
	mealCards []MealCard,
	now time.Time,
) (Restaurant, error) {
	if strings.TrimSpace(id) == "" {
		return Restaurant{}, ErrInvalidRestaurantID
	}

	if strings.TrimSpace(name) == "" {
		return Restaurant{}, ErrInvalidRestaurantName
	}

	if strings.TrimSpace(city) == "" {
		return Restaurant{}, ErrInvalidRestaurantCity
	}

	if strings.TrimSpace(district) == "" {
		return Restaurant{}, ErrInvalidRestaurantDistrict
	}

	normalizedMealCards, err := normalizeMealCards(mealCards)
	if err != nil {
		return Restaurant{}, err
	}

	return Restaurant{
		ID:                 strings.TrimSpace(id),
		Name:               strings.TrimSpace(name),
		Address:            strings.TrimSpace(address),
		City:               strings.TrimSpace(city),
		District:           strings.TrimSpace(district),
		SupportedMealCards: normalizedMealCards,
		IsActive:           true,
		CreatedAt:          now,
	}, nil
}

func normalizeMealCards(cards []MealCard) ([]MealCard, error) {
	if len(cards) == 0 {
		return nil, ErrInvalidMealCard
	}

	seen := make(map[MealCard]struct{}, len(cards))
	normalized := make([]MealCard, 0, len(cards))

	for _, card := range cards {
		if !card.IsValid() {
			return nil, ErrInvalidMealCard
		}

		if _, ok := seen[card]; ok {
			return nil, ErrDuplicateMealCard
		}

		seen[card] = struct{}{}
		normalized = append(normalized, card)
	}

	return normalized, nil
}
