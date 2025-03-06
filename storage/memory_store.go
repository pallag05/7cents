package storage

import (
	"allen_hackathon/models"
	"sync"
	"time"
)

type MemoryStore struct {
	users              map[string]*models.User
	streaks            map[string]*models.Streak
	streakItems        map[string]*models.StreakItem
	streakToUsers      map[string]*models.StreakToUser
	rewards            map[string]*models.Reward
	userRewardMappings map[string]*models.UserRewardMapping
	userActivities     map[string]*models.UserActivity
	mu                 sync.RWMutex
}

var store *MemoryStore

func GetStore() *MemoryStore {
	if store == nil {
		store = &MemoryStore{
			users:              make(map[string]*models.User),
			streaks:            make(map[string]*models.Streak),
			streakItems:        make(map[string]*models.StreakItem),
			streakToUsers:      make(map[string]*models.StreakToUser),
			rewards:            make(map[string]*models.Reward),
			userRewardMappings: make(map[string]*models.UserRewardMapping),
			userActivities:     make(map[string]*models.UserActivity),
		}
	}
	return store
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

func (s *MemoryStore) GetAllUsers() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
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

func (s *MemoryStore) GetAllStreaks() []*models.Streak {
	s.mu.RLock()
	defer s.mu.RUnlock()

	streaks := make([]*models.Streak, 0, len(s.streaks))
	for _, streak := range s.streaks {
		streaks = append(streaks, streak)
	}
	return streaks
}

// StreakItem operations
func (s *MemoryStore) GetStreakItem(itemID string) (*models.StreakItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, exists := s.streakItems[itemID]
	return item, exists
}

func (s *MemoryStore) SaveStreakItem(item *models.StreakItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streakItems[item.ID] = item
}

func (s *MemoryStore) GetStreakItemsByStreakID(streakID string) []*models.StreakItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []*models.StreakItem
	for _, item := range s.streakItems {
		if item.StreakID == streakID {
			items = append(items, item)
		}
	}
	return items
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

func (s *MemoryStore) GetStreakToUsersByBatchID(batchID string) []*models.StreakToUser {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var streakToUsers []*models.StreakToUser
	for _, streakToUser := range s.streakToUsers {
		if streakToUser.BatchID == batchID {
			streakToUsers = append(streakToUsers, streakToUser)
		}
	}
	return streakToUsers
}

func (s *MemoryStore) GetAllStreakToUsers() []*models.StreakToUser {
	s.mu.RLock()
	defer s.mu.RUnlock()

	streakToUsers := make([]*models.StreakToUser, 0, len(s.streakToUsers))
	for _, streakToUser := range s.streakToUsers {
		streakToUsers = append(streakToUsers, streakToUser)
	}
	return streakToUsers
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

func (s *MemoryStore) GetAllRewards() []*models.Reward {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rewards := make([]*models.Reward, 0, len(s.rewards))
	for _, reward := range s.rewards {
		rewards = append(rewards, reward)
	}
	return rewards
}

func (s *MemoryStore) GetRewardsByRating(rating int) []*models.Reward {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rewards []*models.Reward
	for _, reward := range s.rewards {
		if reward.Rating == rating {
			rewards = append(rewards, reward)
		}
	}
	return rewards
}

// UserRewardMapping operations
func (s *MemoryStore) GetUserRewardMapping(userID, rewardID string) (*models.UserRewardMapping, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	mapping, exists := s.userRewardMappings[userID+"_"+rewardID]
	return mapping, exists
}

func (s *MemoryStore) SaveUserRewardMapping(mapping *models.UserRewardMapping) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userRewardMappings[mapping.UserID+"_"+mapping.RewardID] = mapping
}

func (s *MemoryStore) GetUserRewards(userID string) []*models.Reward {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rewards []*models.Reward
	for _, mapping := range s.userRewardMappings {
		if mapping.UserID == userID {
			if reward, exists := s.rewards[mapping.RewardID]; exists {
				rewards = append(rewards, reward)
			}
		}
	}
	return rewards
}

// UserActivity operations
func (s *MemoryStore) GetUserActivity(userID, activityType string) (*models.UserActivity, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	activity, exists := s.userActivities[userID+"_"+activityType]
	return activity, exists
}

func (s *MemoryStore) SaveUserActivity(activity *models.UserActivity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userActivities[activity.UserID+"_"+activity.ActivityType] = activity
}

func (s *MemoryStore) GetUserActivities(userID string) []*models.UserActivity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activities []*models.UserActivity
	for _, activity := range s.userActivities {
		if activity.UserID == userID {
			activities = append(activities, activity)
		}
	}
	return activities
}

func (s *MemoryStore) UpdateUserActivityValue(userID, activityType string, value int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := userID + "_" + activityType
	activity, exists := s.userActivities[key]
	if !exists {
		return nil // Activity doesn't exist, nothing to update
	}

	activity.Value = value
	activity.UpdatedAt = time.Now()
	s.userActivities[key] = activity
	return nil
}
