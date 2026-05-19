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

type mockQueueService struct {
	queues []domain.Queue
	err    error
}

func (m *mockQueueService) BookQueue(_ context.Context, userID, zoneID uint) (*domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	queue := domain.Queue{
		ID:          uint(len(m.queues) + 1),
		QueueNumber: 100 + len(m.queues) + 1,
		ZoneID:      zoneID,
		UserID:      userID,
		Status:      "waiting",
	}
	m.queues = append(m.queues, queue)
	return &queue, nil
}

func (m *mockQueueService) GetQueueHistory(_ context.Context, userID uint, page, limit int) ([]domain.Queue, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	result := []domain.Queue{}
	for _, queue := range m.queues {
		if queue.UserID == userID {
			result = append(result, queue)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockQueueService) GetQueue(queueNumber int) (*domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, queue := range m.queues {
		if queue.QueueNumber == queueNumber {
			return &queue, nil
		}
	}
	return nil, service.ErrQueueNotFound
}

func (m *mockQueueService) CancelQueue(_ context.Context, id, userID uint) error {
	if m.err != nil {
		return m.err
	}
	for i := range m.queues {
		if m.queues[i].ID == id {
			m.queues[i].Status = "cancelled"
			return nil
		}
	}
	return service.ErrQueueNotFound
}

func (m *mockQueueService) GetQueuesByZone(_ context.Context, zoneID uint) ([]domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := []domain.Queue{}
	for _, queue := range m.queues {
		if queue.ZoneID == zoneID {
			result = append(result, queue)
		}
	}
	return result, nil
}

func (m *mockQueueService) CallQueue(_ context.Context, id uint) (*domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.queues {
		if m.queues[i].ID == id {
			m.queues[i].Status = "called"
			return &m.queues[i], nil
		}
	}
	return nil, service.ErrQueueNotFound
}

func (m *mockQueueService) CompleteQueue(_ context.Context, id uint) (*domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.queues {
		if m.queues[i].ID == id {
			m.queues[i].Status = "completed"
			return &m.queues[i], nil
		}
	}
	return nil, service.ErrQueueNotFound
}

func (m *mockQueueService) SkipQueue(_ context.Context, id uint) (*domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.queues {
		if m.queues[i].ID == id {
			m.queues[i].Status = "skipped"
			return &m.queues[i], nil
		}
	}
	return nil, service.ErrQueueNotFound
}

func (m *mockQueueService) GetQueueByNumber(_ context.Context, queueNumber int, userID uint) (*domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, queue := range m.queues {
		if queue.QueueNumber == queueNumber && queue.UserID == userID {
			return &queue, nil
		}
	}
	return nil, service.ErrQueueNotFound
}

func TestBookQueue(t *testing.T) {
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPost, "/api/queues/book", `{"zone_id":1}`)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["zone_id"] != float64(1) {
		t.Fatalf("unexpected zone_id: %v", response["zone_id"])
	}
	if response["status"] != "waiting" {
		t.Fatalf("unexpected status: %v", response["status"])
	}
}

func TestBookQueueWithInvalidFormat(t *testing.T) {
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPost, "/api/queues/book", `{"invalid":"json"}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestGetHistory(t *testing.T) {
	router, svc := setupQueueTestRouter()
	// pre-populate a queue for userID=1
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "waiting"},
	}

	res := performQueueRequest(router, http.MethodGet, "/api/queues/history", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response struct {
		Data       []domain.Queue `json:"data"`
		Pagination map[string]any `json:"pagination"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(response.Data) == 0 {
		t.Fatalf("expected at least one queue in response")
	}
	if response.Pagination["total"] == nil {
		t.Fatalf("expected pagination metadata")
	}
}

func TestGetQueue(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "waiting"},
	}

	res := performQueueRequest(router, http.MethodGet, "/api/queues/101", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["queue_number"] == nil {
		t.Fatalf("expected queue_number in response")
	}
}

func TestCancelQueue(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "waiting"},
	}

	res := performQueueRequest(router, http.MethodPatch, "/api/queues/1/cancel", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "queue cancelled" {
		t.Fatalf("unexpected message: %v", response["message"])
	}
}

func TestQueueHandlerReturnsInternalServerError(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.err = errors.New("service error")

	res := performQueueRequest(router, http.MethodGet, "/api/queues/history", "")

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestBookQueueUnauthorized(t *testing.T) {
	router, _ := setupQueueTestRouterWithoutAuth()

	res := performQueueRequest(router, http.MethodPost, "/api/queues/book", `{"zone_id":1}`)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestBookQueueErrorMappings(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"invalid user", service.ErrInvalidUserID, http.StatusBadRequest},
		{"invalid zone", service.ErrInvalidZoneID, http.StatusBadRequest},
		{"zone not found", service.ErrZoneNotFound, http.StatusNotFound},
		{"zone closed", service.ErrZoneClosed, http.StatusConflict},
		{"internal", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, svc := setupQueueTestRouter()
			svc.err = tt.err

			res := performQueueRequest(router, http.MethodPost, "/api/queues/book", `{"zone_id":1}`)
			if res.Code != tt.wantCode {
				t.Fatalf("expected status %d, got %d", tt.wantCode, res.Code)
			}
		})
	}
}

