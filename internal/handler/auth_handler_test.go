package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"qflow/internal/domain"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	otp     *domain.OTP
	user    *domain.User
	token   string
	profile *domain.User
	err     error
}

func (m *mockAuthService) RequestOTP(phone string) (*domain.OTP, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.otp, nil
}

func (m *mockAuthService) VerifyOTP(phone, code string) (*domain.User, string, error) {
	if m.err != nil {
		return nil, "", m.err
	}
	return m.user, m.token, nil
}

func (m *mockAuthService) RegisterUser(phone, name, role, otpCode string) (*domain.User, string, error) {
	if m.err != nil {
		return nil, "", m.err
	}
	return m.user, m.token, nil
}

func (m *mockAuthService) GetUserProfile(userID uint) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.profile, nil
}

func (m *mockAuthService) UpdateUserProfile(userID uint, name, role string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.profile != nil {
		if name != "" {
			m.profile.Name = name
		}
	}
	return m.profile, nil
}

func TestRequestOTP(t *testing.T) {
	router, _ := setupAuthTestRouter()

	res := performAuthRequest(router, http.MethodPost, "/api/auth/request-otp", `{"phone":"0812345678"}`)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "OTP sent successfully" {
		t.Fatalf("unexpected message: %v", response["message"])
	}
}

func TestRequestOTPWithInvalidFormat(t *testing.T) {
	router, _ := setupAuthTestRouter()

	res := performAuthRequest(router, http.MethodPost, "/api/auth/request-otp", `{"invalid":"json"}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestVerifyOTP(t *testing.T) {
	router, _ := setupAuthTestRouter()

	res := performAuthRequest(router, http.MethodPost, "/api/auth/verify-otp", `{"phone":"0812345678","code":"123456"}`)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "OTP verified successfully" {
		t.Fatalf("unexpected message: %v", response["message"])
	}

	if response["token"] == nil {
		t.Fatalf("expected token in response")
	}
}

func TestVerifyOTPWithInvalidFormat(t *testing.T) {
	router, _ := setupAuthTestRouter()

	res := performAuthRequest(router, http.MethodPost, "/api/auth/verify-otp", `{"invalid":"json"}`)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestRegister(t *testing.T) {
	router, _ := setupAuthTestRouter()

	res := performAuthRequest(router, http.MethodPost, "/api/auth/register", `{"phone":"0812345678","name":"Test User","otp_code":"123456"}`)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "User registered successfully" {
		t.Fatalf("unexpected message: %v", response["message"])
	}
}

func TestGetProfile(t *testing.T) {
	router, svc := setupAuthTestRouter()
	svc.profile = &domain.User{ID: 1, Phone: "0812345678", Name: "Test User", Role: "user"}

	res := performAuthRequestWithAuth(router, http.MethodGet, "/api/auth/me", "")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	user := response["user"].(map[string]interface{})
	if user["name"] != "Test User" {
		t.Fatalf("unexpected user name: %v", user["name"])
	}
}

func TestUpdateProfile(t *testing.T) {
	router, svc := setupAuthTestRouter()
	svc.profile = &domain.User{ID: 1, Phone: "0812345678", Name: "Test User", Role: "user"}

	res := performAuthRequestWithAuth(router, http.MethodPut, "/api/auth/me", `{"name":"Updated Name"}`)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response["message"] != "Profile updated successfully" {
		t.Fatalf("unexpected message: %v", response["message"])
	}
}

func TestAuthHandlerReturnsInternalServerError(t *testing.T) {
	router, svc := setupAuthTestRouter()
	svc.err = errors.New("service error")

	res := performAuthRequest(router, http.MethodPost, "/api/auth/request-otp", `{"phone":"0812345678"}`)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
}

func setupAuthTestRouter() (*gin.Engine, *mockAuthService) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	svc := &mockAuthService{
		otp:     &domain.OTP{ID: 1, Phone: "0812345678", Code: "123456"},
		user:    &domain.User{ID: 1, Phone: "0812345678", Name: "Test User", Role: "user"},
		token:   "test-token",
		profile: &domain.User{ID: 1, Phone: "0812345678", Name: "Test User", Role: "user"},
	}
	auth := NewAuthHandler(svc, false)
	api := router.Group("/api")
	api.POST("/auth/request-otp", auth.RequestOTP)
	api.POST("/auth/verify-otp", auth.VerifyOTP)
	api.POST("/auth/register", auth.Register)

	// Protected routes require user_id in context (simulates JWT middleware)
	protected := api.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	protected.GET("/auth/me", auth.GetProfile)
	protected.PUT("/auth/me", auth.UpdateProfile)

	return router, svc
}

func performAuthRequest(router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
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

func performAuthRequestWithAuth(router *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body == "" {
		reqBody = bytes.NewBuffer(nil)
	} else {
		reqBody = bytes.NewBufferString(body)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}
