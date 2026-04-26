package handlers

import (
	"encoding/json"
	"net/http"
	"qflow/models"
	"strconv"
	"strings"
)

var categories = []models.Category{
	{
		ID:          1,
		Name:        "ชาบู",
	},
	{
		ID:          2,
		Name:        "ซูชิ",
	},
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(categories)
		return
	}

	if r.Method == http.MethodPost {
		var newCategory models.Category

		json.NewDecoder(r.Body).Decode(&newCategory)

		newCategory.ID = len(categories) + 1

		categories = append(categories, newCategory)

		json.NewEncoder(w).Encode(newCategory)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(path)

	for _, category := range categories {
		if category.ID == id {
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.NotFound(w, r)
}
