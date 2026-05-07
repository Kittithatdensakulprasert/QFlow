package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// resolveContextUserID extracts the authenticated user's ID from the Gin context.
// It handles multiple type representations that may be set by different JWT middleware implementations.
// Returns (userID, true) on success, or (0, false) if the value is missing or invalid.
func resolveContextUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userIDVal.(type) {
	case uint:
		return v, v > 0
	case int:
		if v > 0 {
			return uint(v), true
		}
	case float64:
		if v > 0 {
			return uint(v), true
		}
	case string:
		uid, err := strconv.ParseUint(v, 10, 64)
		if err == nil && uid > 0 {
			return uint(uid), true
		}
	}

	return 0, false
}

// parseUintParam parses a named URL parameter as uint.
// Writes a 400 Bad Request response and returns (0, false) if parsing fails.
func parseUintParam(c *gin.Context, name string, errorMessage string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return 0, false
	}
	return uint(value), true
}

// parseID parses a URL parameter string into a uint ID.
// Returns an error if the value is not a valid positive integer.
func parseID(idParam string) (uint, error) {
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, strconv.ErrRange
	}
	return uint(id), nil
}

// respondError writes a structured JSON error response with an error code.
func respondError(c *gin.Context, status int, errorCode string, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"error":   errorCode,
		"message": message,
		"path":    c.Request.URL.Path,
	})
}
