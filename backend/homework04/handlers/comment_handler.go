package handlers

import (
	"homework04/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

func (h *CommentHandler) ListComments(ctx *gin.Context) {
	postId := ctx.Param("id")
	var comments []models.Comment
	if err := h.db.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	postId := ctx.Param("id")
	userId := ctx.GetUint("userId")

	var req models.Comment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	if err := h.db.First(&post, postId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment := models.Comment{
		Content: req.Content,
		PostID:  post.ID,
		UserID:  userId,
	}

	if err := h.db.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, comment)

}
