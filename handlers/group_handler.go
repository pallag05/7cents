package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pallag05/7cents/models"
	"github.com/pallag05/7cents/services"
	"github.com/pallag05/7cents/storage"
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
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag is required in request body"})
		return
	}

	groups := h.store.SearchGroupsByTag(request.Tag)
	c.JSON(http.StatusOK, groups)
}
