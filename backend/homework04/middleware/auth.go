package middleware

import (
	"net/http"

	"strings"

	"homework04/utils"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
		if jwtToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		if claims, err := utils.ParseToken(jwtToken, jwtSecret); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		} else {
			c.Set("userId", claims.UserID)
		}
		c.Next()
	}

}
