package handlers

import (
	"net/http"

	"homework04/models"

	"homework04/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db        *gorm.DB
	jwtSecret string
}

func NewUserHandler(db *gorm.DB, jwtSecret string) *UserHandler {
	return &UserHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	if err := h.db.Create(&user).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			utils.Error(c, http.StatusConflict, "User already exists")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}
	utils.Success(c, gin.H{"message": "User registered successfully"})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := h.db.Where("username=?", req.Username).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	utils.Success(c, gin.H{"token": token})
}
