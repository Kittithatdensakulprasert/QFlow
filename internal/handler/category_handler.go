package handler

import (
	"errors"
	"net/http"
	"qflow/internal/domain"
	"qflow/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service domain.CategoryService
}

func NewCategoryHandler(svc domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

type categoryRequest struct {
	Name string `json:"name"`
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to get categories")
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "category id must be a number")
		return
	}

	category, err := h.service.GetCategory(c.Request.Context(), id)
	if errors.Is(err, service.ErrCategoryNotFound) {
		respondError(c, http.StatusNotFound, "CATEGORY_NOT_FOUND", "category not found")
		return
	}

	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to get category")
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req categoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid JSON body")
		return
	}

	category, err := h.service.CreateCategory(c.Request.Context(), req.Name)
	if errors.Is(err, service.ErrCategoryNameRequired) {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "category name is required")
		return
	}

	if errors.Is(err, service.ErrCategoryDuplicate) {
		respondError(c, http.StatusConflict, "CATEGORY_DUPLICATE", "category name already exists")
		return
	}

	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create category")
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "category id must be a number")
		return
	}

	var req categoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid JSON body")
		return
	}

	category, err := h.service.UpdateCategory(c.Request.Context(), id, req.Name)
	if errors.Is(err, service.ErrCategoryNameRequired) {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "category name is required")
		return
	}

	if errors.Is(err, service.ErrCategoryNotFound) {
		respondError(c, http.StatusNotFound, "CATEGORY_NOT_FOUND", "category not found")
		return
	}

	if errors.Is(err, service.ErrCategoryDuplicate) {
		respondError(c, http.StatusConflict, "CATEGORY_DUPLICATE", "category name already exists")
		return
	}

	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to update category")
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "category id must be a number")
		return
	}

	err = h.service.DeleteCategory(c.Request.Context(), id)
	if errors.Is(err, service.ErrCategoryNotFound) {
		respondError(c, http.StatusNotFound, "CATEGORY_NOT_FOUND", "category not found")
		return
	}

	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to delete category")
		return
	}

	c.Status(http.StatusNoContent)
}
