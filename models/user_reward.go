package models

// UserRewardMapping represents the mapping between users and their earned rewards
type UserRewardMapping struct {
	UserID   string `json:"user_id"`
	RewardID string `json:"reward_id"`
}
