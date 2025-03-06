package models

import "time"

type Freezer struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	StreakCount   int       `json:"streak_count"`
	AllowedDays   int       `json:"allowed_days"`
	Duration      string    `json:"duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
} 