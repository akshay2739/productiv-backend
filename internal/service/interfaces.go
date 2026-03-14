package service

import (
	"context"

	"github.com/akshay/productiv-backend/internal/domain"
)

// FastingServiceInterface defines the fasting business operations.
type FastingServiceInterface interface {
	StartFast(ctx context.Context, userID int64, req StartFastRequest) (*domain.FastingSession, error)
	EndFast(ctx context.Context, userID int64) (*domain.FastingSession, error)
	GetStats(ctx context.Context, userID int64) (*domain.FastingStats, error)
	GetProtocols() []domain.FastingProtocol
}

// GymServiceInterface defines the gym business operations.
type GymServiceInterface interface {
	LogWorkout(ctx context.Context, userID int64, req LogWorkoutRequest) (*domain.GymSession, error)
	GetStats(ctx context.Context, userID int64) (*domain.GymStats, error)
}

// MeditationServiceInterface defines the meditation business operations.
type MeditationServiceInterface interface {
	StartSession(ctx context.Context, userID int64, req StartSessionRequest) (*domain.MeditationSession, error)
	EndSession(ctx context.Context, userID int64, req EndSessionRequest) (*domain.MeditationSession, error)
	GetStats(ctx context.Context, userID int64) (*domain.MeditationStats, error)
}

// RetentionServiceInterface defines the retention business operations.
type RetentionServiceInterface interface {
	StartTracking(ctx context.Context, userID int64) (*domain.RetentionStreak, error)
	ResetCounter(ctx context.Context, userID int64, req ResetRequest) (*domain.RetentionStreak, error)
	GetStats(ctx context.Context, userID int64) (*domain.RetentionStats, error)
}

// DashboardServiceInterface defines the dashboard business operations.
type DashboardServiceInterface interface {
	GetDashboard(ctx context.Context, userID int64) (*domain.DashboardData, error)
}
