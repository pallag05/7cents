package models

import "time"

// RewardType represents the type of reward
type RewardType string

const (
	RewardTypeBadge       RewardType = "badge"
	RewardTypePoints      RewardType = "points"
	RewardTypeDiscount    RewardType = "discount"
	RewardTypeCertificate RewardType = "certificate"
)

// RewardLevel represents the level required to earn the reward
type RewardLevel string

const (
	RewardLevelNovice       RewardLevel = "novice"
	RewardLevelBeginner     RewardLevel = "beginner"
	RewardLevelIntermediate RewardLevel = "intermediate"
	RewardLevelAdvanced     RewardLevel = "advanced"
	RewardLevelExpert       RewardLevel = "expert"
)

// Reward represents a reward that can be earned by users
type Reward struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        RewardType  `json:"type"`
	Level       RewardLevel `json:"level"`
	Value       int         `json:"value"` // Points amount, discount percentage, etc.
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// UserReward represents a reward earned by a user
type UserReward struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RewardID  string    `json:"reward_id"`
	EarnedAt  time.Time `json:"earned_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
