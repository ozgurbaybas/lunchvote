package postgres

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ozgurbaybas/lunchvote/modules/restaurant/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, restaurant domain.Restaurant) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin restaurant transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const insertRestaurantQuery = `
		INSERT INTO restaurants (
			id,
			name,
			address,
			city,
			district,
			is_active,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.Exec(
		ctx,
		insertRestaurantQuery,
		restaurant.ID,
		restaurant.Name,
		restaurant.Address,
		restaurant.City,
		restaurant.District,
		restaurant.IsActive,
		restaurant.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert restaurant: %w", err)
	}

	const insertMealCardQuery = `
		INSERT INTO restaurant_meal_cards (restaurant_id, meal_card)
		VALUES ($1, $2)
	`

	for _, mealCard := range restaurant.SupportedMealCards {
		_, err = tx.Exec(ctx, insertMealCardQuery, restaurant.ID, string(mealCard))
		if err != nil {
			return fmt.Errorf("insert restaurant meal card: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit restaurant transaction: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Restaurant, error) {
	const query = `
		SELECT
			r.id,
			r.name,
			r.address,
			r.city,
			r.district,
			r.is_active,
			r.created_at,
			COALESCE(rmc.meal_card, '')
		FROM restaurants r
		LEFT JOIN restaurant_meal_cards rmc ON rmc.restaurant_id = r.id
		ORDER BY r.created_at ASC, rmc.meal_card ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list restaurants: %w", err)
	}
	defer rows.Close()

	restaurantsByID := make(map[string]*domain.Restaurant)

	for rows.Next() {
		var (
			id        string
			name      string
			address   string
			city      string
			district  string
			isActive  bool
			createdAt time.Time
			mealCard  string
		)

		if err := rows.Scan(
			&id,
			&name,
			&address,
			&city,
			&district,
			&isActive,
			&createdAt,
			&mealCard,
		); err != nil {
			return nil, fmt.Errorf("scan restaurant row: %w", err)
		}

		restaurant, ok := restaurantsByID[id]
		if !ok {
			restaurant = &domain.Restaurant{
				ID:                 id,
				Name:               name,
				Address:            address,
				City:               city,
				District:           district,
				SupportedMealCards: make([]domain.MealCard, 0),
				IsActive:           isActive,
				CreatedAt:          createdAt,
			}
			restaurantsByID[id] = restaurant
		}

		if mealCard != "" {
			restaurant.SupportedMealCards = append(restaurant.SupportedMealCards, domain.MealCard(mealCard))
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate restaurant rows: %w", err)
	}

	restaurants := make([]domain.Restaurant, 0, len(restaurantsByID))
	for _, restaurant := range restaurantsByID {
		restaurants = append(restaurants, *restaurant)
	}

	sort.Slice(restaurants, func(i, j int) bool {
		return restaurants[i].CreatedAt.Before(restaurants[j].CreatedAt)
	})

	return restaurants, nil
}
