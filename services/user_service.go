package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	store *storage.MemoryStore
}

func NewUserService() *UserService {
	return &UserService{
		store: storage.GetStore(),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	// Generate a new UUID for the user
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save the user
	s.store.SaveUser(user)

	return user, nil
}

// GetUser returns user details
func (s *UserService) GetUser(userID string) (*models.User, error) {
	user, exists := s.store.GetUser(userID)
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
