package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockCategoryService struct {
	categories []domain.Category
	err        error
}

func (m *mockCategoryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.categories, nil
}

func (m *mockCategoryService) GetCategory(ctx context.Context, id uint) (*domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}

	for i := range m.categories {
		if m.categories[i].ID == id {
			return &m.categories[i], nil
		}
	}

	return nil, service.ErrCategoryNotFound
}

func (m *mockCategoryService) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}

	category := domain.Category{ID: uint(len(m.categories) + 1), Name: name}
	m.categories = append(m.categories, category)

	return &category, nil
}

func (m *mockCategoryService) UpdateCategory(ctx context.Context, id uint, name string) (*domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}

	for i := range m.categories {
		if m.categories[i].ID == id {
			m.categories[i].Name = name
			return &m.categories[i], nil
		}
	}

	return nil, service.ErrCategoryNotFound
}

func (m *mockCategoryService) DeleteCategory(ctx context.Context, id uint) error {
	if m.err != nil {
		return m.err
	}

	for i := range m.categories {
		if m.categories[i].ID == id {
			m.categories = append(m.categories[:i], m.categories[i+1:]...)
			return nil
		}
	}

	return service.ErrCategoryNotFound
}

func TestGetCategories(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.categories = append(svc.categories, domain.Category{ID: 1, Name: "ชาบู"})

	res := performCategoryRequest(router, http.MethodGet, "/api/categories", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var categories []domain.Category
	if err := json.NewDecoder(res.Body).Decode(&categories); err != nil {
		t.Fatalf("decode categories: %v", err)
	}

	if len(categories) != 1 || categories[0].Name != "ชาบู" {
		t.Fatalf("unexpected categories: %+v", categories)
	}
}

func TestGetCategory(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.categories = append(svc.categories, domain.Category{ID: 1, Name: "ชาบู"})

	res := performCategoryRequest(router, http.MethodGet, "/api/categories/1", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var category domain.Category
	if err := json.NewDecoder(res.Body).Decode(&category); err != nil {
		t.Fatalf("decode category: %v", err)
	}

	if category.ID != 1 || category.Name != "ชาบู" {
		t.Fatalf("unexpected category: %+v", category)
	}
}

func TestGetCategoryNotFound(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodGet, "/api/categories/99", "")

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestGetCategoryInternalServerError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = errors.New("db down")

	res := performCategoryRequest(router, http.MethodGet, "/api/categories/1", "")

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestGetCategoryRejectsInvalidID(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodGet, "/api/categories/0", "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestGetCategoryRejectsNonNumericID(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodGet, "/api/categories/abc", "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestCreateCategory(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodPost, "/api/categories", `{"name":"ชาบู"}`)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var category domain.Category
	if err := json.NewDecoder(res.Body).Decode(&category); err != nil {
		t.Fatalf("decode category: %v", err)
	}

	if category.ID != 1 || category.Name != "ชาบู" {
		t.Fatalf("unexpected category: %+v", category)
	}
}

func TestCreateCategoryValidationError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = service.ErrCategoryNameRequired

	res := performCategoryRequest(router, http.MethodPost, "/api/categories", `{"name":""}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestCreateCategoryInvalidJSON(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodPost, "/api/categories", `{`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestCreateCategoryDuplicate(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = service.ErrCategoryDuplicate

	res := performCategoryRequest(router, http.MethodPost, "/api/categories", `{"name":"ชาบู"}`)

	if res.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, res.Code)
	}
}

func TestCreateCategoryInternalServerError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = errors.New("db down")

	res := performCategoryRequest(router, http.MethodPost, "/api/categories", `{"name":"ชาบู"}`)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestUpdateCategory(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.categories = append(svc.categories, domain.Category{ID: 1, Name: "ชาบู"})

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/1", `{"name":"ปิ้งย่าง"}`)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var category domain.Category
	if err := json.NewDecoder(res.Body).Decode(&category); err != nil {
		t.Fatalf("decode category: %v", err)
	}

	if category.ID != 1 || category.Name != "ปิ้งย่าง" {
		t.Fatalf("unexpected category: %+v", category)
	}
}

func TestUpdateCategoryNotFound(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/99", `{"name":"ปิ้งย่าง"}`)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestUpdateCategoryRejectsInvalidID(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/0", `{"name":"ปิ้งย่าง"}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestUpdateCategoryRejectsNonNumericID(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/abc", `{"name":"ปิ้งย่าง"}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestUpdateCategoryInvalidJSON(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/1", `{`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestUpdateCategoryValidationError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = service.ErrCategoryNameRequired

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/1", `{"name":""}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestUpdateCategoryDuplicate(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = service.ErrCategoryDuplicate

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/1", `{"name":"ซ้ำ"}`)

	if res.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, res.Code)
	}
}

func TestUpdateCategoryInternalServerError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = errors.New("db down")

	res := performCategoryRequest(router, http.MethodPut, "/api/categories/1", `{"name":"ปิ้งย่าง"}`)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestDeleteCategory(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.categories = append(svc.categories, domain.Category{ID: 1, Name: "ชาบู"})

	res := performCategoryRequest(router, http.MethodDelete, "/api/categories/1", "")

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, res.Code)
	}
}

func TestDeleteCategoryNotFound(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodDelete, "/api/categories/99", "")

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestDeleteCategoryRejectsInvalidID(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodDelete, "/api/categories/0", "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestDeleteCategoryRejectsNonNumericID(t *testing.T) {
	router, _ := setupCategoryTestRouter()

	res := performCategoryRequest(router, http.MethodDelete, "/api/categories/abc", "")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestDeleteCategoryInternalServerError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = errors.New("db down")

	res := performCategoryRequest(router, http.MethodDelete, "/api/categories/1", "")

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestCategoryHandlerReturnsInternalServerError(t *testing.T) {
	router, svc := setupCategoryTestRouter()
	svc.err = errors.New("db down")

	res := performCategoryRequest(router, http.MethodGet, "/api/categories", "")

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestNewCategoryHandler(t *testing.T) {
	h := NewCategoryHandler()
	if h == nil {
		t.Fatal("expected handler to be non-nil")
	}
}

func setupCategoryTestRouter() (*gin.Engine, *mockCategoryService) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	svc := &mockCategoryService{}
	category := NewCategoryHandlerWithService(svc)
	api := router.Group("/api")
	api.GET("/categories", category.GetCategories)
	api.GET("/categories/:id", category.GetCategory)
	api.POST("/categories", category.CreateCategory)
	api.PUT("/categories/:id", category.UpdateCategory)
	api.DELETE("/categories/:id", category.DeleteCategory)

	return router, svc
}

func performCategoryRequest(router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body == "" {
		reqBody = bytes.NewBuffer(nil)
	} else {
		reqBody = bytes.NewBufferString(body)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}
