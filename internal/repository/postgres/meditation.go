package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// MeditationRepo implements repository.MeditationRepository using GORM.
type MeditationRepo struct {
	db *gorm.DB
}

// NewMeditationRepo creates a new MeditationRepo.
func NewMeditationRepo(db *gorm.DB) *MeditationRepo {
	return &MeditationRepo{db: db}
}

// Create inserts a new meditation session.
func (r *MeditationRepo) Create(ctx context.Context, session *domain.MeditationSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("creating meditation session: %w", err)
	}
	return nil
}

// GetActive returns the currently active meditation session for a user, or nil if none.
func (r *MeditationRepo) GetActive(ctx context.Context, userID int64) (*domain.MeditationSession, error) {
	var session domain.MeditationSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time IS NULL", userID).
		Order("start_time DESC").
		First(&session).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting active meditation session: %w", err)
	}
	return &session, nil
}

// EndSession marks a meditation session as complete.
func (r *MeditationRepo) EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, moodAfter *int) error {
	err := r.db.WithContext(ctx).
		Model(&domain.MeditationSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"end_time":                endTime,
			"actual_duration_minutes": actualDuration,
			"mood_after":              moodAfter,
		}).Error
	if err != nil {
		return fmt.Errorf("ending meditation session: %w", err)
	}
	return nil
}

// GetCompletedByDateRange returns completed meditation sessions within a date range.
func (r *MeditationRepo) GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.MeditationSession, error) {
	var sessions []domain.MeditationSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time IS NOT NULL AND start_time >= ? AND start_time < ?", userID, start, end).
		Order("start_time DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("getting completed meditation sessions: %w", err)
	}
	return sessions, nil
}

// TotalMinutes returns the total minutes meditated for a user.
func (r *MeditationRepo) TotalMinutes(ctx context.Context, userID int64) (float64, error) {
	var result struct{ Total *float64 }
	err := r.db.WithContext(ctx).
		Model(&domain.MeditationSession{}).
		Select("SUM(actual_duration_minutes) as total").
		Where("user_id = ? AND end_time IS NOT NULL", userID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("totaling meditation minutes: %w", err)
	}
	if result.Total == nil {
		return 0, nil
	}
	return *result.Total, nil
}

// CountCompleted returns the total number of completed meditation sessions.
func (r *MeditationRepo) CountCompleted(ctx context.Context, userID int64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MeditationSession{}).
		Where("user_id = ? AND end_time IS NOT NULL", userID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting completed meditations: %w", err)
	}
	return int(count), nil
}

// AverageDuration returns the average session length in minutes.
func (r *MeditationRepo) AverageDuration(ctx context.Context, userID int64) (float64, error) {
	var result struct{ Avg *float64 }
	err := r.db.WithContext(ctx).
		Model(&domain.MeditationSession{}).
		Select("AVG(actual_duration_minutes) as avg").
		Where("user_id = ? AND end_time IS NOT NULL", userID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("averaging meditation duration: %w", err)
	}
	if result.Avg == nil {
		return 0, nil
	}
	return *result.Avg, nil
}

// HasCompletedOnDate checks if a completed meditation exists on a given date.
func (r *MeditationRepo) HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MeditationSession{}).
		Where("user_id = ? AND end_time IS NOT NULL AND end_time >= ? AND end_time < ?", userID, dayStart, dayEnd).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("checking meditation on date: %w", err)
	}
	return count > 0, nil
}
