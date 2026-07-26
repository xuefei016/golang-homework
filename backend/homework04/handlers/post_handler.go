package handlers

import (
	"net/http"

	"homework04/models"
	"homework04/utils"

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
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetUint("userId")

	newPost := models.Post{UserID: userId, Title: post.Title, Content: post.Content}
	if err := h.db.Create(&newPost).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Created(ctx, newPost)
}

func (h *PostHandler) ListPosts(ctx *gin.Context) {
	var posts []models.Post
	if err := h.db.Find(&posts).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, posts)
}

func (h *PostHandler) GetPost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post
	if err := h.db.Where("id = ?", id).First(&post).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "post not found")
		return
	}

	utils.Success(ctx, post)
}

func (h *PostHandler) UpdatePost(ctx *gin.Context) {
	var input models.Post
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetUint("userId")
	var post models.Post
	id := ctx.Param("id")
	if err := h.db.Where("id = ?", id).First(&post).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "post not found")
		return
	}

	if post.UserID != userId {
		utils.Error(ctx, http.StatusForbidden, "forbidden")
		return
	}

	post.Title = input.Title
	post.Content = input.Content
	if err := h.db.Save(&post).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, post)
}

func (h *PostHandler) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.GetUint("userId")

	var post models.Post
	if err := h.db.Where("id = ?", id).First(&post).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "post not found")
		return
	}

	if post.UserID != userId {
		utils.Error(ctx, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.db.Delete(&models.Post{}, id).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "post deleted"})
}
