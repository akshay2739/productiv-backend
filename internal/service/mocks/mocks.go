package mocks

import (
	"context"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// MockFastingService mocks service.FastingServiceInterface.
type MockFastingService struct{ mock.Mock }

func (m *MockFastingService) StartFast(ctx context.Context, userID int64, req service.StartFastRequest) (*domain.FastingSession, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FastingSession), args.Error(1)
}

func (m *MockFastingService) EndFast(ctx context.Context, userID int64) (*domain.FastingSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FastingSession), args.Error(1)
}

func (m *MockFastingService) GetStats(ctx context.Context, userID int64) (*domain.FastingStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FastingStats), args.Error(1)
}

func (m *MockFastingService) GetProtocols() []domain.FastingProtocol {
	args := m.Called()
	return args.Get(0).([]domain.FastingProtocol)
}

// MockGymService mocks service.GymServiceInterface.
type MockGymService struct{ mock.Mock }

func (m *MockGymService) LogWorkout(ctx context.Context, userID int64, req service.LogWorkoutRequest) (*domain.GymSession, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GymSession), args.Error(1)
}

func (m *MockGymService) GetStats(ctx context.Context, userID int64) (*domain.GymStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GymStats), args.Error(1)
}

// MockMeditationService mocks service.MeditationServiceInterface.
type MockMeditationService struct{ mock.Mock }

func (m *MockMeditationService) StartSession(ctx context.Context, userID int64, req service.StartSessionRequest) (*domain.MeditationSession, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MeditationSession), args.Error(1)
}

func (m *MockMeditationService) EndSession(ctx context.Context, userID int64, req service.EndSessionRequest) (*domain.MeditationSession, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MeditationSession), args.Error(1)
}

func (m *MockMeditationService) GetStats(ctx context.Context, userID int64) (*domain.MeditationStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MeditationStats), args.Error(1)
}

// MockRetentionService mocks service.RetentionServiceInterface.
type MockRetentionService struct{ mock.Mock }

func (m *MockRetentionService) StartTracking(ctx context.Context, userID int64) (*domain.RetentionStreak, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RetentionStreak), args.Error(1)
}

func (m *MockRetentionService) ResetCounter(ctx context.Context, userID int64, req service.ResetRequest) (*domain.RetentionStreak, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RetentionStreak), args.Error(1)
}

func (m *MockRetentionService) GetStats(ctx context.Context, userID int64) (*domain.RetentionStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RetentionStats), args.Error(1)
}

// MockDashboardService mocks service.DashboardServiceInterface.
type MockDashboardService struct{ mock.Mock }

func (m *MockDashboardService) GetDashboard(ctx context.Context, userID int64) (*domain.DashboardData, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DashboardData), args.Error(1)
}
