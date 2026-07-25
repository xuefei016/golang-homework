package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Email     string `gorm:"size:255;uniqueIndex"`
	Age       uint   `gorm:"default:18"`
	Posts     []Post
	PostCount int `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"size:255"`
	Content       string `gorm:"type:text"`
	UserID        uint
	Comments      []Comment
	CommentStatus string `gorm:"size:32;default:'pending'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text"`
	PostID    uint
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
