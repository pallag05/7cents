package models

type StreakToUser struct {
	UserID            string  `json:"user_id" gorm:"primaryKey"`
	StreakCount       int     `json:"streak_count"`
	CurrentStreakID   string  `json:"current_streak_id"`
	CurrentRating     float64 `json:"current_rating"`
	MaxRating         float64 `json:"max_rating"`
	LastStreakUpdated string  `json:"last_streak_updated_at"`
	FreezesUsed       int     `json:"freezes_used"`     // Number of freezes used in current period
	LastFreezeUsed    string  `json:"last_freeze_used"` // Last time a freeze was used
	IsFrozen          bool    `json:"is_frozen"`        // Whether the streak is currently frozen
	FreezeEndTime     string  `json:"freeze_end_time"`  // When the current freeze ends
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
