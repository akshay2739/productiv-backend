package mocks

import (
	"context"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks repository.UserRepository.
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) GetDefault(ctx context.Context) (*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockPillarRepository mocks repository.PillarRepository.
type MockPillarRepository struct{ mock.Mock }

func (m *MockPillarRepository) ListByUser(ctx context.Context, userID int64) ([]domain.Pillar, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Pillar), args.Error(1)
}

func (m *MockPillarRepository) GetByType(ctx context.Context, userID int64, pillarType domain.PillarType) (*domain.Pillar, error) {
	args := m.Called(ctx, userID, pillarType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pillar), args.Error(1)
}

// MockFastingRepository mocks repository.FastingRepository.
type MockFastingRepository struct{ mock.Mock }

func (m *MockFastingRepository) Create(ctx context.Context, session *domain.FastingSession) error {
	return m.Called(ctx, session).Error(0)
}

func (m *MockFastingRepository) GetActive(ctx context.Context, userID int64) (*domain.FastingSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FastingSession), args.Error(1)
}

func (m *MockFastingRepository) EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, targetReached bool) error {
	return m.Called(ctx, id, endTime, actualDuration, targetReached).Error(0)
}

func (m *MockFastingRepository) GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.FastingSession, error) {
	args := m.Called(ctx, userID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FastingSession), args.Error(1)
}

func (m *MockFastingRepository) CountCompleted(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockFastingRepository) AverageDuration(ctx context.Context, userID int64) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockFastingRepository) HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	args := m.Called(ctx, userID, date)
	return args.Bool(0), args.Error(1)
}

func (m *MockFastingRepository) GetCompletedDurationsForMonth(ctx context.Context, userID int64, year int, month time.Month, loc *time.Location) (map[string]float64, error) {
	args := m.Called(ctx, userID, year, month, loc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]float64), args.Error(1)
}

// MockGymRepository mocks repository.GymRepository.
type MockGymRepository struct{ mock.Mock }

func (m *MockGymRepository) Create(ctx context.Context, session *domain.GymSession) error {
	return m.Called(ctx, session).Error(0)
}

func (m *MockGymRepository) HasLoggedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	args := m.Called(ctx, userID, date)
	return args.Bool(0), args.Error(1)
}

func (m *MockGymRepository) CountByWeek(ctx context.Context, userID int64, weekStart time.Time) (int, error) {
	args := m.Called(ctx, userID, weekStart)
	return args.Int(0), args.Error(1)
}

func (m *MockGymRepository) CountByMonth(ctx context.Context, userID int64, year int, month time.Month) (int, error) {
	args := m.Called(ctx, userID, year, month)
	return args.Int(0), args.Error(1)
}

func (m *MockGymRepository) GetByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.GymSession, error) {
	args := m.Called(ctx, userID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.GymSession), args.Error(1)
}

// MockMeditationRepository mocks repository.MeditationRepository.
type MockMeditationRepository struct{ mock.Mock }

func (m *MockMeditationRepository) Create(ctx context.Context, session *domain.MeditationSession) error {
	return m.Called(ctx, session).Error(0)
}

func (m *MockMeditationRepository) GetActive(ctx context.Context, userID int64) (*domain.MeditationSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MeditationSession), args.Error(1)
}

func (m *MockMeditationRepository) EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, moodAfter *int) error {
	return m.Called(ctx, id, endTime, actualDuration, moodAfter).Error(0)
}

func (m *MockMeditationRepository) GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.MeditationSession, error) {
	args := m.Called(ctx, userID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MeditationSession), args.Error(1)
}

func (m *MockMeditationRepository) TotalMinutes(ctx context.Context, userID int64) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockMeditationRepository) CountCompleted(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockMeditationRepository) AverageDuration(ctx context.Context, userID int64) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockMeditationRepository) HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	args := m.Called(ctx, userID, date)
	return args.Bool(0), args.Error(1)
}

// MockRetentionRepository mocks repository.RetentionRepository.
type MockRetentionRepository struct{ mock.Mock }

func (m *MockRetentionRepository) Create(ctx context.Context, streak *domain.RetentionStreak) error {
	return m.Called(ctx, streak).Error(0)
}

func (m *MockRetentionRepository) GetActive(ctx context.Context, userID int64) (*domain.RetentionStreak, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RetentionStreak), args.Error(1)
}

func (m *MockRetentionRepository) EndStreak(ctx context.Context, id int64, endDate time.Time, daysCount int, reason *string) error {
	return m.Called(ctx, id, endDate, daysCount, reason).Error(0)
}

func (m *MockRetentionRepository) GetBestStreak(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockRetentionRepository) ListPast(ctx context.Context, userID int64, limit int) ([]domain.RetentionStreak, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RetentionStreak), args.Error(1)
}
