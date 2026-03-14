package service

import (
	"context"
	"math"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/repository"
)

// RetentionService handles retention business logic.
type RetentionService struct {
	repo     repository.RetentionRepository
	userRepo repository.UserRepository
}

// NewRetentionService creates a new RetentionService.
func NewRetentionService(repo repository.RetentionRepository, userRepo repository.UserRepository) *RetentionService {
	return &RetentionService{repo: repo, userRepo: userRepo}
}

// StartTracking begins a new retention streak.
func (s *RetentionService) StartTracking(ctx context.Context, userID int64) (*domain.RetentionStreak, error) {
	active, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, domain.ErrActiveStreakExists
	}

	streak := &domain.RetentionStreak{
		UserID:    userID,
		StartDate: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, streak); err != nil {
		return nil, err
	}
	return streak, nil
}

// ResetRequest holds data for resetting the retention counter.
type ResetRequest struct {
	Reason *string `json:"reason,omitempty"`
}

// ResetCounter ends the active streak and starts fresh.
func (s *RetentionService) ResetCounter(ctx context.Context, userID int64, req ResetRequest) (*domain.RetentionStreak, error) {
	active, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return nil, domain.ErrNoActiveStreak
	}

	user, err := s.userRepo.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	loc := getUserTimezone(user.Timezone)

	now := time.Now().In(loc)
	daysCount := computeDayCount(active.StartDate, now, loc)

	if err := s.repo.EndStreak(ctx, active.ID, time.Now().UTC(), daysCount, req.Reason); err != nil {
		return nil, err
	}

	active.EndDate = ptrTime(time.Now().UTC())
	active.DaysCount = daysCount
	active.Reason = req.Reason
	return active, nil
}

// GetStats returns retention statistics.
func (s *RetentionService) GetStats(ctx context.Context, userID int64) (*domain.RetentionStats, error) {
	user, err := s.userRepo.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	loc := getUserTimezone(user.Timezone)

	activeStreak, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}

	currentDayCount := 0
	if activeStreak != nil {
		now := time.Now().In(loc)
		currentDayCount = computeDayCount(activeStreak.StartDate, now, loc)
		activeStreak.DaysCount = currentDayCount
	}

	bestStreak, err := s.repo.GetBestStreak(ctx, userID)
	if err != nil {
		return nil, err
	}
	if currentDayCount > bestStreak {
		bestStreak = currentDayCount
	}

	milestones := buildMilestones(currentDayCount)
	nextMilestone := findNextMilestone(currentDayCount)

	pastStreaks, err := s.repo.ListPast(ctx, userID, 10)
	if err != nil {
		return nil, err
	}

	return &domain.RetentionStats{
		CurrentDayCount: currentDayCount,
		BestStreak:      bestStreak,
		NextMilestone:   nextMilestone,
		Milestones:      milestones,
		ActiveStreak:    activeStreak,
		PastStreaks:      pastStreaks,
	}, nil
}

func computeDayCount(startDate time.Time, now time.Time, loc *time.Location) int {
	start := startDate.In(loc)
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	days := int(math.Floor(nowDay.Sub(startDay).Hours() / 24))
	if days < 0 {
		return 0
	}
	return days
}

func buildMilestones(currentDays int) []domain.Milestone {
	milestones := make([]domain.Milestone, len(domain.RetentionMilestones))
	for i, m := range domain.RetentionMilestones {
		milestones[i] = domain.Milestone{
			Days:     m,
			Achieved: currentDays >= m,
		}
	}
	return milestones
}

func findNextMilestone(currentDays int) *int {
	for _, m := range domain.RetentionMilestones {
		if currentDays < m {
			return &m
		}
	}
	return nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
