package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/modules/identity/domain"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Save(ctx context.Context, team domain.Team) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin team transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const insertTeamQuery = `
		INSERT INTO teams (id, name, owner_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			owner_id = EXCLUDED.owner_id
	`

	_, err = tx.Exec(ctx, insertTeamQuery, team.ID, team.Name, team.OwnerID, team.CreatedAt)
	if err != nil {
		return fmt.Errorf("upsert team: %w", err)
	}

	const deleteMembershipsQuery = `
		DELETE FROM team_memberships
		WHERE team_id = $1
	`

	_, err = tx.Exec(ctx, deleteMembershipsQuery, team.ID)
	if err != nil {
		return fmt.Errorf("delete team memberships: %w", err)
	}

	const insertMembershipQuery = `
		INSERT INTO team_memberships (team_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)
	`

	for _, member := range team.Members {
		_, err = tx.Exec(ctx, insertMembershipQuery, team.ID, member.UserID, string(member.Role), member.JoinedAt)
		if err != nil {
			return fmt.Errorf("insert team membership: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit team transaction: %w", err)
	}

	return nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id string) (domain.Team, error) {
	const teamQuery = `
		SELECT id, name, owner_id, created_at
		FROM teams
		WHERE id = $1
	`

	var team domain.Team
	err := r.db.QueryRow(ctx, teamQuery, id).Scan(
		&team.ID,
		&team.Name,
		&team.OwnerID,
		&team.CreatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return domain.Team{}, domain.ErrTeamNotFound
		}
		return domain.Team{}, fmt.Errorf("select team by id: %w", err)
	}

	const membersQuery = `
		SELECT user_id, role, joined_at
		FROM team_memberships
		WHERE team_id = $1
		ORDER BY joined_at ASC
	`

	rows, err := r.db.Query(ctx, membersQuery, id)
	if err != nil {
		return domain.Team{}, fmt.Errorf("select team memberships: %w", err)
	}
	defer rows.Close()

	members := make([]domain.Membership, 0)
	for rows.Next() {
		var member domain.Membership
		var role string

		if err := rows.Scan(&member.UserID, &role, &member.JoinedAt); err != nil {
			return domain.Team{}, fmt.Errorf("scan team membership: %w", err)
		}

		member.Role = domain.MembershipRole(role)
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return domain.Team{}, fmt.Errorf("iterate team memberships: %w", err)
	}

	team.Members = members
	return team, nil
}
