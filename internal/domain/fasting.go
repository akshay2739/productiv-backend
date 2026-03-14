package domain

import "time"

// FastingProtocol defines a fasting protocol with its window.
type FastingProtocol struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	FastingHours int    `json:"fasting_hours"`
	EatingHours  int    `json:"eating_hours"`
}

// AvailableFastingProtocols returns the supported protocols.
func AvailableFastingProtocols() []FastingProtocol {
	return []FastingProtocol{
		{Name: "16:8", Label: "Most popular", FastingHours: 16, EatingHours: 8},
		{Name: "18:6", Label: "Fat burn", FastingHours: 18, EatingHours: 6},
		{Name: "20:4", Label: "Warrior", FastingHours: 20, EatingHours: 4},
		{Name: "OMAD", Label: "One meal", FastingHours: 23, EatingHours: 1},
	}
}

// FastingSession represents a single fasting session.
type FastingSession struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	UserID         int64      `json:"user_id" gorm:"index"`
	Protocol       string     `json:"protocol"`
	TargetHours    int        `json:"target_hours"`
	StartTime      time.Time  `json:"start_time" gorm:"index"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	ActualDuration *float64   `json:"actual_duration_hours,omitempty" gorm:"column:actual_duration_hours"`
	TargetReached  bool       `json:"target_reached" gorm:"default:false"`
	CreatedAt      time.Time  `json:"created_at"`
}

// IsActive returns true if the fasting session is still ongoing.
func (f *FastingSession) IsActive() bool {
	return f.EndTime == nil
}

// FastingStats holds aggregated fasting statistics.
type FastingStats struct {
	CurrentStreak   int                  `json:"current_streak"`
	LongestStreak   int                  `json:"longest_streak"`
	AverageDuration float64              `json:"average_duration_hours"`
	TotalFasts      int                  `json:"total_fasts"`
	CalendarDays    []FastingCalendarDay `json:"calendar_days"`
	ActiveSession   *FastingSession      `json:"active_session,omitempty"`
}

// CalendarDay represents a single day in the 14-day calendar view.
type CalendarDay struct {
	Date        string `json:"date"`
	HasActivity bool   `json:"has_activity"`
	IsToday     bool   `json:"is_today"`
}

// FastingCalendarDay represents a single day in the fasting month calendar.
type FastingCalendarDay struct {
	Date          string   `json:"date"`
	HasActivity   bool     `json:"has_activity"`
	IsToday       bool     `json:"is_today"`
	DurationHours *float64 `json:"duration_hours,omitempty"`
}
