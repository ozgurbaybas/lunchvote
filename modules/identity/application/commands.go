package application

type CreateUserCommand struct {
	ID    string
	Name  string
	Email string
}

type CreateTeamCommand struct {
	ID      string
	Name    string
	OwnerID string
}

type AddTeamMemberCommand struct {
	TeamID string
	UserID string
}
