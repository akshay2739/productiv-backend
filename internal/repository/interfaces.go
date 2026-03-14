package repository

import (
	"context"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
)

// UserRepository defines data access for users.
type UserRepository interface {
	GetDefault(ctx context.Context) (*domain.User, error)
}

// PillarRepository defines data access for pillars.
type PillarRepository interface {
	ListByUser(ctx context.Context, userID int64) ([]domain.Pillar, error)
	GetByType(ctx context.Context, userID int64, pillarType domain.PillarType) (*domain.Pillar, error)
}

// FastingRepository defines data access for fasting sessions.
type FastingRepository interface {
	Create(ctx context.Context, session *domain.FastingSession) error
	GetActive(ctx context.Context, userID int64) (*domain.FastingSession, error)
	EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, targetReached bool) error
	GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.FastingSession, error)
	CountCompleted(ctx context.Context, userID int64) (int, error)
	AverageDuration(ctx context.Context, userID int64) (float64, error)
	HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error)
	GetCompletedDurationsForMonth(ctx context.Context, userID int64, year int, month time.Month, loc *time.Location) (map[string]float64, error)
}

// GymRepository defines data access for gym sessions.
type GymRepository interface {
	Create(ctx context.Context, session *domain.GymSession) error
	HasLoggedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error)
	CountByWeek(ctx context.Context, userID int64, weekStart time.Time) (int, error)
	CountByMonth(ctx context.Context, userID int64, year int, month time.Month) (int, error)
	GetByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.GymSession, error)
}

// MeditationRepository defines data access for meditation sessions.
type MeditationRepository interface {
	Create(ctx context.Context, session *domain.MeditationSession) error
	GetActive(ctx context.Context, userID int64) (*domain.MeditationSession, error)
	EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, moodAfter *int) error
	GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.MeditationSession, error)
	TotalMinutes(ctx context.Context, userID int64) (float64, error)
	CountCompleted(ctx context.Context, userID int64) (int, error)
	AverageDuration(ctx context.Context, userID int64) (float64, error)
	HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error)
}

// RetentionRepository defines data access for retention streaks.
type RetentionRepository interface {
	Create(ctx context.Context, streak *domain.RetentionStreak) error
	GetActive(ctx context.Context, userID int64) (*domain.RetentionStreak, error)
	EndStreak(ctx context.Context, id int64, endDate time.Time, daysCount int, reason *string) error
	GetBestStreak(ctx context.Context, userID int64) (int, error)
	ListPast(ctx context.Context, userID int64, limit int) ([]domain.RetentionStreak, error)
}
