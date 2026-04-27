package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"qflow/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var users []models.User
var otps []models.OTP
var jwtSecret = []byte("your-secret-key-change-in-production")

// generateOTP สร้างหมายเลข OTP แบบสุ่ม 6 หลัก
func generateOTP() string {
	digits := "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		b := make([]byte, 1)
		rand.Read(b)
		otp[i] = digits[b[0]%byte(len(digits))]
	}
	return string(otp)
}

// generateJWT สร้าง JWT token สำหรับผู้ใช้ มีอายุ 24 ชั่วโมง
func generateJWT(userID int, phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"phone":  phone,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWT ตรวจสอบความถูกต้องของ JWT token
func validateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RequestOTP จัดการการขอ OTP สำหรับยืนยันตัวตน
// Method: POST
// Role: Guest
func RequestOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// แปลง JSON request body เป็น struct
	var req models.OTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	// สร้าง OTP ใหม่และกำหนดอายุ 5 นาที
	otpCode := generateOTP()
	otp := models.OTP{
		Phone:     req.Phone,
		Code:      otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// ลบ OTP เก่าของเบอร์โทรศัพท์เดิม (ถ้ามี)
	for i, existingOTP := range otps {
		if existingOTP.Phone == req.Phone {
			otps = append(otps[:i], otps[i+1:]...)
			break
		}
	}

	// เพิ่ม OTP ลงในระบบ
	otps = append(otps, otp)

	// ส่ง response กลับไป (ใน production ไม่ควรส่ง OTP กลับไปด้วย)
	response := map[string]string{
		"message": "OTP sent successfully",
		"otp":     otpCode, // ใน production ไม่ควรคืนค่า OTP
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VerifyOTP ตรวจสอบ OTP และสร้าง JWT token
// Method: POST
// Role: Guest
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.OTPVerify
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Phone == "" || req.OTP == "" {
		http.Error(w, "Phone and OTP are required", http.StatusBadRequest)
		return
	}

	// ค้นหาและตรวจสอบความถูกต้องของ OTP
	var validOTP *models.OTP
	for i, otp := range otps {
		if otp.Phone == req.Phone && otp.Code == req.OTP && otp.ExpiresAt.After(time.Now()) {
			validOTP = &otp
			// ลบ OTP ที่ใช้แล้ว
			otps = append(otps[:i], otps[i+1:]...)
			break
		}
	}

	if validOTP == nil {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// ตรวจสอบว่ามีผู้ใช้ในระบบหรือไม่
	var user *models.User
	for i, u := range users {
		if u.Phone == req.Phone {
			user = &users[i]
			break
		}
	}

	// ถ้าไม่มีผู้ใช้ สร้างผู้ใช้ชั่วคราวสำหรับการยืนยัน OTP
	if user == nil {
		tempUser := models.User{
			ID:        len(users) + 1000, // ID ชั่วคราว
			Phone:     req.Phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user = &tempUser
	}

	token, err := generateJWT(user.ID, user.Phone)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register ลงทะเบียนผู้ใช้ใหม่ด้วยการยืนยัน OTP
// Method: POST
// Role: Guest
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ตรวจสอบฟิลด์ที่จำเป็นต้องกรอก
	if req.Phone == "" || req.OTP == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Phone, OTP, firstName, and lastName are required", http.StatusBadRequest)
		return
	}

	// ตรวจสอบ OTP
	var validOTP *models.OTP
	for i, otp := range otps {
		if otp.Phone == req.Phone && otp.Code == req.OTP && otp.ExpiresAt.After(time.Now()) {
			validOTP = &otp
			// ลบ OTP ที่ใช้แล้ว
			otps = append(otps[:i], otps[i+1:]...)
			break
		}
	}

	if validOTP == nil {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// ตรวจสอบว่ามีผู้ใช้อยู่แล้วหรือไม่
	for _, user := range users {
		if user.Phone == req.Phone {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
	}

	// สร้างผู้ใช้ใหม่
	newUser := models.User{
		ID:        len(users) + 1,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, newUser)

	token, err := generateJWT(newUser.ID, newUser.Phone)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  newUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProfileHandler จัดการการดูและแก้ไขโปรไฟล์ผู้ใช้
// GET /api/auth/me - ดูโปรไฟล์ (User)
// PUT /api/auth/me - แก้ไขโปรไฟล์ (User)
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := validateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// ดึงข้อมูลผู้ใช้จาก JWT token
	userID := (*claims)["userId"].(float64)

	// GET /api/auth/me - ดูข้อมูลโปรไฟล์ผู้ใช้
	if r.Method == http.MethodGet {
		for _, user := range users {
			if user.ID == int(userID) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// PUT /api/auth/me - แก้ไขข้อมูลโปรไฟล์ผู้ใช้
	if r.Method == http.MethodPut {
		var req models.UpdateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// อัปเดตข้อมูลผู้ใช้
		for i, user := range users {
			if user.ID == int(userID) {
				if req.FirstName != "" {
					users[i].FirstName = req.FirstName
				}
				if req.LastName != "" {
					users[i].LastName = req.LastName
				}
				if req.Email != "" {
					users[i].Email = req.Email
				}
				users[i].UpdatedAt = time.Now()

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(users[i])
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
