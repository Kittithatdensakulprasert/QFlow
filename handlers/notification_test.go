package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetNotifications() {
	notifications = []Notification{}
}

// --- GetNotifications ---

func TestGetNotifications_Empty(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	w := httptest.NewRecorder()
	GetNotifications(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []Notification
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestGetNotifications_WithData(t *testing.T) {
	resetNotifications()
	notifications = []Notification{
		{ID: 1, Message: "test", Read: false},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	w := httptest.NewRecorder()
	GetNotifications(w, req)

	var result []Notification
	json.NewDecoder(w.Body).Decode(&result)
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

// --- SendNotification ---

func TestSendNotification_Success(t *testing.T) {
	resetNotifications()
	body := `{"message":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/send", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	SendNotification(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	var result Notification
	json.NewDecoder(w.Body).Decode(&result)
	if result.Message != "hello" {
		t.Errorf("expected message 'hello', got '%s'", result.Message)
	}
	if result.Read != false {
		t.Error("expected Read to be false")
	}
}

func TestSendNotification_EmptyMessage(t *testing.T) {
	resetNotifications()
	body := `{"message":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/send", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	SendNotification(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSendNotification_InvalidBody(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/send", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	SendNotification(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSendNotification_WrongMethod(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications/send", nil)
	w := httptest.NewRecorder()
	SendNotification(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// --- PATCH /api/notifications/:id/read ---

func TestNotificationHandler_ReadSuccess(t *testing.T) {
	resetNotifications()
	notifications = []Notification{{ID: 1, Message: "msg", Read: false}}
	req := httptest.NewRequest(http.MethodPatch, "/api/notifications/1/read", nil)
	w := httptest.NewRecorder()
	NotificationHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result Notification
	json.NewDecoder(w.Body).Decode(&result)
	if !result.Read {
		t.Error("expected Read to be true")
	}
}

func TestNotificationHandler_ReadNotFound(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodPatch, "/api/notifications/999/read", nil)
	w := httptest.NewRecorder()
	NotificationHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestNotificationHandler_ReadInvalidID(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodPatch, "/api/notifications/abc/read", nil)
	w := httptest.NewRecorder()
	NotificationHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- DELETE /api/notifications/:id ---

func TestNotificationHandler_DeleteSuccess(t *testing.T) {
	resetNotifications()
	notifications = []Notification{{ID: 1, Message: "msg", Read: false}}
	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/1", nil)
	w := httptest.NewRecorder()
	NotificationHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if len(notifications) != 0 {
		t.Error("expected notifications to be empty after delete")
	}
}

func TestNotificationHandler_DeleteNotFound(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/999", nil)
	w := httptest.NewRecorder()
	NotificationHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestNotificationHandler_DeleteInvalidID(t *testing.T) {
	resetNotifications()
	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/abc", nil)
	w := httptest.NewRecorder()
	NotificationHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
