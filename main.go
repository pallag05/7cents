package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pallag05/7cents/handlers"
	"github.com/pallag05/7cents/services"
	"github.com/pallag05/7cents/storage"
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
		}
	}

	r.Run(":8080")
}
