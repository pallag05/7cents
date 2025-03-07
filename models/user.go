package models

type User struct {
	ID    string  `json:"id"`
	Email string  `json:"email"`
	Score []Score `json:"score"`
	Name  string  `json:"name"`
}

type Score struct {
	Subject string `json:"subject"`
	Score   int    `json:"score"`
}
