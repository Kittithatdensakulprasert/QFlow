package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	queues = nil
	r := gin.Default()
	api := r.Group("/api/queues")
	{
		api.POST("/book", BookQueue)
		api.GET("/history", GetHistory)
		api.GET("/:queueNumber", GetQueue)
		api.PATCH("/:id/cancel", CancelQueue)
	}
	return r
}

func TestBookQueue_Success(t *testing.T) {
	router := setupRouter()
	body := `{"name":"test"}`
	req, _ := http.NewRequest("POST", "/api/queues/book", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 got %d", w.Code)
	}
}

func TestBookQueue_InvalidBody(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("POST", "/api/queues/book", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}

func TestGetHistory_Empty(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/queues/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", w.Code)
	}
}

func TestGetHistory_WithData(t *testing.T) {
	router := setupRouter()
	bookReq, _ := http.NewRequest("POST", "/api/queues/book", strings.NewReader(`{"name":"A"}`))
	bookReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), bookReq)
	req, _ := http.NewRequest("GET", "/api/queues/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", w.Code)
	}
}

func TestGetQueue_Found(t *testing.T) {
	router := setupRouter()
	bookReq, _ := http.NewRequest("POST", "/api/queues/book", strings.NewReader(`{"name":"test"}`))
	bookReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), bookReq)
	req, _ := http.NewRequest("GET", "/api/queues/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", w.Code)
	}
}

func TestGetQueue_NotFound(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/queues/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 got %d", w.Code)
	}
}

func TestGetQueue_InvalidNumber(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/queues/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}

func TestCancelQueue_Success(t *testing.T) {
	router := setupRouter()
	bookReq, _ := http.NewRequest("POST", "/api/queues/book", strings.NewReader(`{"name":"test"}`))
	bookReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), bookReq)
	req, _ := http.NewRequest("PATCH", "/api/queues/1/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", w.Code)
	}
}

func TestCancelQueue_NotFound(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("PATCH", "/api/queues/999/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 got %d", w.Code)
	}
}

func TestCancelQueue_AlreadyCancelled(t *testing.T) {
	router := setupRouter()
	bookReq, _ := http.NewRequest("POST", "/api/queues/book", strings.NewReader(`{"name":"test"}`))
	bookReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), bookReq)
	req1, _ := http.NewRequest("PATCH", "/api/queues/1/cancel", nil)
	router.ServeHTTP(httptest.NewRecorder(), req1)
	req2, _ := http.NewRequest("PATCH", "/api/queues/1/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req2)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}

func TestCancelQueue_InvalidID(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("PATCH", "/api/queues/abc/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d", w.Code)
	}
}
