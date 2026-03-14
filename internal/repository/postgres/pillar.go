package postgres

import (
	"context"
	"fmt"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// PillarRepo implements repository.PillarRepository using GORM.
type PillarRepo struct {
	db *gorm.DB
}

// NewPillarRepo creates a new PillarRepo.
func NewPillarRepo(db *gorm.DB) *PillarRepo {
	return &PillarRepo{db: db}
}

// ListByUser returns all active pillars for a user ordered by display_order.
func (r *PillarRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Pillar, error) {
	var pillars []domain.Pillar
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("display_order").
		Find(&pillars).Error
	if err != nil {
		return nil, fmt.Errorf("listing pillars: %w", err)
	}
	return pillars, nil
}

// GetByType returns a specific pillar by type for a user.
func (r *PillarRepo) GetByType(ctx context.Context, userID int64, pillarType domain.PillarType) (*domain.Pillar, error) {
	var pillar domain.Pillar
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, pillarType).
		First(&pillar).Error
	if err != nil {
		return nil, fmt.Errorf("getting pillar by type: %w", err)
	}
	return &pillar, nil
}
