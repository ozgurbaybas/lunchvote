package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	polldomain "github.com/ozgurbaybas/lunchvote/modules/poll/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, poll polldomain.Poll) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin poll transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const upsertPollQuery = `
		INSERT INTO polls (id, team_id, title, status, created_at, closed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id)
		DO UPDATE SET
			team_id = EXCLUDED.team_id,
			title = EXCLUDED.title,
			status = EXCLUDED.status,
			closed_at = EXCLUDED.closed_at
	`

	_, err = tx.Exec(
		ctx,
		upsertPollQuery,
		poll.ID,
		poll.TeamID,
		poll.Title,
		string(poll.Status),
		poll.CreatedAt,
		poll.ClosedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert poll: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM poll_options WHERE poll_id = $1`, poll.ID)
	if err != nil {
		return fmt.Errorf("delete poll options: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM poll_votes WHERE poll_id = $1`, poll.ID)
	if err != nil {
		return fmt.Errorf("delete poll votes: %w", err)
	}

	const insertOptionQuery = `
		INSERT INTO poll_options (poll_id, restaurant_id)
		VALUES ($1, $2)
	`

	for _, option := range poll.Options {
		_, err = tx.Exec(ctx, insertOptionQuery, poll.ID, option.RestaurantID)
		if err != nil {
			return fmt.Errorf("insert poll option: %w", err)
		}
	}

	const insertVoteQuery = `
		INSERT INTO poll_votes (poll_id, user_id, restaurant_id, voted_at)
		VALUES ($1, $2, $3, $4)
	`

	for _, vote := range poll.Votes {
		_, err = tx.Exec(ctx, insertVoteQuery, poll.ID, vote.UserID, vote.RestaurantID, vote.VotedAt)
		if err != nil {
			return fmt.Errorf("insert poll vote: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit poll transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (polldomain.Poll, error) {
	const pollQuery = `
		SELECT id, team_id, title, status, created_at, closed_at
		FROM polls
		WHERE id = $1
	`

	var (
		poll     polldomain.Poll
		status   string
		closedAt *time.Time
	)

	err := r.db.QueryRow(ctx, pollQuery, id).Scan(
		&poll.ID,
		&poll.TeamID,
		&poll.Title,
		&status,
		&poll.CreatedAt,
		&closedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return polldomain.Poll{}, polldomain.ErrPollNotFound
		}
		return polldomain.Poll{}, fmt.Errorf("select poll: %w", err)
	}

	poll.Status = polldomain.PollStatus(status)
	poll.ClosedAt = closedAt
	poll.Options = make([]polldomain.PollOption, 0)
	poll.Votes = make([]polldomain.Vote, 0)

	optionRows, err := r.db.Query(ctx, `
		SELECT restaurant_id
		FROM poll_options
		WHERE poll_id = $1
		ORDER BY restaurant_id ASC
	`, id)
	if err != nil {
		return polldomain.Poll{}, fmt.Errorf("select poll options: %w", err)
	}
	defer optionRows.Close()

	for optionRows.Next() {
		var option polldomain.PollOption
		if err := optionRows.Scan(&option.RestaurantID); err != nil {
			return polldomain.Poll{}, fmt.Errorf("scan poll option: %w", err)
		}
		poll.Options = append(poll.Options, option)
	}
	if err := optionRows.Err(); err != nil {
		return polldomain.Poll{}, fmt.Errorf("iterate poll options: %w", err)
	}

	voteRows, err := r.db.Query(ctx, `
		SELECT user_id, restaurant_id, voted_at
		FROM poll_votes
		WHERE poll_id = $1
		ORDER BY voted_at ASC
	`, id)
	if err != nil {
		return polldomain.Poll{}, fmt.Errorf("select poll votes: %w", err)
	}
	defer voteRows.Close()

	for voteRows.Next() {
		var vote polldomain.Vote
		if err := voteRows.Scan(&vote.UserID, &vote.RestaurantID, &vote.VotedAt); err != nil {
			return polldomain.Poll{}, fmt.Errorf("scan poll vote: %w", err)
		}
		poll.Votes = append(poll.Votes, vote)
	}
	if err := voteRows.Err(); err != nil {
		return polldomain.Poll{}, fmt.Errorf("iterate poll votes: %w", err)
	}

	return poll, nil
}

func (r *Repository) ListByTeamID(ctx context.Context, teamID string) ([]polldomain.Poll, error) {
	const query = `
		SELECT id
		FROM polls
		WHERE team_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("list polls by team: %w", err)
	}
	defer rows.Close()

	polls := make([]polldomain.Poll, 0)
	for rows.Next() {
		var pollID string
		if err := rows.Scan(&pollID); err != nil {
			return nil, fmt.Errorf("scan poll id: %w", err)
		}

		poll, err := r.GetByID(ctx, pollID)
		if err != nil {
			return nil, fmt.Errorf("get poll by id from list: %w", err)
		}

		polls = append(polls, poll)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate polls by team: %w", err)
	}

	return polls, nil
}
