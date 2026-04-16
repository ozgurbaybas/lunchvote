package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	recommendationapp "github.com/ozgurbaybas/lunchvote/modules/recommendation/application"
)

type Handler struct {
	service *recommendationapp.Service
}

func NewHandler(service *recommendationapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RecommendRestaurants(w http.ResponseWriter, r *http.Request) {
	teamID := strings.TrimSpace(r.PathValue("id"))
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "team id is required")
		return
	}

	limit := 0
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed < 0 {
			writeError(w, http.StatusBadRequest, "limit must be a non-negative integer")
			return
		}
		limit = parsed
	}

	items, err := h.service.RecommendRestaurants(r.Context(), recommendationapp.RecommendRestaurantsQuery{
		TeamID: teamID,
		Limit:  limit,
	})
	if err != nil {
		switch {
		case errors.Is(err, identitydomain.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response := make([]recommendationResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toRecommendationResponse(item))
	}

	writeJSON(w, http.StatusOK, response)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
