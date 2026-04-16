package http

import (
	"bytes"
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	identitydomain "github.com/ozgurbaybas/lunchvote/modules/identity/domain"
	pollapp "github.com/ozgurbaybas/lunchvote/modules/poll/application"
	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

type inMemoryPollRepository struct {
	byID map[string]polldomain.Poll
}

func newInMemoryPollRepository() *inMemoryPollRepository {
	return &inMemoryPollRepository{
		byID: make(map[string]polldomain.Poll),
	}
}

func (r *inMemoryPollRepository) Save(_ context.Context, poll polldomain.Poll) error {
	r.byID[poll.ID] = poll
	return nil
}

func (r *inMemoryPollRepository) GetByID(_ context.Context, id string) (polldomain.Poll, error) {
	poll, ok := r.byID[id]
	if !ok {
		return polldomain.Poll{}, polldomain.ErrPollNotFound
	}
	return poll, nil
}

type inMemoryTeamRepository struct {
	byID map[string]identitydomain.Team
}

func newInMemoryTeamRepository() *inMemoryTeamRepository {
	return &inMemoryTeamRepository{
		byID: make(map[string]identitydomain.Team),
	}
}

func (r *inMemoryTeamRepository) Save(_ context.Context, team identitydomain.Team) error {
	r.byID[team.ID] = team
	return nil
}

func (r *inMemoryTeamRepository) GetByID(_ context.Context, id string) (identitydomain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return identitydomain.Team{}, identitydomain.ErrTeamNotFound
	}
	return team, nil
}

func newTestMux() *nethttp.ServeMux {
	pollRepo := newInMemoryPollRepository()
	teamRepo := newInMemoryTeamRepository()
	now := time.Date(2026, time.April, 12, 12, 0, 0, 0, time.UTC)

	team, _ := identitydomain.NewTeam("team-1", "Backend Team", "user-1", now)
	_ = team.AddMember("user-2", now)
	_ = teamRepo.Save(context.Background(), team)

	service := pollapp.NewService(pollRepo, teamRepo, func() time.Time { return now })
	handler := NewHandler(service)

	mux := nethttp.NewServeMux()
	RegisterRoutes(mux, handler)
	return mux
}

func TestCreatePoll_ReturnsCreated(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls",
		bytes.NewReader([]byte(`{
			"id":"poll-1",
			"team_id":"team-1",
			"title":"Friday Lunch",
			"restaurant_ids":["restaurant-1","restaurant-2"],
			"creator_user_id":"user-1"
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestCreatePoll_ReturnsBadRequestWhenBodyInvalid(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(nethttp.MethodPost, "/v1/polls", bytes.NewReader([]byte(`{`)))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestCreatePoll_ReturnsForbiddenWhenCreatorNotTeamMember(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls",
		bytes.NewReader([]byte(`{
			"id":"poll-1",
			"team_id":"team-1",
			"title":"Friday Lunch",
			"restaurant_ids":["restaurant-1","restaurant-2"],
			"creator_user_id":"user-x"
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
}

func TestVote_ReturnsOK(t *testing.T) {
	mux := newTestMux()

	createReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls",
		bytes.NewReader([]byte(`{
			"id":"poll-1",
			"team_id":"team-1",
			"title":"Friday Lunch",
			"restaurant_ids":["restaurant-1","restaurant-2"],
			"creator_user_id":"user-1"
		}`)),
	)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	req := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls/poll-1/votes",
		bytes.NewReader([]byte(`{
			"user_id":"user-2",
			"restaurant_id":"restaurant-1"
		}`)),
	)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestVote_ReturnsConflictWhenDuplicate(t *testing.T) {
	mux := newTestMux()

	createReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls",
		bytes.NewReader([]byte(`{
			"id":"poll-1",
			"team_id":"team-1",
			"title":"Friday Lunch",
			"restaurant_ids":["restaurant-1","restaurant-2"],
			"creator_user_id":"user-1"
		}`)),
	)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	firstVoteReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls/poll-1/votes",
		bytes.NewReader([]byte(`{
			"user_id":"user-2",
			"restaurant_id":"restaurant-1"
		}`)),
	)
	firstVoteRec := httptest.NewRecorder()
	mux.ServeHTTP(firstVoteRec, firstVoteReq)

	secondVoteReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls/poll-1/votes",
		bytes.NewReader([]byte(`{
			"user_id":"user-2",
			"restaurant_id":"restaurant-2"
		}`)),
	)
	secondVoteRec := httptest.NewRecorder()
	mux.ServeHTTP(secondVoteRec, secondVoteReq)

	if secondVoteRec.Code != nethttp.StatusConflict {
		t.Fatalf("expected status 409, got %d", secondVoteRec.Code)
	}
}

func TestResults_ReturnsOK(t *testing.T) {
	mux := newTestMux()

	createReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls",
		bytes.NewReader([]byte(`{
			"id":"poll-1",
			"team_id":"team-1",
			"title":"Friday Lunch",
			"restaurant_ids":["restaurant-1","restaurant-2"],
			"creator_user_id":"user-1"
		}`)),
	)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	voteReq := httptest.NewRequest(
		nethttp.MethodPost,
		"/v1/polls/poll-1/votes",
		bytes.NewReader([]byte(`{
			"user_id":"user-2",
			"restaurant_id":"restaurant-1"
		}`)),
	)
	voteRec := httptest.NewRecorder()
	mux.ServeHTTP(voteRec, voteReq)

	req := httptest.NewRequest(nethttp.MethodGet, "/v1/polls/poll-1/results", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response["poll_id"] != "poll-1" {
		t.Fatalf("expected poll_id poll-1, got %v", response["poll_id"])
	}
}
