package main

import (
	"Gin/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"Gin/middleware/grpc_auth"
)

func main() {
	// Initialize gRPC connection once
	handlers.InitGRPC()

	router := gin.Default()

	// Recovery middleware
	router.Use(gin.Recovery())

	// Public routes
	router.POST("/register", handlers.CreateUser)
	router.POST("/login", auth.Login) // Fixed: POST instead of GET
	// router.POST("/auth", test.Login)    // optional test endpoint

	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
