package handlers

import (
    "Gin/database"
    "Gin/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "Gin/middleware"
    "golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{
            "error":   "Invalid input",
            "details": err.Error(),
        })
        return
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(400, gin.H{
            "error":   "Invalid input",
            "details": err.Error(),
        })
        return
    }
    user.Pass = string(hash)
    // Save the user to the database
    db := database.GetDB()
    result := db.Create(&user)
    if result.Error != nil {
        c.JSON(500, gin.H{"error": "Failed to create user", "details": result.Error.Error()})
        return
    }

}
func LoginUser (c *gin.Context){
    var users models.User 
    var dbUser models.User
    db := database.GetDB()

    if err := c.ShouldBindBodyWithJSON(&users); err != nil {
        c.JSON(400, gin.H{
            "error":   "Invalid input",
            "details": err.Error(),
        })
        return
    }

    if err := db.Where("email = ?", users.Email).First(&dbUser).Error; err != nil{
        c.JSON(401, gin.H{
            "error": "Invalid credentials",
        })
        return

    }

    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Pass), []byte(users.Pass)); err != nil {
        c.JSON(401, gin.H{
            "error": "Invalid credentials",
        })
        return
    }
    token, err := middleware.CreateToken(users)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate token"})
        return
    }

    // Set the token in a cookie
    middleware.SetTokenCookie(c, token)
}
    

// ListUsers retrieves all users (excluding sensitive data)
func ListUsers(c *gin.Context) {
    var users []models.User
    
    // Get database instance
    db := database.GetDB()
    
    // Retrieve only necessary fields and exclude sensitive data
    result := db.Select("ID", "Email", "CreatedAt").Find(&users)
    if result.Error != nil {
        c.JSON(500, gin.H{
            "error": "Failed to retrieve users",
            "details": "Database error occurred",
        })
        return
    }

    // Return clean response with status code
    c.JSON(200, gin.H{
        "count": len(users),
        "users": users,
    })
}
// ReadUser retrieves a specific user by email
func ReadUser(c *gin.Context) {
    email := c.Query("email")
    if email == "" {
        c.JSON(400, gin.H{"error": "Email is required"})
        return
    }   

    var user models.User
    db := database.GetDB()  
    result := db.Select("ID", "Email", "CreatedAt").Where("email = ?", email).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(404, gin.H{"error": "User not found"})
        } else {
            c.JSON(500, gin.H{"error": "Failed to retrieve user", "details": result.Error.Error()})
        }
        return
    }

    c.JSON(200, user)
}
func UpdateUser(c *gin.Context) {
    email := c.Param("email")
    if email == "" {
        c.JSON(400, gin.H{"error": "Email is required"})
        return
    }

    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    db := database.GetDB()
    result := db.Model(&models.User{}).Where("email = ?", email).Update("Pass", user.Pass)
    if result.Error != nil {
        c.JSON(500, gin.H{"error": "Failed to update user", "details": result.Error.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "User updated successfully"})
}
func DeleteUser(c *gin.Context) {
    email := c.Param("email")
    if email == "" {
        c.JSON(400, gin.H{"error": "Email is required"})
        return
    }

    db := database.GetDB()
    result := db.Where("email = ?", email).Delete(&models.User{})
    if result.Error != nil {
        c.JSON(500, gin.H{"error": "Failed to delete user", "details": result.Error.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "User deleted successfully"})
}

