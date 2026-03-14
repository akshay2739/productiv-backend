package domain

import "time"

// MeditationDurationOptions returns the available duration goals in minutes.
func MeditationDurationOptions() []int {
	return []int{5, 10, 15, 20, 30, 45, 60}
}

// MeditationSession represents a single meditation session.
type MeditationSession struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	UserID         int64      `json:"user_id" gorm:"index"`
	TargetMinutes  int        `json:"target_minutes"`
	StartTime      time.Time  `json:"start_time" gorm:"index"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	ActualDuration *float64   `json:"actual_duration_minutes,omitempty" gorm:"column:actual_duration_minutes"`
	MoodBefore     *int       `json:"mood_before,omitempty"`
	MoodAfter      *int       `json:"mood_after,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// IsActive returns true if the meditation session is still ongoing.
func (m *MeditationSession) IsActive() bool {
	return m.EndTime == nil
}

// MeditationStats holds aggregated meditation statistics.
type MeditationStats struct {
	CurrentStreak  int                `json:"current_streak"`
	LongestStreak  int                `json:"longest_streak"`
	TotalMinutes   float64            `json:"total_minutes"`
	TotalSessions  int                `json:"total_sessions"`
	AverageSession float64            `json:"average_session_minutes"`
	CalendarDays   []CalendarDay      `json:"calendar_days"`
	ActiveSession  *MeditationSession `json:"active_session,omitempty"`
}
