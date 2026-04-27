package handlers

import (
	"encoding/json"
	"net/http"
	"qflow/models"
	"strconv"
	"strings"
	"sync"
)

var (
	providerMu     sync.Mutex
	providers      []models.Provider
	zones          []models.Zone
	nextProviderID = 1
	nextZoneID     = 1
)

type providerRequest struct {
	Name string `json:"name"`
}

type zoneRequest struct {
	Name string `json:"name"`
}

func ProvidersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createProvider(w, r)
	case http.MethodGet:
		listProviders(w)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func ProviderHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/providers/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) != 2 || parts[1] != "zones" {
		http.NotFound(w, r)
		return
	}

	providerID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid provider id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		createZone(w, r, providerID)
	case http.MethodGet:
		listZones(w, r, providerID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func ZoneHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/zones/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(parts) != 2 || parts[1] != "toggle" {
		http.NotFound(w, r)
		return
	}

	zoneID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid zone id", http.StatusBadRequest)
		return
	}

	providerMu.Lock()
	defer providerMu.Unlock()

	for i := range zones {
		if zones[i].ID == zoneID {
			zones[i].IsOpen = !zones[i].IsOpen
			writeJSON(w, http.StatusOK, zones[i])
			return
		}
	}

	http.NotFound(w, r)
}

func createProvider(w http.ResponseWriter, r *http.Request) {
	var req providerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	providerMu.Lock()
	defer providerMu.Unlock()

	provider := models.Provider{
		ID:   nextProviderID,
		Name: req.Name,
	}
	nextProviderID++
	providers = append(providers, provider)

	writeJSON(w, http.StatusCreated, provider)
}

func listProviders(w http.ResponseWriter) {
	providerMu.Lock()
	defer providerMu.Unlock()

	providerList := make([]models.Provider, len(providers))
	copy(providerList, providers)

	writeJSON(w, http.StatusOK, providerList)
}

func createZone(w http.ResponseWriter, r *http.Request, providerID int) {
	var req zoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	providerMu.Lock()
	defer providerMu.Unlock()

	if !providerExists(providerID) {
		http.NotFound(w, r)
		return
	}

	zone := models.Zone{
		ID:         nextZoneID,
		ProviderID: providerID,
		Name:       req.Name,
		IsOpen:     true,
		QueueCount: 0,
	}
	nextZoneID++
	zones = append(zones, zone)

	writeJSON(w, http.StatusCreated, zone)
}

func listZones(w http.ResponseWriter, r *http.Request, providerID int) {
	providerMu.Lock()
	defer providerMu.Unlock()

	if !providerExists(providerID) {
		http.NotFound(w, r)
		return
	}

	providerZones := make([]models.Zone, 0)
	for _, zone := range zones {
		if zone.ProviderID == providerID {
			providerZones = append(providerZones, zone)
		}
	}

	writeJSON(w, http.StatusOK, providerZones)
}

func providerExists(providerID int) bool {
	for _, provider := range providers {
		if provider.ID == providerID {
			return true
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func resetProviderData() {
	providerMu.Lock()
	defer providerMu.Unlock()

	providers = nil
	zones = nil
	nextProviderID = 1
	nextZoneID = 1
}
