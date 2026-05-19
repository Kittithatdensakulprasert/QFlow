package jwt

import (
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

func TestNewJWTManager(t *testing.T) {
	m := NewJWTManager("mysecret")
	if m == nil {
		t.Error("expected non-nil JWTManager")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	m := NewJWTManager("supersecretkey")

	token, err := m.GenerateToken(1, "0812345678", "user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := m.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", claims.UserID)
	}
	if claims.Phone != "0812345678" {
		t.Errorf("expected phone 0812345678, got %s", claims.Phone)
	}
	if claims.Role != "user" {
		t.Errorf("expected role user, got %s", claims.Role)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	m := NewJWTManager("supersecretkey")

	_, err := m.ValidateToken("not.a.valid.token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	m1 := NewJWTManager("secret1")
	m2 := NewJWTManager("secret2")

	token, _ := m1.GenerateToken(1, "0812345678", "user")
	_, err := m2.ValidateToken(token)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	m := NewJWTManager("supersecretkey")

	claims := &Claims{
		UserID: 1,
		Phone:  "0812345678",
		Role:   "user",
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  gojwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("supersecretkey"))

	_, err := m.ValidateToken(tokenStr)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateToken_Malformed(t *testing.T) {
	m := NewJWTManager("supersecretkey")

	_, err := m.ValidateToken("aaa.bbb.ccc")
	if err == nil {
		t.Error("expected error for malformed token")
	}
}
