package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /v1/teams/{id}/recommendations", handler.RecommendRestaurants)
}
