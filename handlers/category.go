package handlers

import (
	"encoding/json"
	"net/http"
	"qflow/models"
)

var categories []models.Category

func GetCategories(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(categories)
}
