package domain

import "time"

type MembershipRole string

const (
	MembershipRoleOwner  MembershipRole = "owner"
	MembershipRoleMember MembershipRole = "member"
)

type Membership struct {
	UserID   string
	Role     MembershipRole
	JoinedAt time.Time
}

func NewMembership(userID string, role MembershipRole, joinedAt time.Time) (Membership, error) {
	if userID == "" {
		return Membership{}, ErrInvalidUserID
	}

	if !role.IsValid() {
		return Membership{}, ErrInvalidMembershipRole
	}

	return Membership{
		UserID:   userID,
		Role:     role,
		JoinedAt: joinedAt,
	}, nil
}

func (r MembershipRole) IsValid() bool {
	return r == MembershipRoleOwner || r == MembershipRoleMember
}
