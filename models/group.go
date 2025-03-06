package models

import "time"

type Group struct {
	ID            string    `json:"id"`
	Title         string    `json:"Title"`
	Description   string    `json:"description"`
	Members       []string  `json:"members"`
	Tags          []string  `json:"tags"`
	Type          string    `json:"type"`
	Private       bool      `json:"private"`
	Messages      []Message `json:"messages"`
	CreateBy      string    `json:"createBy"`
	Capacity      int       `json:"capacity"`
	ActivityScore int       `json:"activityScore"`
}

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Sender    User      `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}
