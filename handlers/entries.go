package handlers

import (
	"Gin/database"
	"Gin/middleware"
	//"Gin/middleware"
	"Gin/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	// Initialize the database connection
	db = database.GetDB()
}

// CreateEntry handles the creation of a new journal entry
func CreateEntry(c *gin.Context) {
    var input models.EntryIn
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input format",
            "details": err.Error(),
        })
        return
    }
    userID := middleware.ParseToken(c)
    // 3. Create entry with authenticated user
    entry := models.Entry{
        Title:     input.Title,
        UserID:    uint(userID),     
        Content:   input.Content,
        Mood:      input.Mood,
        Tags:      input.Tags,
        IsPrivate: input.IsPrivate,
    }

    // 4. Save to database
    db := database.GetDB()
    if err := db.Create(&entry).Error; err != nil {
        log.Printf("Database error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create entry",
        })
        return
    }

    // 5. Return clean response
    c.JSON(http.StatusCreated, gin.H{
        "id":         entry.ID,
        "title":      entry.Title,
        "created_at": entry.CreatedAt.Format(time.RFC3339),
        "is_private": entry.IsPrivate,
    })
}

func ListEntries(c *gin.Context) {
    // Get user ID from token (just 1 line!)
    userID := middleware.ParseToken(c)
    
    // Now use the userID to fetch entries
    var entries []models.Entry
    database.GetDB().Where("user_id = ?", userID).Find(&entries)
    
    c.JSON(200, gin.H{
        "entries": entries,
        "UserID": userID,
    })
}
func ReadEntry(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	var entry models.Entry
	result := db.First(&entry, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Entry not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve entry", "details": result.Error.Error()})
		}
		return
	}

	c.JSON(200, gin.H{"entry": entry})
}

// UpdateEntry updates an existing journal entry by ID
func UpdateEntry(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	var entry models.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	result := db.Model(&models.Entry{}).Where("id = ?", id).Updates(entry)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to update entry", "details": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Entry updated successfully"})
}

// DeleteEntry deletes a journal entry by ID
func DeleteEntry(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	result := db.Delete(&models.Entry{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete entry", "details": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Entry not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Entry deleted successfully"})
}

// SearchEntries searches for journal entries by a query string
func SearchEntries(c *gin.Context) {
	db := database.GetDB()
	query := c.Query("query")
	if query == "" {
		c.JSON(400, gin.H{"error": "Query parameter is required"})
		return
	}

	var entries []models.Entry
	result := db.Where("content LIKE ?", "%"+query+"%").Find(&entries)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to search entries", "details": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"entries": entries})
}

// ListEntriesByMood retrieves journal entries filtered by mood
func ListEntriesByMood(c *gin.Context) {
	db := database.GetDB()
	mood := c.Param("mood")
	if mood == "" {
		c.JSON(400, gin.H{"error": "Mood is required"})
		return
	}

	var entries []models.Entry
	result := db.Where("mood = ?", mood).Find(&entries)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve entries by mood", "details": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"entries": entries})
}

// ListTags retrieves a list of unique tags from journal entries
func ListTags(c *gin.Context) {
	db := database.GetDB()
	var entries []models.Entry
	result := db.Select("DISTINCT tags").Find(&entries)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve tags", "details": result.Error.Error()})
		return
	}

	var tags []string
	for _, entry := range entries {
		if entry.Tags != "" {
			tags = append(tags, entry.Tags)
		}
	}

	c.JSON(200, gin.H{"tags": tags})
}

// ListPrivateEntries retrieves all private journal entries
func ListPrivateEntries(c *gin.Context) {
	db := database.GetDB()
	var entries []models.Entry
	result := db.Where("is_private = ?", true).Find(&entries)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve private entries", "details": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"private_entries": entries})
}
