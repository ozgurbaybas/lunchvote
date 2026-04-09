package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ozgurbaybas/lunchvote/modules/identity/application"
	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		writeError(w, http.StatusBadRequest, "id, name and email are required")
		return
	}

	user, err := h.service.CreateUser(r.Context(), application.CreateUserCommand{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrUserEmailAlreadyExists):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrInvalidUserID),
			errors.Is(err, domain.ErrInvalidUserName),
			errors.Is(err, domain.ErrInvalidUserEmail):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req createTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.OwnerID) == "" {
		writeError(w, http.StatusBadRequest, "id, name and owner_id are required")
		return
	}

	team, err := h.service.CreateTeam(r.Context(), application.CreateTeamCommand{
		ID:      req.ID,
		Name:    req.Name,
		OwnerID: req.OwnerID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrInvalidTeamID),
			errors.Is(err, domain.ErrInvalidTeamName),
			errors.Is(err, domain.ErrInvalidOwnerID):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, toTeamResponse(team))
}

func (h *Handler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	teamID := strings.TrimSpace(r.PathValue("id"))
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "team id is required")
		return
	}

	var req addTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	team, err := h.service.AddTeamMember(r.Context(), application.AddTeamMemberCommand{
		TeamID: teamID,
		UserID: req.UserID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTeamNotFound), errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrMemberAlreadyExists):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrInvalidUserID):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, toTeamResponse(team))
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
