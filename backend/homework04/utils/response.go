package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 返回 200 + 统一成功结构：{"code":0,"data":...}
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

// Created 返回 201 + 统一成功结构，用于资源创建
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"code": 0, "data": data})
}

// Error 返回指定状态码 + 统一错误结构：{"code":<status>,"message":...}
func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"code": status, "message": msg})
}
