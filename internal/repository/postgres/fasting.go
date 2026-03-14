package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FastingRepo implements repository.FastingRepository using PostgreSQL.
type FastingRepo struct {
	pool *pgxpool.Pool
}

// NewFastingRepo creates a new FastingRepo.
func NewFastingRepo(pool *pgxpool.Pool) *FastingRepo {
	return &FastingRepo{pool: pool}
}

// Create inserts a new fasting session.
func (r *FastingRepo) Create(ctx context.Context, session *domain.FastingSession) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO fasting_sessions (user_id, protocol, target_hours, start_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		session.UserID, session.Protocol, session.TargetHours, session.StartTime,
	).Scan(&session.ID, &session.CreatedAt)
	if err != nil {
		return fmt.Errorf("creating fasting session: %w", err)
	}
	return nil
}

// GetActive returns the currently active fasting session for a user, or nil if none.
func (r *FastingRepo) GetActive(ctx context.Context, userID int64) (*domain.FastingSession, error) {
	var s domain.FastingSession
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, protocol, target_hours, start_time, end_time, actual_duration_hours, target_reached, created_at
		 FROM fasting_sessions
		 WHERE user_id = $1 AND end_time IS NULL
		 ORDER BY start_time DESC LIMIT 1`, userID,
	).Scan(&s.ID, &s.UserID, &s.Protocol, &s.TargetHours, &s.StartTime,
		&s.EndTime, &s.ActualDuration, &s.TargetReached, &s.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting active fasting session: %w", err)
	}
	return &s, nil
}

// EndSession marks a fasting session as complete.
func (r *FastingRepo) EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, targetReached bool) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE fasting_sessions
		 SET end_time = $2, actual_duration_hours = $3, target_reached = $4
		 WHERE id = $1`,
		id, endTime, actualDuration, targetReached)
	if err != nil {
		return fmt.Errorf("ending fasting session: %w", err)
	}
	return nil
}

// GetCompletedByDateRange returns completed fasting sessions within a date range.
func (r *FastingRepo) GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.FastingSession, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, protocol, target_hours, start_time, end_time, actual_duration_hours, target_reached, created_at
		 FROM fasting_sessions
		 WHERE user_id = $1 AND end_time IS NOT NULL AND start_time >= $2 AND start_time < $3
		 ORDER BY start_time DESC`, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("getting completed fasting sessions: %w", err)
	}
	defer rows.Close()

	var sessions []domain.FastingSession
	for rows.Next() {
		var s domain.FastingSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.Protocol, &s.TargetHours, &s.StartTime,
			&s.EndTime, &s.ActualDuration, &s.TargetReached, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning fasting session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// CountCompleted returns the total number of completed fasts for a user.
func (r *FastingRepo) CountCompleted(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM fasting_sessions WHERE user_id = $1 AND end_time IS NOT NULL`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting completed fasts: %w", err)
	}
	return count, nil
}

// AverageDuration returns the average duration of completed fasts in hours.
func (r *FastingRepo) AverageDuration(ctx context.Context, userID int64) (float64, error) {
	var avg *float64
	err := r.pool.QueryRow(ctx,
		`SELECT AVG(actual_duration_hours) FROM fasting_sessions WHERE user_id = $1 AND end_time IS NOT NULL`, userID,
	).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("averaging fasting duration: %w", err)
	}
	if avg == nil {
		return 0, nil
	}
	return *avg, nil
}

// HasCompletedOnDate checks if a completed fast exists on a given date in the user's timezone.
func (r *FastingRepo) HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM fasting_sessions
			WHERE user_id = $1 AND end_time IS NOT NULL AND end_time >= $2 AND end_time < $3
		)`, userID, dayStart, dayEnd,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking fasting on date: %w", err)
	}
	return exists, nil
}
