package handler

import (
	"fmt"
	"net/http"
	"sanjay/api/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT from the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("run middlewere")

		authHeader := c.GetHeader("Authorization")
		fmt.Println("authHeader", authHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		tokenString := parts[1]
		claims, err := service.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		phone, ok := claims["phone"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token payload"})
			c.Abort()
			return
		}

		c.Set("phone", phone)
		c.Next()
	}
}
func QueryAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token in query"})
			c.Abort()
			return
		}
		claims, err := service.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		phone, ok := claims["phone"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token payload"})
			c.Abort()
			return
		}

		c.Set("phone", phone)
		c.Next()
	}
}
