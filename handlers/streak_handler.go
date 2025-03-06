package handlers

import (
	"allen_hackathon/models"
	"allen_hackathon/services"
	"net/http"
	"time"

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

// StreakInfoResponse represents the response for streak information
type StreakInfoResponse struct {
	UserID             string                   `json:"user_id"`
	CurrentStreakType  models.StreakType        `json:"current_streak_type"`
	CurrentStreakCount int                      `json:"current_streak_count"`
	LastStreakUpdated  time.Time                `json:"last_streak_updated"`
	IsFrozen           bool                     `json:"is_frozen"`
	FreezeEndTime      time.Time                `json:"freeze_end_time,omitempty"`
	FreezesUsed        int                      `json:"freezes_used"`
	MaxFreezesAllowed  int                      `json:"max_freezes_allowed"`
	Requirements       models.StreakRequirement `json:"requirements"`
	NextMilestone      int                      `json:"next_milestone"`   // Next milestone to achieve
	ProgressToNext     float64                  `json:"progress_to_next"` // Progress percentage to next milestone
}

// GetStreakInfo handles the request to get detailed streak information
func (h *StreakHandler) GetStreakInfo(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Get user's streak info
	streakToUser, err := h.streakService.GetUserStreakInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get current streak
	streak, exists := h.streakService.GetStreak(streakToUser.CurrentStreakID)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streak not found"})
		return
	}

	// Get freeze config
	freezeConfig, exists := h.streakService.GetFreezeConfig()
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "freeze configuration not found"})
		return
	}

	// Get streak requirements
	requirements := models.GetStreakRequirement(streak.Type)

	// Calculate next milestone and progress
	nextMilestone := calculateNextMilestone(streakToUser.StreakCount)
	progressToNext := calculateProgressToNext(streakToUser.StreakCount, nextMilestone)

	response := StreakInfoResponse{
		UserID:             userID,
		CurrentStreakType:  streak.Type,
		CurrentStreakCount: streakToUser.StreakCount,
		LastStreakUpdated:  streakToUser.LastStreakUpdated,
		IsFrozen:           streakToUser.IsFrozen,
		FreezeEndTime:      streakToUser.FreezeEndTime,
		FreezesUsed:        streakToUser.FreezesUsed,
		MaxFreezesAllowed:  freezeConfig.MaxFreezes,
		Requirements:       requirements,
		NextMilestone:      nextMilestone,
		ProgressToNext:     progressToNext,
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to calculate the next milestone
func calculateNextMilestone(currentStreak int) int {
	milestones := []int{7, 14, 21, 30, 60, 90, 180, 365}
	for _, milestone := range milestones {
		if currentStreak < milestone {
			return milestone
		}
	}
	return currentStreak + 30 // Default to next 30 days if all milestones achieved
}

// Helper function to calculate progress to next milestone
func calculateProgressToNext(currentStreak, nextMilestone int) float64 {
	if currentStreak == 0 {
		return 0
	}
	progress := float64(currentStreak) / float64(nextMilestone) * 100
	if progress > 100 {
		return 100
	}
	return progress
}
