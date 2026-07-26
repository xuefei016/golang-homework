package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:64;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:128;uniqueIndex;not null"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"size:255;not null" binding:"required"`
	Content string `gorm:"type:text;not null" binding:"required"`
	UserID  uint   `gorm:"not null"`
	User    User   `json:"-"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" binding:"required"`
	PostID  uint   `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	User    User   `json:"-"`
	Post    Post   `json:"-"`
}
