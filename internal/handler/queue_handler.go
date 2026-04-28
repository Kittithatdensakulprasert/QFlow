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
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := resolveRequiredUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
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
	userID, ok := resolveRequiredUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
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

	userID, ok := resolveRequiredUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	queue, err := h.svc.GetQueueByNumber(queueNumber, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrQueueNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrForbiddenQueue):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

	userID, ok := resolveRequiredUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
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

func resolveRequiredUserID(c *gin.Context) (uint, bool) {
	if userIDVal, exists := c.Get("user_id"); exists {
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
			userID, err := strconv.ParseUint(v, 10, 64)
			if err == nil && userID > 0 {
				return uint(userID), true
			}
		}
	}
	return 0, false
}
