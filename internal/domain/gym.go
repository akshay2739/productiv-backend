package domain

import "time"

// WorkoutType defines the kind of workout.
type WorkoutType string

const (
	WorkoutStrength WorkoutType = "strength"
	WorkoutCardio   WorkoutType = "cardio"
	WorkoutHIIT     WorkoutType = "hiit"
	WorkoutYoga     WorkoutType = "yoga"
	WorkoutSports   WorkoutType = "sports"
	WorkoutSwimming WorkoutType = "swimming"
)

// ValidWorkoutTypes returns all valid workout types.
func ValidWorkoutTypes() []WorkoutType {
	return []WorkoutType{
		WorkoutStrength, WorkoutCardio, WorkoutHIIT,
		WorkoutYoga, WorkoutSports, WorkoutSwimming,
	}
}

// IsValidWorkoutType checks if a workout type is valid.
func IsValidWorkoutType(wt WorkoutType) bool {
	for _, valid := range ValidWorkoutTypes() {
		if valid == wt {
			return true
		}
	}
	return false
}

// ValidDurationOptions returns all valid workout duration options in minutes.
func ValidDurationOptions() []int {
	return []int{30, 45, 60, 90, 120}
}

// IsValidDurationMin checks if a workout duration is valid.
func IsValidDurationMin(d int) bool {
	for _, valid := range ValidDurationOptions() {
		if valid == d {
			return true
		}
	}
	return false
}

// GymSession represents a single gym workout log.
type GymSession struct {
	ID          int64       `json:"id" gorm:"primaryKey"`
	UserID      int64       `json:"user_id" gorm:"index"`
	WorkoutType WorkoutType `json:"workout_type"`
	DurationMin *int        `json:"duration_min,omitempty"`
	EnergyLevel *int        `json:"energy_level,omitempty"`
	LoggedAt    time.Time   `json:"logged_at" gorm:"index"`
	CreatedAt   time.Time   `json:"created_at"`
}

// GymStats holds aggregated gym statistics.
type GymStats struct {
	CurrentStreak  int           `json:"current_streak"`
	LongestStreak  int           `json:"longest_streak"`
	WorkoutsWeek   int           `json:"workouts_this_week"`
	WorkoutsMonth  int           `json:"workouts_this_month"`
	WeeklyGoal     int           `json:"weekly_goal"`
	CalendarDays   []CalendarDay `json:"calendar_days"`
	LoggedToday    bool          `json:"logged_today"`
}
