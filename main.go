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
	groupHandler := handlers.NewGroupHandler(groupService, store)

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// Group routes
	api := r.Group("/api")
	{
		groups := api.Group("/groups")
		{
			groups.POST("", groupHandler.CreateGroup)
			groups.GET("/user/:user_id", groupHandler.GetGroupsPage)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.POST("/:id/join/:user_id", groupHandler.JoinGroup)
			groups.PUT("/:id", groupHandler.UpdateGroup)
			groups.POST("/:id/leave/:user_id", groupHandler.LeaveGroup)
			groups.POST("/search", groupHandler.SearchGroupsByTag)
			groups.POST("/:id/reject/:user_id", groupHandler.RejectGroupRecommendation)
		}
	}

	r.Run(":96")
}
