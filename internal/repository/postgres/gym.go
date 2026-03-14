package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// GymRepo implements repository.GymRepository using GORM.
type GymRepo struct {
	db *gorm.DB
}

// NewGymRepo creates a new GymRepo.
func NewGymRepo(db *gorm.DB) *GymRepo {
	return &GymRepo{db: db}
}

// Create inserts a new gym session.
func (r *GymRepo) Create(ctx context.Context, session *domain.GymSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("creating gym session: %w", err)
	}
	return nil
}

// HasLoggedOnDate checks if a workout has already been logged on the given date.
func (r *GymRepo) HasLoggedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.GymSession{}).
		Where("user_id = ? AND logged_at >= ? AND logged_at < ?", userID, dayStart, dayEnd).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("checking gym log on date: %w", err)
	}
	return count > 0, nil
}

// CountByWeek returns the number of workouts in the week starting from weekStart.
func (r *GymRepo) CountByWeek(ctx context.Context, userID int64, weekStart time.Time) (int, error) {
	weekEnd := weekStart.Add(7 * 24 * time.Hour)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.GymSession{}).
		Where("user_id = ? AND logged_at >= ? AND logged_at < ?", userID, weekStart, weekEnd).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting weekly workouts: %w", err)
	}
	return int(count), nil
}

// CountByMonth returns the number of workouts in the given month.
func (r *GymRepo) CountByMonth(ctx context.Context, userID int64, year int, month time.Month) (int, error) {
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.GymSession{}).
		Where("user_id = ? AND logged_at >= ? AND logged_at < ?", userID, monthStart, monthEnd).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting monthly workouts: %w", err)
	}
	return int(count), nil
}

// GetByDateRange returns gym sessions within a date range.
func (r *GymRepo) GetByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.GymSession, error) {
	var sessions []domain.GymSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND logged_at >= ? AND logged_at < ?", userID, start, end).
		Order("logged_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("getting gym sessions by date range: %w", err)
	}
	return sessions, nil
}
