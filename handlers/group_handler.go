package handlers

import (
	"allen_hackathon/models"
	"allen_hackathon/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
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
