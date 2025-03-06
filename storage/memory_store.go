package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/pallag05/7cents/models"
)

type MemoryStore struct {
	users  map[string]*models.User
	groups map[string]*models.Group
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users:  make(map[string]*models.User),
		groups: make(map[string]*models.Group),
	}

	// Add dummy groups
	subjects := []string{"physics", "chemistry", "maths"}
	for i, subject := range subjects {
		group := &models.Group{
			ID:            uuid.New().String(),
			Title:         subject + " Study Group",
			Description:   "A group for studying " + subject,
			Members:       []string{},
			Tag:           subject,
			Type:          "study",
			Private:       false,
			Messages:      []models.Message{},
			CreateBy:      uuid.New().String(),
			Capacity:      10,
			ActivityScore: 100 - (i * 25), // 100, 75, 50
		}
		store.groups[group.ID] = group
	}

	// Add one more physics group with different activity score
	physicsGroup2 := &models.Group{
		ID:            uuid.New().String(),
		Title:         "Advanced Physics Group",
		Description:   "Advanced physics study group",
		Members:       []string{},
		Tag:           "physics",
		Type:          "study",
		Private:       false,
		Messages:      []models.Message{},
		CreateBy:      uuid.New().String(),
		Capacity:      10,
		ActivityScore: 85,
	}
	store.groups[physicsGroup2.ID] = physicsGroup2

	// Add some dummy messages to groups
	for _, group := range store.groups {
		message := models.Message{
			ID:        uuid.New().String(),
			Content:   "Welcome to " + group.Title,
			SenderId:  group.CreateBy,
			Timestamp: time.Now(),
		}
		group.Messages = append(group.Messages, message)
	}

	return store
}

// User operations
func (s *MemoryStore) GetUser(id string) (*models.User, error) {
	if user, exists := s.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (s *MemoryStore) CreateUser(user *models.User) error {
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) UpdateUser(user *models.User) error {
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) DeleteUser(id string) error {
	delete(s.users, id)
	return nil
}

// Group operations
func (s *MemoryStore) GetGroup(id string) (*models.Group, error) {
	if group, exists := s.groups[id]; exists {
		return group, nil
	}
	return nil, nil
}

func (s *MemoryStore) CreateGroup(group *models.Group) error {
	s.groups[group.ID] = group
	return nil
}

func (s *MemoryStore) UpdateGroup(group *models.Group) error {
	s.groups[group.ID] = group
	return nil
}

func (s *MemoryStore) DeleteGroup(id string) error {
	delete(s.groups, id)
	return nil
}

func (s *MemoryStore) GetGroupsByUser(userID string) ([]*models.Group, error) {
	var userGroups []*models.Group
	for _, group := range s.groups {
		for _, memberID := range group.Members {
			if memberID == userID {
				userGroups = append(userGroups, group)
				break
			}
		}
	}
	return userGroups, nil
}

func (s *MemoryStore) AddMemberToGroup(groupID string, userID string) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	// Check if user is already a member
	for _, memberID := range group.Members {
		if memberID == userID {
			return nil
		}
	}

	group.Members = append(group.Members, userID)
	return s.UpdateGroup(group)
}

func (s *MemoryStore) RemoveMemberFromGroup(groupID string, userID string) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	// Remove user from members list
	for i, memberID := range group.Members {
		if memberID == userID {
			group.Members = append(group.Members[:i], group.Members[i+1:]...)
			return s.UpdateGroup(group)
		}
	}
	return nil
}

func (s *MemoryStore) AddMessageToGroup(groupID string, message *models.Message) error {
	group, err := s.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	group.Messages = append(group.Messages, *message)
	return s.UpdateGroup(group)
}

func (s *MemoryStore) SearchGroupsByTag(tag string) []*models.Group {
	var matchingGroups []*models.Group
	for _, group := range s.groups {
		if group != nil && group.Tag == tag {
			matchingGroups = append(matchingGroups, group)
		}
	}
	return matchingGroups
}
