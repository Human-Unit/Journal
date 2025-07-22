package main

import (
	"Gin/database"
	"Gin/handlers"
	"Gin/middleware"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	if err := database.Connect(); err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }

	// Create a new Gin router
	router := gin.Default()

	// Add recovery middleware to handle panics
	router.Use(gin.Recovery())

	// Public routes
	router.POST("/register", handlers.CreateUser)
	router.GET("/login", handlers.LoginUser) // Assuming this is a login endpoint
	//router.POST("/login", handlers.LoginUser)  // You should add this

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.Auth())
	{
		protected.GET("/users", handlers.ListUsers)
		protected.GET("/rusers", handlers.ReadUser)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.DELETE("/users/:id", handlers.DeleteUser)
	}

	// Journal routes
	journal := router.Group("/entries")
	journal.Use(middleware.Auth())
	{
		journal.POST("/", handlers.CreateEntry)
		journal.GET("/list", handlers.ListEntries)
		journal.GET("/:id", handlers.ReadEntry)
		journal.PUT("/:id", handlers.UpdateEntry)
		journal.DELETE("/:id", handlers.DeleteEntry)
		journal.GET("/search", handlers.SearchEntries)
		journal.GET("/mood/:mood", handlers.ListEntriesByMood)
		journal.GET("/tags", handlers.ListTags)
		journal.GET("/private", handlers.ListPrivateEntries)
	}

	// Start server
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
//github added