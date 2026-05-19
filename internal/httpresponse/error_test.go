package httpresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/test", nil)

	Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "phone is required")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["error"] != "VALIDATION_ERROR" {
		t.Errorf("expected error VALIDATION_ERROR, got %v", body["error"])
	}
	if body["message"] != "phone is required" {
		t.Errorf("expected message 'phone is required', got %v", body["message"])
	}
	if body["path"] != "/api/test" {
		t.Errorf("expected path /api/test, got %v", body["path"])
	}
}

func TestError_StatusCodes(t *testing.T) {
	cases := []struct {
		status int
		code   string
	}{
		{http.StatusUnauthorized, "UNAUTHORIZED"},
		{http.StatusForbidden, "FORBIDDEN"},
		{http.StatusNotFound, "NOT_FOUND"},
		{http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		Error(c, tc.status, tc.code, "test message")

		if w.Code != tc.status {
			t.Errorf("expected %d, got %d", tc.status, w.Code)
		}
	}
}
