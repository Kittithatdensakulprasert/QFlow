package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func ManageHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/manage/queues/")
	parts := strings.Split(path, "/")

	id, _ := strconv.Atoi(parts[0])

	for i := range queues {

		if queues[i].ID == id {

			if parts[1] == "call" {
				queues[i].Status = "called"
			} else if parts[1] == "complete" {
				queues[i].Status = "completed"
			} else if parts[1] == "skip" {
				queues[i].Status = "skipped"
			}

			json.NewEncoder(w).Encode(queues[i])
			return
		}
	}

	http.NotFound(w, r)
}
