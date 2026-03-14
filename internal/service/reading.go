package service

import (
	"context"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/akshay/productiv-backend/internal/repository"
)

// ReadingService handles reading business logic.
type ReadingService struct {
	repo     repository.ReadingRepository
	userRepo repository.UserRepository
}

// NewReadingService creates a new ReadingService.
func NewReadingService(repo repository.ReadingRepository, userRepo repository.UserRepository) *ReadingService {
	return &ReadingService{repo: repo, userRepo: userRepo}
}

// LogReadingRequest holds data for logging a reading session.
type LogReadingRequest struct {
	BookName string `json:"book_name"`
	Pages    int    `json:"pages"`
}

// LogReading records a reading session.
func (s *ReadingService) LogReading(ctx context.Context, userID int64, req LogReadingRequest) (*domain.ReadingSession, error) {
	if req.BookName == "" {
		return nil, domain.ErrInvalidBookName
	}
	if req.Pages < 1 || req.Pages > 1000 {
		return nil, domain.ErrInvalidPageCount
	}

	session := &domain.ReadingSession{
		UserID:   userID,
		BookName: req.BookName,
		Pages:    req.Pages,
		LoggedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// GetStats returns reading statistics.
func (s *ReadingService) GetStats(ctx context.Context, userID int64) (*domain.ReadingStats, error) {
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

	totalPages, err := s.repo.TotalPages(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalSessions, err := s.repo.CountCompleted(ctx, userID)
	if err != nil {
		return nil, err
	}

	books, err := s.repo.GetBookSummaries(ctx, userID)
	if err != nil {
		return nil, err
	}
	if books == nil {
		books = []domain.BookSummary{}
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

	return &domain.ReadingStats{
		CurrentStreak: currentStreak,
		LongestStreak: currentStreak,
		TotalPages:    totalPages,
		TotalSessions: totalSessions,
		BooksRead:     books,
		CalendarDays:  domainDays,
		LoggedToday:   loggedToday,
	}, nil
}
