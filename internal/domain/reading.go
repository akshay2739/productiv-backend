package domain

import "time"

// ReadingSession represents a single reading log entry.
type ReadingSession struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index"`
	BookName  string    `json:"book_name"`
	Pages     int       `json:"pages"`
	LoggedAt  time.Time `json:"logged_at" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

// ReadingStats holds aggregated reading statistics.
type ReadingStats struct {
	CurrentStreak int            `json:"current_streak"`
	LongestStreak int            `json:"longest_streak"`
	TotalPages    int            `json:"total_pages"`
	TotalSessions int            `json:"total_sessions"`
	BooksRead     []BookSummary  `json:"books_read"`
	CalendarDays  []CalendarDay  `json:"calendar_days"`
	LoggedToday   bool           `json:"logged_today"`
}

// BookSummary holds aggregated data for a single book.
type BookSummary struct {
	BookName   string `json:"book_name"`
	TotalPages int    `json:"total_pages"`
	Sessions   int    `json:"sessions"`
}
