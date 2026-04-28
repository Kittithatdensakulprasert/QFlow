package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockProviderService struct {
	providers []domain.Provider
	zones     []domain.Zone
	err       error
}

func (m *mockProviderService) CreateProvider(name string, categoryID uint) (*domain.Provider, error) {
	if m.err != nil {
		return nil, m.err
	}
	provider := domain.Provider{ID: uint(len(m.providers) + 1), Name: name}
	if categoryID > 0 {
		provider.CategoryID = &categoryID
	}
	m.providers = append(m.providers, provider)
	return &provider, nil
}

func (m *mockProviderService) GetProviders() ([]domain.Provider, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.providers, nil
}

func (m *mockProviderService) CreateZone(providerID uint, name string) (*domain.Zone, error) {
	if m.err != nil {
		return nil, m.err
	}
	zone := domain.Zone{
		ID:         uint(len(m.zones) + 1),
		ProviderID: providerID,
		Name:       name,
		IsOpen:     true,
		QueueCount: 0,
	}
	m.zones = append(m.zones, zone)
	return &zone, nil
}

func (m *mockProviderService) GetZones(providerID uint) ([]domain.Zone, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := []domain.Zone{}
	for _, zone := range m.zones {
		if zone.ProviderID == providerID {
			result = append(result, zone)
		}
	}
	return result, nil
}

func (m *mockProviderService) ToggleZone(id uint) (*domain.Zone, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.zones {
		if m.zones[i].ID == id {
			m.zones[i].IsOpen = !m.zones[i].IsOpen
			return &m.zones[i], nil
		}
	}
	return nil, service.ErrProviderZoneNotFound
}

func TestCreateProvider(t *testing.T) {
	router, _ := setupProviderTestRouter()

	res := performProviderRequest(router, http.MethodPost, "/api/providers", `{"name":"Bangkok Clinic"}`)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var provider domain.Provider
	if err := json.NewDecoder(res.Body).Decode(&provider); err != nil {
		t.Fatalf("decode provider: %v", err)
	}

	if provider.ID != 1 || provider.Name != "Bangkok Clinic" {
		t.Fatalf("unexpected provider: %+v", provider)
	}
}

func TestGetProviders(t *testing.T) {
	router, svc := setupProviderTestRouter()
	svc.providers = append(svc.providers, domain.Provider{ID: 1, Name: "Bangkok Clinic"})

	res := performProviderRequest(router, http.MethodGet, "/api/providers", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var providers []domain.Provider
	if err := json.NewDecoder(res.Body).Decode(&providers); err != nil {
		t.Fatalf("decode providers: %v", err)
	}

	if len(providers) != 1 || providers[0].Name != "Bangkok Clinic" {
		t.Fatalf("unexpected providers: %+v", providers)
	}
}

func TestCreateAndGetZones(t *testing.T) {
	router, svc := setupProviderTestRouter()

	res := performProviderRequest(router, http.MethodPost, "/api/providers/1/zones", `{"name":"Counter A"}`)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var zone domain.Zone
	if err := json.NewDecoder(res.Body).Decode(&zone); err != nil {
		t.Fatalf("decode zone: %v", err)
	}

	if zone.ID != 1 || zone.ProviderID != 1 || zone.Name != "Counter A" || !zone.IsOpen || zone.QueueCount != 0 {
		t.Fatalf("unexpected zone: %+v", zone)
	}

	svc.zones[0].QueueCount = 2
	res = performProviderRequest(router, http.MethodGet, "/api/providers/1/zones", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var zones []domain.Zone
	if err := json.NewDecoder(res.Body).Decode(&zones); err != nil {
		t.Fatalf("decode zones: %v", err)
	}

	if len(zones) != 1 || zones[0].Name != "Counter A" || zones[0].QueueCount != 2 {
		t.Fatalf("unexpected zones: %+v", zones)
	}
}

func TestToggleZone(t *testing.T) {
	router, svc := setupProviderTestRouter()
	svc.zones = append(svc.zones, domain.Zone{ID: 1, ProviderID: 1, Name: "Counter A", IsOpen: true})

	res := performProviderRequest(router, http.MethodPatch, "/api/zones/1/toggle", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var zone domain.Zone
	if err := json.NewDecoder(res.Body).Decode(&zone); err != nil {
		t.Fatalf("decode zone: %v", err)
	}

	if zone.IsOpen {
		t.Fatalf("expected zone to be closed after toggle: %+v", zone)
	}
}

func TestCreateZoneRequiresExistingProvider(t *testing.T) {
	router, svc := setupProviderTestRouter()
	svc.err = service.ErrProviderNotFound

	res := performProviderRequest(router, http.MethodPost, "/api/providers/99/zones", `{"name":"Counter A"}`)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestProviderHandlerReturnsInternalServerError(t *testing.T) {
	router, svc := setupProviderTestRouter()
	svc.err = errors.New("db down")

	res := performProviderRequest(router, http.MethodGet, "/api/providers", "")

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func setupProviderTestRouter() (*gin.Engine, *mockProviderService) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	svc := &mockProviderService{}
	provider := NewProviderHandler(svc)
	api := router.Group("/api")
	api.POST("/providers", provider.CreateProvider)
	api.GET("/providers", provider.GetProviders)
	api.POST("/providers/:id/zones", provider.CreateZone)
	api.GET("/providers/:id/zones", provider.GetZones)
	api.PATCH("/zones/:id/toggle", provider.ToggleZone)

	return router, svc
}

func performProviderRequest(router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
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
