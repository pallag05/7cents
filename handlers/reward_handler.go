package handlers

import (
	"allen_hackathon/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	rewardService *services.RewardService
}

func NewRewardHandler(rewardService *services.RewardService) *RewardHandler {
	return &RewardHandler{
		rewardService: rewardService,
	}
}

// GetUserRewards returns all rewards earned by a user
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

// GetRewardDetails returns detailed information about a specific reward
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

// GetAvailableRewards returns all rewards available for a specific rating level
func (h *RewardHandler) GetAvailableRewards(c *gin.Context) {
	rating := c.Param("rating")
	if rating == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating is required"})
		return
	}

	rewards, err := h.rewardService.GetAvailableRewards(rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// GetRewardProgress returns information about a user's progress towards earning rewards
func (h *RewardHandler) GetRewardProgress(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	rating := c.Query("rating")
	if rating == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating is required"})
		return
	}

	progress, err := h.rewardService.GetRewardProgress(userID, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetUserRewardProgress handles the request to get user's reward progress
func (h *RewardHandler) GetUserRewardProgress(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	progress, err := h.rewardService.GetUserRewardProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// CheckAndAwardRewards handles the request to check and award rewards to a user
func (h *RewardHandler) CheckAndAwardRewards(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	currentRating := c.Query("rating")
	if currentRating == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating is required"})
		return
	}

	err := h.rewardService.CheckAndAwardRewards(userID, currentRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rewards checked and awarded successfully"})
}
