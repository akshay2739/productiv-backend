package service

import (
	"context"
	"math"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/repository"
)

// MeditationService handles meditation business logic.
type MeditationService struct {
	repo     repository.MeditationRepository
	userRepo repository.UserRepository
}

// NewMeditationService creates a new MeditationService.
func NewMeditationService(repo repository.MeditationRepository, userRepo repository.UserRepository) *MeditationService {
	return &MeditationService{repo: repo, userRepo: userRepo}
}

// StartSessionRequest holds data for starting a meditation session.
type StartSessionRequest struct {
	TargetMinutes int  `json:"target_minutes"`
	MoodBefore    *int `json:"mood_before,omitempty"`
}

// StartSession begins a new meditation session.
func (s *MeditationService) StartSession(ctx context.Context, userID int64, req StartSessionRequest) (*domain.MeditationSession, error) {
	if !isValidDuration(req.TargetMinutes) {
		return nil, domain.ErrInvalidDuration
	}

	if req.MoodBefore != nil && (*req.MoodBefore < 1 || *req.MoodBefore > 5) {
		return nil, domain.ErrInvalidMoodValue
	}

	active, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, domain.ErrActiveSessionExists
	}

	session := &domain.MeditationSession{
		UserID:        userID,
		TargetMinutes: req.TargetMinutes,
		StartTime:     time.Now().UTC(),
		MoodBefore:    req.MoodBefore,
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// EndSessionRequest holds data for ending a meditation session.
type EndSessionRequest struct {
	MoodAfter *int `json:"mood_after,omitempty"`
}

// EndSession completes the active meditation session.
func (s *MeditationService) EndSession(ctx context.Context, userID int64, req EndSessionRequest) (*domain.MeditationSession, error) {
	if req.MoodAfter != nil && (*req.MoodAfter < 1 || *req.MoodAfter > 5) {
		return nil, domain.ErrInvalidMoodValue
	}

	active, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return nil, domain.ErrNoActiveSession
	}

	endTime := time.Now().UTC()
	durationMinutes := endTime.Sub(active.StartTime).Minutes()
	roundedDuration := math.Round(durationMinutes*100) / 100

	if err := s.repo.EndSession(ctx, active.ID, endTime, roundedDuration, req.MoodAfter); err != nil {
		return nil, err
	}

	active.EndTime = &endTime
	active.ActualDuration = &roundedDuration
	active.MoodAfter = req.MoodAfter
	return active, nil
}

// GetStats returns meditation statistics.
func (s *MeditationService) GetStats(ctx context.Context, userID int64) (*domain.MeditationStats, error) {
	user, err := s.userRepo.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	loc := getUserTimezone(user.Timezone)

	activeSession, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}

	checker := func(ctx context.Context, uid int64, date time.Time) (bool, error) {
		return s.repo.HasCompletedOnDate(ctx, uid, date)
	}

	currentStreak, err := calculateCurrentStreak(ctx, userID, loc, checker)
	if err != nil {
		return nil, err
	}

	totalMinutes, err := s.repo.TotalMinutes(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalSessions, err := s.repo.CountCompleted(ctx, userID)
	if err != nil {
		return nil, err
	}

	avgDuration, err := s.repo.AverageDuration(ctx, userID)
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

	return &domain.MeditationStats{
		CurrentStreak:  currentStreak,
		LongestStreak:  currentStreak,
		TotalMinutes:   math.Round(totalMinutes*100) / 100,
		TotalSessions:  totalSessions,
		AverageSession: math.Round(avgDuration*100) / 100,
		CalendarDays:   domainDays,
		ActiveSession:  activeSession,
	}, nil
}

func isValidDuration(minutes int) bool {
	for _, d := range domain.MeditationDurationOptions() {
		if d == minutes {
			return true
		}
	}
	return false
}
