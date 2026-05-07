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
	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	notifications, err := h.svc.GetNotifications(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var body struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	n, err := h.svc.SendNotification(userID, body.Message)
	if err != nil {
		respondError(c, http.StatusBadRequest, "SEND_FAILED", err.Error())
		return
	}
	c.JSON(http.StatusCreated, n)
}

func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return
	}

	if err := h.svc.MarkNotificationRead(uint(id), userID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotificationNotFound):
			respondError(c, http.StatusNotFound, "NOTIFICATION_NOT_FOUND", err.Error())
		case errors.Is(err, service.ErrNotificationForbidden):
			respondError(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return
	}

	if err := h.svc.DeleteNotification(uint(id), userID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotificationNotFound):
			respondError(c, http.StatusNotFound, "NOTIFICATION_NOT_FOUND", err.Error())
		case errors.Is(err, service.ErrNotificationForbidden):
			respondError(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
