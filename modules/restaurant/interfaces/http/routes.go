package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /v1/restaurants", handler.CreateRestaurant)
	mux.HandleFunc("GET /v1/restaurants", handler.ListRestaurants)
}
