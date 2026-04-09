package domain

import (
	"strings"
	"time"
)

type Team struct {
	ID        string
	Name      string
	OwnerID   string
	Members   []Membership
	CreatedAt time.Time
}

func NewTeam(id, name, ownerID string, now time.Time) (Team, error) {
	if strings.TrimSpace(id) == "" {
		return Team{}, ErrInvalidTeamID
	}

	if strings.TrimSpace(name) == "" {
		return Team{}, ErrInvalidTeamName
	}

	if strings.TrimSpace(ownerID) == "" {
		return Team{}, ErrInvalidOwnerID
	}

	ownerMembership, err := NewMembership(strings.TrimSpace(ownerID), MembershipRoleOwner, now)
	if err != nil {
		return Team{}, err
	}

	return Team{
		ID:        strings.TrimSpace(id),
		Name:      strings.TrimSpace(name),
		OwnerID:   strings.TrimSpace(ownerID),
		Members:   []Membership{ownerMembership},
		CreatedAt: now,
	}, nil
}

func (t *Team) AddMember(userID string, now time.Time) error {
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidUserID
	}

	if t.HasMember(userID) {
		return ErrMemberAlreadyExists
	}

	member, err := NewMembership(strings.TrimSpace(userID), MembershipRoleMember, now)
	if err != nil {
		return err
	}

	t.Members = append(t.Members, member)
	return nil
}

func (t *Team) RemoveMember(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidUserID
	}

	for i, member := range t.Members {
		if member.UserID != strings.TrimSpace(userID) {
			continue
		}

		if member.Role == MembershipRoleOwner {
			return ErrOwnerCannotBeRemoved
		}

		t.Members = append(t.Members[:i], t.Members[i+1:]...)
		return nil
	}

	return ErrMemberNotFound
}

func (t Team) HasMember(userID string) bool {
	for _, member := range t.Members {
		if member.UserID == strings.TrimSpace(userID) {
			return true
		}
	}

	return false
}
