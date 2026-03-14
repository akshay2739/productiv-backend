package postgres

import (
	"context"
	"fmt"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepo implements repository.UserRepository using PostgreSQL.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// GetDefault returns the default user (id=1).
func (r *UserRepo) GetDefault(ctx context.Context) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, timezone, created_at, updated_at FROM users WHERE id = 1`,
	).Scan(&u.ID, &u.Name, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting default user: %w", err)
	}
	return &u, nil
}
