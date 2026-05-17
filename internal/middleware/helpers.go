package middleware

import "github.com/gin-gonic/gin"

// respondError writes the same structured error shape used by HTTP handlers.
func respondError(c *gin.Context, status int, errorCode string, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"error":   errorCode,
		"message": message,
		"path":    c.Request.URL.Path,
	})
}
