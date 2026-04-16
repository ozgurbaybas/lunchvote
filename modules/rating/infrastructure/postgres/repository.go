package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	ratingdomain "github.com/ozgurbaybas/lunchvote/modules/rating/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, rating ratingdomain.Rating) error {
	const query = `
		INSERT INTO ratings (
			id,
			restaurant_id,
			user_id,
			score,
			comment,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		rating.ID,
		rating.RestaurantID,
		rating.UserID,
		rating.Score,
		rating.Comment,
		rating.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		// 23505: unique constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ratingdomain.ErrRatingAlreadyExists
		}
		return fmt.Errorf("insert rating: %w", err)
	}

	return nil
}

func (r *Repository) GetByRestaurantAndUser(
	ctx context.Context,
	restaurantID string,
	userID string,
) (ratingdomain.Rating, error) {
	const query = `
		SELECT id, restaurant_id, user_id, score, comment, created_at
		FROM ratings
		WHERE restaurant_id = $1 AND user_id = $2
	`

	var rating ratingdomain.Rating
	err := r.db.QueryRow(ctx, query, restaurantID, userID).Scan(
		&rating.ID,
		&rating.RestaurantID,
		&rating.UserID,
		&rating.Score,
		&rating.Comment,
		&rating.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ratingdomain.Rating{}, ratingdomain.ErrRatingNotFound
		}
		return ratingdomain.Rating{}, fmt.Errorf("select rating: %w", err)
	}

	return rating, nil
}

func (r *Repository) ListByRestaurantID(
	ctx context.Context,
	restaurantID string,
) ([]ratingdomain.Rating, error) {
	const query = `
		SELECT id, restaurant_id, user_id, score, comment, created_at
		FROM ratings
		WHERE restaurant_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list ratings: %w", err)
	}
	defer rows.Close()

	var ratings []ratingdomain.Rating
	for rows.Next() {
		var rating ratingdomain.Rating
		if err := rows.Scan(
			&rating.ID,
			&rating.RestaurantID,
			&rating.UserID,
			&rating.Score,
			&rating.Comment,
			&rating.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rating: %w", err)
		}
		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ratings: %w", err)
	}

	return ratings, nil
}
