package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qflow/models"
	"testing"
)

func TestCreateAndListProviders(t *testing.T) {
	resetProviderData()

	req := httptest.NewRequest(http.MethodPost, "/api/providers", bytes.NewBufferString(`{"name":"Bangkok Clinic"}`))
	res := httptest.NewRecorder()

	ProvidersHandler(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var provider models.Provider
	if err := json.NewDecoder(res.Body).Decode(&provider); err != nil {
		t.Fatalf("decode provider: %v", err)
	}

	if provider.ID != 1 || provider.Name != "Bangkok Clinic" {
		t.Fatalf("unexpected provider: %+v", provider)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/providers", nil)
	res = httptest.NewRecorder()

	ProvidersHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var list []models.Provider
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatalf("decode providers: %v", err)
	}

	if len(list) != 1 || list[0].Name != "Bangkok Clinic" {
		t.Fatalf("unexpected providers: %+v", list)
	}
}

func TestCreateListAndToggleZones(t *testing.T) {
	resetProviderData()

	createProviderForTest(t, "Bangkok Clinic")

	req := httptest.NewRequest(http.MethodPost, "/api/providers/1/zones", bytes.NewBufferString(`{"name":"Counter A"}`))
	res := httptest.NewRecorder()

	ProviderHandler(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var zone models.Zone
	if err := json.NewDecoder(res.Body).Decode(&zone); err != nil {
		t.Fatalf("decode zone: %v", err)
	}

	if zone.ID != 1 || zone.ProviderID != 1 || zone.Name != "Counter A" || !zone.IsOpen || zone.QueueCount != 0 {
		t.Fatalf("unexpected zone: %+v", zone)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/providers/1/zones", nil)
	res = httptest.NewRecorder()

	ProviderHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var list []models.Zone
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatalf("decode zones: %v", err)
	}

	if len(list) != 1 || list[0].Name != "Counter A" {
		t.Fatalf("unexpected zones: %+v", list)
	}

	req = httptest.NewRequest(http.MethodPatch, "/api/zones/1/toggle", nil)
	res = httptest.NewRecorder()

	ZoneHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	if err := json.NewDecoder(res.Body).Decode(&zone); err != nil {
		t.Fatalf("decode toggled zone: %v", err)
	}

	if zone.IsOpen {
		t.Fatalf("expected zone to be closed after toggle: %+v", zone)
	}
}

func TestCreateZoneRequiresExistingProvider(t *testing.T) {
	resetProviderData()

	req := httptest.NewRequest(http.MethodPost, "/api/providers/99/zones", bytes.NewBufferString(`{"name":"Counter A"}`))
	res := httptest.NewRecorder()

	ProviderHandler(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func createProviderForTest(t *testing.T, name string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/providers", bytes.NewBufferString(`{"name":"`+name+`"}`))
	res := httptest.NewRecorder()

	ProvidersHandler(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("create provider: expected status %d, got %d", http.StatusCreated, res.Code)
	}
}
