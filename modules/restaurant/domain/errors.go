package domain

import "errors"

var (
	ErrInvalidRestaurantID       = errors.New("invalid restaurant id")
	ErrInvalidRestaurantName     = errors.New("invalid restaurant name")
	ErrInvalidRestaurantCity     = errors.New("invalid restaurant city")
	ErrInvalidRestaurantDistrict = errors.New("invalid restaurant district")
	ErrInvalidMealCard           = errors.New("invalid meal card")
	ErrDuplicateMealCard         = errors.New("duplicate meal card")
)
