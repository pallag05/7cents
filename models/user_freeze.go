package models

import "time"

type UserFreeze struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	UserID     string    `json:"user_id"`
	FreezeUsed int       `json:"freeze_used"`
	Month      string    `json:"month"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
