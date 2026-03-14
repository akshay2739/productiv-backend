package domain

import "time"

// RetentionMilestones defines the milestone day markers.
var RetentionMilestones = []int{7, 14, 30, 60, 90, 180, 365}

// RetentionStreak represents a retention tracking period.
type RetentionStreak struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	DaysCount int        `json:"days_count"`
	Reason    *string    `json:"reason,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// IsActive returns true if the retention streak is still ongoing.
func (r *RetentionStreak) IsActive() bool {
	return r.EndDate == nil
}

// Milestone represents a milestone in the retention journey.
type Milestone struct {
	Days     int  `json:"days"`
	Achieved bool `json:"achieved"`
}

// RetentionStats holds aggregated retention statistics.
type RetentionStats struct {
	CurrentDayCount int              `json:"current_day_count"`
	BestStreak      int              `json:"best_streak"`
	NextMilestone   *int             `json:"next_milestone,omitempty"`
	Milestones      []Milestone      `json:"milestones"`
	ActiveStreak    *RetentionStreak `json:"active_streak,omitempty"`
	PastStreaks     []RetentionStreak `json:"past_streaks"`
}
