package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /v1/users", handler.CreateUser)
	mux.HandleFunc("POST /v1/teams", handler.CreateTeam)
	mux.HandleFunc("POST /v1/teams/{id}/members", handler.AddTeamMember)
}
