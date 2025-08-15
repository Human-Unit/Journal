package handlers

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "Gin/proto/gen/go"
)

type UserLog struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// This will hold a persistent gRPC client
var grpcClient pb.AuthServiceClient

func InitGRPC() {
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }

    grpcClient = pb.NewAuthServiceClient(conn)
    log.Println("Connected to gRPC server")
}

func GetData(c *gin.Context) (UserLog, error) {
    var userLog UserLog
    if err := c.ShouldBindJSON(&userLog); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid input format",
            "details": err.Error(),
        })
        return userLog, err
    }
    return userLog, nil
}

func CreateUser(c *gin.Context) {
    start := time.Now()

    data, err := GetData(c)
    if err != nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    callStart := time.Now()
    response, err := grpcClient.SaveUser(ctx, &pb.SaveUserRequest{
        Email:    data.Email,
        Password: data.Password,
    })
    if err != nil {
        log.Printf("gRPC call failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to communicate with auth service",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Status": response.Message,
    })

    log.Printf("gRPC call took: %v", time.Since(callStart))
    log.Printf("Total request took: %v", time.Since(start))
}
