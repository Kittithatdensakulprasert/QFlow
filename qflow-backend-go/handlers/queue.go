package handlers

import (
	"encoding/json"
	"net/http"
	"qflow/models"
	"strconv"
	"strings"
)

var queues []models.Queue

func BookQueue(w http.ResponseWriter, r *http.Request) {
	q := models.Queue{
		ID:          len(queues) + 1,
		QueueNumber: len(queues) + 1,
		Status:      "waiting",
	}
	queues = append(queues, q)

	json.NewEncoder(w).Encode(q)
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(queues)
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/queues/")
	parts := strings.Split(path, "/")

	// GET /api/queues/:queueNumber
	if r.Method == http.MethodGet {
		num, _ := strconv.Atoi(parts[0])
		for _, q := range queues {
			if q.QueueNumber == num {
				json.NewEncoder(w).Encode(q)
				return
			}
		}
	}

	// PATCH /api/queues/:id/cancel
	if len(parts) == 2 && parts[1] == "cancel" {
		id, _ := strconv.Atoi(parts[0])
		for i := range queues {
			if queues[i].ID == id {
				queues[i].Status = "cancelled"
				json.NewEncoder(w).Encode(queues[i])
				return
			}
		}
	}

	http.NotFound(w, r)
}
