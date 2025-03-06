package models

import "time"

// UserFreeze represents a freeze period for a user's streak
type UserFreeze struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	FreezeCount int       `json:"freeze_count"` // Number of freezes used
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FreezeConfig represents the configuration for freeze functionality
type FreezeConfig struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	MinStreakCount  int       `json:"min_streak_count"`  // Minimum streak count required to use freeze
	MaxFreezes      int       `json:"max_freezes"`       // Maximum number of freezes allowed
	MaxDurationDays int       `json:"max_duration_days"` // Maximum duration of a single freeze in days
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
