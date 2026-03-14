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

func newMeditationService() (*MeditationService, *repomocks.MockMeditationRepository, *repomocks.MockUserRepository) {
	medRepo := new(repomocks.MockMeditationRepository)
	userRepo := new(repomocks.MockUserRepository)
	return NewMeditationService(medRepo, userRepo), medRepo, userRepo
}

func TestMeditationService_StartSession_Success(t *testing.T) {
	svc, medRepo, _ := newMeditationService()

	medRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	medRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.MeditationSession")).Return(nil)

	mood := 3
	session, err := svc.StartSession(context.Background(), 1, StartSessionRequest{
		TargetMinutes: 10,
		MoodBefore:    &mood,
	})

	require.NoError(t, err)
	assert.Equal(t, 10, session.TargetMinutes)
	assert.Equal(t, &mood, session.MoodBefore)
	assert.Nil(t, session.EndTime)
	medRepo.AssertExpectations(t)
}

func TestMeditationService_StartSession_InvalidDuration(t *testing.T) {
	svc, _, _ := newMeditationService()

	_, err := svc.StartSession(context.Background(), 1, StartSessionRequest{TargetMinutes: 7})

	assert.ErrorIs(t, err, domain.ErrInvalidDuration)
}

func TestMeditationService_StartSession_ValidDurations(t *testing.T) {
	durations := []int{5, 10, 15, 20}
	for _, d := range durations {
		t.Run(string(rune('0'+d/5)), func(t *testing.T) {
			svc, medRepo, _ := newMeditationService()
			medRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
			medRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.MeditationSession")).Return(nil)

			session, err := svc.StartSession(context.Background(), 1, StartSessionRequest{TargetMinutes: d})
			require.NoError(t, err)
			assert.Equal(t, d, session.TargetMinutes)
		})
	}
}

func TestMeditationService_StartSession_InvalidMoodBefore(t *testing.T) {
	svc, _, _ := newMeditationService()

	mood := 6
	_, err := svc.StartSession(context.Background(), 1, StartSessionRequest{
		TargetMinutes: 10,
		MoodBefore:    &mood,
	})

	assert.ErrorIs(t, err, domain.ErrInvalidMoodValue)
}

func TestMeditationService_StartSession_InvalidMoodBeforeZero(t *testing.T) {
	svc, _, _ := newMeditationService()

	mood := 0
	_, err := svc.StartSession(context.Background(), 1, StartSessionRequest{
		TargetMinutes: 10,
		MoodBefore:    &mood,
	})

	assert.ErrorIs(t, err, domain.ErrInvalidMoodValue)
}

func TestMeditationService_StartSession_ActiveSessionExists(t *testing.T) {
	svc, medRepo, _ := newMeditationService()

	existing := &domain.MeditationSession{ID: 1, UserID: 1}
	medRepo.On("GetActive", mock.Anything, int64(1)).Return(existing, nil)

	_, err := svc.StartSession(context.Background(), 1, StartSessionRequest{TargetMinutes: 10})

	assert.ErrorIs(t, err, domain.ErrActiveSessionExists)
	medRepo.AssertExpectations(t)
}

func TestMeditationService_EndSession_Success(t *testing.T) {
	svc, medRepo, _ := newMeditationService()

	active := &domain.MeditationSession{
		ID:            1,
		UserID:        1,
		TargetMinutes: 10,
		StartTime:     time.Now().UTC().Add(-12 * time.Minute),
	}
	medRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	medRepo.On("EndSession", mock.Anything, int64(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("float64"), mock.Anything).Return(nil)

	mood := 5
	session, err := svc.EndSession(context.Background(), 1, EndSessionRequest{MoodAfter: &mood})

	require.NoError(t, err)
	assert.NotNil(t, session.EndTime)
	assert.NotNil(t, session.ActualDuration)
	assert.Equal(t, &mood, session.MoodAfter)
	medRepo.AssertExpectations(t)
}

func TestMeditationService_EndSession_NoActiveSession(t *testing.T) {
	svc, medRepo, _ := newMeditationService()

	medRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)

	_, err := svc.EndSession(context.Background(), 1, EndSessionRequest{})

	assert.ErrorIs(t, err, domain.ErrNoActiveSession)
}

func TestMeditationService_EndSession_InvalidMoodAfter(t *testing.T) {
	svc, _, _ := newMeditationService()

	mood := 7
	_, err := svc.EndSession(context.Background(), 1, EndSessionRequest{MoodAfter: &mood})

	assert.ErrorIs(t, err, domain.ErrInvalidMoodValue)
}

func TestMeditationService_GetStats_Success(t *testing.T) {
	svc, medRepo, userRepo := newMeditationService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	medRepo.On("GetActive", mock.Anything, int64(1)).Return(nil, nil)
	medRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	medRepo.On("TotalMinutes", mock.Anything, int64(1)).Return(150.5, nil)
	medRepo.On("CountCompleted", mock.Anything, int64(1)).Return(12, nil)
	medRepo.On("AverageDuration", mock.Anything, int64(1)).Return(12.54, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 150.5, stats.TotalMinutes)
	assert.Equal(t, 12, stats.TotalSessions)
	assert.Equal(t, 12.54, stats.AverageSession)
	assert.Len(t, stats.CalendarDays, 14)
	assert.Nil(t, stats.ActiveSession)
}

func TestMeditationService_GetStats_WithActiveSession(t *testing.T) {
	svc, medRepo, userRepo := newMeditationService()

	active := &domain.MeditationSession{ID: 1, TargetMinutes: 15, StartTime: time.Now().UTC()}
	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	medRepo.On("GetActive", mock.Anything, int64(1)).Return(active, nil)
	medRepo.On("HasCompletedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	medRepo.On("TotalMinutes", mock.Anything, int64(1)).Return(0.0, nil)
	medRepo.On("CountCompleted", mock.Anything, int64(1)).Return(0, nil)
	medRepo.On("AverageDuration", mock.Anything, int64(1)).Return(0.0, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.NotNil(t, stats.ActiveSession)
	assert.Equal(t, 15, stats.ActiveSession.TargetMinutes)
}
