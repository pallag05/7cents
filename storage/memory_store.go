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
	mu            sync.RWMutex
}

var store *MemoryStore

func GetStore() *MemoryStore {
	if store == nil {
		store = &MemoryStore{
			users:         make(map[string]*models.User),
			streaks:       make(map[string]*models.Streak),
			streakItems:   make(map[string]*models.StreakItem),
			streakToUsers: make(map[string]*models.StreakToUser),
		}
		// Populate with dummy data
		store.populateDummyData()
	}
	return store
}

// populateDummyData adds some initial test data to the store
func (s *MemoryStore) populateDummyData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create dummy users
	dummyUsers := []*models.User{
		{
			ID:        "user1",
			Name:      "John Doe",
			Phone:     "+1234567890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user2",
			Name:      "Jane Smith",
			Phone:     "+1987654321",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user3",
			Name:      "Bob Johnson",
			Phone:     "+1122334455",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Create dummy streaks
	dummyStreaks := []*models.Streak{
		{
			ID:                "streak1",
			Type:              models.StreakTypeBeginner,
			ThresholdDuration: 30,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                "streak2",
			Type:              models.StreakTypeIntermediate,
			ThresholdDuration: 45,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                "streak3",
			Type:              models.StreakTypeAdvanced,
			ThresholdDuration: 60,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	// Create dummy streak items
	dummyStreakItems := []*models.StreakItem{
		{
			ID:        "item1",
			Type:      models.StreakItemTypeVideo,
			StreakID:  "streak1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "item2",
			Type:      models.StreakItemTypeQuestion,
			StreakID:  "streak1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "item3",
			Type:      models.StreakItemTypeFlash,
			StreakID:  "streak2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Create dummy streak to user mappings
	dummyStreakToUsers := []*models.StreakToUser{
		{
			UserID:            "user1",
			StreakCount:       5,
			CurrentStreakID:   "streak1",
			CurrentRating:     "beginner",
			MaxRating:         "beginner",
			LastStreakUpdated: time.Now(),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			UserID:            "user2",
			StreakCount:       12,
			CurrentStreakID:   "streak2",
			CurrentRating:     "intermediate",
			MaxRating:         "intermediate",
			LastStreakUpdated: time.Now(),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			UserID:            "user3",
			StreakCount:       25,
			CurrentStreakID:   "streak3",
			CurrentRating:     "advanced",
			MaxRating:         "advanced",
			LastStreakUpdated: time.Now(),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
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
