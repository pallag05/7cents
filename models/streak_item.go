package models

import "time"

type StreakItemType string

const (
	StreakItemTypeVideo    StreakItemType = "video"
	StreakItemTypeFlash    StreakItemType = "flash"
	StreakItemTypeQuestion StreakItemType = "question"
)

type StreakItem struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Type      StreakItemType `json:"type"`
	StreakID  string         `json:"streak_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
