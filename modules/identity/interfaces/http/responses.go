package http

import "github.com/ozgurbaybas/lunchvote/modules/identity/domain"

type userResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type membershipResponse struct {
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

type teamResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	OwnerID   string               `json:"owner_id"`
	Members   []membershipResponse `json:"members"`
	CreatedAt string               `json:"created_at"`
}

func toUserResponse(user domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func toTeamResponse(team domain.Team) teamResponse {
	members := make([]membershipResponse, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, membershipResponse{
			UserID:   member.UserID,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}

	return teamResponse{
		ID:        team.ID,
		Name:      team.Name,
		OwnerID:   team.OwnerID,
		Members:   members,
		CreatedAt: team.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
