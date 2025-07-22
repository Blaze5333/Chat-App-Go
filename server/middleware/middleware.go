package middleware

import (
	"chat-server/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "No Authorization token provided"})
			c.Abort()
			return
		}
		token = token[len("Bearer "):] // Remove "Bearer " prefix
		claims, err := tokens.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Next()
	}
}
