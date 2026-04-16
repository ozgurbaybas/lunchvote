package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	pollapp "github.com/ozgurbaybas/lunchvote/modules/poll/application"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

type Handler struct {
	service *pollapp.Service
}

func NewHandler(service *pollapp.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req createPollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" ||
		strings.TrimSpace(req.TeamID) == "" ||
		strings.TrimSpace(req.Title) == "" ||
		strings.TrimSpace(req.CreatorUserID) == "" {
		writeError(w, http.StatusBadRequest, "id, team_id, title and creator_user_id are required")
		return
	}

	poll, err := h.service.CreatePoll(r.Context(), pollapp.CreatePollCommand{
		ID:            req.ID,
		TeamID:        req.TeamID,
		Title:         req.Title,
		RestaurantIDs: req.RestaurantIDs,
		CreatorUserID: req.CreatorUserID,
	})
	if err != nil {
		switch {
		case errors.Is(err, identitydomain.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, polldomain.ErrUserNotTeamMember):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, polldomain.ErrInvalidPollID),
			errors.Is(err, polldomain.ErrInvalidTeamID),
			errors.Is(err, polldomain.ErrInvalidPollTitle),
			errors.Is(err, polldomain.ErrNotEnoughPollOptions),
			errors.Is(err, polldomain.ErrDuplicatePollOption),
			errors.Is(err, polldomain.ErrInvalidRestaurantID):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, toPollResponse(poll))
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	pollID := strings.TrimSpace(r.PathValue("id"))
	if pollID == "" {
		writeError(w, http.StatusBadRequest, "poll id is required")
		return
	}

	var req voteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.RestaurantID) == "" {
		writeError(w, http.StatusBadRequest, "user_id and restaurant_id are required")
		return
	}

	poll, err := h.service.Vote(r.Context(), pollapp.VoteCommand{
		PollID:       pollID,
		UserID:       req.UserID,
		RestaurantID: req.RestaurantID,
	})
	if err != nil {
		switch {
		case errors.Is(err, polldomain.ErrPollNotFound),
			errors.Is(err, identitydomain.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, polldomain.ErrUserNotTeamMember):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, polldomain.ErrVoteAlreadyExists),
			errors.Is(err, polldomain.ErrPollClosed):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, polldomain.ErrPollOptionNotFound):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, toPollResponse(poll))
}

func (h *Handler) Results(w http.ResponseWriter, r *http.Request) {
	pollID := strings.TrimSpace(r.PathValue("id"))
	if pollID == "" {
		writeError(w, http.StatusBadRequest, "poll id is required")
		return
	}

	results, err := h.service.Results(r.Context(), pollID)
	if err != nil {
		switch {
		case errors.Is(err, polldomain.ErrPollNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, pollResultsResponse{
		PollID:  pollID,
		Results: results,
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
