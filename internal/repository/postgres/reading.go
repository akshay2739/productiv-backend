package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/gorm"
)

// ReadingRepo implements repository.ReadingRepository using GORM.
type ReadingRepo struct {
	db *gorm.DB
}

// NewReadingRepo creates a new ReadingRepo.
func NewReadingRepo(db *gorm.DB) *ReadingRepo {
	return &ReadingRepo{db: db}
}

// Create inserts a new reading session.
func (r *ReadingRepo) Create(ctx context.Context, session *domain.ReadingSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("creating reading session: %w", err)
	}
	return nil
}

// HasLoggedOnDate checks if a reading session exists on a given date.
func (r *ReadingRepo) HasLoggedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.ReadingSession{}).
		Where("user_id = ? AND logged_at >= ? AND logged_at < ?", userID, dayStart, dayEnd).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("checking reading on date: %w", err)
	}
	return count > 0, nil
}

// TotalPages returns the total number of pages read.
func (r *ReadingRepo) TotalPages(ctx context.Context, userID int64) (int, error) {
	var result struct{ Total *int64 }
	err := r.db.WithContext(ctx).
		Model(&domain.ReadingSession{}).
		Select("SUM(pages) as total").
		Where("user_id = ?", userID).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("totaling reading pages: %w", err)
	}
	if result.Total == nil {
		return 0, nil
	}
	return int(*result.Total), nil
}

// CountCompleted returns the total number of reading sessions.
func (r *ReadingRepo) CountCompleted(ctx context.Context, userID int64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.ReadingSession{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting reading sessions: %w", err)
	}
	return int(count), nil
}

// GetBookSummaries returns aggregated reading data per book.
func (r *ReadingRepo) GetBookSummaries(ctx context.Context, userID int64) ([]domain.BookSummary, error) {
	var summaries []domain.BookSummary
	err := r.db.WithContext(ctx).
		Model(&domain.ReadingSession{}).
		Select("book_name, SUM(pages) as total_pages, COUNT(*) as sessions").
		Where("user_id = ?", userID).
		Group("book_name").
		Order("MAX(logged_at) DESC").
		Scan(&summaries).Error
	if err != nil {
		return nil, fmt.Errorf("getting book summaries: %w", err)
	}
	return summaries, nil
}
