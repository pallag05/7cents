package main

import (
	"allen_hackathon/handlers"
	"allen_hackathon/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.

	// Create a default gin router
	router := gin.Default()

	// Initialize services
	streakService := services.NewStreakService()

	// Initialize handlers
	streakHandler := handlers.NewStreakHandler(streakService)
	ratingHandler := handlers.NewRatingHandler(streakService)
	leaderboardHandler := handlers.NewLeaderboardHandler(streakService)

	// API routes
	api := router.Group("/api")
	{
		// Streak routes
		streaks := api.Group("/streaks")
		{
			streaks.POST("/activity", streakHandler.RecordActivity)
			streaks.GET("/user/:user_id", streakHandler.GetUserStreakInfo)
		}

		// Rating routes
		ratings := api.Group("/ratings")
		{
			ratings.GET("/user/:user_id", ratingHandler.CalculateRating)
			ratings.GET("/user/:user_id/breakdown", ratingHandler.GetRatingBreakdown)
		}

		// Leaderboard routes
		leaderboards := api.Group("/leaderboards")
		{
			leaderboards.GET("/batch/:batch_id", leaderboardHandler.GetBatchLeaderboard)
			leaderboards.GET("/top", leaderboardHandler.GetTopPerformers)
			leaderboards.GET("/batch/:batch_id/stats", leaderboardHandler.GetLeaderboardStats)
			leaderboards.GET("/batch/:batch_id/rating-distribution", leaderboardHandler.GetRatingDistribution)
			leaderboards.GET("/batch/:batch_id/streak-distribution", leaderboardHandler.GetStreakDistribution)
		}
	}

	// Start the server on port 96
	fmt.Println("Server starting on http://localhost:96")
	router.Run(":96")
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
