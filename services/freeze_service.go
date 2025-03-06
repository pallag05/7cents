package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FreezeService handles streak freeze operations
type FreezeService struct {
	store *storage.MemoryStore
}

// NewFreezeService creates a new freeze service
func NewFreezeService() *FreezeService {
	return &FreezeService{
		store: storage.GetStore(),
	}
}

// FreezeStreak freezes a user's streak for a specified duration
func (s *FreezeService) FreezeStreak(userID string, durationDays int) error {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return fmt.Errorf("user streak not found")
	}

	// Get freeze configuration
	freezeConfig, exists := s.store.GetFreezeConfig()
	if !exists {
		return fmt.Errorf("freeze configuration not found")
	}

	// Check if user has already used maximum freezes
	if streakToUser.FreezesUsed >= freezeConfig.MaxFreezes {
		return fmt.Errorf("maximum number of freezes used")
	}

	// Check if user has minimum required streak count
	if streakToUser.StreakCount < freezeConfig.MinStreakCount {
		return fmt.Errorf("insufficient streak count")
	}

	// Check if duration is within allowed range
	if durationDays > freezeConfig.MaxDurationDays {
		return fmt.Errorf("freeze duration exceeds maximum allowed")
	}

	now := time.Now()
	endTime := now.AddDate(0, 0, durationDays)

	// Create freeze record
	freeze := &models.UserFreeze{
		ID:        uuid.New().String(),
		UserID:    userID,
		StartTime: now.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}

	// Update streak to user
	streakToUser.IsFrozen = true
	streakToUser.FreezeEndTime = endTime.Format(time.RFC3339)
	streakToUser.LastFreezeUsed = now.Format(time.RFC3339)
	streakToUser.FreezesUsed++
	streakToUser.UpdatedAt = now.Format(time.RFC3339)

	s.store.SaveUserFreeze(freeze)
	s.store.SaveStreakToUser(streakToUser)

	return nil
}

// UnfreezeStreak unfreezes a user's streak
func (s *FreezeService) UnfreezeStreak(userID string) error {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return fmt.Errorf("user streak not found")
	}

	if !streakToUser.IsFrozen {
		return fmt.Errorf("streak is not frozen")
	}

	// Check if freeze has expired
	endTime, err := time.Parse(time.RFC3339, streakToUser.FreezeEndTime)
	if err != nil {
		return fmt.Errorf("invalid freeze end time")
	}

	if time.Now().After(endTime) {
		streakToUser.IsFrozen = false
		streakToUser.FreezeEndTime = ""
		streakToUser.UpdatedAt = time.Now().Format(time.RFC3339)
		s.store.SaveStreakToUser(streakToUser)
		return nil
	}

	return fmt.Errorf("freeze has not expired yet")
}

// GetFreezeStatus returns the current freeze status for a user
func (s *FreezeService) GetFreezeStatus(userID string) (*models.UserFreeze, error) {
	streakToUser, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return nil, fmt.Errorf("user streak not found")
	}

	if !streakToUser.IsFrozen {
		return nil, nil
	}

	freeze, exists := s.store.GetActiveUserFreeze(userID)
	if !exists {
		return nil, fmt.Errorf("active freeze not found")
	}

	return freeze, nil
}

// CheckAndAutoUnfreeze checks if any frozen streaks need to be unfrozen
func (s *FreezeService) CheckAndAutoUnfreeze() error {
	// Get all frozen streaks
	frozenStreaks := s.store.GetAllFrozenStreaks()
	for _, streak := range frozenStreaks {
		freezeEndTime, _ := time.Parse(time.RFC3339, streak.FreezeEndTime)
		if time.Now().After(freezeEndTime) {
			if err := s.UnfreezeStreak(streak.UserID); err != nil {
				// Log error but continue with other streaks
				continue
			}
		}
	}
	return nil
}
