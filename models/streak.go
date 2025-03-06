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
