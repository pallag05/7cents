package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"errors"
	"time"

	"github.com/google/uuid"
)

type StreakService struct {
	store *storage.MemoryStore
}

// RatingBreakdown provides detailed information about rating calculation
type RatingBreakdown struct {
	BaseScore           float64 `json:"base_score"`
	StreakMultiplier    float64 `json:"streak_multiplier"`
	TypeMultiplier      float64 `json:"type_multiplier"`
	PenaltyPoints       float64 `json:"penalty_points"`
	FinalScore          float64 `json:"final_score"`
	CurrentRating       string  `json:"current_rating"`
	StreakCount         int     `json:"streak_count"`
	StreakType          string  `json:"streak_type"`
	LastStreakUpdated   string  `json:"last_streak_updated"`
	DaysSinceLastStreak int     `json:"days_since_last_streak"`
}

func NewStreakService() *StreakService {
	return &StreakService{
		store: storage.GetStore(),
	}
}

// CalculateUserRating calculates the user's rating based on their streak performance
func (s *StreakService) CalculateUserRating(userID string) (string, error) {
	breakdown, err := s.GetRatingBreakdown(userID)
	if err != nil {
		return "", err
	}

	return breakdown.CurrentRating, nil
}

// GetRatingBreakdown provides detailed information about the rating calculation
func (s *StreakService) GetRatingBreakdown(userID string) (*RatingBreakdown, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return nil, errors.New("user not found")
	}

	streak, exists := s.store.GetStreak(streakToUser.CurrentStreakID)
	if !exists {
		return nil, errors.New("streak not found")
	}

	// Calculate base score (0-100)
	baseScore := calculateBaseScore(streakToUser.StreakCount)

	// Calculate streak multiplier (1.0-2.0)
	streakMultiplier := calculateStreakMultiplier(streakToUser.StreakCount)

	// Calculate type multiplier based on streak type
	typeMultiplier := calculateTypeMultiplier(streak.Type)

	// Calculate penalty points for missed streaks
	penaltyPoints := calculatePenaltyPoints(streakToUser.LastStreakUpdated)

	// Calculate final score
	finalScore := (baseScore * streakMultiplier * typeMultiplier) - penaltyPoints

	// Ensure final score is within bounds
	if finalScore < 0 {
		finalScore = 0
	}
	if finalScore > 100 {
		finalScore = 100
	}

	// Calculate current rating based on final score
	currentRating := calculateRatingFromScore(finalScore)

	// Calculate days since last streak
	daysSinceLastStreak := int(time.Since(streakToUser.LastStreakUpdated).Hours() / 24)

	return &RatingBreakdown{
		BaseScore:           baseScore,
		StreakMultiplier:    streakMultiplier,
		TypeMultiplier:      typeMultiplier,
		PenaltyPoints:       penaltyPoints,
		FinalScore:          finalScore,
		CurrentRating:       currentRating,
		StreakCount:         streakToUser.StreakCount,
		StreakType:          string(streak.Type),
		LastStreakUpdated:   streakToUser.LastStreakUpdated.Format(time.RFC3339),
		DaysSinceLastStreak: daysSinceLastStreak,
	}, nil
}

// Helper functions for rating calculation
func calculateBaseScore(streakCount int) float64 {
	// Base score increases with streak count, but with diminishing returns
	if streakCount <= 0 {
		return 0
	}
	if streakCount >= 30 {
		return 100
	}
	return float64(streakCount) * 3.33 // Linear scaling up to 30 days
}

func calculateStreakMultiplier(streakCount int) float64 {
	// Multiplier increases with streak count, but with diminishing returns
	if streakCount <= 0 {
		return 1.0
	}
	if streakCount >= 30 {
		return 2.0
	}
	return 1.0 + (float64(streakCount) * 0.033) // Linear scaling up to 30 days
}

func calculateTypeMultiplier(streakType models.StreakType) float64 {
	switch streakType {
	case models.StreakTypeBeginner:
		return 1.0
	case models.StreakTypeIntermediate:
		return 1.2
	case models.StreakTypeAdvanced:
		return 1.5
	default:
		return 1.0
	}
}

