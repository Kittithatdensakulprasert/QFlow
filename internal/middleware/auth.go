package middleware

import (
	"github.com/gin-gonic/gin"
)

// JWTAuth validates the JWT token in the Authorization header.
// TODO: implement JWT validation
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
