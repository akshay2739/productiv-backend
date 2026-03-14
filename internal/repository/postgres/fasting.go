package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// FastingRepo implements repository.FastingRepository using GORM.
type FastingRepo struct {
	db *gorm.DB
}

// NewFastingRepo creates a new FastingRepo.
func NewFastingRepo(db *gorm.DB) *FastingRepo {
	return &FastingRepo{db: db}
}

// Create inserts a new fasting session.
func (r *FastingRepo) Create(ctx context.Context, session *domain.FastingSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("creating fasting session: %w", err)
	}
	return nil
}

// GetActive returns the currently active fasting session for a user, or nil if none.
func (r *FastingRepo) GetActive(ctx context.Context, userID int64) (*domain.FastingSession, error) {
	var session domain.FastingSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time IS NULL", userID).
		Order("start_time DESC").
		First(&session).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting active fasting session: %w", err)
	}
	return &session, nil
}

// EndSession marks a fasting session as complete.
func (r *FastingRepo) EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, targetReached bool) error {
	err := r.db.WithContext(ctx).
		Model(&domain.FastingSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"end_time":              endTime,
			"actual_duration_hours": actualDuration,
			"target_reached":        targetReached,
		}).Error
	if err != nil {
		return fmt.Errorf("ending fasting session: %w", err)
	}
	return nil
}

// GetCompletedByDateRange returns completed fasting sessions within a date range.
func (r *FastingRepo) GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.FastingSession, error) {
	var sessions []domain.FastingSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time IS NOT NULL AND start_time >= ? AND start_time < ?", userID, start, end).
		Order("start_time DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("getting completed fasting sessions: %w", err)
	}
	return sessions, nil
}

// CountCompleted returns the total number of completed fasts for a user.
func (r *FastingRepo) CountCompleted(ctx context.Context, userID int64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.FastingSession{}).
		Where("user_id = ? AND end_time IS NOT NULL", userID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting completed fasts: %w", err)
	}
	return int(count), nil
}

// AverageDuration returns the average duration of completed fasts in hours.
func (r *FastingRepo) AverageDuration(ctx context.Context, userID int64) (float64, error) {
	var result struct{ Avg *float64 }
	err := r.db.WithContext(ctx).
		Model(&domain.FastingSession{}).
		Select("AVG(actual_duration_hours) as avg").
		Where("user_id = ? AND end_time IS NOT NULL", userID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("averaging fasting duration: %w", err)
	}
	if result.Avg == nil {
		return 0, nil
	}
	return *result.Avg, nil
}

// HasCompletedOnDate checks if a completed fast exists on a given date in the user's timezone.
func (r *FastingRepo) HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.FastingSession{}).
		Where("user_id = ? AND end_time IS NOT NULL AND end_time >= ? AND end_time < ?", userID, dayStart, dayEnd).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("checking fasting on date: %w", err)
	}
	return count > 0, nil
}
