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

	// API routes
	api := router.Group("/api")
	{
		// Streak routes
		streaks := api.Group("/streaks")
		{
			streaks.POST("/activity", streakHandler.RecordActivity)
			streaks.GET("/user/:user_id", streakHandler.GetUserStreakInfo)
		}
	}

	// Start the server on port 96
	fmt.Println("Server starting on http://localhost:96")
	router.Run(":96")
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
