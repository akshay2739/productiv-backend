package service

import (
	"context"
	"testing"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	repomocks "github.com/akshay/productiv-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func defaultUser() *domain.User {
	return &domain.User{ID: 1, Name: "Test", Timezone: "Asia/Kolkata"}
}

func newFastingService() (*FastingService, *repomocks.MockFastingRepository, *repomocks.MockUserRepository) {
	fastingRepo := new(repomocks.MockFastingRepository)
	userRepo := new(repomocks.MockUserRepository)
	return NewFastingService(fastingRepo, userRepo), fastingRepo, userRepo
}

func TestFastingService_StartFast_Success(t *testing.T) {
	svc, fastingRepo, _ := newFastingService()

	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	fastingRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.FastingSession")).Return(nil)

	session, err := svc.StartFast(context.Background(), 1, StartFastRequest{Protocol: "16:8"})

	require.NoError(t, err)
	assert.Equal(t, "16:8", session.Protocol)
	assert.Equal(t, 16, session.TargetHours)
	assert.Equal(t, int64(1), session.UserID)
	assert.Nil(t, session.EndTime)
	fastingRepo.AssertExpectations(t)
}

func TestFastingService_StartFast_InvalidProtocol(t *testing.T) {
	svc, _, _ := newFastingService()

	_, err := svc.StartFast(context.Background(), 1, StartFastRequest{Protocol: "invalid"})

	assert.ErrorIs(t, err, domain.ErrInvalidDuration)
}

func TestFastingService_StartFast_ActiveFastExists(t *testing.T) {
	svc, fastingRepo, _ := newFastingService()

	existing := &domain.FastingSession{ID: 1, UserID: 1, Protocol: "16:8"}
	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(existing, nil)

	_, err := svc.StartFast(context.Background(), 1, StartFastRequest{Protocol: "16:8"})

	assert.ErrorIs(t, err, domain.ErrActiveFastExists)
	fastingRepo.AssertExpectations(t)
}

func TestFastingService_StartFast_AllProtocols(t *testing.T) {
	protocols := []struct {
		name  string
		hours int
	}{
		{"16:8", 16},
		{"18:6", 18},
		{"20:4", 20},
		{"OMAD", 23},
	}

	for _, p := range protocols {
		t.Run(p.name, func(t *testing.T) {
			svc, fastingRepo, _ := newFastingService()
			fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
			fastingRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.FastingSession")).Return(nil)

			session, err := svc.StartFast(context.Background(), 1, StartFastRequest{Protocol: p.name})

			require.NoError(t, err)
			assert.Equal(t, p.name, session.Protocol)
			assert.Equal(t, p.hours, session.TargetHours)
		})
	}
}

func TestFastingService_EndFast_Success(t *testing.T) {
	svc, fastingRepo, _ := newFastingService()

	active := &domain.FastingSession{
		ID:          1,
		UserID:      1,
		Protocol:    "16:8",
		TargetHours: 16,
		StartTime:   time.Now().UTC().Add(-17 * time.Hour),
	}
	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	fastingRepo.On("EndSession", mock.Anything, int64(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("float64"), true).Return(nil)

	session, err := svc.EndFast(context.Background(), 1)

	require.NoError(t, err)
	assert.NotNil(t, session.EndTime)
	assert.NotNil(t, session.ActualDuration)
	assert.True(t, session.TargetReached)
	fastingRepo.AssertExpectations(t)
}

func TestFastingService_EndFast_TargetNotReached(t *testing.T) {
	svc, fastingRepo, _ := newFastingService()

	active := &domain.FastingSession{
		ID:          1,
		UserID:      1,
		Protocol:    "16:8",
		TargetHours: 16,
		StartTime:   time.Now().UTC().Add(-2 * time.Hour),
	}
	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	fastingRepo.On("EndSession", mock.Anything, int64(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("float64"), false).Return(nil)

	session, err := svc.EndFast(context.Background(), 1)

	require.NoError(t, err)
	assert.False(t, session.TargetReached)
	fastingRepo.AssertExpectations(t)
}

func TestFastingService_EndFast_NoActiveFast(t *testing.T) {
	svc, fastingRepo, _ := newFastingService()

	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)

	_, err := svc.EndFast(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrNoActiveFast)
	fastingRepo.AssertExpectations(t)
}

func TestFastingService_GetStats_Success(t *testing.T) {
	svc, fastingRepo, userRepo := newFastingService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	fastingRepo.On("CountCompleted", mock.Anything, int64(1)).Return(5, nil)
	fastingRepo.On("AverageDuration", mock.Anything, int64(1)).Return(16.5, nil)
	fastingRepo.On("GetCompletedDurationsForMonth", mock.Anything, int64(1), mock.AnythingOfType("int"), mock.AnythingOfType("time.Month"), mock.Anything).Return(map[string]float64{}, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 5, stats.TotalFasts)
	assert.Equal(t, 16.5, stats.AverageDuration)
	assert.True(t, len(stats.CalendarDays) >= 28) // full month calendar
	assert.Nil(t, stats.ActiveSession)
	userRepo.AssertExpectations(t)
}

func TestFastingService_GetStats_WithActiveSession(t *testing.T) {
	svc, fastingRepo, userRepo := newFastingService()

	active := &domain.FastingSession{ID: 1, Protocol: "18:6", StartTime: time.Now().UTC()}
	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	fastingRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	fastingRepo.On("CountCompleted", mock.Anything, int64(1)).Return(0, nil)
	fastingRepo.On("AverageDuration", mock.Anything, int64(1)).Return(0.0, nil)
	fastingRepo.On("GetCompletedDurationsForMonth", mock.Anything, int64(1), mock.AnythingOfType("int"), mock.AnythingOfType("time.Month"), mock.Anything).Return(map[string]float64{}, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.NotNil(t, stats.ActiveSession)
	assert.Equal(t, "18:6", stats.ActiveSession.Protocol)
}

func TestFastingService_GetProtocols(t *testing.T) {
	svc, _, _ := newFastingService()

	protocols := svc.GetProtocols()

	assert.Len(t, protocols, 4)
	assert.Equal(t, "16:8", protocols[0].Name)
	assert.Equal(t, "OMAD", protocols[3].Name)
}
