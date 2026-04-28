package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/application"
	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
	"github.com/ozgurbaybas/lunchvote/platform/httpserver"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var req createRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" ||
		strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.City) == "" ||
		strings.TrimSpace(req.District) == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "id, name, city and district are required")
		return
	}

	restaurant, err := h.service.CreateRestaurant(r.Context(), application.CreateRestaurantCommand{
		ID:                 req.ID,
		Name:               req.Name,
		Address:            req.Address,
		City:               req.City,
		District:           req.District,
		SupportedMealCards: req.SupportedMealCards,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRestaurantID),
			errors.Is(err, domain.ErrInvalidRestaurantName),
			errors.Is(err, domain.ErrInvalidRestaurantCity),
			errors.Is(err, domain.ErrInvalidRestaurantDistrict),
			errors.Is(err, domain.ErrInvalidMealCard),
			errors.Is(err, domain.ErrDuplicateMealCard):
			httpserver.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, toRestaurantResponse(restaurant))
}

func (h *Handler) ListRestaurants(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.service.ListRestaurants(r.Context())
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := make([]restaurantResponse, 0, len(restaurants))
	for _, restaurant := range restaurants {
		response = append(response, toRestaurantResponse(restaurant))
	}

	httpserver.WriteJSON(w, http.StatusOK, response)
}
