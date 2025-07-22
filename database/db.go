package database

import (
	"Gin/models"
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Connect establishes the database connection
func Connect() error {
	var err error
	db, err = gorm.Open(sqlite.Open("storage.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return err
	}

	// Migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Entry{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

// Get returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Close cleans up the database connection
func Close() {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}