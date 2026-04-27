package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"qflow/models"
)

var categories []models.Category

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	// GET /api/categories
	case http.MethodGet:
		json.NewEncoder(w).Encode(categories)

	// POST /api/categories
	case http.MethodPost:
		var newCategory models.Category

		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if newCategory.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		newCategory.ID = len(categories) + 1
		categories = append(categories, newCategory)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCategory)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}