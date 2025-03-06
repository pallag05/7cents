package models

import "time"

// UserActivity represents a user's daily activity record
type UserActivity struct {
	UserID       string    `json:"user_id"`
	Value        int       `json:"value"`
	ActivityType string    `json:"activity_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
