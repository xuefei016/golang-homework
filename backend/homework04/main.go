package main

import (
	"homework04/config"
	"homework04/handlers"
	"homework04/middleware"
	"homework04/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	if cfg.DSN == "" {
		log.Fatal("BLOG_MYSQL_DSN is not set")
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	userHandler := handlers.NewUserHandler(db, cfg.JWTSecret)

	user_api := r.Group("/api/user")
	{
		// Here you would register your user-related routes, for example:
		user_api.POST("/register", userHandler.Register)
		user_api.POST("/login", userHandler.Login)
	}

	auth := r.Group("/api")
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	auth.Use(middleware.Auth(cfg.JWTSecret))
	{
		auth.POST("/posts", postHandler.CreatePost)
		auth.PUT("/posts/:id", postHandler.UpdatePost)
		auth.DELETE("/posts/:id", postHandler.DeletePost)
		auth.POST("/posts/:id/comments", commentHandler.CreateComment)
	}

	post_api := r.Group("/api")
	{
		post_api.GET("/posts", postHandler.ListPosts)
		post_api.GET("/posts/:id", postHandler.GetPost)
		post_api.GET("/posts/:id/comments", commentHandler.ListComments)

	}

	log.Printf("Starting server on : %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
