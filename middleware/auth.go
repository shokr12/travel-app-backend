package middleware

import (
	"Visa/pkg"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware() gin.HandlerFunc {
	godotenv.Load()
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

		// ✅ Store values in context (IMPORTANT)
		c.Set("userId", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(jwt.MapClaims)
		if user["role"] != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthinticateUser(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(jwt.MapClaims)
		if user["role"] != role {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
