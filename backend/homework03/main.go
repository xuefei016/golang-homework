package main

import (
	"fmt"
	"homework03/models"
	"homework03/query"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getConnect() (*gorm.DB, error) {
	dsn := os.Getenv("BLOG_MYSQL_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("BLOG_MYSQL_DSN is not set")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	return db, nil
}

func initData(db *gorm.DB) ([]models.User, error) {
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	users := []models.User{
		{
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			Age:       30,
			PostCount: 0,
			Posts: []models.Post{
				{
					Title:         "First Post",
					Content:       "This is the content of the first post.",
					CommentStatus: "pending",
					Comments: []models.Comment{
						{Content: "写得不错", UserID: 2},
						{Content: "学到了", UserID: 3},
					},
				},
				{
					Title:         "Second Post",
					Content:       "This is the content of the second post.",
					CommentStatus: "pending",
					Comments: []models.Comment{
						{Content: "channel 讲得好", UserID: 2},
					},
				},
			}}, {
			Name:  "Alice",
			Email: "alice@example.com",
			Posts: []models.Post{
				{
					Title:   "GORM 入门",
					Content: "这是第一篇文章",
					Comments: []models.Comment{
						{Content: "垃圾", UserID: 2},
						{Content: "太辣鸡", UserID: 3},
						{Content: "啥玩意啊", UserID: 3},
					},
				},
				{
					Title:   "Go 并发",
					Content: "聊聊 goroutine",
					Comments: []models.Comment{
						{Content: "channel 讲得好", UserID: 2},
					},
				},
			},
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			fmt.Println("Failed to create user:", err)
			return users, err
		}
	}
	return users, nil
}

func main() {
	db, err := getConnect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		fmt.Println("Failed to migrate database:", err)
		return
	}
	fmt.Println("Database migrated successfully")

	userList, err := initData(db)
	if err != nil {
		fmt.Println("Failed to initialize data:", err)
		return
	}
	fmt.Println("Data initialized successfully")

	user, err := query.QueryUsersWithPostsAndComments(db, userList[0].ID)
	if err != nil {
		fmt.Println("Failed to query user with posts and comments:", err)
		return
	}
	fmt.Printf("Queried user: %s(ID:%d), Posts: %d\n", user.Name, user.ID, len(user.Posts))

	post, err := query.QueryMostCommentedPost(db)
	if err != nil {
		fmt.Println("Failed to query most commented post:", err)
		return
	}
	fmt.Printf("Most commented post: %s(ID:%d)\n", post.Title, post.ID)
	for _, comment := range post.Comments {
		fmt.Printf("\n\t评论: %s(ID:%d)", comment.Content, comment.ID)
	}

	for comment := range post.Comments {
		if err := query.DeleteCommentById(db, post.Comments[comment].ID); err != nil {
			fmt.Println("Failed to delete comment:", err)
			return
		}
	}
	fmt.Println("All comments for the post deleted successfully")

	queriedPost, err := query.QueryPostById(db, post.ID)
	if err != nil {
		fmt.Println("Failed to query post by ID:", err)
		return
	}
	fmt.Printf("Queried post: %s(ID:%d), Comments: %d, CommentStatus:%s \n", queriedPost.Title, queriedPost.ID, len(queriedPost.Comments), queriedPost.CommentStatus)
}
