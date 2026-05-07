package handler

import (
	"errors"
	"net/http"
	"strconv"

	"qflow/internal/domain"
	"qflow/internal/service"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	svc domain.QueueService
}

func NewQueueHandler(svc domain.QueueService) *QueueHandler {
	return &QueueHandler{svc: svc}
}

// ===================== Queue Booking (User) =====================

// POST /api/queues/book
func (h *QueueHandler) BookQueue(c *gin.Context) {
	var body struct {
		ZoneID uint `json:"zone_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}
	queue, err := h.svc.BookQueue(userID, body.ZoneID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID),
			errors.Is(err, service.ErrInvalidZoneID):
			respondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		case errors.Is(err, service.ErrZoneNotFound):
			respondError(c, http.StatusNotFound, "ZONE_NOT_FOUND", err.Error())
		case errors.Is(err, service.ErrZoneClosed):
			respondError(c, http.StatusConflict, "ZONE_CLOSED", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
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

// GET /api/queues/history
func (h *QueueHandler) GetHistory(c *gin.Context) {
	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}
	queues, err := h.svc.GetQueueHistory(userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			respondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}
	c.JSON(http.StatusOK, queues)
}

// GET /api/queues/:queueNumber
func (h *QueueHandler) GetQueue(c *gin.Context) {
	queueNumber, err := strconv.Atoi(c.Param("queueNumber"))
	if err != nil || queueNumber <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_QUEUE_NUMBER", "invalid queue number")
		return
	}

	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}
	queue, err := h.svc.GetQueueByNumber(queueNumber, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			respondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		case errors.Is(err, service.ErrQueueNotFound):
			respondError(c, http.StatusNotFound, "QUEUE_NOT_FOUND", err.Error())
		case errors.Is(err, service.ErrForbiddenQueue):
			respondError(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, queue)
}

// PATCH /api/queues/:id/cancel
func (h *QueueHandler) CancelQueue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return
	}

	userID, ok := resolveContextUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}
	err = h.svc.CancelQueue(uint(id), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			respondError(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		case errors.Is(err, service.ErrQueueNotFound):
			respondError(c, http.StatusNotFound, "QUEUE_NOT_FOUND", err.Error())
		case errors.Is(err, service.ErrForbiddenQueue):
			respondError(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		case errors.Is(err, service.ErrQueueFinalized),
			errors.Is(err, service.ErrQueueCancelled):
			respondError(c, http.StatusConflict, "QUEUE_ALREADY_FINALIZED", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "queue cancelled"})
}

// ===================== Queue Management (Provider) =====================

// GET /api/manage/queues/:zoneId
func (h *QueueHandler) GetQueuesByZone(c *gin.Context) {
	zoneID, err := strconv.ParseUint(c.Param("zoneId"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ZONE_ID", "invalid zone id")
		return
	}

	queues, err := h.svc.GetQueuesByZone(uint(zoneID))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, queues)
}

// PATCH /api/manage/queues/:id/call
func (h *QueueHandler) CallQueue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return
	}

	queue, err := h.svc.CallQueue(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrQueueNotFound):
			respondError(c, http.StatusNotFound, "QUEUE_NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrQueueCannotBeCalled):
			respondError(c, http.StatusConflict, "QUEUE_INVALID_STATE", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, queue)
}

// PATCH /api/manage/queues/:id/complete
func (h *QueueHandler) CompleteQueue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return
	}

	queue, err := h.svc.CompleteQueue(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrQueueNotFound):
			respondError(c, http.StatusNotFound, "QUEUE_NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrQueueCannotBeCompleted):
			respondError(c, http.StatusConflict, "QUEUE_INVALID_STATE", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, queue)
}

// PATCH /api/manage/queues/:id/skip
func (h *QueueHandler) SkipQueue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "invalid id")
		return
	}

	queue, err := h.svc.SkipQueue(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrQueueNotFound):
			respondError(c, http.StatusNotFound, "QUEUE_NOT_FOUND", err.Error())
		case errors.Is(err, domain.ErrQueueCannotBeSkipped):
			respondError(c, http.StatusConflict, "QUEUE_INVALID_STATE", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, queue)
}
