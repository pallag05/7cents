package main

import (
	"allen_hackathon/handlers"
	"allen_hackathon/services"
	"allen_hackathon/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize store
	store := storage.NewMemoryStore()

	// Initialize services
	groupService := services.NewGroupService(store)

	// Initialize handlers
	groupHandler := handlers.NewGroupHandler(groupService)

	// Group routes
	api := r.Group("/api")
	{
		groups := api.Group("/groups")
		{
			groups.POST("", groupHandler.CreateGroup)
			groups.GET("/user/:user_id", groupHandler.GetGroupsPage)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.POST("/:id/join/:user_id", groupHandler.JoinGroup)
		}
	}

	r.Run(":8080")
}
