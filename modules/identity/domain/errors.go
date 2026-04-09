package domain

import "errors"

var (
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrInvalidUserName       = errors.New("invalid user name")
	ErrInvalidUserEmail      = errors.New("invalid user email")
	ErrInvalidTeamID         = errors.New("invalid team id")
	ErrInvalidTeamName       = errors.New("invalid team name")
	ErrInvalidOwnerID        = errors.New("invalid owner id")
	ErrMemberAlreadyExists   = errors.New("member already exists")
	ErrMemberNotFound        = errors.New("member not found")
	ErrOwnerCannotBeRemoved  = errors.New("owner cannot be removed")
	ErrInvalidMembershipRole = errors.New("invalid membership role")
)
