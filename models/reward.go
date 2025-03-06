package models

// RewardType represents the type of reward
type RewardType string

const (
	RewardTypeGoodies      RewardType = "GOODIES"
	RewardTypeFreezeUnlock RewardType = "FREEZE_UNLOCK"
)

// Reward represents a reward definition
type Reward struct {
	ID          string     `json:"id"`
	Rating      int        `json:"rating"`
	Title       string     `json:"title"`
	Type        RewardType `json:"type"`
	Description string     `json:"description"`
}
