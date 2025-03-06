package services

import (
	"allen_hackathon/models"
	"allen_hackathon/storage"
	"fmt"

	"github.com/google/uuid"
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

func (s *GroupService) GetGroupsPage(userID string) (*models.GroupsPageResponse, error) {
	// Get user's group data
	userGroup, err := s.store.GetUserGroup(userID)
	if err != nil {
		return nil, err
	}

	// If user has no group data yet, return empty response
	if userGroup == nil {
		return &models.GroupsPageResponse{
			SystemRecommendedGroups: []models.Group{},
			UserActiveGroups:        []models.Group{},
		}, nil
	}

	// Get active groups
	activeGroups, err := s.store.GetGroupsByIDs(userGroup.ActiveGroups)
	if err != nil {
		return nil, err
	}

	// Get recommended groups
	recommendedGroups, err := s.store.GetGroupsByIDs(userGroup.RecommendedGroups)
	if err != nil {
		return nil, err
	}

	// Convert []*models.Group to []models.Group
	activeGroupsList := make([]models.Group, len(activeGroups))
	for i, group := range activeGroups {
		activeGroupsList[i] = *group
	}

	recommendedGroupsList := make([]models.Group, len(recommendedGroups))
	for i, group := range recommendedGroups {
		recommendedGroupsList[i] = *group
	}

	return &models.GroupsPageResponse{
		SystemRecommendedGroups: recommendedGroupsList,
		UserActiveGroups:        activeGroupsList,
	}, nil
}

func (s *GroupService) GetGroup(id string) (*models.Group, error) {
	group, err := s.store.GetGroup(id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) JoinGroup(groupID string, userID string) error {
	// Get the group
	group, err := s.store.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Check capacity
	if len(group.Members) >= group.Capacity {
		return fmt.Errorf("group has reached maximum capacity")
	}

	// Check if user is already a member
	for _, memberID := range group.Members {
		if memberID == userID {
			return fmt.Errorf("user is already a member of this group")
		}
	}

	// Add user to group members
	if err := s.store.AddMemberToGroup(groupID, userID); err != nil {
		return err
	}

	// Get user's group data
	userGroup, err := s.store.GetUserGroup(userID)
	if err != nil {
		return err
	}

	// If user has no group data yet, create it
	if userGroup == nil {
		userGroup = &models.UserGroup{
			ID:                uuid.New().String(),
			UserID:            userID,
			ActiveGroups:      []string{},
			RecommendedGroups: []string{},
		}
	}

	// Add group to user's active groups if not already present
	isActive := false
	for _, activeGroupID := range userGroup.ActiveGroups {
		if activeGroupID == groupID {
			isActive = true
			break
		}
	}
	if !isActive {
		userGroup.ActiveGroups = append(userGroup.ActiveGroups, groupID)
	}

	// Remove from recommended groups if present
	recommendedGroups := []string{}
	for _, recGroupID := range userGroup.RecommendedGroups {
		if recGroupID != groupID {
			recommendedGroups = append(recommendedGroups, recGroupID)
		}
	}
	userGroup.RecommendedGroups = recommendedGroups

	// Save or update user group data
	if userGroup.ID == "" {
		return s.store.CreateUserGroup(userGroup)
	}
	return s.store.UpdateUserGroup(userGroup)
}

func (s *GroupService) LeaveGroup(groupID string, userID string) error {
	// Get the group
	group, err := s.store.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Check if user is a member
	isMember := false
	for _, memberID := range group.Members {
		if memberID == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return fmt.Errorf("user is not a member of this group")
	}

	// Remove user from group members
	if err := s.store.RemoveMemberFromGroup(groupID, userID); err != nil {
		return err
	}

	// Get user's group data
	userGroup, err := s.store.GetUserGroup(userID)
	if err != nil {
		return err
	}
	if userGroup == nil {
		return fmt.Errorf("user group data not found")
	}

	// Remove group from user's active groups
	activeGroups := []string{}
	for _, activeGroupID := range userGroup.ActiveGroups {
		if activeGroupID != groupID {
			activeGroups = append(activeGroups, activeGroupID)
		}
	}
	userGroup.ActiveGroups = activeGroups

	// Update user group data
	return s.store.UpdateUserGroup(userGroup)
}

func (s *GroupService) UpdateGroup(groupID string, update *models.GroupUpdateRequest) error {
	// Get the group
	group, err := s.store.GetGroup(groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	// Handle message update
	if update.Message != nil {
		message := models.Message{
			ID:        uuid.New().String(),
			Content:   update.Message.Content,
			SenderId:  update.Message.SenderID,
			Timestamp: update.Message.Timestamp,
		}
		if err := s.store.AddMessageToGroup(groupID, &message); err != nil {
			return err
		}
	}

	// Handle action update
	if update.Action != nil {
		// Create action
		action := models.Action{
			ID:        uuid.New().String(),
			Type:      update.Action.Type,
			Content:   update.Action.Content,
			SenderId:  "system",
			Timestamp: update.Action.Timestamp,
		}

		// Add action to group
		if err := s.store.AddActionToGroup(groupID, &action); err != nil {
			return err
		}

		// Also create a message for this action
		actionMessage := models.Message{
			ID:        uuid.New().String(),
			Content:   fmt.Sprintf("[%s] %s", update.Action.Type, update.Action.Content),
			SenderId:  "system",
			Timestamp: update.Action.Timestamp,
		}

		// Add the action message
		if err := s.store.AddMessageToGroup(groupID, &actionMessage); err != nil {
			return err
		}
	}

	return nil
}
