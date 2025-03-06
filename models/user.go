package models

type User struct {
	ID    string  `json:"id"`
	Stream string  `json:"stream"`
	Score []Score `json:"score"`
}

type Score struct {
	Subject string `json:"subject"`
	Score   int    `json:"score"`
}
