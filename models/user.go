package models

type User struct {
	ID    string  `json:"id"`
	Email string  `json:"email"`
	Score []Score `json:"score"`
}

type Score struct {
	Subject string `json:"subject"`
	Score   int    `json:"score"`
}
