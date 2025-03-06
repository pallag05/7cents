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

// UserRewardProgress represents the user's reward progress
type UserRewardProgress struct {
	PreviousReward *Reward `json:"previous_reward"`
	CurrentReward  *Reward `json:"current_reward"`
	NextReward     *Reward `json:"next_reward"`
	MaxRating      float64 `json:"max_rating"`
	CurrentRating  float64 `json:"current_rating"`
	ImageURL       string  `json:"image_url"`
}
