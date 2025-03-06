package models

import "time"

// StreakToUser represents the relationship between a user and their streak
type StreakToUser struct {
	UserID            string    `json:"user_id"`
	StreakCount       int       `json:"streak_count"`
	CurrentStreakID   string    `json:"current_streak_id"`
	CurrentRating     int       `json:"current_rating"`
	MaxRating         int       `json:"max_rating"`
	LastStreakUpdated time.Time `json:"last_streak_updated_at"`
	BatchID           string    `json:"batch_id"`
	IsFrozen          bool      `json:"is_frozen"`
	FrozenAt          time.Time `json:"frozen_at"`
}
