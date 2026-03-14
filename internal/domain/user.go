package domain

import "time"

// User represents a user of the application.
// Currently single-user, but designed for multi-user expansion.
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone" gorm:"default:Asia/Kolkata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
