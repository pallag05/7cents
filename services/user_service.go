package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"time"

	"github.com/google/uuid"
)

// UserService handles user-related operations
type UserService struct {
	store *storage.MemoryStore
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		store: storage.GetStore(),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	user.CreatedAt = now
	user.UpdatedAt = now
	s.store.SaveUser(user)
	return nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID string) (*models.User, error) {
	user, exists := s.store.GetUser(userID)
	if !exists {
		return nil, nil
	}
	return user, nil
}
