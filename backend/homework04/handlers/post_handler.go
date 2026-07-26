package handlers

import (
	"homework04/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db: db}
}

func (h *PostHandler) CreatePost(ctx *gin.Context) {
	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := ctx.GetUint("userId")

	newPost := models.Post{UserID: userId, Title: post.Title, Content: post.Content}
	if err := h.db.Create(&newPost).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, newPost)
}

func (h *PostHandler) ListPosts(ctx *gin.Context) {
	var posts []models.Post
	if err := h.db.Find(&posts).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post
	if err := h.db.Where("id = ?", id).First(&post).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(ctx *gin.Context) {
	var input models.Post
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := ctx.GetUint("userId")
	var post models.Post
	id := ctx.Param("id")
	if err := h.db.Where("id = ?", id).First(&post).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if post.UserID != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	post.Title = input.Title
	post.Content = input.Content
	if err := h.db.Save(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.GetUint("userId")

	var post models.Post
	if err := h.db.Where("id = ?", id).First(&post).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if post.UserID != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.db.Delete(&models.Post{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
