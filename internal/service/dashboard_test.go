package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	repomocks "github.com/akshay/productiv-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newDashboardService() (
	*DashboardService,
	*repomocks.MockPillarRepository,
	*repomocks.MockFastingRepository,
	*repomocks.MockGymRepository,
	*repomocks.MockMeditationRepository,
	*repomocks.MockRetentionRepository,
	*repomocks.MockUserRepository,
) {
	pillarRepo := new(repomocks.MockPillarRepository)
	fastingRepo := new(repomocks.MockFastingRepository)
	gymRepo := new(repomocks.MockGymRepository)
	medRepo := new(repomocks.MockMeditationRepository)
	retRepo := new(repomocks.MockRetentionRepository)
	userRepo := new(repomocks.MockUserRepository)
	svc := NewDashboardService(pillarRepo, fastingRepo, gymRepo, medRepo, retRepo, userRepo)
	return svc, pillarRepo, fastingRepo, gymRepo, medRepo, retRepo, userRepo
}

// todayMatcher returns a MatchedBy that matches today's date in Asia/Kolkata timezone.
func todayMatcher() interface{} {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return mock.MatchedBy(func(t time.Time) bool {
		tDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return tDate.Equal(today)
	})
}

// notTodayMatcher returns a MatchedBy that matches any date that is NOT today.
func notTodayMatcher() interface{} {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return mock.MatchedBy(func(t time.Time) bool {
		tDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return !tDate.Equal(today)
	})
}

func TestDashboardService_GetDashboard_AllActive(t *testing.T) {
	svc, pillarRepo, fastingRepo, gymRepo, medRepo, retRepo, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)

	pillars := []domain.Pillar{
		{Type: domain.PillarFasting, Name: "Fasting", Icon: "🍽️", Color: "#e94560"},
		{Type: domain.PillarGym, Name: "Gym", Icon: "💪", Color: "#4a9eff"},
		{Type: domain.PillarMeditation, Name: "Meditation", Icon: "🧘", Color: "#9b59b6"},
		{Type: domain.PillarRetention, Name: "Retention", Icon: "🔥", Color: "#2ecc71"},
	}
	pillarRepo.On("ListByUser", mock.Anything, int64(1)).Return(pillars, nil)

	// Today has activity, other days do not (keeps streak at 1, fast lookback)
	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), todayMatcher()).Return(true, nil)
	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), notTodayMatcher()).Return(false, nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), todayMatcher()).Return(true, nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), notTodayMatcher()).Return(false, nil)
	medRepo.On("HasCompletedOnDate", mock.Anything, int64(1), todayMatcher()).Return(true, nil)
	medRepo.On("HasCompletedOnDate", mock.Anything, int64(1), notTodayMatcher()).Return(false, nil)

	activeStreak := &domain.RetentionStreak{
		ID:        1,
		StartDate: time.Now().UTC().AddDate(0, 0, -5),
	}
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(activeStreak, nil)

	data, err := svc.GetDashboard(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 4, len(data.Pillars))
	assert.Equal(t, "All systems active. Stay the course.", data.TodaysFocus)
	for _, p := range data.Pillars {
		assert.True(t, p.HasActivityToday, "pillar %s should have activity today", p.Name)
	}
}

func TestDashboardService_GetDashboard_NothingDone(t *testing.T) {
	svc, pillarRepo, fastingRepo, gymRepo, medRepo, retRepo, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)

	pillars := []domain.Pillar{
		{Type: domain.PillarFasting, Name: "Fasting", Icon: "🍽️", Color: "#e94560"},
		{Type: domain.PillarGym, Name: "Gym", Icon: "💪", Color: "#4a9eff"},
		{Type: domain.PillarMeditation, Name: "Meditation", Icon: "🧘", Color: "#9b59b6"},
		{Type: domain.PillarRetention, Name: "Retention", Icon: "🔥", Color: "#2ecc71"},
	}
	pillarRepo.On("ListByUser", mock.Anything, int64(1)).Return(pillars, nil)

	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	medRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)

	data, err := svc.GetDashboard(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 0, data.DisciplineScore)
	assert.Equal(t, "Begin your journey. Pick any pillar and start.", data.TodaysFocus)
}

func TestDashboardService_GetDashboard_PartialActivity(t *testing.T) {
	svc, pillarRepo, fastingRepo, gymRepo, medRepo, retRepo, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)

	pillars := []domain.Pillar{
		{Type: domain.PillarFasting, Name: "Fasting", Icon: "🍽️", Color: "#e94560"},
		{Type: domain.PillarGym, Name: "Gym", Icon: "💪", Color: "#4a9eff"},
		{Type: domain.PillarMeditation, Name: "Meditation", Icon: "🧘", Color: "#9b59b6"},
		{Type: domain.PillarRetention, Name: "Retention", Icon: "🔥", Color: "#2ecc71"},
	}
	pillarRepo.On("ListByUser", mock.Anything, int64(1)).Return(pillars, nil)

	// Fasting done today, others not
	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), todayMatcher()).Return(true, nil)
	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), notTodayMatcher()).Return(false, nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	medRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	retRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)

	data, err := svc.GetDashboard(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, "Hit the gym today.", data.TodaysFocus)
}

func TestDashboardService_GetDashboard_DisciplineScore(t *testing.T) {
	svc, pillarRepo, fastingRepo, gymRepo, _, _, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)

	pillars := []domain.Pillar{
		{Type: domain.PillarFasting, Name: "Fasting", Icon: "🍽️", Color: "#e94560"},
		{Type: domain.PillarGym, Name: "Gym", Icon: "💪", Color: "#4a9eff"},
	}
	pillarRepo.On("ListByUser", mock.Anything, int64(1)).Return(pillars, nil)

	fastingRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)

	data, err := svc.GetDashboard(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 0, data.DisciplineScore)
	assert.Len(t, data.Pillars, 2)
}

func TestDashboardService_GetDashboard_EmptyPillars(t *testing.T) {
	svc, pillarRepo, _, _, _, _, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	pillarRepo.On("ListByUser", mock.Anything, int64(1)).Return([]domain.Pillar{}, nil)

	data, err := svc.GetDashboard(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 0, data.DisciplineScore)
	assert.Empty(t, data.Pillars)
}

func TestDashboardService_GetDashboard_UserRepoError(t *testing.T) {
	svc, _, _, _, _, _, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(nil, errors.New("db error"))

	_, err := svc.GetDashboard(context.Background(), 1)

	assert.Error(t, err)
}

func TestDashboardService_GetDashboard_PillarRepoError(t *testing.T) {
	svc, pillarRepo, _, _, _, _, userRepo := newDashboardService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	pillarRepo.On("ListByUser", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetDashboard(context.Background(), 1)

	assert.Error(t, err)
}
