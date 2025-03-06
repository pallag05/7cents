package main

import (
	"allen_hackathon/handlers"
	"allen_hackathon/services"
	"allen_hackathon/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize storage
	store := storage.GetStore()

	// Initialize services
	streakService := services.NewStreakService()

	// Initialize handlers
	streakHandler := handlers.NewStreakHandler(streakService)

	// Set up Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api")
	{
		// Streak routes
		streaks := api.Group("/streaks")
		{
			streaks.POST("/activity", streakHandler.RecordActivity)
			streaks.GET("/user/:user_id", streakHandler.GetUserStreakInfo)
			streaks.GET("/info/:user_id", streakHandler.GetStreakInfo)
		}
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "96"
	}

	// Start server
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
