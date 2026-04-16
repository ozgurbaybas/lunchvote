package domain

import "errors"

var (
	ErrInvalidPollID        = errors.New("invalid poll id")
	ErrInvalidTeamID        = errors.New("invalid team id")
	ErrInvalidPollTitle     = errors.New("invalid poll title")
	ErrNotEnoughPollOptions = errors.New("not enough poll options")
	ErrDuplicatePollOption  = errors.New("duplicate poll option")
	ErrInvalidRestaurantID  = errors.New("invalid restaurant id")
	ErrPollNotFound         = errors.New("poll not found")
	ErrPollClosed           = errors.New("poll is closed")
	ErrVoteAlreadyExists    = errors.New("vote already exists")
	ErrPollOptionNotFound   = errors.New("poll option not found")
	ErrUserNotTeamMember    = errors.New("user is not a team member")
)
