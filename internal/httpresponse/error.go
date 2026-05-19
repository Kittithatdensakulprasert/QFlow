package httpresponse

import "github.com/gin-gonic/gin"

// Error writes a structured JSON error response.
func Error(c *gin.Context, status int, errorCode string, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"error":   errorCode,
		"message": message,
		"path":    c.Request.URL.Path,
	})
}
