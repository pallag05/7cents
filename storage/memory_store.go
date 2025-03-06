package storage

import "github.com/pallag05/7cents/models"

type MemoryStore struct {
	users      map[string]*models.User
	groups     map[string]*models.Group
	userGroups map[string]*models.UserGroup
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:      make(map[string]*models.User),
		groups:     make(map[string]*models.Group),
		userGroups: make(map[string]*models.UserGroup),
	}
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

// UserGroup operations
func (s *MemoryStore) GetUserGroup(userID string) (*models.UserGroup, error) {
	if userGroup, exists := s.userGroups[userID]; exists {
		return userGroup, nil
	}
	return nil, nil
}

func (s *MemoryStore) CreateUserGroup(userGroup *models.UserGroup) error {
	s.userGroups[userGroup.UserID] = userGroup
	return nil
}

func (s *MemoryStore) UpdateUserGroup(userGroup *models.UserGroup) error {
	s.userGroups[userGroup.UserID] = userGroup
	return nil
}

func (s *MemoryStore) GetGroupsByIDs(groupIDs []string) ([]*models.Group, error) {
	var groups []*models.Group
	for _, id := range groupIDs {
		if group, exists := s.groups[id]; exists {
			groups = append(groups, group)
		}
	}
	return groups, nil
}
