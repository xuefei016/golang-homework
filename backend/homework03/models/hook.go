package models

import "gorm.io/gorm"

func (p *Post) AfterCreate(db *gorm.DB) error {
	return db.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

func (c *Comment) AfterDelete(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Comment{}).Where("post_id=?", c.PostID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return db.Model(&Post{}).Where("id=?", c.PostID).UpdateColumn("comment_status", "no_comments").Error
	}
	return nil
}
