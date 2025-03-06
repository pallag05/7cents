package storage

import (
	"github.com/anshparmar/7cents/models"
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
}
