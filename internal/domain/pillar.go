package domain

import "time"

// PillarType identifies the kind of pillar.
type PillarType string

const (
	PillarFasting    PillarType = "fasting"
	PillarGym        PillarType = "gym"
	PillarMeditation PillarType = "meditation"
	PillarRetention  PillarType = "retention"
	PillarReading    PillarType = "reading"
)

// Pillar represents a self-improvement tracking area.
type Pillar struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	UserID       int64      `json:"user_id" gorm:"index"`
	Type         PillarType `json:"type" gorm:"uniqueIndex:idx_user_pillar"`
	Name         string     `json:"name"`
	Icon         string     `json:"icon"`
	Color        string     `json:"color"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	DisplayOrder int        `json:"display_order" gorm:"default:0"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// PillarSummary holds streak data for the dashboard.
type PillarSummary struct {
	Type             PillarType `json:"type"`
	Name             string     `json:"name"`
	Icon             string     `json:"icon"`
	Color            string     `json:"color"`
	CurrentStreak    int        `json:"current_streak"`
	HasActivityToday bool       `json:"has_activity_today"`
}

// DashboardData holds the full dashboard response.
type DashboardData struct {
	DisciplineScore int             `json:"discipline_score"`
	Pillars         []PillarSummary `json:"pillars"`
	TodaysFocus     string          `json:"todays_focus"`
}
