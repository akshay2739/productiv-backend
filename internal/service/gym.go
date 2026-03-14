package service

import (
	"context"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/repository"
)

// GymService handles gym business logic.
type GymService struct {
	repo     repository.GymRepository
	userRepo repository.UserRepository
}

// NewGymService creates a new GymService.
func NewGymService(repo repository.GymRepository, userRepo repository.UserRepository) *GymService {
	return &GymService{repo: repo, userRepo: userRepo}
}

// LogWorkoutRequest holds data for logging a workout.
type LogWorkoutRequest struct {
	WorkoutType string `json:"workout_type"`
	DurationMin *int   `json:"duration_min,omitempty"`
	EnergyLevel *int   `json:"energy_level,omitempty"`
}

// LogWorkout records a gym workout for today.
func (s *GymService) LogWorkout(ctx context.Context, userID int64, req LogWorkoutRequest) (*domain.GymSession, error) {
	wt := domain.WorkoutType(req.WorkoutType)
	if !domain.IsValidWorkoutType(wt) {
		return nil, domain.ErrInvalidWorkoutType
	}

	if req.EnergyLevel != nil && (*req.EnergyLevel < 1 || *req.EnergyLevel > 5) {
		return nil, domain.ErrInvalidEnergyLevel
	}

	user, err := s.userRepo.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	loc := getUserTimezone(user.Timezone)
	now := time.Now().In(loc)

	logged, err := s.repo.HasLoggedOnDate(ctx, userID, now)
	if err != nil {
		return nil, err
	}
	if logged {
		return nil, domain.ErrAlreadyLoggedToday
	}

	session := &domain.GymSession{
		UserID:      userID,
		WorkoutType: wt,
		DurationMin: req.DurationMin,
		EnergyLevel: req.EnergyLevel,
		LoggedAt:    time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// GetStats returns gym statistics.
func (s *GymService) GetStats(ctx context.Context, userID int64) (*domain.GymStats, error) {
	user, err := s.userRepo.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	loc := getUserTimezone(user.Timezone)
	now := time.Now().In(loc)

	checker := func(ctx context.Context, uid int64, date time.Time) (bool, error) {
		return s.repo.HasLoggedOnDate(ctx, uid, date)
	}

	currentStreak, err := calculateCurrentStreak(ctx, userID, loc, checker)
	if err != nil {
		return nil, err
	}

	// Calculate week start (Monday)
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-int(weekday-time.Monday), 0, 0, 0, 0, loc)

	workoutsWeek, err := s.repo.CountByWeek(ctx, userID, weekStart)
	if err != nil {
		return nil, err
	}

	workoutsMonth, err := s.repo.CountByMonth(ctx, userID, now.Year(), now.Month())
	if err != nil {
		return nil, err
	}

	loggedToday, err := s.repo.HasLoggedOnDate(ctx, userID, now)
	if err != nil {
		return nil, err
	}

	calendarDays, err := buildCalendarDays(ctx, userID, loc, checker)
	if err != nil {
		return nil, err
	}

	domainDays := make([]domain.CalendarDay, len(calendarDays))
	for i, d := range calendarDays {
		domainDays[i] = domain.CalendarDay{
			Date:        d.Date,
			HasActivity: d.HasActivity,
			IsToday:     d.IsToday,
		}
	}

	return &domain.GymStats{
		CurrentStreak:  currentStreak,
		LongestStreak:  currentStreak,
		WorkoutsWeek:   workoutsWeek,
		WorkoutsMonth:  workoutsMonth,
		WeeklyGoal:     5,
		CalendarDays:   domainDays,
		LoggedToday:    loggedToday,
	}, nil
}
