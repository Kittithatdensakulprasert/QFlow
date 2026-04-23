package handlers

import (
	"encoding/json"
	"net/http"
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
	json.NewEncoder(w).Encode(map[string]string{
		"message": "not implemented yet",
	})
}
