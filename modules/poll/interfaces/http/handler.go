package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	pollapp "github.com/ozgurbaybas/lunchvote/modules/poll/application"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
	"github.com/ozgurbaybas/lunchvote/platform/httpserver"
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
		httpserver.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" ||
		strings.TrimSpace(req.TeamID) == "" ||
		strings.TrimSpace(req.Title) == "" ||
		strings.TrimSpace(req.CreatorUserID) == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "id, team_id, title and creator_user_id are required")
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
			httpserver.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, polldomain.ErrUserNotTeamMember):
			httpserver.WriteError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, polldomain.ErrInvalidPollID),
			errors.Is(err, polldomain.ErrInvalidTeamID),
			errors.Is(err, polldomain.ErrInvalidPollTitle),
			errors.Is(err, polldomain.ErrNotEnoughPollOptions),
			errors.Is(err, polldomain.ErrDuplicatePollOption),
			errors.Is(err, polldomain.ErrInvalidRestaurantID):
			httpserver.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, toPollResponse(poll))
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	pollID := strings.TrimSpace(r.PathValue("id"))
	if pollID == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "poll id is required")
		return
	}

	var req voteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.RestaurantID) == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "user_id and restaurant_id are required")
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
			httpserver.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, polldomain.ErrUserNotTeamMember):
			httpserver.WriteError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, polldomain.ErrVoteAlreadyExists),
			errors.Is(err, polldomain.ErrPollClosed):
			httpserver.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, polldomain.ErrPollOptionNotFound):
			httpserver.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, toPollResponse(poll))
}

func (h *Handler) Results(w http.ResponseWriter, r *http.Request) {
	pollID := strings.TrimSpace(r.PathValue("id"))
	if pollID == "" {
		httpserver.WriteError(w, http.StatusBadRequest, "poll id is required")
		return
	}

	results, err := h.service.Results(r.Context(), pollID)
	if err != nil {
		switch {
		case errors.Is(err, polldomain.ErrPollNotFound):
			httpserver.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpserver.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, pollResultsResponse{
		PollID:  pollID,
		Results: results,
	})
}
