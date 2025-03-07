package storage

import (
	"allen_hackathon/models"
)

type Store interface {
	// User operations
	GetUser(id string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id string) error

	// Group operations
	GetGroup(id string) (*models.Group, error)
	CreateGroup(group *models.Group) error
	UpdateGroup(group *models.Group) error
	DeleteGroup(id string) error
	GetGroupsByUser(userID string) ([]*models.Group, error)
	AddMemberToGroup(groupID string, userID string) error
	RemoveMemberFromGroup(groupID string, userID string) error
	AddMessageToGroup(groupID string, message *models.Message) error
	GetGroupsByIDs(groupIDs []string) ([]*models.Group, error)
	AddActionToGroup(groupID string, action *models.Action) error
	SearchGroupsByTag(tag string, userID string) []*models.Group

	// UserGroup operations
	GetUserGroup(userID string) (*models.UserGroup, error)
	CreateUserGroup(userGroup *models.UserGroup) error
	UpdateUserGroup(userGroup *models.UserGroup) error

	// Match operations
	GetMatches(userID string) []*models.UserPair
	GetAllMatches() []*models.UserPair
	CreateMatch(user1ID, user2ID string) (*models.UserPair, error)
	DeleteMatch(matchID string) error
}
