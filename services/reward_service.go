package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RewardService struct {
	store *storage.MemoryStore
}

func NewRewardService() *RewardService {
	return &RewardService{
		store: storage.GetStore(),
	}
}

// CheckAndAwardRewards checks if a user is eligible for any rewards based on their rating
func (s *RewardService) CheckAndAwardRewards(userID string, currentRating string) error {
	// Get user's existing rewards
	existingRewards := s.store.GetUserRewards(userID)
	existingRewardIDs := make(map[string]bool)
	for _, reward := range existingRewards {
		existingRewardIDs[reward.RewardID] = true
	}

	// Get rewards for the current rating level
	levelRewards := s.store.GetRewardsByLevel(models.RewardLevel(currentRating))
	for _, reward := range levelRewards {
		// Skip if user already has this reward
		if existingRewardIDs[reward.ID] {
			continue
		}

		// Create user reward
		userReward := &models.UserReward{
			ID:        uuid.New().String(),
			UserID:    userID,
			RewardID:  reward.ID,
			EarnedAt:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save the user reward
		s.store.SaveUserReward(userReward)
	}

	return nil
}

// GetUserRewards returns all rewards earned by a user
func (s *RewardService) GetUserRewards(userID string) ([]models.Reward, error) {
	userRewards := s.store.GetUserRewards(userID)
	if len(userRewards) == 0 {
		return nil, nil
	}

	var rewards []models.Reward
	for _, userReward := range userRewards {
		reward, exists := s.store.GetReward(userReward.RewardID)
		if !exists {
			continue
		}
		rewards = append(rewards, *reward)
	}

	return rewards, nil
}

// GetRewardDetails returns detailed information about a specific reward
func (s *RewardService) GetRewardDetails(rewardID string) (*models.Reward, error) {
	reward, exists := s.store.GetReward(rewardID)
	if !exists {
		return nil, errors.New("reward not found")
	}
	return reward, nil
}

// GetAvailableRewards returns all rewards available for a specific rating level
func (s *RewardService) GetAvailableRewards(rating string) ([]models.Reward, error) {
	rewards := s.store.GetRewardsByLevel(models.RewardLevel(rating))
	if len(rewards) == 0 {
		return nil, nil
	}

	var result []models.Reward
	for _, reward := range rewards {
		result = append(result, *reward)
	}

	return result, nil
}

// GetRewardProgress returns information about a user's progress towards earning rewards
func (s *RewardService) GetRewardProgress(userID string, rating string) (map[string]interface{}, error) {
	// Get user's current rewards
	userRewards, err := s.GetUserRewards(userID)
	if err != nil {
		return nil, err
	}

	// Get available rewards for current rating
	availableRewards, err := s.GetAvailableRewards(rating)
	if err != nil {
		return nil, err
	}

	// Create a map of earned reward IDs
	earnedRewardIDs := make(map[string]bool)
	for _, reward := range userRewards {
		earnedRewardIDs[reward.ID] = true
	}

	// Calculate progress
	progress := make(map[string]interface{})
	progress["total_available"] = len(availableRewards)
	progress["earned"] = len(userRewards)
	progress["remaining"] = len(availableRewards) - len(userRewards)

	// Add details about available rewards
	var availableRewardDetails []map[string]interface{}
	for _, reward := range availableRewards {
		if !earnedRewardIDs[reward.ID] {
			availableRewardDetails = append(availableRewardDetails, map[string]interface{}{
				"id":          reward.ID,
				"name":        reward.Name,
				"description": reward.Description,
				"type":        reward.Type,
				"value":       reward.Value,
			})
		}
	}
	progress["available_rewards"] = availableRewardDetails

	return progress, nil
}
