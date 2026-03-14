package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MeditationRepo implements repository.MeditationRepository using PostgreSQL.
type MeditationRepo struct {
	pool *pgxpool.Pool
}

// NewMeditationRepo creates a new MeditationRepo.
func NewMeditationRepo(pool *pgxpool.Pool) *MeditationRepo {
	return &MeditationRepo{pool: pool}
}

// Create inserts a new meditation session.
func (r *MeditationRepo) Create(ctx context.Context, session *domain.MeditationSession) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO meditation_sessions (user_id, target_minutes, start_time, mood_before)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		session.UserID, session.TargetMinutes, session.StartTime, session.MoodBefore,
	).Scan(&session.ID, &session.CreatedAt)
	if err != nil {
		return fmt.Errorf("creating meditation session: %w", err)
	}
	return nil
}

// GetActive returns the currently active meditation session for a user, or nil if none.
func (r *MeditationRepo) GetActive(ctx context.Context, userID int64) (*domain.MeditationSession, error) {
	var s domain.MeditationSession
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, target_minutes, start_time, end_time, actual_duration_minutes, mood_before, mood_after, created_at
		 FROM meditation_sessions
		 WHERE user_id = $1 AND end_time IS NULL
		 ORDER BY start_time DESC LIMIT 1`, userID,
	).Scan(&s.ID, &s.UserID, &s.TargetMinutes, &s.StartTime, &s.EndTime,
		&s.ActualDuration, &s.MoodBefore, &s.MoodAfter, &s.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting active meditation session: %w", err)
	}
	return &s, nil
}

// EndSession marks a meditation session as complete.
func (r *MeditationRepo) EndSession(ctx context.Context, id int64, endTime time.Time, actualDuration float64, moodAfter *int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE meditation_sessions
		 SET end_time = $2, actual_duration_minutes = $3, mood_after = $4
		 WHERE id = $1`,
		id, endTime, actualDuration, moodAfter)
	if err != nil {
		return fmt.Errorf("ending meditation session: %w", err)
	}
	return nil
}

// GetCompletedByDateRange returns completed meditation sessions within a date range.
func (r *MeditationRepo) GetCompletedByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.MeditationSession, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, target_minutes, start_time, end_time, actual_duration_minutes, mood_before, mood_after, created_at
		 FROM meditation_sessions
		 WHERE user_id = $1 AND end_time IS NOT NULL AND start_time >= $2 AND start_time < $3
		 ORDER BY start_time DESC`, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("getting completed meditation sessions: %w", err)
	}
	defer rows.Close()

	var sessions []domain.MeditationSession
	for rows.Next() {
		var s domain.MeditationSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.TargetMinutes, &s.StartTime, &s.EndTime,
			&s.ActualDuration, &s.MoodBefore, &s.MoodAfter, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning meditation session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

// TotalMinutes returns the total minutes meditated for a user.
func (r *MeditationRepo) TotalMinutes(ctx context.Context, userID int64) (float64, error) {
	var total *float64
	err := r.pool.QueryRow(ctx,
		`SELECT SUM(actual_duration_minutes) FROM meditation_sessions WHERE user_id = $1 AND end_time IS NOT NULL`, userID,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("totaling meditation minutes: %w", err)
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

// CountCompleted returns the total number of completed meditation sessions.
func (r *MeditationRepo) CountCompleted(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM meditation_sessions WHERE user_id = $1 AND end_time IS NOT NULL`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting completed meditations: %w", err)
	}
	return count, nil
}

// AverageDuration returns the average session length in minutes.
func (r *MeditationRepo) AverageDuration(ctx context.Context, userID int64) (float64, error) {
	var avg *float64
	err := r.pool.QueryRow(ctx,
		`SELECT AVG(actual_duration_minutes) FROM meditation_sessions WHERE user_id = $1 AND end_time IS NOT NULL`, userID,
	).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("averaging meditation duration: %w", err)
	}
	if avg == nil {
		return 0, nil
	}
	return *avg, nil
}

// HasCompletedOnDate checks if a completed meditation exists on a given date.
func (r *MeditationRepo) HasCompletedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM meditation_sessions
			WHERE user_id = $1 AND end_time IS NOT NULL AND end_time >= $2 AND end_time < $3
		)`, userID, dayStart, dayEnd,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking meditation on date: %w", err)
	}
	return exists, nil
}
