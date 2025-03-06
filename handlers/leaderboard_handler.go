package handlers

import (
	"allen_hackathon/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	streakService *services.StreakService
}

func NewLeaderboardHandler(streakService *services.StreakService) *LeaderboardHandler {
	return &LeaderboardHandler{
		streakService: streakService,
	}
}

// GetBatchLeaderboard returns the leaderboard for a specific batch/class
func (h *LeaderboardHandler) GetBatchLeaderboard(c *gin.Context) {
	batchID := c.Param("batch_id")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id is required"})
		return
	}

	// Query parameters
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	rating := c.Query("rating")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parse dates if provided
	var startTime, endTime time.Time
	var err error
	if startDate != "" {
		startTime, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
			return
		}
	}
	if endDate != "" {
		endTime, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
			return
		}
	}

	leaderboard, err := h.streakService.GetBatchLeaderboard(batchID, limit, offset, rating, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetTopPerformers returns the top performers across all batches
func (h *LeaderboardHandler) GetTopPerformers(c *gin.Context) {
	// Query parameters
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	rating := c.Query("rating")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parse dates if provided
	var startTime, endTime time.Time
	var err error
	if startDate != "" {
		startTime, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
			return
		}
	}
	if endDate != "" {
		endTime, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
			return
		}
	}

	performers, err := h.streakService.GetTopPerformers(limit, offset, rating, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, performers)
}

// GetLeaderboardStats returns statistics for the leaderboard
func (h *LeaderboardHandler) GetLeaderboardStats(c *gin.Context) {
	batchID := c.Param("batch_id")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id is required"})
		return
	}

	stats, err := h.streakService.GetLeaderboardStats(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRatingDistribution returns the distribution of ratings in the leaderboard
func (h *LeaderboardHandler) GetRatingDistribution(c *gin.Context) {
	batchID := c.Param("batch_id")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id is required"})
		return
	}

	distribution, err := h.streakService.GetRatingDistribution(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, distribution)
}

// GetStreakDistribution returns the distribution of streaks in the leaderboard
func (h *LeaderboardHandler) GetStreakDistribution(c *gin.Context) {
	batchID := c.Param("batch_id")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id is required"})
		return
	}

	distribution, err := h.streakService.GetStreakDistribution(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, distribution)
}
