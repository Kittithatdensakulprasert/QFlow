package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Notification struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
}

var notifications []Notification

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(notifications)
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	parts := strings.Split(path, "/")

	// PATCH /api/notifications/:id/read
	if r.Method == http.MethodPatch && len(parts) == 2 && parts[1] == "read" {
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		for i, n := range notifications {
			if n.ID == id {
				notifications[i].Read = true
				json.NewEncoder(w).Encode(notifications[i])
				return
			}
		}
		http.NotFound(w, r)
		return
	}

	// DELETE /api/notifications/:id
	if r.Method == http.MethodDelete && len(parts) == 1 && parts[0] != "" {
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		for i, n := range notifications {
			if n.ID == id {
				notifications = append(notifications[:i], notifications[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.NotFound(w, r)
		return
	}

	http.NotFound(w, r)
}

func SendNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	n := Notification{
		ID:      len(notifications) + 1,
		Message: body.Message,
		Read:    false,
	}
	notifications = append(notifications, n)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}
