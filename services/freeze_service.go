package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FreezeService struct {
	store *storage.MemoryStore
}

func NewFreezeService() *FreezeService {
	return &FreezeService{
		store: storage.GetStore(),
	}
}

// FreezeStreak freezes a user's streak for a specified duration
func (s *FreezeService) FreezeStreak(userID string, durationDays int) error {
	// Get user's streak info
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return errors.New("user streak not found")
	}

	// Get freeze configuration
	freezeConfig, exists := s.store.GetFreezeConfig()
	if !exists {
		return errors.New("freeze configuration not found")
	}

	// Validate streak count requirement
	if streakToUser.StreakCount < freezeConfig.MinStreakCount {
		return errors.New("insufficient streak count to use freeze")
	}

	// Validate maximum freezes
	if streakToUser.FreezesUsed >= freezeConfig.MaxFreezes {
		return errors.New("maximum number of freezes reached")
	}

	// Validate duration
	if durationDays > freezeConfig.MaxDurationDays {
		return errors.New("freeze duration exceeds maximum allowed")
	}

	// Check if streak is already frozen
	if streakToUser.IsFrozen {
		return errors.New("streak is already frozen")
	}

	// Create freeze record
	freeze := &models.UserFreeze{
		ID:          uuid.New().String(),
		UserID:      userID,
		StartTime:   time.Now(),
		EndTime:     time.Now().AddDate(0, 0, durationDays),
		FreezeCount: streakToUser.FreezesUsed + 1,
	}

	// Update streak info
	streakToUser.IsFrozen = true
	streakToUser.FreezeEndTime = freeze.EndTime
	streakToUser.FreezesUsed++
	streakToUser.LastFreezeUsed = time.Now()

	// Save changes
	s.store.SaveUserFreeze(freeze)
	s.store.SaveStreakToUser(streakToUser)

	return nil
}

// UnfreezeStreak unfreezes a user's streak
func (s *FreezeService) UnfreezeStreak(userID string) error {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return errors.New("user streak not found")
	}

	if !streakToUser.IsFrozen {
		return errors.New("streak is not frozen")
	}

	// Update streak info
	streakToUser.IsFrozen = false
	streakToUser.FreezeEndTime = time.Time{} // Reset freeze end time

	// Save changes
	s.store.SaveStreakToUser(streakToUser)

	return nil
}

// GetFreezeStatus returns the current freeze status for a user
func (s *FreezeService) GetFreezeStatus(userID string) (*models.UserFreeze, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return nil, errors.New("user streak not found")
	}

	if !streakToUser.IsFrozen {
		return nil, nil
	}

	// Get the active freeze record
	freeze, exists := s.store.GetActiveUserFreeze(userID)
	if !exists {
		return nil, errors.New("active freeze record not found")
	}

	return freeze, nil
}

// CheckAndAutoUnfreeze checks if any frozen streaks need to be unfrozen
func (s *FreezeService) CheckAndAutoUnfreeze() error {
	// Get all frozen streaks
	frozenStreaks := s.store.GetAllFrozenStreaks()
	for _, streak := range frozenStreaks {
		if time.Now().After(streak.FreezeEndTime) {
			if err := s.UnfreezeStreak(streak.UserID); err != nil {
				// Log error but continue with other streaks
				continue
			}
		}
	}
	return nil
}
