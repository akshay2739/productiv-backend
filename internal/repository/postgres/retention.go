package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// RetentionRepo implements repository.RetentionRepository using GORM.
type RetentionRepo struct {
	db *gorm.DB
}

// NewRetentionRepo creates a new RetentionRepo.
func NewRetentionRepo(db *gorm.DB) *RetentionRepo {
	return &RetentionRepo{db: db}
}

// Create inserts a new retention streak.
func (r *RetentionRepo) Create(ctx context.Context, streak *domain.RetentionStreak) error {
	if err := r.db.WithContext(ctx).Create(streak).Error; err != nil {
		return fmt.Errorf("creating retention streak: %w", err)
	}
	return nil
}

// GetActive returns the currently active retention streak for a user, or nil if none.
func (r *RetentionRepo) GetActive(ctx context.Context, userID int64) (*domain.RetentionStreak, error) {
	var streak domain.RetentionStreak
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_date IS NULL", userID).
		Order("start_date DESC").
		First(&streak).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting active retention streak: %w", err)
	}
	return &streak, nil
}

// EndStreak marks a retention streak as ended.
func (r *RetentionRepo) EndStreak(ctx context.Context, id int64, endDate time.Time, daysCount int, reason *string) error {
	err := r.db.WithContext(ctx).
		Model(&domain.RetentionStreak{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"end_date":   endDate,
			"days_count": daysCount,
			"reason":     reason,
		}).Error
	if err != nil {
		return fmt.Errorf("ending retention streak: %w", err)
	}
	return nil
}

// GetBestStreak returns the highest days_count across all streaks for a user.
func (r *RetentionRepo) GetBestStreak(ctx context.Context, userID int64) (int, error) {
	var result struct{ Best *int }
	err := r.db.WithContext(ctx).
		Model(&domain.RetentionStreak{}).
		Select("MAX(days_count) as best").
		Where("user_id = ?", userID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("getting best retention streak: %w", err)
	}
	if result.Best == nil {
		return 0, nil
	}
	return *result.Best, nil
}

// ListPast returns past (ended) retention streaks, most recent first.
func (r *RetentionRepo) ListPast(ctx context.Context, userID int64, limit int) ([]domain.RetentionStreak, error) {
	var streaks []domain.RetentionStreak
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_date IS NOT NULL", userID).
		Order("end_date DESC").
		Limit(limit).
		Find(&streaks).Error
	if err != nil {
		return nil, fmt.Errorf("listing past retention streaks: %w", err)
	}
	return streaks, nil
}
