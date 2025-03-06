package storage

import (
	"allen_hackathon/models"
	"sync"
	"time"
)

type MemoryStore struct {
	users         map[string]*models.User
	streaks       map[string]*models.Streak
	streakItems   map[string]*models.StreakItem
	streakToUsers map[string]*models.StreakToUser
	rewards       map[string]*models.Reward
	userRewards   map[string]*models.UserReward
	userFreezes   map[string]*models.UserFreeze
	freezeConfig  *models.FreezeConfig
	mu            sync.RWMutex
}

var store *MemoryStore

func GetStore() *MemoryStore {
	if store == nil {
		now := time.Now().Format(time.RFC3339)
		store = &MemoryStore{
			users:         make(map[string]*models.User),
			streaks:       make(map[string]*models.Streak),
			streakItems:   make(map[string]*models.StreakItem),
			streakToUsers: make(map[string]*models.StreakToUser),
			rewards:       make(map[string]*models.Reward),
			userRewards:   make(map[string]*models.UserReward),
			userFreezes:   make(map[string]*models.UserFreeze),
			freezeConfig: &models.FreezeConfig{
				ID:              "default",
				MinStreakCount:  7, // Minimum 7 days streak required
				MaxFreezes:      3, // Maximum 3 freezes allowed
				MaxDurationDays: 7, // Maximum 7 days per freeze
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		}
		// Initialize rewards
		store.InitializeRewards()
		// Populate with dummy data
		store.populateDummyData()
	}
	return store
}

// populateDummyData adds some initial test data to the store
func (s *MemoryStore) populateDummyData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Format(time.RFC3339)

	// Create dummy users
	dummyUsers := []*models.User{
		{
			ID:        "user1",
			Name:      "John Doe",
			Phone:     "+1234567890",
			BatchID:   "batch1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "user2",
			Name:      "Jane Smith",
			Phone:     "+1987654321",
			BatchID:   "batch1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "user3",
			Name:      "Bob Johnson",
			Phone:     "+1122334455",
			BatchID:   "batch2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Create dummy streaks
	dummyStreaks := []*models.Streak{
		{
			ID:                "streak1",
			Type:              models.StreakTypeBeginner,
			ThresholdDuration: 30,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		{
			ID:                "streak2",
			Type:              models.StreakTypeIntermediate,
			ThresholdDuration: 45,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		{
			ID:                "streak3",
			Type:              models.StreakTypeAdvanced,
			ThresholdDuration: 60,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
	}

	// Create dummy streak items
	dummyStreakItems := []*models.StreakItem{
		{
			ID:        "item1",
			Type:      models.StreakItemTypeVideo,
			StreakID:  "streak1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "item2",
			Type:      models.StreakItemTypeQuestion,
			StreakID:  "streak1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "item3",
			Type:      models.StreakItemTypeFlash,
			StreakID:  "streak2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Create dummy streak to user mappings
	dummyStreakToUsers := []*models.StreakToUser{
		{
			UserID:            "user1",
			StreakCount:       5,
			CurrentStreakID:   "streak1",
			CurrentRating:     16.65, // 5 * 3.33 (base score)
			MaxRating:         16.65, // 5 * 3.33 (base score)
			LastStreakUpdated: now,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		{
			UserID:            "user2",
			StreakCount:       12,
			CurrentStreakID:   "streak2",
			CurrentRating:     39.96, // 12 * 3.33 (base score)
			MaxRating:         39.96, // 12 * 3.33 (base score)
			LastStreakUpdated: now,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		{
			UserID:            "user3",
			StreakCount:       25,
			CurrentStreakID:   "streak3",
			CurrentRating:     83.25, // 25 * 3.33 (base score)
			MaxRating:         83.25, // 25 * 3.33 (base score)
			LastStreakUpdated: now,
			CreatedAt:         now,
			UpdatedAt:         now,
		},
	}

	// Create dummy rewards
	dummyRewards := []*models.Reward{
		{
			ID:          "reward1",
			Name:        "Novice Badge",
			Description: "Achieved novice level",
			Type:        models.RewardTypeBadge,
			Level:       models.RewardLevelNovice,
			Value:       "0",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "reward2",
			Name:        "Beginner Points",
			Description: "Earned 100 points for reaching beginner level",
			Type:        models.RewardTypePoints,
			Level:       models.RewardLevelBeginner,
			Value:       "100",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "reward3",
			Name:        "Intermediate Discount",
			Description: "10% discount on next course",
			Type:        models.RewardTypeDiscount,
			Level:       models.RewardLevelIntermediate,
			Value:       "10",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "reward4",
			Name:        "Advanced Certificate",
			Description: "Certificate of Advanced Achievement",
			Type:        models.RewardTypeCertificate,
			Level:       models.RewardLevelAdvanced,
			Value:       "0",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "reward5",
			Name:        "Expert Badge",
			Description: "Achieved expert level",
			Type:        models.RewardTypeBadge,
			Level:       models.RewardLevelExpert,
			Value:       "0",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Add all dummy data to the store
	for _, user := range dummyUsers {
		s.users[user.ID] = user
	}
	for _, streak := range dummyStreaks {
		s.streaks[streak.ID] = streak
	}
	for _, item := range dummyStreakItems {
		s.streakItems[item.ID] = item
	}
	for _, streakToUser := range dummyStreakToUsers {
		s.streakToUsers[streakToUser.UserID] = streakToUser
	}
	for _, reward := range dummyRewards {
		s.rewards[reward.ID] = reward
	}
}

// User operations
func (s *MemoryStore) GetUser(userID string) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	return user, exists
}

func (s *MemoryStore) SaveUser(user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
}

// Streak operations
func (s *MemoryStore) GetStreak(streakID string) (*models.Streak, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	streak, exists := s.streaks[streakID]
	return streak, exists
}

func (s *MemoryStore) SaveStreak(streak *models.Streak) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streaks[streak.ID] = streak
}

// StreakItem operations
func (s *MemoryStore) SaveStreakItem(item *models.StreakItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streakItems[item.ID] = item
}

// StreakToUser operations
func (s *MemoryStore) GetStreakToUser(userID string) (*models.StreakToUser, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	streakToUser, exists := s.streakToUsers[userID]
	return streakToUser, exists
}

func (s *MemoryStore) SaveStreakToUser(streakToUser *models.StreakToUser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streakToUsers[streakToUser.UserID] = streakToUser
}

// GetUsersByBatch returns all users in a specific batch
func (s *MemoryStore) GetUsersByBatch(batchID string) []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var batchUsers []*models.User
	for _, user := range s.users {
		if user.BatchID == batchID {
			batchUsers = append(batchUsers, user)
		}
	}
	return batchUsers
}

// GetAllUsers returns all users in the store
func (s *MemoryStore) GetAllUsers() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// Reward operations
func (s *MemoryStore) GetReward(rewardID string) (*models.Reward, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	reward, exists := s.rewards[rewardID]
	return reward, exists
}

func (s *MemoryStore) SaveReward(reward *models.Reward) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rewards[reward.ID] = reward
}

func (s *MemoryStore) GetUserRewards(userID string) []*models.UserReward {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userRewards []*models.UserReward
	for _, reward := range s.userRewards {
		if reward.UserID == userID {
			userRewards = append(userRewards, reward)
		}
	}
	return userRewards
}

func (s *MemoryStore) SaveUserReward(userReward *models.UserReward) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userRewards[userReward.ID] = userReward
}

func (s *MemoryStore) GetRewardsByLevel(level models.RewardLevel) []*models.Reward {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rewards []*models.Reward
	for _, reward := range s.rewards {
		if reward.Level == level {
			rewards = append(rewards, reward)
		}
	}
	return rewards
}

// Freeze operations
func (s *MemoryStore) GetFreezeConfig() (*models.FreezeConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.freezeConfig, true
}

func (s *MemoryStore) SaveUserFreeze(freeze *models.UserFreeze) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userFreezes[freeze.ID] = freeze
}

func (s *MemoryStore) GetActiveUserFreeze(userID string) (*models.UserFreeze, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	for _, freeze := range s.userFreezes {
		if freeze.UserID == userID {
			endTime, err := time.Parse(time.RFC3339, freeze.EndTime)
			if err != nil {
				continue
			}
			if now.Before(endTime) {
				return freeze, true
			}
		}
	}
	return nil, false
}

func (s *MemoryStore) GetAllFrozenStreaks() []*models.StreakToUser {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var frozenStreaks []*models.StreakToUser
	for _, streak := range s.streakToUsers {
		if streak.IsFrozen {
			frozenStreaks = append(frozenStreaks, streak)
		}
	}
	return frozenStreaks
}

// GetAllRewards returns all rewards in the store
func (s *MemoryStore) GetAllRewards() []*models.Reward {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rewards := make([]*models.Reward, 0, len(s.rewards))
	for _, reward := range s.rewards {
		rewards = append(rewards, reward)
	}
	return rewards
}

// GetAllStreakToUsers returns all streak-to-user mappings
func (s *MemoryStore) GetAllStreakToUsers() []*models.StreakToUser {
	var streaks []*models.StreakToUser
	for _, streak := range s.streakToUsers {
		streaks = append(streaks, streak)
	}
	return streaks
}

// InitializeRewards initializes default rewards
func (s *MemoryStore) InitializeRewards() {
	now := time.Now().Format(time.RFC3339)
	s.rewards = map[string]*models.Reward{
		"bronze": {
			ID:          "bronze",
			Name:        "Bronze Achievement",
			Description: "Achieve a rating of 50 or higher",
			Rating:      50,
			Type:        models.RewardTypeBadge,
			Level:       models.RewardLevelNovice,
			Value:       "bronze_badge",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		"silver": {
			ID:          "silver",
			Name:        "Silver Achievement",
			Description: "Achieve a rating of 75 or higher",
			Rating:      75,
			Type:        models.RewardTypeBadge,
			Level:       models.RewardLevelBeginner,
			Value:       "silver_badge",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		"gold": {
			ID:          "gold",
			Name:        "Gold Achievement",
			Description: "Achieve a rating of 90 or higher",
			Rating:      90,
			Type:        models.RewardTypeBadge,
			Level:       models.RewardLevelAdvanced,
			Value:       "gold_badge",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		"points_100": {
			ID:          "points_100",
			Name:        "100 Points",
			Description: "Earn 100 points",
			Rating:      50,
			Type:        models.RewardTypePoints,
			Level:       models.RewardLevelBeginner,
			Value:       "100",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		"points_500": {
			ID:          "points_500",
			Name:        "500 Points",
			Description: "Earn 500 points",
			Rating:      75,
			Type:        models.RewardTypePoints,
			Level:       models.RewardLevelAdvanced,
			Value:       "500",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

// GetAllStreaks returns all streaks in the store
func (s *MemoryStore) GetAllStreaks() []*models.Streak {
	s.mu.RLock()
	defer s.mu.RUnlock()

	streaks := make([]*models.Streak, 0, len(s.streaks))
	for _, streak := range s.streaks {
		streaks = append(streaks, streak)
	}
	return streaks
}
