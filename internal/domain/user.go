package domain

import "time"

// User represents a user of the application.
// Currently single-user, but designed for multi-user expansion.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
