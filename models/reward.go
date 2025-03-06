package models

// RewardType represents the type of a reward
type RewardType string

const (
	RewardTypeBadge       RewardType = "badge"
	RewardTypePoints      RewardType = "points"
	RewardTypeDiscount    RewardType = "discount"
	RewardTypeCertificate RewardType = "certificate"
)

// RewardLevel represents the level of a reward
type RewardLevel string

const (
	RewardLevelNovice       RewardLevel = "novice"
	RewardLevelBeginner     RewardLevel = "beginner"
	RewardLevelIntermediate RewardLevel = "intermediate"
	RewardLevelAdvanced     RewardLevel = "advanced"
	RewardLevelExpert       RewardLevel = "expert"
	RewardLevelMaster       RewardLevel = "master"
)

// Reward represents a reward that can be earned by users
type Reward struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Rating      float64     `json:"rating"`
	Type        RewardType  `json:"type"`
	Level       RewardLevel `json:"level"`
	Value       string      `json:"value"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

// UserReward represents a reward earned by a user
type UserReward struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	RewardID  string `json:"reward_id"`
	EarnedAt  string `json:"earned_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// RewardProgress represents a user's progress towards earning rewards
type RewardProgress struct {
	UserID           string   `json:"user_id"`
	CurrentRating    float64  `json:"current_rating"`
	AvailableRewards []Reward `json:"available_rewards"`
	EarnedRewards    []Reward `json:"earned_rewards"`
	LastUpdated      string   `json:"last_updated"`
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
