package models

import "time"

type Group struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Members       []string  `json:"members"`
	Tag           string    `json:"tag"`
	Type          string    `json:"type"`
	Private       bool      `json:"private"`
	Messages      []Message `json:"messages"`
	Actions       []Action  `json:"actions"`
	CreateBy      string    `json:"createBy"`
	Capacity      int       `json:"capacity"`
	ActivityScore int       `json:"activityScore"`
}

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	SenderId  string    `json:"senderId"`
	Timestamp time.Time `json:"timestamp"`
}

type Action struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	SenderId  string    `json:"senderId"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	ActionTypeCall = "CALL"
	ActionTypeTest = "TEST"
)

type GroupUpdateRequest struct {
	Message *MessageUpdate `json:"message,omitempty"`
	Action  *ActionUpdate  `json:"action,omitempty"`
}

type MessageUpdate struct {
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ActionUpdate struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
