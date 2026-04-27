package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
    "strconv"
	"qflow/models"
)

var categories = map[int]models.Category{}
var categoryOrder = []int{}
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

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := parseCategoryID(r)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		handleGetCategory(w, id)

	case http.MethodPut:
		handleUpdateCategory(w, r, id)

	case http.MethodDelete:
		handleDeleteCategory(w, id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /api/categories
func getCategories(w http.ResponseWriter) {
	var result []models.Category

	for _, id := range categoryOrder {
		if c, ok := categories[id]; ok {
			result = append(result, c)
		}
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

	if isDuplicateCategoryName(newCategory.Name) {
		http.Error(w, "Category already exists", http.StatusConflict)
		return
	}

	newCategory.ID = nextCategoryID
	nextCategoryID++

	categories[newCategory.ID] = newCategory

	categoryOrder = append(categoryOrder, newCategory.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func isDuplicateCategoryName(name string) bool {
	for _, c := range categories {
		if strings.EqualFold(c.Name, name) {
			return true
		}
	}
	return false
}

func getCategory(id int) (models.Category, bool) {
	c, ok := categories[id]
	return c, ok
}

// GET /api/categories/{id}
func handleGetCategory(w http.ResponseWriter, id int) {
	category, exists := getCategory(id)

	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(category)
}

// PUT /api/categories/{id}
func handleUpdateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var update models.Category

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if update.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	category, exists := getCategory(id)
	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category.Name = update.Name
	categories[id] = category

	json.NewEncoder(w).Encode(category)
}

// DELETE /api/categories/{id}
func handleDeleteCategory(w http.ResponseWriter, id int) {
    _, exists := categories[id]
    if !exists {
        http.Error(w, "Category not found", http.StatusNotFound)
        return
    }

    delete(categories, id)

    for i, v := range categoryOrder {
        if v == id {
            categoryOrder = append(categoryOrder[:i], categoryOrder[i+1:]...)
            break
        }
    }

    w.WriteHeader(http.StatusNoContent)
}

func parseCategoryID(r *http.Request) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	return strconv.Atoi(idStr)
}