package models

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model `json:"-"`
	ID         uint   `gorm:"primaryKey;autoIncrement"` // ID is auto-incremented
	UserID     uint   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	Mood       string `gorm:"size:20;index"` // Added index if you'll query by mood
	Tags       string `gorm:"size:255"`      // Kept as string but consider adding index
	IsPrivate  bool   `gorm:"default:true;not null"`
}
type EntryIn struct {          // ID is optional for creation, but useful for updates
	UserID     uint   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	Mood       string `gorm:"size:20;index"` // Added index if you'll query by mood
	Tags       string `gorm:"size:255"`      // Kept as string but consider adding index
	IsPrivate  bool   `gorm:"default:true;not null"`
}
