package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ozgurbaybas/lunchvote/modules/identity/application"
	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

type inMemoryUserRepository struct {
	byID    map[string]domain.User
	byEmail map[string]domain.User
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		byID:    make(map[string]domain.User),
		byEmail: make(map[string]domain.User),
	}
}

func (r *inMemoryUserRepository) Save(_ context.Context, user domain.User) error {
	r.byID[user.ID] = user
	r.byEmail[strings.ToLower(user.Email)] = user
	return nil
}

func (r *inMemoryUserRepository) GetByID(_ context.Context, id string) (domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepository) GetByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

type inMemoryTeamRepository struct {
	byID map[string]domain.Team
}

func newInMemoryTeamRepository() *inMemoryTeamRepository {
	return &inMemoryTeamRepository{
		byID: make(map[string]domain.Team),
	}
}

func (r *inMemoryTeamRepository) Save(_ context.Context, team domain.Team) error {
	r.byID[team.ID] = team
	return nil
}

func (r *inMemoryTeamRepository) GetByID(_ context.Context, id string) (domain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return domain.Team{}, domain.ErrTeamNotFound
	}
	return team, nil
}

func newTestMux() *nethttp.ServeMux {
	users := newInMemoryUserRepository()
	teams := newInMemoryTeamRepository()
	now := func() time.Time {
		return time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)
	}

	service := application.NewService(users, teams, now)
	handler := NewHandler(service)

	mux := nethttp.NewServeMux()
	RegisterRoutes(mux, handler)
	return mux
}

func TestCreateUser_ReturnsCreated(t *testing.T) {
	mux := newTestMux()

	body := []byte(`{"id":"user-1","name":"Ozgur","email":"ozgur@example.com"}`)
	req := httptest.NewRequest(nethttp.MethodPost, "/v1/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response["id"] != "user-1" {
		t.Fatalf("expected id user-1, got %v", response["id"])
	}
}

func TestCreateUser_ReturnsBadRequestWhenBodyInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/users", bytes.NewReader([]byte(`{`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateUser_ReturnsConflictWhenEmailExists(t *testing.T) {
	mux := newTestMux()

	first := httptest.NewRequest(nethttp.MethodPost, "/v1/users", bytes.NewReader([]byte(`{"id":"user-1","name":"Ozgur","email":"ozgur@example.com"}`)))
	firstRec := httptest.NewRecorder()
	mux.ServeHTTP(firstRec, first)

	second := httptest.NewRequest(nethttp.MethodPost, "/v1/users", bytes.NewReader([]byte(`{"id":"user-2","name":"Another","email":"ozgur@example.com"}`)))
	secondRec := httptest.NewRecorder()
	mux.ServeHTTP(secondRec, second)

	if secondRec.Code != nethttp.StatusConflict {
		t.Fatalf("expected status 409, got %d", secondRec.Code)
	}
}

func TestCreateTeam_ReturnsCreated(t *testing.T) {
	mux := newTestMux()

	createUserReq := httptest.NewRequest(nethttp.MethodPost, "/v1/users", bytes.NewReader([]byte(`{"id":"user-1","name":"Owner","email":"owner@example.com"}`)))
	createUserRec := httptest.NewRecorder()
	mux.ServeHTTP(createUserRec, createUserReq)

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/teams", bytes.NewReader([]byte(`{"id":"team-1","name":"Backend Team","owner_id":"user-1"}`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestCreateTeam_ReturnsNotFoundWhenOwnerMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/teams", bytes.NewReader([]byte(`{"id":"team-1","name":"Backend Team","owner_id":"missing-user"}`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestAddTeamMember_ReturnsOK(t *testing.T) {
	mux := newTestMux()

	requests := []struct {
		method string
		path   string
		body   string
	}{
		{method: nethttp.MethodPost, path: "/v1/users", body: `{"id":"user-1","name":"Owner","email":"owner@example.com"}`},
		{method: nethttp.MethodPost, path: "/v1/users", body: `{"id":"user-2","name":"Member","email":"member@example.com"}`},
		{method: nethttp.MethodPost, path: "/v1/teams", body: `{"id":"team-1","name":"Backend Team","owner_id":"user-1"}`},
	}

	for _, item := range requests {
		req := httptest.NewRequest(item.method, item.path, bytes.NewReader([]byte(item.body)))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/teams/team-1/members", bytes.NewReader([]byte(`{"user_id":"user-2"}`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestAddTeamMember_ReturnsNotFoundWhenTeamMissing(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/teams/missing-team/members", bytes.NewReader([]byte(`{"user_id":"user-2"}`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
