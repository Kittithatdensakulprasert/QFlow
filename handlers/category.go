package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"qflow/models"
)

var categories []models.Category

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newCategory models.Category
	json.NewDecoder(r.Body).Decode(&newCategory)

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCategory)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	json.NewEncoder(w).Encode(categories)
}