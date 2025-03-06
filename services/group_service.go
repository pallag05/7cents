package services

import (
	"github.com/google/uuid"
	"github.com/pallag05/7cents/models"
	"github.com/pallag05/7cents/storage"
)

type GroupService struct {
	store storage.Store
}

func NewGroupService(store storage.Store) *GroupService {
	return &GroupService{
		store: store,
	}
}

func (s *GroupService) CreateGroup(group *models.Group) error {
	// Generate a new UUID for the group
	group.ID = uuid.New().String()

	// Set initial values
	group.Messages = make([]models.Message, 0)
	group.ActivityScore = 0

	// Store the group
	return s.store.CreateGroup(group)
}