func TestGetHistoryUnauthorized(t *testing.T) {
	router, _ := setupQueueTestRouterWithoutAuth()

	res := performQueueRequest(router, http.MethodGet, "/api/queues/history", "")
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestGetHistoryInvalidUserID(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.err = service.ErrInvalidUserID

	res := performQueueRequest(router, http.MethodGet, "/api/queues/history", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestGetQueueInvalidQueueNumber(t *testing.T) {
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodGet, "/api/queues/abc", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}

	res = performQueueRequest(router, http.MethodGet, "/api/queues/0", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestGetQueueUnauthorized(t *testing.T) {
	router, _ := setupQueueTestRouterWithoutAuth()

	res := performQueueRequest(router, http.MethodGet, "/api/queues/101", "")
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestGetQueueErrorMappings(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"invalid user", service.ErrInvalidUserID, http.StatusBadRequest},
		{"not found", service.ErrQueueNotFound, http.StatusNotFound},
		{"forbidden", service.ErrForbiddenQueue, http.StatusForbidden},
		{"internal", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, svc := setupQueueTestRouter()
			svc.err = tt.err

			res := performQueueRequest(router, http.MethodGet, "/api/queues/101", "")
			if res.Code != tt.wantCode {
				t.Fatalf("expected status %d, got %d", tt.wantCode, res.Code)
			}
		})
	}
}

func TestCancelQueueInvalidID(t *testing.T) {
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPatch, "/api/queues/abc/cancel", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}

	res = performQueueRequest(router, http.MethodPatch, "/api/queues/0/cancel", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestCancelQueueUnauthorized(t *testing.T) {
	router, _ := setupQueueTestRouterWithoutAuth()

	res := performQueueRequest(router, http.MethodPatch, "/api/queues/1/cancel", "")
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestCancelQueueErrorMappings(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"invalid user", service.ErrInvalidUserID, http.StatusBadRequest},
		{"not found", service.ErrQueueNotFound, http.StatusNotFound},
		{"forbidden", service.ErrForbiddenQueue, http.StatusForbidden},
		{"finalized", service.ErrQueueFinalized, http.StatusConflict},
		{"cancelled", service.ErrQueueCancelled, http.StatusConflict},
		{"internal", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, svc := setupQueueTestRouter()
			svc.err = tt.err

			res := performQueueRequest(router, http.MethodPatch, "/api/queues/1/cancel", "")
			if res.Code != tt.wantCode {
				t.Fatalf("expected status %d, got %d", tt.wantCode, res.Code)
			}
		})
	}
}

func TestGetQueuesByZoneBranches(t *testing.T) {
	router, svc := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodGet, "/api/manage/queues/abc", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}

	svc.err = errors.New("boom")
	res = performQueueRequest(router, http.MethodGet, "/api/manage/queues/1", "")
	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestCallQueueBranches(t *testing.T) {
	router, svc := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPatch, "/api/manage/queues/abc/call", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}

	svc.err = service.ErrQueueNotFound
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/call", "")
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}

	svc.err = domain.ErrQueueCannotBeCalled
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/call", "")
	if res.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, res.Code)
	}

	svc.err = errors.New("boom")
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/call", "")
	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestCompleteQueueBranches(t *testing.T) {
	router, svc := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPatch, "/api/manage/queues/abc/complete", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}

	svc.err = service.ErrQueueNotFound
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/complete", "")
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}

	svc.err = domain.ErrQueueCannotBeCompleted
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/complete", "")
	if res.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, res.Code)
	}

	svc.err = errors.New("boom")
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/complete", "")
	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestSkipQueueBranches(t *testing.T) {
	router, svc := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPatch, "/api/manage/queues/abc/skip", "")
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}

	svc.err = service.ErrQueueNotFound
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/skip", "")
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}

	svc.err = domain.ErrQueueCannotBeSkipped
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/skip", "")
	if res.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, res.Code)
	}

	svc.err = errors.New("boom")
	res = performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/skip", "")
	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestGetQueuesByZoneSuccess(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "waiting"},
		{ID: 2, QueueNumber: 102, ZoneID: 2, UserID: 1, Status: "waiting"},
	}

	res := performQueueRequest(router, http.MethodGet, "/api/manage/queues/1", "")
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestCallQueueSuccess(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "waiting"},
	}

	res := performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/call", "")
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestCompleteQueueSuccess(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "called"},
	}

	res := performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/complete", "")
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestSkipQueueSuccess(t *testing.T) {
	router, svc := setupQueueTestRouter()
	svc.queues = []domain.Queue{
		{ID: 1, QueueNumber: 101, ZoneID: 1, UserID: 1, Status: "called"},
	}

	res := performQueueRequest(router, http.MethodPatch, "/api/manage/queues/1/skip", "")
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func setupQueueTestRouter() (*gin.Engine, *mockQueueService) {
	return setupQueueTestRouterWithAuth(true)
}

func setupQueueTestRouterWithoutAuth() (*gin.Engine, *mockQueueService) {
	return setupQueueTestRouterWithAuth(false)
}

func setupQueueTestRouterWithAuth(withAuth bool) (*gin.Engine, *mockQueueService) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	svc := &mockQueueService{}
	handler := NewQueueHandler(svc)

	if withAuth {
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Next()
		})
	}

	api := router.Group("/api")
	api.POST("/queues/book", handler.BookQueue)
	api.GET("/queues/history", handler.GetHistory)
	api.GET("/queues/:queueNumber", handler.GetQueue)
	api.PATCH("/queues/:id/cancel", handler.CancelQueue)
	api.GET("/manage/queues/:zoneId", handler.GetQueuesByZone)
	api.PATCH("/manage/queues/:id/call", handler.CallQueue)
	api.PATCH("/manage/queues/:id/complete", handler.CompleteQueue)
	api.PATCH("/manage/queues/:id/skip", handler.SkipQueue)

	return router, svc
}

func performQueueRequest(router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
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
