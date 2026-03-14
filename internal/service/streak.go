package service

import (
	"context"
	"time"
)

// dateChecker is a function type that checks if an activity exists on a given date.
type dateChecker func(ctx context.Context, userID int64, date time.Time) (bool, error)

// calculateStreak computes the current streak and longest streak
// by walking backward from today using a date-checking function.
func calculateStreak(ctx context.Context, userID int64, loc *time.Location, checker dateChecker) (current int, longest int, err error) {
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// Check today first
	hasToday, err := checker(ctx, userID, today)
	if err != nil {
		return 0, 0, err
	}

	streak := 0
	if hasToday {
		streak = 1
	}

	// Walk backward from yesterday
	checkDate := today.AddDate(0, 0, -1)
	if !hasToday {
		// If today has no activity, the streak might still be alive
		// (grace period: check from yesterday)
		hasYesterday, err := checker(ctx, userID, checkDate)
		if err != nil {
			return 0, 0, err
		}
		if !hasYesterday {
			return 0, 0, nil // streak is broken
		}
		streak = 1
		checkDate = checkDate.AddDate(0, 0, -1)
	}

	// Continue walking backward (max 365 days to prevent unbounded lookback)
	const maxLookback = 365
	for i := 0; i < maxLookback; i++ {
		has, err := checker(ctx, userID, checkDate)
		if err != nil {
			return 0, 0, err
		}
		if !has {
			break
		}
		streak++
		checkDate = checkDate.AddDate(0, 0, -1)
	}

	return streak, streak, nil
}

// calculateStreakWithHistory computes streak and tracks the longest across all history.
// This is used when we also have a separate "best streak" stored/calculated.
func calculateCurrentStreak(ctx context.Context, userID int64, loc *time.Location, checker dateChecker) (int, error) {
	current, _, err := calculateStreak(ctx, userID, loc, checker)
	return current, err
}

// buildCalendarDays creates the 14-day calendar data.
func buildCalendarDays(ctx context.Context, userID int64, loc *time.Location, checker dateChecker) ([]CalendarDay, error) {
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	days := make([]CalendarDay, 14)
	for i := 13; i >= 0; i-- {
		date := today.AddDate(0, 0, -(13 - i))
		has, err := checker(ctx, userID, date)
		if err != nil {
			return nil, err
		}
		days[i] = CalendarDay{
			Date:        date.Format("2006-01-02"),
			HasActivity: has,
			IsToday:     date.Equal(today),
		}
	}
	return days, nil
}

// CalendarDay represents a single day in the 14-day calendar view.
type CalendarDay struct {
	Date        string `json:"date"`
	HasActivity bool   `json:"has_activity"`
	IsToday     bool   `json:"is_today"`
}

// getUserTimezone loads the timezone for the user. Falls back to Asia/Kolkata.
func getUserTimezone(tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc, _ = time.LoadLocation("Asia/Kolkata")
	}
	return loc
}
