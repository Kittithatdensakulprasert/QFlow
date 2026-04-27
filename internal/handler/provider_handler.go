package handler

import (
	"net/http"
	"qflow/internal/domain"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	mu             sync.Mutex
	providers      []domain.Provider
	zones          []domain.Zone
	nextProviderID uint
	nextZoneID     uint
}

type createProviderRequest struct {
	Name string `json:"name"`
}

type createZoneRequest struct {
	Name string `json:"name"`
}

func NewProviderHandler() *ProviderHandler {
	return &ProviderHandler{
		providers:      make([]domain.Provider, 0),
		zones:          make([]domain.Zone, 0),
		nextProviderID: 1,
		nextZoneID:     1,
	}
}

func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req createProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name is required"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	provider := domain.Provider{
		ID:   h.nextProviderID,
		Name: req.Name,
	}
	h.nextProviderID++
	h.providers = append(h.providers, provider)

	c.JSON(http.StatusCreated, provider)
}

func (h *ProviderHandler) GetProviders(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	providers := make([]domain.Provider, len(h.providers))
	copy(providers, h.providers)

	c.JSON(http.StatusOK, providers)
}

func (h *ProviderHandler) CreateZone(c *gin.Context) {
	providerID, ok := parseUintParam(c, "id", "invalid provider id")
	if !ok {
		return
	}

	var req createZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name is required"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.providerExists(providerID) {
		c.JSON(http.StatusNotFound, gin.H{"message": "provider not found"})
		return
	}

	zone := domain.Zone{
		ID:         h.nextZoneID,
		ProviderID: providerID,
		Name:       req.Name,
		IsOpen:     true,
		QueueCount: 0,
	}
	h.nextZoneID++
	h.zones = append(h.zones, zone)

	c.JSON(http.StatusCreated, zone)
}

func (h *ProviderHandler) GetZones(c *gin.Context) {
	providerID, ok := parseUintParam(c, "id", "invalid provider id")
	if !ok {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.providerExists(providerID) {
		c.JSON(http.StatusNotFound, gin.H{"message": "provider not found"})
		return
	}

	providerZones := make([]domain.Zone, 0)
	for _, zone := range h.zones {
		if zone.ProviderID == providerID {
			providerZones = append(providerZones, zone)
		}
	}

	c.JSON(http.StatusOK, providerZones)
}

func (h *ProviderHandler) ToggleZone(c *gin.Context) {
	zoneID, ok := parseUintParam(c, "id", "invalid zone id")
	if !ok {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.zones {
		if h.zones[i].ID == zoneID {
			h.zones[i].IsOpen = !h.zones[i].IsOpen
			c.JSON(http.StatusOK, h.zones[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "zone not found"})
}

func (h *ProviderHandler) providerExists(id uint) bool {
	for _, provider := range h.providers {
		if provider.ID == id {
			return true
		}
	}
	return false
}

func parseUintParam(c *gin.Context, name string, errorMessage string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": errorMessage})
		return 0, false
	}

	return uint(value), true
}
