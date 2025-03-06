package models

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
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}