func calculatePenaltyPoints(lastStreakUpdated time.Time) float64 {
	daysSinceLastStreak := time.Since(lastStreakUpdated).Hours() / 24
	if daysSinceLastStreak <= 1 {
		return 0
	}

	// Exponential penalty for missed streaks
	penalty := 0.0
	for i := 1.0; i <= daysSinceLastStreak; i++ {
		penalty += i * 2 // Each day adds more penalty
	}
	return penalty
}

func calculateRatingFromScore(score float64) string {
	switch {
	case score >= 90:
		return "expert"
	case score >= 75:
		return "advanced"
	case score >= 50:
		return "intermediate"
	case score >= 25:
		return "beginner"
	default:
		return "novice"
	}
}

// RecordActivity records a learning activity and updates the user's streak
func (s *StreakService) RecordActivity(userID string, activityType models.StreakItemType) error {
	// 1. Get user's current streak info
	streakToUser, err := s.getUserStreak(userID)
	if err != nil {
		return err
	}

	// 2. Check if the activity is within the same day
	if !isSameDay(streakToUser.LastStreakUpdated, time.Now()) {
		// 3. Check if streak should be broken
		if shouldBreakStreak(streakToUser.LastStreakUpdated) {
			streakToUser.StreakCount = 0
			streakToUser.CurrentStreakID = ""
		}
	}

	// 4. Get or create streak based on user's level
	streak, err := s.getOrCreateStreak(userID)
	if err != nil {
		return err
	}

	// 5. Create streak item
	streakItem := &models.StreakItem{
		ID:       generateID(),
		Type:     activityType,
		StreakID: streak.ID,
	}

	// 6. Update streak count and last updated time
	streakToUser.StreakCount++
	streakToUser.CurrentStreakID = streak.ID
	streakToUser.LastStreakUpdated = time.Now()

	// 7. Update user's rating if needed
	if shouldUpdateRating(streakToUser) {
		streakToUser.CurrentRating = calculateNewRating(streakToUser)
		if streakToUser.CurrentRating > streakToUser.MaxRating {
			streakToUser.MaxRating = streakToUser.CurrentRating
		}
	}

	// 8. Save all changes
	return s.saveChanges(streakToUser, streakItem)
}

// GetUserStreakInfo returns the user's current streak information
func (s *StreakService) GetUserStreakInfo(userID string) (*models.StreakToUser, error) {
	return s.getUserStreak(userID)
}

// Helper functions
func (s *StreakService) getUserStreak(userID string) (*models.StreakToUser, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		// Create new streak for user if not exists
		streakToUser = &models.StreakToUser{
			UserID:            userID,
			StreakCount:       0,
			CurrentStreakID:   "",
			CurrentRating:     "beginner",
			MaxRating:         "beginner",
			LastStreakUpdated: time.Now(),
		}
		s.store.SaveStreakToUser(streakToUser)
	}
	return streakToUser, nil
}

func (s *StreakService) getOrCreateStreak(userID string) (*models.Streak, error) {
	streak, exists := s.store.GetStreak(userID)
	if !exists {
		// Create new streak with default settings
		streak = &models.Streak{
			ID:                generateID(),
			Type:              models.StreakTypeBeginner,
			ThresholdDuration: 30, // 30 minutes default
		}
		s.store.SaveStreak(streak)
	}
	return streak, nil
}

func (s *StreakService) saveChanges(streakToUser *models.StreakToUser, streakItem *models.StreakItem) error {
	// Save streak item
	s.store.SaveStreakItem(streakItem)

	// Update streak to user
	s.store.SaveStreakToUser(streakToUser)

	return nil
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func shouldBreakStreak(lastUpdated time.Time) bool {
	return time.Since(lastUpdated) > 24*time.Hour
}

func shouldUpdateRating(streakToUser *models.StreakToUser) bool {
	return streakToUser.StreakCount%5 == 0
}

func calculateNewRating(streakToUser *models.StreakToUser) string {
	switch streakToUser.CurrentRating {
	case "beginner":
		if streakToUser.StreakCount >= 10 {
			return "intermediate"
		}
	case "intermediate":
		if streakToUser.StreakCount >= 20 {
			return "advanced"
		}
	}
	return streakToUser.CurrentRating
}

func generateID() string {
	return uuid.New().String()
}
