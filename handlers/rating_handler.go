package handlers

import (
	"allen_hackathon/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	streakService *services.StreakService
}

func NewRatingHandler(streakService *services.StreakService) *RatingHandler {
	return &RatingHandler{
		streakService: streakService,
	}
}

// CalculateRating calculates the student's rating based on their streak performance
func (h *RatingHandler) CalculateRating(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	rating, err := h.streakService.CalculateUserRating(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rating)
}

// GetRatingBreakdown provides detailed information about the rating calculation
func (h *RatingHandler) GetRatingBreakdown(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	breakdown, err := h.streakService.GetRatingBreakdown(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}
