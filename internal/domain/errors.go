package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrActiveFastExists    = errors.New("an active fasting session already exists")
	ErrNoActiveFast        = errors.New("no active fasting session to end")
	ErrAlreadyLoggedToday  = errors.New("workout already logged for today")
	ErrActiveSessionExists = errors.New("an active meditation session already exists")
	ErrNoActiveSession     = errors.New("no active meditation session to end")
	ErrActiveStreakExists   = errors.New("an active retention streak already exists")
	ErrNoActiveStreak      = errors.New("no active retention streak to reset")
	ErrInvalidWorkoutType  = errors.New("invalid workout type")
	ErrInvalidDuration     = errors.New("invalid duration option")
	ErrInvalidMoodValue    = errors.New("mood value must be between 1 and 5")
	ErrInvalidEnergyLevel  = errors.New("energy level must be between 1 and 5")
	ErrInvalidBookName     = errors.New("book name is required")
	ErrInvalidPageCount    = errors.New("pages must be between 1 and 1000")
)
