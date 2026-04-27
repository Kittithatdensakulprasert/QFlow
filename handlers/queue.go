package handlers

import (
	"net/http"
	"qflow/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var queues []models.Queue

// POST /api/queues/book
func BookQueue(c *gin.Context) {
	var req models.CreateQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	q := models.Queue{
		ID:          len(queues) + 1,
		QueueNumber: len(queues) + 1,
		Status:      "waiting",
	}
	queues = append(queues, q)

	c.JSON(http.StatusCreated, q)
}

// GET /api/queues/history
func GetHistory(c *gin.Context) {
	c.JSON(http.StatusOK, queues)
}

// GET /api/queues/:queueNumber
func GetQueue(c *gin.Context) {
	num, err := strconv.Atoi(c.Param("queueNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid queue number"})
		return
	}

	for _, q := range queues {
		if q.QueueNumber == num {
			c.JSON(http.StatusOK, q)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "queue not found"})
}

// PATCH /api/queues/:id/cancel
func CancelQueue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	for i := range queues {
		if queues[i].ID == id {
			if queues[i].Status == "cancelled" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "queue already cancelled"})
				return
			}
			queues[i].Status = "cancelled"
			c.JSON(http.StatusOK, queues[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "queue not found"})
}
