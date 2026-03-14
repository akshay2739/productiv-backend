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

func newRetentionService() (*RetentionService, *repomocks.MockRetentionRepository, *repomocks.MockUserRepository) {
	retRepo := new(repomocks.MockRetentionRepository)
	userRepo := new(repomocks.MockUserRepository)
	return NewRetentionService(retRepo, userRepo), retRepo, userRepo
}

func TestRetentionService_StartTracking_Success(t *testing.T) {
	svc, retRepo, _ := newRetentionService()

	retRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	retRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RetentionStreak")).Return(nil)

	streak, err := svc.StartTracking(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), streak.UserID)
	assert.Nil(t, streak.EndDate)
	retRepo.AssertExpectations(t)
}

func TestRetentionService_StartTracking_ActiveStreakExists(t *testing.T) {
	svc, retRepo, _ := newRetentionService()

	existing := &domain.RetentionStreak{ID: 1, UserID: 1}
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(existing, nil)

	_, err := svc.StartTracking(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrActiveStreakExists)
	retRepo.AssertExpectations(t)
}

func TestRetentionService_ResetCounter_Success(t *testing.T) {
	svc, retRepo, userRepo := newRetentionService()

	active := &domain.RetentionStreak{
		ID:        1,
		UserID:    1,
		StartDate: time.Now().UTC().AddDate(0, 0, -10),
	}
	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	retRepo.On("EndStreak", mock.Anything, int64(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("int"), mock.Anything).Return(nil)

	reason := "testing"
	streak, err := svc.ResetCounter(context.Background(), 1, ResetRequest{Reason: &reason})

	require.NoError(t, err)
	assert.NotNil(t, streak.EndDate)
	assert.Equal(t, &reason, streak.Reason)
	assert.GreaterOrEqual(t, streak.DaysCount, 9) // at least 9 days (timezone may shift by 1)
	retRepo.AssertExpectations(t)
}

func TestRetentionService_ResetCounter_NoActiveStreak(t *testing.T) {
	svc, retRepo, _ := newRetentionService()

	retRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)

	_, err := svc.ResetCounter(context.Background(), 1, ResetRequest{})

	assert.ErrorIs(t, err, domain.ErrNoActiveStreak)
}

func TestRetentionService_ResetCounter_WithoutReason(t *testing.T) {
	svc, retRepo, userRepo := newRetentionService()

	active := &domain.RetentionStreak{
		ID:        1,
		UserID:    1,
		StartDate: time.Now().UTC().AddDate(0, 0, -5),
	}
	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	retRepo.On("EndStreak", mock.Anything, int64(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("int"), mock.Anything).Return(nil)

	streak, err := svc.ResetCounter(context.Background(), 1, ResetRequest{})

	require.NoError(t, err)
	assert.Nil(t, streak.Reason)
}

func TestRetentionService_GetStats_NoActiveStreak(t *testing.T) {
	svc, retRepo, userRepo := newRetentionService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	retRepo.On("GetBestStreak", mock.Anything, int64(1)).Return(30, nil)
	retRepo.On("ListPast", mock.Anything, int64(1), 10).Return([]domain.RetentionStreak{}, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 0, stats.CurrentDayCount)
	assert.Equal(t, 30, stats.BestStreak)
	assert.Nil(t, stats.ActiveStreak)
	assert.Len(t, stats.Milestones, 7)
}

func TestRetentionService_GetStats_WithActiveStreak(t *testing.T) {
	svc, retRepo, userRepo := newRetentionService()

	active := &domain.RetentionStreak{
		ID:        1,
		UserID:    1,
		StartDate: time.Now().UTC().AddDate(0, 0, -15),
	}
	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	retRepo.On("GetBestStreak", mock.Anything, int64(1)).Return(10, nil)
	retRepo.On("ListPast", mock.Anything, int64(1), 10).Return([]domain.RetentionStreak{}, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.CurrentDayCount, 14)
	assert.GreaterOrEqual(t, stats.BestStreak, stats.CurrentDayCount) // current > stored best, so best = current
	assert.NotNil(t, stats.ActiveStreak)

	// Check milestones: 7 and 14 should be achieved
	assert.True(t, stats.Milestones[0].Achieved)  // 7 days
	assert.True(t, stats.Milestones[1].Achieved)  // 14 days
	assert.False(t, stats.Milestones[2].Achieved) // 30 days

	// Next milestone should be 30
	require.NotNil(t, stats.NextMilestone)
	assert.Equal(t, 30, *stats.NextMilestone)
}

func TestRetentionService_GetStats_PastStreaks(t *testing.T) {
	svc, retRepo, userRepo := newRetentionService()

	past := []domain.RetentionStreak{
		{ID: 1, DaysCount: 20},
		{ID: 2, DaysCount: 10},
	}
	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	retRepo.On("GetBestStreak", mock.Anything, int64(1)).Return(20, nil)
	retRepo.On("ListPast", mock.Anything, int64(1), 10).Return(past, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, stats.PastStreaks, 2)
	assert.Equal(t, 20, stats.PastStreaks[0].DaysCount)
}
