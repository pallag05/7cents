package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	BatchID   string    `json:"batch_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StreakToUser struct {
	UserID            string    `json:"user_id"`
	StreakCount       int       `json:"streak_count"`
	CurrentStreakID   string    `json:"current_streak_id"`
	CurrentRating     float64   `json:"current_rating"`
	MaxRating         float64   `json:"max_rating"`
	LastStreakUpdated time.Time `json:"last_streak_updated"`
	IsFrozen          bool      `json:"is_frozen"`
	FreezeEndTime     time.Time `json:"freeze_end_time"`
	FreezesUsed       int       `json:"freezes_used"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type FreezeConfig struct {
	MaxFreezes int `json:"max_freezes"`
	Duration   int `json:"duration"` // in hours
}
