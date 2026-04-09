package http

type createUserRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type createTeamRequest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type addTeamMemberRequest struct {
	UserID string `json:"user_id"`
}
