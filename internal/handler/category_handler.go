package handler

import (
	"errors"
	"net/http"
	"qflow/db"
	"qflow/internal/domain"
	"qflow/internal/repository"
	"qflow/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service domain.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	repo := repository.NewCategoryGormRepository(db.DB)
	svc := service.NewCategoryService(repo)

	return &CategoryHandler{
		service: svc,
	}
}

type categoryRequest struct {
	Name string `json:"name"`
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	category, err := h.service.GetCategory(c.Request.Context(), uint(id))
	if errors.Is(err, service.ErrCategoryNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req categoryRequest
	c.ShouldBindJSON(&req)

	category, err := h.service.CreateCategory(c.Request.Context(), req.Name)

	if errors.Is(err, service.ErrCategoryNameRequired) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name required"})
		return
	}

	if errors.Is(err, service.ErrCategoryDuplicate) {
		c.JSON(http.StatusConflict, gin.H{"message": "duplicate"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req categoryRequest
	c.ShouldBindJSON(&req)

	category, err := h.service.UpdateCategory(c.Request.Context(), uint(id), req.Name)

	if errors.Is(err, service.ErrCategoryNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.DeleteCategory(c.Request.Context(), uint(id))
	if errors.Is(err, service.ErrCategoryNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.Status(http.StatusNoContent)
}