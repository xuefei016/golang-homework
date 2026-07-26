package handlers

import (
	"net/http"

	"homework04/models"
	"homework04/utils"

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
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, comments)
}

func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	postId := ctx.Param("id")
	userId := ctx.GetUint("userId")

	var req models.Comment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var post models.Post
	if err := h.db.First(&post, postId).Error; err != nil {
		utils.Error(ctx, http.StatusNotFound, "Post not found")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		PostID:  post.ID,
		UserID:  userId,
	}

	if err := h.db.Create(&comment).Error; err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Created(ctx, comment)
}
