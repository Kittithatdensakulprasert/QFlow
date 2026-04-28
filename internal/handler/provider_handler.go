package handler

import (
	"errors"
	"net/http"
	"qflow/internal/domain"
	"qflow/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	svc domain.ProviderService
}

type createProviderRequest struct {
	Name       string `json:"name" binding:"required"`
	CategoryID uint   `json:"category_id"`
}

type createZoneRequest struct {
	Name string `json:"name" binding:"required"`
}

func NewProviderHandler(svc domain.ProviderService) *ProviderHandler {
	return &ProviderHandler{svc: svc}
}

func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req createProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider, err := h.svc.CreateProvider(req.Name, req.CategoryID)
	if err != nil {
		if errors.Is(err, service.ErrProviderNameRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, provider)
}

func (h *ProviderHandler) GetProviders(c *gin.Context) {
	providers, err := h.svc.GetProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, providers)
}

func (h *ProviderHandler) CreateZone(c *gin.Context) {
	providerID, ok := parseUintParam(c, "id", "invalid provider id")
	if !ok {
		return
	}

	var req createZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zone, err := h.svc.CreateZone(providerID, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrZoneNameRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, zone)
}

func (h *ProviderHandler) GetZones(c *gin.Context) {
	providerID, ok := parseUintParam(c, "id", "invalid provider id")
	if !ok {
		return
	}

	zones, err := h.svc.GetZones(providerID)
	if err != nil {
		if errors.Is(err, service.ErrProviderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, zones)
}

func (h *ProviderHandler) ToggleZone(c *gin.Context) {
	zoneID, ok := parseUintParam(c, "id", "invalid zone id")
	if !ok {
		return
	}

	zone, err := h.svc.ToggleZone(zoneID)
	if err != nil {
		if errors.Is(err, service.ErrZoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, zone)
}

func parseUintParam(c *gin.Context, name string, errorMessage string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return 0, false
	}

	return uint(value), true
}
