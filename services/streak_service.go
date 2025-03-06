package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"time"

	"github.com/google/uuid"
)

type StreakService struct {
	store *storage.MemoryStore
}

func NewStreakService() *StreakService {
	return &StreakService{
		store: storage.GetStore(),
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
