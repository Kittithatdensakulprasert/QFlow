package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"qflow/internal/domain"
	"qflow/internal/service"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockNotificationService struct {
	notifications []domain.Notification
	getErr        error
	sendErr       error
	markErr       error
	deleteErr     error
}

func (m *mockNotificationService) GetNotifications(userID uint) ([]domain.Notification, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}

	result := make([]domain.Notification, 0)
	for _, n := range m.notifications {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	return result, nil
}

func (m *mockNotificationService) SendNotification(userID uint, message string) (*domain.Notification, error) {
	if m.sendErr != nil {
		return nil, m.sendErr
	}

	n := domain.Notification{
		ID:      uint(len(m.notifications) + 1),
		UserID:  userID,
		Message: message,
		IsRead:  false,
	}
	m.notifications = append(m.notifications, n)
	return &n, nil
}

func (m *mockNotificationService) MarkNotificationRead(id, userID uint) error {
	if m.markErr != nil {
		return m.markErr
	}

	for i := range m.notifications {
		if m.notifications[i].ID == id {
			if m.notifications[i].UserID != userID {
				return service.ErrNotificationForbidden
			}
			m.notifications[i].IsRead = true
			return nil
		}
	}

	return service.ErrNotificationNotFound
}

func (m *mockNotificationService) DeleteNotification(id, userID uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}

	for i := range m.notifications {
		if m.notifications[i].ID == id {
			if m.notifications[i].UserID != userID {
				return service.ErrNotificationForbidden
			}
			m.notifications = append(m.notifications[:i], m.notifications[i+1:]...)
			return nil
		}
	}

	return service.ErrNotificationNotFound
}

func TestGetNotifications(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.notifications = []domain.Notification{
		{ID: 1, UserID: 7, Message: "A"},
		{ID: 2, UserID: 8, Message: "B"},
	}

	res := performNotificationRequest(router, http.MethodGet, "/api/notifications", "", 7)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var notifications []domain.Notification
	if err := json.NewDecoder(res.Body).Decode(&notifications); err != nil {
		t.Fatalf("decode notifications: %v", err)
	}

	if len(notifications) != 1 || notifications[0].ID != 1 {
		t.Fatalf("unexpected notifications: %+v", notifications)
	}
}

func TestGetNotificationsUnauthorized(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodGet, "/api/notifications", "", nil)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestGetNotificationsInternalServerError(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.getErr = errors.New("db down")

	res := performNotificationRequest(router, http.MethodGet, "/api/notifications", "", 7)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestSendNotification(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodPost, "/api/notifications/send", `{"message":"hello"}`, 7)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var n domain.Notification
	if err := json.NewDecoder(res.Body).Decode(&n); err != nil {
		t.Fatalf("decode notification: %v", err)
	}

	if n.ID != 1 || n.UserID != 7 || n.Message != "hello" {
		t.Fatalf("unexpected notification: %+v", n)
	}
}

func TestSendNotificationInvalidBody(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodPost, "/api/notifications/send", `{`, 7)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestSendNotificationUnauthorized(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodPost, "/api/notifications/send", `{"message":"hello"}`, nil)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestSendNotificationServiceError(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.sendErr = errors.New("message is required")

	res := performNotificationRequest(router, http.MethodPost, "/api/notifications/send", `{"message":"hello"}`, 7)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestMarkNotificationRead(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.notifications = []domain.Notification{{ID: 1, UserID: 7, Message: "A"}}

	res := performNotificationRequest(router, http.MethodPatch, "/api/notifications/1/read", "", 7)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestMarkNotificationReadUnauthorized(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodPatch, "/api/notifications/1/read", "", nil)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestMarkNotificationReadInvalidID(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodPatch, "/api/notifications/abc/read", "", 7)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestMarkNotificationReadNotFound(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.markErr = service.ErrNotificationNotFound

	res := performNotificationRequest(router, http.MethodPatch, "/api/notifications/1/read", "", 7)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestMarkNotificationReadForbidden(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.markErr = service.ErrNotificationForbidden

	res := performNotificationRequest(router, http.MethodPatch, "/api/notifications/1/read", "", 7)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
}

func TestMarkNotificationReadInternalServerError(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.markErr = errors.New("db down")

	res := performNotificationRequest(router, http.MethodPatch, "/api/notifications/1/read", "", 7)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestDeleteNotification(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.notifications = []domain.Notification{{ID: 1, UserID: 7, Message: "A"}}

	res := performNotificationRequest(router, http.MethodDelete, "/api/notifications/1", "", 7)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestDeleteNotificationUnauthorized(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodDelete, "/api/notifications/1", "", nil)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, res.Code)
	}
}

func TestDeleteNotificationInvalidID(t *testing.T) {
	router, _ := setupNotificationTestRouter()

	res := performNotificationRequest(router, http.MethodDelete, "/api/notifications/abc", "", 7)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestDeleteNotificationNotFound(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.deleteErr = service.ErrNotificationNotFound

	res := performNotificationRequest(router, http.MethodDelete, "/api/notifications/1", "", 7)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestDeleteNotificationForbidden(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.deleteErr = service.ErrNotificationForbidden

	res := performNotificationRequest(router, http.MethodDelete, "/api/notifications/1", "", 7)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, res.Code)
	}
}

func TestDeleteNotificationInternalServerError(t *testing.T) {
	router, svc := setupNotificationTestRouter()
	svc.deleteErr = errors.New("db down")

	res := performNotificationRequest(router, http.MethodDelete, "/api/notifications/1", "", 7)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func TestResolveNotificationUserID(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		ok     bool
		userID uint
	}{
		{name: "uint", value: uint(7), ok: true, userID: 7},
		{name: "int positive", value: 9, ok: true, userID: 9},
		{name: "int zero", value: 0, ok: false},
		{name: "float64 positive", value: float64(11), ok: true, userID: 11},
		{name: "float64 zero", value: float64(0), ok: false},
		{name: "string numeric", value: "7", ok: true, userID: 7},
		{name: "unsupported type", value: struct{}{}, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user_id", tt.value)

			got, ok := resolveContextUserID(c)
			if ok != tt.ok {
				t.Fatalf("expected ok=%v, got %v", tt.ok, ok)
			}
			if got != tt.userID {
				t.Fatalf("expected userID=%d, got %d", tt.userID, got)
			}
		})
	}
}

func setupNotificationTestRouter() (*gin.Engine, *mockNotificationService) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	svc := &mockNotificationService{}
	h := NewNotificationHandler(svc)

	api := router.Group("/api")
	api.Use(func(c *gin.Context) {
		if header := c.GetHeader("X-Test-User-ID"); header != "" {
			id, err := strconv.Atoi(header)
			if err == nil {
				c.Set("user_id", id)
			}
		}
		c.Next()
	})

	api.GET("/notifications", h.GetNotifications)
	api.POST("/notifications/send", h.SendNotification)
	api.PATCH("/notifications/:id/read", h.MarkNotificationRead)
	api.DELETE("/notifications/:id", h.DeleteNotification)

	return router, svc
}

func performNotificationRequest(router *gin.Engine, method string, path string, body string, userID interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body == "" {
		reqBody = bytes.NewBuffer(nil)
	} else {
		reqBody = bytes.NewBufferString(body)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if userID != nil {
		req.Header.Set("X-Test-User-ID", strconv.Itoa(userID.(int)))
	}
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}
