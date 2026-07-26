package utils

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any)           { c.JSON(200, gin.H{"code": 0, "data": data}) }
func Error(c *gin.Context, code int, msg string) { c.JSON(code, gin.H{"code": code, "message": msg}) }
