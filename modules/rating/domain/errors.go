package domain

import "errors"

var (
	ErrInvalidRatingID     = errors.New("invalid rating id")
	ErrInvalidRestaurantID = errors.New("invalid restaurant id")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidRatingScore  = errors.New("invalid rating score")
	ErrRatingAlreadyExists = errors.New("rating already exists")
	ErrRatingNotFound      = errors.New("rating not found")
)
