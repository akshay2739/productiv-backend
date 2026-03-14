package service

import (
	"context"
	"errors"
	"testing"

	"github.com/akshay/productiv-backend/internal/domain"
	repomocks "github.com/akshay/productiv-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newGymService() (*GymService, *repomocks.MockGymRepository, *repomocks.MockUserRepository) {
	gymRepo := new(repomocks.MockGymRepository)
	userRepo := new(repomocks.MockUserRepository)
	return NewGymService(gymRepo, userRepo), gymRepo, userRepo
}

func TestGymService_LogWorkout_Success(t *testing.T) {
	svc, gymRepo, userRepo := newGymService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	gymRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.GymSession")).Return(nil)

	dur := 60
	energy := 4
	session, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{
		WorkoutType: "strength",
		DurationMin: &dur,
		EnergyLevel: &energy,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.WorkoutStrength, session.WorkoutType)
	assert.Equal(t, &dur, session.DurationMin)
	assert.Equal(t, &energy, session.EnergyLevel)
	gymRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestGymService_LogWorkout_InvalidType(t *testing.T) {
	svc, _, _ := newGymService()

	_, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{WorkoutType: "boxing"})

	assert.ErrorIs(t, err, domain.ErrInvalidWorkoutType)
}

func TestGymService_LogWorkout_InvalidEnergyLevel(t *testing.T) {
	svc, _, _ := newGymService()

	energy := 6
	_, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{
		WorkoutType: "strength",
		EnergyLevel: &energy,
	})

	assert.ErrorIs(t, err, domain.ErrInvalidEnergyLevel)
}

func TestGymService_LogWorkout_InvalidEnergyLevelZero(t *testing.T) {
	svc, _, _ := newGymService()

	energy := 0
	_, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{
		WorkoutType: "strength",
		EnergyLevel: &energy,
	})

	assert.ErrorIs(t, err, domain.ErrInvalidEnergyLevel)
}

func TestGymService_LogWorkout_AlreadyLoggedToday(t *testing.T) {
	svc, gymRepo, userRepo := newGymService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(true, nil)

	_, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{WorkoutType: "cardio"})

	assert.ErrorIs(t, err, domain.ErrAlreadyLoggedToday)
	gymRepo.AssertExpectations(t)
}

func TestGymService_LogWorkout_AllWorkoutTypes(t *testing.T) {
	types := []string{"strength", "cardio", "hiit", "yoga", "sports", "swimming"}

	for _, wt := range types {
		t.Run(wt, func(t *testing.T) {
			svc, gymRepo, userRepo := newGymService()
			userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
			gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
			gymRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.GymSession")).Return(nil)

			session, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{WorkoutType: wt})
			require.NoError(t, err)
			assert.Equal(t, domain.WorkoutType(wt), session.WorkoutType)
		})
	}
}

func TestGymService_LogWorkout_RepoError(t *testing.T) {
	svc, gymRepo, userRepo := newGymService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	gymRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.GymSession")).Return(errors.New("db connection lost"))

	_, err := svc.LogWorkout(context.Background(), 1, LogWorkoutRequest{WorkoutType: "strength"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db connection lost")
}

func TestGymService_GetStats_Success(t *testing.T) {
	svc, gymRepo, userRepo := newGymService()

	userRepo.On("GetDefault", mock.Anything).Return(defaultUser(), nil)
	gymRepo.On("HasLoggedOnDate", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(false, nil)
	gymRepo.On("CountByWeek", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(3, nil)
	gymRepo.On("CountByMonth", mock.Anything, int64(1), mock.AnythingOfType("int"), mock.AnythingOfType("time.Month")).Return(12, nil)

	stats, err := svc.GetStats(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 3, stats.WorkoutsWeek)
	assert.Equal(t, 12, stats.WorkoutsMonth)
	assert.Equal(t, 5, stats.WeeklyGoal)
	assert.Len(t, stats.CalendarDays, 14)
}

func TestGymService_GetStats_UserRepoError(t *testing.T) {
	svc, _, userRepo := newGymService()

	userRepo.On("GetDefault", mock.Anything).Return(nil, errors.New("user not found"))

	_, err := svc.GetStats(context.Background(), 1)

	assert.Error(t, err)
}
