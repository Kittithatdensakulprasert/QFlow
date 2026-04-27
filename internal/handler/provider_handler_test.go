package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qflow/internal/domain"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateProvider(t *testing.T) {
	router := setupProviderTestRouter()

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
	router := setupProviderTestRouter()
	performProviderRequest(router, http.MethodPost, "/api/providers", `{"name":"Bangkok Clinic"}`)

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
	router := setupProviderTestRouter()
	performProviderRequest(router, http.MethodPost, "/api/providers", `{"name":"Bangkok Clinic"}`)

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

	res = performProviderRequest(router, http.MethodGet, "/api/providers/1/zones", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var zones []domain.Zone
	if err := json.NewDecoder(res.Body).Decode(&zones); err != nil {
		t.Fatalf("decode zones: %v", err)
	}

	if len(zones) != 1 || zones[0].Name != "Counter A" || zones[0].QueueCount != 0 {
		t.Fatalf("unexpected zones: %+v", zones)
	}
}

func TestToggleZone(t *testing.T) {
	router := setupProviderTestRouter()
	performProviderRequest(router, http.MethodPost, "/api/providers", `{"name":"Bangkok Clinic"}`)
	performProviderRequest(router, http.MethodPost, "/api/providers/1/zones", `{"name":"Counter A"}`)

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
	router := setupProviderTestRouter()

	res := performProviderRequest(router, http.MethodPost, "/api/providers/99/zones", `{"name":"Counter A"}`)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func setupProviderTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	provider := NewProviderHandler()
	api := router.Group("/api")
	api.POST("/providers", provider.CreateProvider)
	api.GET("/providers", provider.GetProviders)
	api.POST("/providers/:id/zones", provider.CreateZone)
	api.GET("/providers/:id/zones", provider.GetZones)
	api.PATCH("/zones/:id/toggle", provider.ToggleZone)

	return router
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
