package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"merch-shop/pkg/jwtAuth"
	"net/http"
	"strings"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"}) //TODO Неавторизован.
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"}) //TODO Неавторизован.
			c.Abort()
			return
		}

		tokenString := parts[1]

		userID, username, err := jwtAuth.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)}) //TODO Неавторизован.
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Next()
	}
}
