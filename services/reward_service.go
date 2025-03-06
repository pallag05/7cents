package services

import (
	"fmt"
	"time"

	"allen_hackathon/models"
	"allen_hackathon/storage"

	"github.com/google/uuid"
)

// RewardService handles reward-related operations
type RewardService struct {
	store *storage.MemoryStore
}

// NewRewardService creates a new reward service
func NewRewardService() *RewardService {
	return &RewardService{
		store: storage.GetStore(),
	}
}

// GetUserRewards returns all rewards earned by a user
func (s *RewardService) GetUserRewards(userID string) ([]models.Reward, error) {
	userRewards := s.store.GetUserRewards(userID)
	if userRewards == nil {
		return []models.Reward{}, nil
	}

	var rewards []models.Reward
	for _, ur := range userRewards {
		reward, exists := s.store.GetReward(ur.RewardID)
		if exists {
			rewards = append(rewards, *reward)
		}
	}

	return rewards, nil
}

// GetRewardDetails returns details of a specific reward
func (s *RewardService) GetRewardDetails(rewardID string) (*models.Reward, error) {
	reward, exists := s.store.GetReward(rewardID)
	if !exists {
		return nil, fmt.Errorf("reward not found")
	}

	return reward, nil
}

// GetAvailableRewards returns all rewards available for a user's rating
func (s *RewardService) GetAvailableRewards(rating float64) ([]models.Reward, error) {
	allRewards := s.store.GetAllRewards()
	var availableRewards []models.Reward

	for _, reward := range allRewards {
		if rating >= reward.Rating {
			availableRewards = append(availableRewards, *reward)
		}
	}

	return availableRewards, nil
}

// GetRewardProgress returns a user's progress towards earning rewards
func (s *RewardService) GetRewardProgress(userID string) (*models.RewardProgress, error) {
	streak, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return nil, fmt.Errorf("user streak not found")
	}

	availableRewards, err := s.GetAvailableRewards(streak.CurrentRating)
	if err != nil {
		return nil, err
	}

	earnedRewards, err := s.GetUserRewards(userID)
	if err != nil {
		return nil, err
	}

	return &models.RewardProgress{
		UserID:           userID,
		CurrentRating:    streak.CurrentRating,
		AvailableRewards: availableRewards,
		EarnedRewards:    earnedRewards,
		LastUpdated:      time.Now().Format(time.RFC3339),
	}, nil
}

// EarnReward marks a reward as earned by a user
func (s *RewardService) EarnReward(userID, rewardID string) error {
	reward, exists := s.store.GetReward(rewardID)
	if !exists {
		return fmt.Errorf("reward not found")
	}

	streak, exists := s.store.GetStreakToUser(userID)
	if !exists {
		return fmt.Errorf("user streak not found")
	}

	if streak.CurrentRating < reward.Rating {
		return fmt.Errorf("user rating not high enough to earn this reward")
	}

	now := time.Now().Format(time.RFC3339)
	userReward := &models.UserReward{
		ID:        uuid.New().String(),
		UserID:    userID,
		RewardID:  rewardID,
		EarnedAt:  now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.store.SaveUserReward(userReward)
	return nil
}

// CheckAndAwardRewards checks if a user is eligible for any rewards based on their rating
func (s *RewardService) CheckAndAwardRewards(userID string, currentRating float64) error {
	// Get user's existing rewards
	existingRewards := s.store.GetUserRewards(userID)
	existingRewardIDs := make(map[string]bool)
	for _, reward := range existingRewards {
		existingRewardIDs[reward.RewardID] = true
	}

	// Get all rewards
	allRewards := s.store.GetAllRewards()
	for _, reward := range allRewards {
		// Skip if user already has this reward
		if existingRewardIDs[reward.ID] {
			continue
		}

		// Check if user's rating meets the reward requirement
		if currentRating >= reward.Rating {
			now := time.Now().Format(time.RFC3339)
			userReward := &models.UserReward{
				ID:        uuid.New().String(),
				UserID:    userID,
				RewardID:  reward.ID,
				EarnedAt:  now,
				CreatedAt: now,
				UpdatedAt: now,
			}
			s.store.SaveUserReward(userReward)
		}
	}

	return nil
}

// GetAllRewards returns all available rewards
func (s *RewardService) GetAllRewards() ([]models.Reward, error) {
	allRewards := s.store.GetAllRewards()
	var rewards []models.Reward
	for _, reward := range allRewards {
		rewards = append(rewards, *reward)
	}
	return rewards, nil
}
