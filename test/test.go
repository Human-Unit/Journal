package test

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "Gin/proto/gen/go" // Ensure this matches your actual generated protobuf package
)

type Token struct {
	Token string `json:"token"` // Field must be exported (capitalized)
}

func GetData(c *gin.Context) (string, error) {
	var token Token
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"details": err.Error(),
		})
		return "", err
	}
	return token.Token, nil
}

func SendToken(c *gin.Context) {
	// Get token from request
	token, err := GetData(c)
	if err != nil {
		return // Error response already handled in GetData
	}

	// Set up gRPC connection
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to auth service",
		})
		return
	}
	defer conn.Close()

	// Create gRPC client
	client := pb.NewAuthServiceClient(conn)

	// Call gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.SendString(ctx, &pb.SendStringRequest{
		Value: token, // Pass the extracted token string
	})
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to communicate with auth service",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"acknowledgment": response.Acknowledgment,
	})
}