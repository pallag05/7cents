package storage

import (
	"allen_hackathon/models"
	"sync"
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
