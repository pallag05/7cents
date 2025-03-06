package handlers

import (
	"allen_hackathon/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RewardHandler handles reward-related HTTP requests
type RewardHandler struct {
	rewardService *services.RewardService
}

// NewRewardHandler creates a new reward handler
func NewRewardHandler(rewardService *services.RewardService) *RewardHandler {
	return &RewardHandler{
		rewardService: rewardService,
	}
}

// GetUserRewards handles requests to get all rewards earned by a user
func (h *RewardHandler) GetUserRewards(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	rewards, err := h.rewardService.GetUserRewards(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// GetRewardDetails handles requests to get details of a specific reward
func (h *RewardHandler) GetRewardDetails(c *gin.Context) {
	rewardID := c.Param("reward_id")
	if rewardID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reward_id is required"})
		return
	}

	reward, err := h.rewardService.GetRewardDetails(rewardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reward)
}

// GetAvailableRewards returns all rewards available for a user's rating
func (h *RewardHandler) GetAvailableRewards(c *gin.Context) {
	var rating float64
	var err error

	// Try to get rating from URL parameter first
	ratingStr := c.Param("rating")
	if ratingStr != "" {
		rating, err = strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating parameter"})
			return
		}
	} else {
		// If no rating provided, get all rewards
		rewards, err := h.rewardService.GetAllRewards()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rewards)
		return
	}

	rewards, err := h.rewardService.GetAvailableRewards(rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rewards)
}

// GetRewardProgress handles requests to get a user's progress towards earning rewards
func (h *RewardHandler) GetRewardProgress(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	progress, err := h.rewardService.GetRewardProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// EarnReward handles requests to mark a reward as earned by a user
func (h *RewardHandler) EarnReward(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		RewardID string `json:"reward_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.rewardService.EarnReward(req.UserID, req.RewardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reward earned successfully"})
}
