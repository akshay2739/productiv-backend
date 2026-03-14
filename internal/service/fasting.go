package service

import (
	"context"
	"math"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/repository"
)

// FastingService handles fasting business logic.
type FastingService struct {
	repo     repository.FastingRepository
	userRepo repository.UserRepository
}

// NewFastingService creates a new FastingService.
func NewFastingService(repo repository.FastingRepository, userRepo repository.UserRepository) *FastingService {
	return &FastingService{repo: repo, userRepo: userRepo}
}

// StartFastRequest holds data for starting a fast.
type StartFastRequest struct {
	Protocol string `json:"protocol"`
}

// StartFast begins a new fasting session.
func (s *FastingService) StartFast(ctx context.Context, userID int64, req StartFastRequest) (*domain.FastingSession, error) {
	// Validate protocol
	protocol := findProtocol(req.Protocol)
	if protocol == nil {
		return nil, domain.ErrInvalidDuration
	}

	// Check for existing active session
	active, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, domain.ErrActiveFastExists
	}

	session := &domain.FastingSession{
		UserID:      userID,
		Protocol:    protocol.Name,
		TargetHours: protocol.FastingHours,
		StartTime:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// EndFast completes the active fasting session.
func (s *FastingService) EndFast(ctx context.Context, userID int64) (*domain.FastingSession, error) {
	active, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return nil, domain.ErrNoActiveFast
	}

	endTime := time.Now().UTC()
	durationHours := endTime.Sub(active.StartTime).Hours()
	roundedDuration := math.Round(durationHours*100) / 100
	targetReached := durationHours >= float64(active.TargetHours)

	if err := s.repo.EndSession(ctx, active.ID, endTime, roundedDuration, targetReached); err != nil {
		return nil, err
	}

	active.EndTime = &endTime
	active.ActualDuration = &roundedDuration
	active.TargetReached = targetReached
	return active, nil
}

// GetStats returns fasting statistics and the active session.
func (s *FastingService) GetStats(ctx context.Context, userID int64) (*domain.FastingStats, error) {
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

	totalFasts, err := s.repo.CountCompleted(ctx, userID)
	if err != nil {
		return nil, err
	}

	avgDuration, err := s.repo.AverageDuration(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build full month calendar with duration data
	now := time.Now().In(loc)
	durations, err := s.repo.GetCompletedDurationsForMonth(ctx, userID, now.Year(), now.Month(), loc)
	if err != nil {
		return nil, err
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	monthEnd := monthStart.AddDate(0, 1, 0)
	totalDays := int(monthEnd.Sub(monthStart).Hours() / 24)

	calendarDays := make([]domain.FastingCalendarDay, totalDays)
	for i := 0; i < totalDays; i++ {
		date := monthStart.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		dur, hasFast := durations[dateStr]
		day := domain.FastingCalendarDay{
			Date:        dateStr,
			HasActivity: hasFast,
			IsToday:     date.Equal(today),
		}
		if hasFast {
			rounded := math.Round(dur*10) / 10
			day.DurationHours = &rounded
		}
		calendarDays[i] = day
	}

	return &domain.FastingStats{
		CurrentStreak:   currentStreak,
		LongestStreak:   currentStreak,
		AverageDuration: math.Round(avgDuration*100) / 100,
		TotalFasts:      totalFasts,
		CalendarDays:    calendarDays,
		ActiveSession:   activeSession,
	}, nil
}

// GetProtocols returns available fasting protocols.
func (s *FastingService) GetProtocols() []domain.FastingProtocol {
	return domain.AvailableFastingProtocols()
}

func findProtocol(name string) *domain.FastingProtocol {
	for _, p := range domain.AvailableFastingProtocols() {
		if p.Name == name {
			return &p
		}
	}
	return nil
}
