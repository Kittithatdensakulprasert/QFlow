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

type mockQueueService struct {
	queues []domain.Queue
	err    error
}

func (m *mockQueueService) BookQueue(userID, zoneID uint) (*domain.Queue, error) {
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

func (m *mockQueueService) GetQueueHistory(userID uint) ([]domain.Queue, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := []domain.Queue{}
	for _, queue := range m.queues {
		if queue.UserID == userID {
			result = append(result, queue)
		}
	}
	return result, nil
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

func (m *mockQueueService) CancelQueue(id, userID uint) error {
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

func (m *mockQueueService) GetQueuesByZone(zoneID uint) ([]domain.Queue, error) {
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

func (m *mockQueueService) CallQueue(id uint) (*domain.Queue, error) {
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

func (m *mockQueueService) CompleteQueue(id uint) (*domain.Queue, error) {
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

func (m *mockQueueService) SkipQueue(id uint) (*domain.Queue, error) {
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

func (m *mockQueueService) GetQueueByNumber(queueNumber int, userID uint) (*domain.Queue, error) {
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

	if response["message"] != "Queue booked successfully" {
		t.Fatalf("unexpected message: %v", response["message"])
	}

	// Check that queue was created
	queue := response["queue"].(map[string]interface{})
	if queue["zone_id"] != float64(1) {
		t.Fatalf("unexpected zone_id: %v", queue["zone_id"])
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
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodGet, "/api/queues/history", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	queues, ok := response["queues"].([]interface{})
	if !ok {
		t.Fatalf("expected queues array in response")
	}

	if len(queues) == 0 {
		t.Fatalf("expected at least one queue in response")
	}
}

func TestGetQueue(t *testing.T) {
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodGet, "/api/queues/101", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["queue"] == nil {
		t.Fatalf("expected queue in response")
	}

	queue := response["queue"].(map[string]interface{})
	if queue["queue_number"] == nil {
		t.Fatalf("expected queue_number in response")
	}
}

func TestCancelQueue(t *testing.T) {
	router, _ := setupQueueTestRouter()

	res := performQueueRequest(router, http.MethodPatch, "/api/queues/1/cancel", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "Queue cancelled successfully" {
		t.Fatalf("unexpected message: %v", response["message"])
	}

	queue := response["queue"].(map[string]interface{})
	if queue["status"] != "cancelled" {
		t.Fatalf("expected queue status to be cancelled: %v", queue["status"])
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

func setupQueueTestRouter() (*gin.Engine, *mockQueueService) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	svc := &mockQueueService{}
	handler := NewQueueHandler(svc)

	// Mock user ID in context
	router.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	api := router.Group("/api")
	api.POST("/queues/book", handler.BookQueue)
	api.GET("/queues/history", handler.GetHistory)
	api.GET("/queues/:queueNumber", handler.GetQueue)
	api.PATCH("/queues/:id/cancel", handler.CancelQueue)

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
