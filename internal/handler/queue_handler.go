package handler

import (
	"errors"
	"net/http"
	"qflow/internal/domain"
	"qflow/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	svc domain.QueueService
}

func NewQueueHandler(svc domain.QueueService) *QueueHandler {
	return &QueueHandler{svc: svc}
}

func (h *QueueHandler) BookQueue(c *gin.Context) {
	var body struct {
		ZoneID uint `json:"zone_id" binding:"required"`
		UserID uint `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := resolveUserID(c, body.UserID)
	queue, err := h.svc.BookQueue(userID, body.ZoneID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID),
			errors.Is(err, service.ErrInvalidZoneID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrZoneNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrZoneClosed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           queue.ID,
		"queue_number": queue.QueueNumber,
		"zone_id":      queue.ZoneID,
		"status":       queue.Status,
		"created_at":   queue.CreatedAt,
	})
}

func (h *QueueHandler) GetHistory(c *gin.Context) {
	userID := resolveUserID(c, parseOptionalUint(c.Query("user_id")))
	queues, err := h.svc.GetQueueHistory(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, queues)
}

func (h *QueueHandler) GetQueue(c *gin.Context) {
	queueNumber, err := strconv.Atoi(c.Param("queueNumber"))
	if err != nil || queueNumber <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue number"})
		return
	}

	queue, err := h.svc.GetQueueByNumber(queueNumber)
	if err != nil {
		if errors.Is(err, service.ErrQueueNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	requestUserID := resolveUserID(c, parseOptionalUint(c.Query("user_id")))
	if requestUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if queue.UserID != requestUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "queue does not belong to user"})
		return
	}
	c.JSON(http.StatusOK, queue)
}

func (h *QueueHandler) CancelQueue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := resolveUserID(c, parseOptionalUint(c.Query("user_id")))
	err = h.svc.CancelQueue(uint(id), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrQueueNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrForbiddenQueue):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrQueueFinalized),
			errors.Is(err, service.ErrQueueCancelled):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "queue cancelled"})
}

func (h *QueueHandler) GetQueuesByZone(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *QueueHandler) CallQueue(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *QueueHandler) CompleteQueue(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *QueueHandler) SkipQueue(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func resolveUserID(c *gin.Context, fallback uint) uint {
	if userIDVal, exists := c.Get("user_id"); exists {
		switch v := userIDVal.(type) {
		case uint:
			return v
		case int:
			if v > 0 {
				return uint(v)
			}
		case float64:
			if v > 0 {
				return uint(v)
			}
		}
	}
	return fallback
}

func parseOptionalUint(value string) uint {
	if value == "" {
		return 0
	}
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}
	return uint(v)
}
