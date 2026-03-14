package service

import (
	"context"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/repository"
)

// DashboardService aggregates data across all pillars.
type DashboardService struct {
	pillarRepo     repository.PillarRepository
	fastingRepo    repository.FastingRepository
	gymRepo        repository.GymRepository
	meditationRepo repository.MeditationRepository
	retentionRepo  repository.RetentionRepository
	userRepo       repository.UserRepository
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(
	pillarRepo repository.PillarRepository,
	fastingRepo repository.FastingRepository,
	gymRepo repository.GymRepository,
	meditationRepo repository.MeditationRepository,
	retentionRepo repository.RetentionRepository,
	userRepo repository.UserRepository,
) *DashboardService {
	return &DashboardService{
		pillarRepo:     pillarRepo,
		fastingRepo:    fastingRepo,
		gymRepo:        gymRepo,
		meditationRepo: meditationRepo,
		retentionRepo:  retentionRepo,
		userRepo:       userRepo,
	}
}

// GetDashboard returns the complete dashboard data.
func (s *DashboardService) GetDashboard(ctx context.Context, userID int64) (*domain.DashboardData, error) {
	user, err := s.userRepo.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	loc := getUserTimezone(user.Timezone)

	pillars, err := s.pillarRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	summaries := make([]domain.PillarSummary, 0, len(pillars))
	disciplineScore := 0
	activePillarsToday := 0

	for _, p := range pillars {
		summary, err := s.getPillarSummary(ctx, userID, p, now, loc)
		if err != nil {
			return nil, err
		}
		disciplineScore += summary.CurrentStreak
		if summary.HasActivityToday {
			activePillarsToday++
		}
		summaries = append(summaries, *summary)
	}

	focus := buildTodaysFocus(summaries, activePillarsToday, len(pillars))

	return &domain.DashboardData{
		DisciplineScore: disciplineScore,
		Pillars:         summaries,
		TodaysFocus:     focus,
	}, nil
}

func (s *DashboardService) getPillarSummary(ctx context.Context, userID int64, p domain.Pillar, now time.Time, loc *time.Location) (*domain.PillarSummary, error) {
	summary := &domain.PillarSummary{
		Type:  p.Type,
		Name:  p.Name,
		Icon:  p.Icon,
		Color: p.Color,
	}

	switch p.Type {
	case domain.PillarFasting:
		streak, err := calculateCurrentStreak(ctx, userID, loc, func(ctx context.Context, uid int64, date time.Time) (bool, error) {
			return s.fastingRepo.HasCompletedOnDate(ctx, uid, date)
		})
		if err != nil {
			return nil, err
		}
		hasToday, err := s.fastingRepo.HasCompletedOnDate(ctx, userID, now)
		if err != nil {
			return nil, err
		}
		summary.CurrentStreak = streak
		summary.HasActivityToday = hasToday

	case domain.PillarGym:
		streak, err := calculateCurrentStreak(ctx, userID, loc, func(ctx context.Context, uid int64, date time.Time) (bool, error) {
			return s.gymRepo.HasLoggedOnDate(ctx, uid, date)
		})
		if err != nil {
			return nil, err
		}
		hasToday, err := s.gymRepo.HasLoggedOnDate(ctx, userID, now)
		if err != nil {
			return nil, err
		}
		summary.CurrentStreak = streak
		summary.HasActivityToday = hasToday

	case domain.PillarMeditation:
		streak, err := calculateCurrentStreak(ctx, userID, loc, func(ctx context.Context, uid int64, date time.Time) (bool, error) {
			return s.meditationRepo.HasCompletedOnDate(ctx, uid, date)
		})
		if err != nil {
			return nil, err
		}
		hasToday, err := s.meditationRepo.HasCompletedOnDate(ctx, userID, now)
		if err != nil {
			return nil, err
		}
		summary.CurrentStreak = streak
		summary.HasActivityToday = hasToday

	case domain.PillarRetention:
		activeStreak, err := s.retentionRepo.GetActive(ctx, userID)
		if err != nil {
			return nil, err
		}
		if activeStreak != nil {
			summary.CurrentStreak = computeDayCount(activeStreak.StartDate, now, loc)
			summary.HasActivityToday = true // retention is always "active" when tracking
		}
	}

	return summary, nil
}

func buildTodaysFocus(summaries []domain.PillarSummary, activePillarsToday, totalPillars int) string {
	if activePillarsToday == 0 {
		return "Begin your journey. Pick any pillar and start."
	}
	if activePillarsToday >= totalPillars {
		return "All systems active. Stay the course."
	}

	// Find the first inactive pillar and suggest it
	for _, s := range summaries {
		if !s.HasActivityToday {
			switch s.Type {
			case domain.PillarFasting:
				return "Start a fast today."
			case domain.PillarGym:
				return "Hit the gym today."
			case domain.PillarMeditation:
				return "Take a moment to meditate."
			case domain.PillarRetention:
				return "Start tracking your retention."
			}
		}
	}
	return "Keep going. You're doing great."
}
