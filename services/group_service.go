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
