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
		Name:        "สุกี้ตี๋น้อย",
		Description: "ร้านอาหารสุกี้ยากี้ ชาบู",
	},
	{
		ID:          2,
		Name:        "Sushiro",
		Description: "ร้านซูชิสายพาน",
	},
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(categories)
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
