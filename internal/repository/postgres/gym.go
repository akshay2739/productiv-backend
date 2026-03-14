package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GymRepo implements repository.GymRepository using PostgreSQL.
type GymRepo struct {
	pool *pgxpool.Pool
}

// NewGymRepo creates a new GymRepo.
func NewGymRepo(pool *pgxpool.Pool) *GymRepo {
	return &GymRepo{pool: pool}
}

// Create inserts a new gym session.
func (r *GymRepo) Create(ctx context.Context, session *domain.GymSession) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO gym_sessions (user_id, workout_type, duration_min, energy_level, logged_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		session.UserID, session.WorkoutType, session.DurationMin, session.EnergyLevel, session.LoggedAt,
	).Scan(&session.ID, &session.CreatedAt)
	if err != nil {
		return fmt.Errorf("creating gym session: %w", err)
	}
	return nil
}

// HasLoggedOnDate checks if a workout has already been logged on the given date.
func (r *GymRepo) HasLoggedOnDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM gym_sessions WHERE user_id = $1 AND logged_at >= $2 AND logged_at < $3
		)`, userID, dayStart, dayEnd,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking gym log on date: %w", err)
	}
	return exists, nil
}

// CountByWeek returns the number of workouts in the week starting from weekStart.
func (r *GymRepo) CountByWeek(ctx context.Context, userID int64, weekStart time.Time) (int, error) {
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM gym_sessions WHERE user_id = $1 AND logged_at >= $2 AND logged_at < $3`,
		userID, weekStart, weekEnd,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting weekly workouts: %w", err)
	}
	return count, nil
}

// CountByMonth returns the number of workouts in the given month.
func (r *GymRepo) CountByMonth(ctx context.Context, userID int64, year int, month time.Month) (int, error) {
	loc := time.UTC
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	monthEnd := monthStart.AddDate(0, 1, 0)

	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM gym_sessions WHERE user_id = $1 AND logged_at >= $2 AND logged_at < $3`,
		userID, monthStart, monthEnd,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting monthly workouts: %w", err)
	}
	return count, nil
}

// GetByDateRange returns gym sessions within a date range.
func (r *GymRepo) GetByDateRange(ctx context.Context, userID int64, start, end time.Time) ([]domain.GymSession, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, workout_type, duration_min, energy_level, logged_at, created_at
		 FROM gym_sessions
		 WHERE user_id = $1 AND logged_at >= $2 AND logged_at < $3
		 ORDER BY logged_at DESC`, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("getting gym sessions by date range: %w", err)
	}
	defer rows.Close()

	var sessions []domain.GymSession
	for rows.Next() {
		var s domain.GymSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.WorkoutType, &s.DurationMin,
			&s.EnergyLevel, &s.LoggedAt, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning gym session: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}
