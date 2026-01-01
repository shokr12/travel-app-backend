package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		Method := c.Request.Method
		Path := c.Request.URL.Path
		clientIP := c.ClientIP()

		c.Next()
		statusCode := c.Writer.Status()

		fmt.Printf("Method: %s, Path: %s, ClientIP: %s, StatusCode: %d\n", Method, Path, clientIP, statusCode)

	}
}
