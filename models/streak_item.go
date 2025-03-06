package models

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityTypeVideo     ActivityType = "video"
	ActivityTypeFlashCard ActivityType = "flash_card"
	ActivityTypeQuestion  ActivityType = "question"
)

// StreakItem represents a streak item configuration
type StreakItem struct {
	ID             string       `json:"id"`
	ActivityType   ActivityType `json:"activity_type"`
	ThresholdValue int          `json:"threshold_value"`
	StreakID       string       `json:"streak_id"`
}
