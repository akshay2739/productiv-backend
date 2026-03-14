package postgres

import (
	"context"
	"fmt"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PillarRepo implements repository.PillarRepository using PostgreSQL.
type PillarRepo struct {
	pool *pgxpool.Pool
}

// NewPillarRepo creates a new PillarRepo.
func NewPillarRepo(pool *pgxpool.Pool) *PillarRepo {
	return &PillarRepo{pool: pool}
}

// ListByUser returns all active pillars for a user ordered by display_order.
func (r *PillarRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Pillar, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, type, name, icon, color, is_active, display_order, created_at, updated_at
		 FROM pillars
		 WHERE user_id = $1 AND is_active = TRUE
		 ORDER BY display_order`, userID)
	if err != nil {
		return nil, fmt.Errorf("listing pillars: %w", err)
	}
	defer rows.Close()

	var pillars []domain.Pillar
	for rows.Next() {
		var p domain.Pillar
		if err := rows.Scan(&p.ID, &p.UserID, &p.Type, &p.Name, &p.Icon, &p.Color,
			&p.IsActive, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning pillar: %w", err)
		}
		pillars = append(pillars, p)
	}
	return pillars, rows.Err()
}

// GetByType returns a specific pillar by type for a user.
func (r *PillarRepo) GetByType(ctx context.Context, userID int64, pillarType domain.PillarType) (*domain.Pillar, error) {
	var p domain.Pillar
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, type, name, icon, color, is_active, display_order, created_at, updated_at
		 FROM pillars
		 WHERE user_id = $1 AND type = $2`, userID, pillarType,
	).Scan(&p.ID, &p.UserID, &p.Type, &p.Name, &p.Icon, &p.Color,
		&p.IsActive, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting pillar by type: %w", err)
	}
	return &p, nil
}
