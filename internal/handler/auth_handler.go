package handler

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"os"
	"qflow/db"
	"qflow/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var bigTen = big.NewInt(10)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var body struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := generateOTPCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate otp"})
		return
	}

	otp := domain.OTP{
		Phone:     body.Phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}
	if err := db.DB.Create(&otp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create otp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "otp sent",
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var body struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var otp domain.OTP
	err := db.DB.
		Where("phone = ? AND code = ? AND used = ? AND expires_at > ?", body.Phone, body.Code, false, time.Now()).
		Order("created_at desc").
		First(&otp).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired otp"})
		return
	}

	if err := db.DB.Model(&otp).Where("used = ?", false).Update("used", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update otp state"})
		return
	}

	var user domain.User
	err = db.DB.Where("phone = ?", body.Phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = domain.User{
				Phone: body.Phone,
				Role:  "user",
			}
			if err := db.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			return
		}
	}

	token, err := issueJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Phone string `json:"phone" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user domain.User
	err := db.DB.Where("phone = ?", body.Phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = domain.User{
				Phone: body.Phone,
				Name:  body.Name,
				Role:  "user",
			}
			if err := db.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "phone already registered"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, ok := resolveAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user domain.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := resolveAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user domain.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Name = body.Name
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func issueJWT(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" || secret == "secret" {
		return "", errors.New("jwt secret is not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func resolveAuthUserID(c *gin.Context) (uint, bool) {
	if v, exists := c.Get("user_id"); exists {
		switch id := v.(type) {
		case uint:
			return id, true
		case int:
			if id > 0 {
				return uint(id), true
			}
		case float64:
			if id > 0 {
				return uint(id), true
			}
		}
	}

	return 0, false
}

func generateOTPCode() (string, error) {
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, bigTen)
		if err != nil {
			return "", err
		}
		code[i] = byte('0' + n.Int64())
	}
	return string(code), nil
}
