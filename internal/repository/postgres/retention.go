package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RetentionRepo implements repository.RetentionRepository using PostgreSQL.
type RetentionRepo struct {
	pool *pgxpool.Pool
}

// NewRetentionRepo creates a new RetentionRepo.
func NewRetentionRepo(pool *pgxpool.Pool) *RetentionRepo {
	return &RetentionRepo{pool: pool}
}

// Create inserts a new retention streak.
func (r *RetentionRepo) Create(ctx context.Context, streak *domain.RetentionStreak) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO retention_streaks (user_id, start_date)
		 VALUES ($1, $2)
		 RETURNING id, created_at`,
		streak.UserID, streak.StartDate,
	).Scan(&streak.ID, &streak.CreatedAt)
	if err != nil {
		return fmt.Errorf("creating retention streak: %w", err)
	}
	return nil
}

// GetActive returns the currently active retention streak for a user, or nil if none.
func (r *RetentionRepo) GetActive(ctx context.Context, userID int64) (*domain.RetentionStreak, error) {
	var s domain.RetentionStreak
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, start_date, end_date, days_count, reason, created_at
		 FROM retention_streaks
		 WHERE user_id = $1 AND end_date IS NULL
		 ORDER BY start_date DESC LIMIT 1`, userID,
	).Scan(&s.ID, &s.UserID, &s.StartDate, &s.EndDate, &s.DaysCount, &s.Reason, &s.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting active retention streak: %w", err)
	}
	return &s, nil
}

// EndStreak marks a retention streak as ended.
func (r *RetentionRepo) EndStreak(ctx context.Context, id int64, endDate time.Time, daysCount int, reason *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE retention_streaks
		 SET end_date = $2, days_count = $3, reason = $4
		 WHERE id = $1`,
		id, endDate, daysCount, reason)
	if err != nil {
		return fmt.Errorf("ending retention streak: %w", err)
	}
	return nil
}

// GetBestStreak returns the highest days_count across all streaks for a user.
func (r *RetentionRepo) GetBestStreak(ctx context.Context, userID int64) (int, error) {
	var best *int
	err := r.pool.QueryRow(ctx,
		`SELECT MAX(days_count) FROM retention_streaks WHERE user_id = $1`, userID,
	).Scan(&best)
	if err != nil {
		return 0, fmt.Errorf("getting best retention streak: %w", err)
	}
	if best == nil {
		return 0, nil
	}
	return *best, nil
}

// ListPast returns past (ended) retention streaks, most recent first.
func (r *RetentionRepo) ListPast(ctx context.Context, userID int64, limit int) ([]domain.RetentionStreak, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, start_date, end_date, days_count, reason, created_at
		 FROM retention_streaks
		 WHERE user_id = $1 AND end_date IS NOT NULL
		 ORDER BY end_date DESC
		 LIMIT $2`, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("listing past retention streaks: %w", err)
	}
	defer rows.Close()

	var streaks []domain.RetentionStreak
	for rows.Next() {
		var s domain.RetentionStreak
		if err := rows.Scan(&s.ID, &s.UserID, &s.StartDate, &s.EndDate, &s.DaysCount,
			&s.Reason, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning retention streak: %w", err)
		}
		streaks = append(streaks, s)
	}
	return streaks, rows.Err()
}
