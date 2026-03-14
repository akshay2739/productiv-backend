package postgres

import (
	"context"
	"fmt"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// UserRepo implements repository.UserRepository using GORM.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetDefault returns the default user (first user).
func (r *UserRepo) GetDefault(ctx context.Context) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user).Error; err != nil {
		return nil, fmt.Errorf("getting default user: %w", err)
	}
	return &user, nil
}
