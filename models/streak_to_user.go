package models

import "time"

type StreakToUser struct {
	UserID            string    `json:"user_id" gorm:"primaryKey"`
	StreakCount       int       `json:"streak_count"`
	CurrentStreakID   string    `json:"current_streak_id"`
	CurrentRating     string    `json:"current_rating"`
	MaxRating         string    `json:"max_rating"`
	LastStreakUpdated time.Time `json:"last_streak_updated_at"`
	FreezesUsed       int       `json:"freezes_used"`     // Number of freezes used in current period
	LastFreezeUsed    time.Time `json:"last_freeze_used"` // Last time a freeze was used
	IsFrozen          bool      `json:"is_frozen"`        // Whether the streak is currently frozen
	FreezeEndTime     time.Time `json:"freeze_end_time"`  // When the current freeze ends
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
