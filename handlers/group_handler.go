package handlers

import (
	"net/http"

	"allen_hackathon/models"
	"allen_hackathon/services"
	"allen_hackathon/storage"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService *services.GroupService
	store        *storage.MemoryStore
}

func NewGroupHandler(groupService *services.GroupService, store *storage.MemoryStore) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		store:        store,
	}
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.groupService.CreateGroup(&group); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, group)
}

func (h *GroupHandler) SearchGroupsByTag(c *gin.Context) {
	var request struct {
		Tag string `json:"tag" binding:"required"`
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag is required in request body"})
		return
	}

	groups := h.store.SearchGroupsByTag(request.Tag, request.UserID)
	c.JSON(http.StatusOK, groups)
}

func (h *GroupHandler) GetGroupsPage(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id is required"})
		return
	}

	groupsPage, err := h.groupService.GetGroupsPage(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, groupsPage)
}

// GetGroup handles the GET request for retrieving a group by ID
func (h *GroupHandler) GetGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	group, err := h.groupService.GetGroup(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// JoinGroup handles the POST request for a user to join a group
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if err := h.groupService.JoinGroup(groupID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined the group"})
}

// LeaveGroup handles the POST request for a user to leave a group
func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if err := h.groupService.LeaveGroup(groupID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left the group"})
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	var update models.GroupUpdateRequest
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that exactly one of message or action is provided
	if (update.Message == nil && update.Action == nil) || (update.Message != nil && update.Action != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exactly one of message or action must be provided"})
		return
	}

	if err := h.groupService.UpdateGroup(groupID, &update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group updated successfully"})
}

// RejectGroupRecommendation handles the POST request for rejecting a group recommendation
func (h *GroupHandler) RejectGroupRecommendation(c *gin.Context) {
	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if err := h.groupService.RejectGroupRecommendation(groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group recommendation rejected successfully"})
}
