package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email    string `gorm:"unique;not null" json:"email"`
    Pass     string `gorm:"not null" json:"pass"`
    Entries  []Entry `gorm:"foreignKey:UserID" json:"entries,omitempty"`
}
