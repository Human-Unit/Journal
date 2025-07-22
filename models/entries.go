package models

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model `json:"-"`
	UserID     uint   `gorm:"not null;index"` // Added index for better query performance
	Title      string `gorm:"size:100"`       // Added reasonable size limit
	Content    string `gorm:"type:text;not null"`
	Mood       string `gorm:"size:20;index"` // Added index if you'll query by mood
	Tags       string `gorm:"size:255"`      // Kept as string but consider adding index
	IsPrivate  bool   `gorm:"default:true;not null"`
}
type EntryIn struct {
	UserID    uint   `gorm:"not null;index"`
	Title     string `json:"title" binding:"required,max=100"`
	Content   string `json:"content" binding:"required"`
	Mood      string `json:"mood" binding:"max=20"`
	Tags      string `json:"tags" binding:"max=255"`
	IsPrivate bool   `json:"is_private"`
}
