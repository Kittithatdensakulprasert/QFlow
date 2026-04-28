package handler

import (
	"errors"
	"net/http"
	"qflow/internal/domain"
	"qflow/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	svc domain.NotificationService
}

func NewNotificationHandler(svc domain.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, ok := resolveNotificationUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	notifications, err := h.svc.GetNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var body struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n, err := h.svc.SendNotification(body.UserID, body.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, n)
}

func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	userID, ok := resolveNotificationUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.MarkNotificationRead(uint(id), userID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotificationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrNotificationForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, ok := resolveNotificationUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteNotification(uint(id), userID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotificationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrNotificationForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func resolveNotificationUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	switch v := userIDVal.(type) {
	case uint:
		return v, true
	case int:
		if v > 0 {
			return uint(v), true
		}
	case float64:
		if v > 0 {
			return uint(v), true
		}
	}
	return 0, false
}
