package models

// StreakType represents the type of streak
type StreakType string

const (
	StreakTypeSteadyStarter       StreakType = "steady_starter"
	StreakTypeCommittedChallenger StreakType = "committed_challenger"
	StreakTypeExamWarrior         StreakType = "exam_warrior"
)

// Streak represents a streak configuration
type Streak struct {
	ID   string     `json:"id"`
	Type StreakType `json:"type"`
}
