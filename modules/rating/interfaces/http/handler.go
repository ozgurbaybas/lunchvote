package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	ratingapp "github.com/ozgurbaybas/lunchvote/modules/rating/application"
	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
	restaurantdomain "github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
	"github.com/ozgurbaybas/lunchvote/platform/httpserver"
)

type Handler struct {
	service *ratingapp.Service
}

func NewHandler(service *ratingapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRating(w http.ResponseWriter, r *http.Request) {
	var req createRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" ||
		strings.TrimSpace(req.RestaurantID) == "" ||
		strings.TrimSpace(req.UserID) == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "id, restaurant_id and user_id are required")
		return
	}

	rating, err := h.service.CreateRating(r.Context(), ratingapp.CreateRatingCommand{
		ID:           req.ID,
		RestaurantID: req.RestaurantID,
		UserID:       req.UserID,
		Score:        req.Score,
		Comment:      req.Comment,
	})
	if err != nil {
		switch {
		case errors.Is(err, identitydomain.ErrUserNotFound),
			errors.Is(err, restaurantdomain.ErrRestaurantNotFound):
			httpserver.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ratingdomain.ErrRatingAlreadyExists):
			httpserver.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, ratingdomain.ErrInvalidRatingID),
			errors.Is(err, ratingdomain.ErrInvalidRestaurantID),
			errors.Is(err, ratingdomain.ErrInvalidUserID),
			errors.Is(err, ratingdomain.ErrInvalidRatingScore):
			httpserver.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, toRatingResponse(rating))
}

func (h *Handler) ListRatingsByRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID := strings.TrimSpace(r.PathValue("id"))
	if restaurantID == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "restaurant id is required")
		return
	}

	ratings, err := h.service.ListRatingsByRestaurant(r.Context(), restaurantID)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := make([]ratingResponse, 0, len(ratings))
	for _, rating := range ratings {
		response = append(response, toRatingResponse(rating))
	}

	httpserver.WriteJSON(w, http.StatusOK, response)
}
