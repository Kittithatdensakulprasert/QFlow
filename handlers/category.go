package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"qflow/models"
)

var categories = map[int]models.Category{}
var nextCategoryID = 1

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	// GET /api/categories
	case http.MethodGet:
        var result []models.Category
        
        for _, c := range categories {
            result = append(result, c)
        }

	json.NewEncoder(w).Encode(result)

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
        
        for _, c := range categories {
            if strings.EqualFold(c.Name, newCategory.Name) {
                http.Error(w, "Category already exists", http.StatusConflict)
                return
            }
        }
        newCategory.ID = nextCategoryID
        nextCategoryID++
        categories[newCategory.ID] = newCategory
        
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(newCategory)
    
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
