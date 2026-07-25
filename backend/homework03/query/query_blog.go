package query

import (
	"fmt"

	"homework03/models"

	"gorm.io/gorm"
)

type PostCommentCount struct {
	PostID       uint
	CommentCount int
}

func QueryUsersWithPostsAndComments(db *gorm.DB, userId uint) (models.User, error) {
	var user models.User
	err := db.Preload("Posts.Comments").First(&user, userId).Error
	if err != nil {
		return models.User{}, err
	}

	fmt.Printf("用户: %s(ID:%d), 发布%d篇文章", user.Name, user.ID, len(user.Posts))

	for _, post := range user.Posts {
		fmt.Printf("\n文章: %s(ID:%d), 评论%d条", post.Title, post.ID, len(post.Comments))
		for _, comment := range post.Comments {
			fmt.Printf("\n\t评论: %s(ID:%d)", comment.Content, comment.ID)
		}
	}

	return user, nil
}

func QueryMostCommentedPost(db *gorm.DB) (models.Post, error) {
	var mostCommented PostCommentCount
	err := db.Model(&models.Comment{}).Select("post_id, count(*) as comment_count").Group("post_id").Order("comment_count DESC").Limit(1).Scan(&mostCommented).Error
	if err != nil {
		return models.Post{}, err
	}

	var post models.Post
	err = db.Preload("Comments").First(&post, mostCommented.PostID).Error
	if err != nil {
		return models.Post{}, err
	}
	fmt.Printf("评论最多的文章: %s(ID:%d), 共评论了:%d条\n", post.Title, post.ID, mostCommented.CommentCount)

	return post, nil
}

func DeleteAllCommentByPostId(db *gorm.DB, postId uint) error {
	return db.Where("post_id = ?", postId).Delete(&models.Comment{}).Error
}

func DeleteCommentById(db *gorm.DB, commentId uint) error {
	var comment models.Comment
	err := db.First(&comment, commentId).Error
	if err != nil {
		return err
	}

	return db.Delete(&comment).Error
}

func QueryPostById(db *gorm.DB, postId uint) (models.Post, error) {
	var post models.Post
	err := db.Preload("Comments").First(&post, postId).Error
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}
