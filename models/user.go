package models

type User struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	BatchID   string `json:"batch_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
