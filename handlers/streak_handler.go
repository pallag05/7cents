package handlers

import (
	"allen_hackathon/models"
	"allen_hackathon/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StreakHandler struct {
	streakService *services.StreakService
}

func NewStreakHandler(streakService *services.StreakService) *StreakHandler {
	return &StreakHandler{
		streakService: streakService,
	}
}

// RecordActivity handles the recording of a learning activity
func (h *StreakHandler) RecordActivity(c *gin.Context) {
	var request struct {
		UserID       string                `json:"user_id" binding:"required"`
		ActivityType models.StreakItemType `json:"activity_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.streakService.RecordActivity(request.UserID, request.ActivityType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity recorded successfully"})
}

// GetUserStreakInfo returns the user's current streak information
func (h *StreakHandler) GetUserStreakInfo(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	streakInfo, err := h.streakService.GetUserStreakInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, streakInfo)
}
