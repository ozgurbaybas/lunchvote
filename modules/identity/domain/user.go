package domain

import (
	"strings"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(id, name, email string, now time.Time) (User, error) {
	if strings.TrimSpace(id) == "" {
		return User{}, ErrInvalidUserID
	}

	if strings.TrimSpace(name) == "" {
		return User{}, ErrInvalidUserName
	}

	if strings.TrimSpace(email) == "" {
		return User{}, ErrInvalidUserEmail
	}

	return User{
		ID:        strings.TrimSpace(id),
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(email),
		CreatedAt: now,
	}, nil
}
