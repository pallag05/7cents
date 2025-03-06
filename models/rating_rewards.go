package models

import "time"

type RatingRewards struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Rating    int       `json:"rating"`
	Reward    string    `json:"reward"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
