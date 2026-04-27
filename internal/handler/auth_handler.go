package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"qflow/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	otp, err := h.authService.RequestOTP(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
		"otp_id":  otp.ID,
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, token, err := h.authService.VerifyOTP(req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"user":    user,
		"token":   token,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, token, err := h.authService.RegisterUser(req.Phone, req.Name, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Try to get user ID from context (JWT middleware) first
	userIDInterface, exists := c.Get("user_id")
	var userIDStr string

	if exists {
		// From context (JWT middleware)
		switch v := userIDInterface.(type) {
		case uint:
			userIDStr = strconv.FormatUint(uint64(v), 10)
		case int:
			userIDStr = strconv.Itoa(v)
		case string:
			userIDStr = v
		case float64:
			userIDStr = strconv.FormatUint(uint64(v), 10)
		default:
			userIDStr = fmt.Sprintf("%v", v)
		}
	} else {
		// Fallback: try header first, then query parameter
		userIDStr = c.GetHeader("X-User-ID")
		if userIDStr == "" {
			userIDStr = c.Query("user_id")
		}
	}

	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetUserProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Try to get user ID from context (JWT middleware) first
	userIDInterface, exists := c.Get("user_id")
	var userIDStr string

	if exists {
		// From context (JWT middleware)
		switch v := userIDInterface.(type) {
		case uint:
			userIDStr = strconv.FormatUint(uint64(v), 10)
		case int:
			userIDStr = strconv.Itoa(v)
		case string:
			userIDStr = v
		case float64:
			userIDStr = strconv.FormatUint(uint64(v), 10)
		default:
			userIDStr = fmt.Sprintf("%v", v)
		}
	} else {
		// Fallback: try header first, then query parameter
		userIDStr = c.GetHeader("X-User-ID")
		if userIDStr == "" {
			userIDStr = c.Query("user_id")
		}
	}

	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.authService.UpdateUserProfile(uint(userID), req.Name, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}
