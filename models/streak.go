package models

import "time"

type StreakType string

const (
	StreakTypeBeginner     StreakType = "beginner"
	StreakTypeIntermediate StreakType = "intermediate"
	StreakTypeAdvanced     StreakType = "advanced"
)

type Streak struct {
	ID                string     `json:"id" gorm:"primaryKey"`
	Type              StreakType `json:"type"`
	ThresholdDuration int        `json:"threshold_duration"` // in minutes
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// StreakRequirement represents the requirements for maintaining a streak
type StreakRequirement struct {
	Type               StreakType `json:"type"`
	MinDailyDuration   int        `json:"min_daily_duration"`  // Minimum duration in minutes required per day
	MaxGapAllowed      int        `json:"max_gap_allowed"`     // Maximum gap in days allowed before streak breaks
	RequiredActivities []string   `json:"required_activities"` // Types of activities required
	BonusMultiplier    float64    `json:"bonus_multiplier"`    // Multiplier for streak count
	Description        string     `json:"description"`         // Description of the streak type
}

// GetStreakRequirement returns the requirements for a specific streak type
func GetStreakRequirement(streakType StreakType) StreakRequirement {
	switch streakType {
	case StreakTypeBeginner:
		return StreakRequirement{
			Type:               StreakTypeBeginner,
			MinDailyDuration:   30, // 30 minutes
			MaxGapAllowed:      1,  // 1 day gap allowed
			RequiredActivities: []string{"video", "question", "flash"},
			BonusMultiplier:    1.0,
			Description:        "Perfect for beginners. Requires 30 minutes of daily activity with any combination of videos, questions, and flashcards.",
		}
	case StreakTypeIntermediate:
		return StreakRequirement{
			Type:               StreakTypeIntermediate,
			MinDailyDuration:   60, // 1 hour
			MaxGapAllowed:      1,  // 1 day gap allowed
			RequiredActivities: []string{"video", "question", "flash", "quiz"},
			BonusMultiplier:    1.2,
			Description:        "For committed learners. Requires 1 hour of daily activity including videos, questions, flashcards, and quizzes.",
		}
	case StreakTypeAdvanced:
		return StreakRequirement{
			Type:               StreakTypeAdvanced,
			MinDailyDuration:   120, // 2 hours
			MaxGapAllowed:      0,   // No gap allowed
			RequiredActivities: []string{"video", "question", "flash", "quiz", "project"},
			BonusMultiplier:    1.5,
			Description:        "For dedicated learners. Requires 2 hours of daily activity including all activity types and at least one project.",
		}
	default:
		return StreakRequirement{
			Type:               StreakTypeBeginner,
			MinDailyDuration:   30,
			MaxGapAllowed:      1,
			RequiredActivities: []string{"video", "question", "flash"},
			BonusMultiplier:    1.0,
			Description:        "Default beginner streak requirements.",
		}
	}
}
