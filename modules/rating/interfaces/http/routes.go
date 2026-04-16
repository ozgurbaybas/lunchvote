package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /v1/ratings", handler.CreateRating)
	mux.HandleFunc("GET /v1/restaurants/{id}/ratings", handler.ListRatingsByRestaurant)
}
