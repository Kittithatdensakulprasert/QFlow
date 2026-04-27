package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qflow/models"
	"strings"
	"testing"
	"time"
)

// ทำความสะอาดข้อมูลก่อนทำแต่ละ test
func setup() {
	users = []models.User{}
	otps = []models.OTP{}
}

func TestGenerateOTP(t *testing.T) {
	otp := generateOTP()
	if len(otp) != 6 {
		t.Errorf("Expected OTP length 6, got %d", len(otp))
	}
	
	// ตรวจสอบว่า OTP ประกอบด้วยตัวเลขเท่านั้น
	for _, char := range otp {
		if char < '0' || char > '9' {
			t.Errorf("OTP should contain only digits, got %c", char)
		}
	}
}

func TestRequestOTP(t *testing.T) {
	setup()
	
	// สร้าง request สำหรับขอ OTP
	requestBody := models.OTPRequest{Phone: "0812345678"}
	body, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/request-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	RequestOTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if response["message"] != "OTP sent successfully" {
		t.Errorf("Expected success message, got %s", response["message"])
	}
	
	if len(response["otp"]) != 6 {
		t.Errorf("Expected OTP length 6, got %d", len(response["otp"]))
	}
}

func TestRequestOTPInvalidMethod(t *testing.T) {
	setup()
	
	req := httptest.NewRequest(http.MethodGet, "/api/auth/request-otp", nil)
	w := httptest.NewRecorder()
	
	RequestOTP(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestRequestOTPInvalidBody(t *testing.T) {
	setup()
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/request-otp", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	RequestOTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRequestOTPMissingPhone(t *testing.T) {
	setup()
	
	requestBody := models.OTPRequest{Phone: ""}
	body, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/request-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	RequestOTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestVerifyOTP(t *testing.T) {
	setup()
	
	// ขอ OTP ก่อน
	otpCode := generateOTP()
	otp := models.OTP{
		Phone:     "0812345678",
		Code:      otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	otps = append(otps, otp)
	
	// ยืนยัน OTP
	requestBody := models.OTPVerify{Phone: "0812345678", OTP: otpCode}
	body, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/verify-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	VerifyOTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if response.Token == "" {
		t.Error("Expected JWT token in response")
	}
	
	if response.User.Phone != "0812345678" {
		t.Errorf("Expected phone 0812345678, got %s", response.User.Phone)
	}
}

func TestVerifyOTPInvalidOTP(t *testing.T) {
	setup()
	
	requestBody := models.OTPVerify{Phone: "0812345678", OTP: "999999"}
	body, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/verify-otp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	VerifyOTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRegister(t *testing.T) {
	setup()
	
	// ขอ OTP ก่อน
	otpCode := generateOTP()
	otp := models.OTP{
		Phone:     "0812345678",
		Code:      otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	otps = append(otps, otp)
	
	// ลงทะเบียนผู้ใช้ใหม่
	requestBody := models.RegisterRequest{
		Phone:     "0812345678",
		OTP:       otpCode,
		FirstName: "สมชาย",
		LastName:  "ใจดี",
		Email:     "somchai@example.com",
	}
	body, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	Register(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if response.User.FirstName != "สมชาย" {
		t.Errorf("Expected firstName สมชาย, got %s", response.User.FirstName)
	}
	
	if response.User.LastName != "ใจดี" {
		t.Errorf("Expected lastName ใจดี, got %s", response.User.LastName)
	}
	
	if len(users) != 1 {
		t.Errorf("Expected 1 user in system, got %d", len(users))
	}
}

func TestRegisterUserExists(t *testing.T) {
	setup()
	
	// เพิ่มผู้ใช้ที่มีอยู่แล้ว
	existingUser := models.User{
		ID:        1,
		Phone:     "0812345678",
		FirstName: "มีอยู่แล้ว",
		LastName:  "ในระบบ",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	users = append(users, existingUser)
	
	// ขอ OTP
	otpCode := generateOTP()
	otp := models.OTP{
		Phone:     "0812345678",
		Code:      otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	otps = append(otps, otp)
	
	// พยายามลงทะเบียนซ้ำ
	requestBody := models.RegisterRequest{
		Phone:     "0812345678",
		OTP:       otpCode,
		FirstName: "สมชาย",
		LastName:  "ใจดี",
	}
	body, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	Register(w, req)
	
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestGenerateJWT(t *testing.T) {
	token, err := generateJWT(1, "0812345678")
	if err != nil {
		t.Errorf("Failed to generate JWT: %v", err)
	}
	
	if token == "" {
		t.Error("Expected non-empty token")
	}
	
	// ตรวจสอบว่า token สามารถ validate ได้
	claims, err := validateJWT(token)
	if err != nil {
		t.Errorf("Failed to validate generated JWT: %v", err)
	}
	
	if (*claims)["userId"] != float64(1) {
		t.Errorf("Expected userId 1, got %v", (*claims)["userId"])
	}
	
	if (*claims)["phone"] != "0812345678" {
		t.Errorf("Expected phone 0812345678, got %v", (*claims)["phone"])
	}
}

func TestValidateJWTInvalid(t *testing.T) {
	invalidToken := "invalid.jwt.token"
	_, err := validateJWT(invalidToken)
	
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestProfileHandlerGet(t *testing.T) {
	setup()
	
	// เพิ่มผู้ใช้
	user := models.User{
		ID:        1,
		Phone:     "0812345678",
		FirstName: "สมชาย",
		LastName:  "ใจดี",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	users = append(users, user)
	
	// สร้าง JWT token
	token, err := generateJWT(1, "0812345678")
	if err != nil {
		t.Errorf("Failed to generate JWT: %v", err)
	}
	
	// ทดสอบ GET /api/auth/me
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	
	ProfileHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var responseUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if responseUser.FirstName != "สมชาย" {
		t.Errorf("Expected firstName สมชาย, got %s", responseUser.FirstName)
	}
}

func TestProfileHandlerUpdate(t *testing.T) {
	setup()
	
	// เพิ่มผู้ใช้
	user := models.User{
		ID:        1,
		Phone:     "0812345678",
		FirstName: "สมชาย",
		LastName:  "ใจดี",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	users = append(users, user)
	
	// สร้าง JWT token
	token, err := generateJWT(1, "0812345678")
	if err != nil {
		t.Errorf("Failed to generate JWT: %v", err)
	}
	
	// ทดสอบ PUT /api/auth/me
	updateRequest := models.UpdateProfileRequest{
		FirstName: "สมศรี",
		LastName:  "รักดี",
		Email:     "somsri@example.com",
	}
	body, _ := json.Marshal(updateRequest)
	
	req := httptest.NewRequest(http.MethodPut, "/api/auth/me", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	ProfileHandler(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var responseUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if responseUser.FirstName != "สมศรี" {
		t.Errorf("Expected updated firstName สมศรี, got %s", responseUser.FirstName)
	}
	
	if responseUser.LastName != "รักดี" {
		t.Errorf("Expected updated lastName รักดี, got %s", responseUser.LastName)
	}
}

func TestProfileHandlerUnauthorized(t *testing.T) {
	setup()
	
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()
	
	ProfileHandler(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProfileHandlerInvalidToken(t *testing.T) {
	setup()
	
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	w := httptest.NewRecorder()
	
	ProfileHandler(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
