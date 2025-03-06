package models

// User represents a user in the education platform
type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	BatchID string `json:"batch_id"`
}
