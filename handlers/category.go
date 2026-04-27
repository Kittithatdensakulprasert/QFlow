package handlers

import (
	"encoding/json"
	"net/http"
	"qflow/models"
)

var categories = map[int]models.Category{}
var nextCategoryID = 1

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		getCategories(w)

	case http.MethodPost:
		createCategory(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


// GET /api/categories
func getCategories(w http.ResponseWriter) {
	var result []models.Category

	for _, c := range categories {
		result = append(result, c)
	}

	json.NewEncoder(w).Encode(result)
}

// POST /api/categories
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category

	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newCategory.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	newCategory.ID = nextCategoryID
	nextCategoryID++

	categories[newCategory.ID] = newCategory

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}