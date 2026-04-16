package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /v1/polls", handler.CreatePoll)
	mux.HandleFunc("POST /v1/polls/{id}/votes", handler.Vote)
	mux.HandleFunc("GET /v1/polls/{id}/results", handler.Results)
}
