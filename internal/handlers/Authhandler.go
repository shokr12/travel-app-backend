// handlers/middleware.go
package handlers

import (
	"Visa/pkg"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &pkg.JWTClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// 🔑 Store values in context
		c.Set("userId", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// AdminMiddleware ensures only admin users can access the endpoint
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
