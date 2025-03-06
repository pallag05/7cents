package handlers

import (
	"allen_hackathon/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FreezeHandler struct {
	freezeService *services.FreezeService
}

func NewFreezeHandler(freezeService *services.FreezeService) *FreezeHandler {
	return &FreezeHandler{
		freezeService: freezeService,
	}
}

// FreezeStreak handles the request to freeze a user's streak
func (h *FreezeHandler) FreezeStreak(c *gin.Context) {
	var request struct {
		UserID       string `json:"user_id" binding:"required"`
		DurationDays int    `json:"duration_days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.freezeService.FreezeStreak(request.UserID, request.DurationDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Streak frozen successfully"})
}

// UnfreezeStreak handles the request to unfreeze a user's streak
func (h *FreezeHandler) UnfreezeStreak(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	err := h.freezeService.UnfreezeStreak(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Streak unfrozen successfully"})
}

// GetFreezeStatus handles the request to get a user's freeze status
func (h *FreezeHandler) GetFreezeStatus(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	freeze, err := h.freezeService.GetFreezeStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, freeze)
}
