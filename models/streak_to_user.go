package models

import "time"

type StreakToUser struct {
	UserID            string    `json:"user_id" gorm:"primaryKey"`
	StreakCount       int       `json:"streak_count"`
	CurrentStreakID   string    `json:"current_streak_id"`
	CurrentRating     string    `json:"current_rating"`
	MaxRating         string    `json:"max_rating"`
	LastStreakUpdated time.Time `json:"last_streak_updated_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
