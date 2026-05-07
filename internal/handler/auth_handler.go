package handler

import (
	"net/http"
	"qflow/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService       domain.AuthService
	exposeOTPResponse bool
}

func NewAuthHandler(authService domain.AuthService, exposeOTPResponse bool) *AuthHandler {
	return &AuthHandler{authService: authService, exposeOTPResponse: exposeOTPResponse}
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request format")
		return
	}

	otp, err := h.authService.RequestOTP(req.Phone)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "OTP_SEND_FAILED", "failed to send OTP")
		return
	}

	response := gin.H{
		"message": "OTP sent successfully",
		"otp_id":  otp.ID,
	}
	if h.exposeOTPResponse {
		response["otp_code"] = otp.Code
	}
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request format")
		return
	}

	user, token, err := h.authService.VerifyOTP(req.Phone, req.Code)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "OTP_INVALID", err.Error())
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
		Phone   string `json:"phone" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Role    string `json:"role"`
		OTPCode string `json:"otp_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request format")
		return
	}

	user, token, err := h.authService.RegisterUser(req.Phone, req.Name, req.Role, req.OTPCode)
	if err != nil {
		respondError(c, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	user, err := h.authService.GetUserProfile(userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var req struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request format")
		return
	}

	user, err := h.authService.UpdateUserProfile(userID, req.Name, req.Role)
	if err != nil {
		respondError(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}
