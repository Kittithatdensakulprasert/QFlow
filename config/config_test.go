package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("OTP_CLEANUP_INTERVAL")

	cfg := Load()

	if cfg.Port != "3000" {
		t.Errorf("expected port 3000, got %s", cfg.Port)
	}
	if cfg.DSN != "" {
		t.Errorf("expected empty DSN, got %s", cfg.DSN)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("expected production, got %s", cfg.AppEnv)
	}
	if cfg.OTPCleanupInterval != "1h" {
		t.Errorf("expected 1h, got %s", cfg.OTPCleanupInterval)
	}
	if cfg.BootstrapAdminName != "Bootstrap Admin" {
		t.Errorf("expected Bootstrap Admin, got %s", cfg.BootstrapAdminName)
	}
	if cfg.BootstrapProviderName != "Bootstrap Provider" {
		t.Errorf("expected Bootstrap Provider, got %s", cfg.BootstrapProviderName)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://localhost/test")
	os.Setenv("JWT_SECRET", "mysecret")
	os.Setenv("APP_ENV", "dev")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("APP_ENV")
	}()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected 8080, got %s", cfg.Port)
	}
	if cfg.DSN != "postgres://localhost/test" {
		t.Errorf("expected postgres URL, got %s", cfg.DSN)
	}
	if cfg.JWTSecret != "mysecret" {
		t.Errorf("expected mysecret, got %s", cfg.JWTSecret)
	}
	if cfg.AppEnv != "dev" {
		t.Errorf("expected dev, got %s", cfg.AppEnv)
	}
}

func TestExposeOTPInResponse(t *testing.T) {
	cases := []struct {
		env      string
		expected bool
	}{
		{"dev", true},
		{"development", true},
		{"test", true},
		{"DEV", true},
		{"DEVELOPMENT", true},
		{"TEST", true},
		{"production", false},
		{"staging", false},
		{"", false},
	}

	for _, tc := range cases {
		cfg := &Config{AppEnv: tc.env}
		got := cfg.ExposeOTPInResponse()
		if got != tc.expected {
			t.Errorf("AppEnv=%q: expected %v, got %v", tc.env, tc.expected, got)
		}
	}
}

func TestParsedOTPCleanupInterval_Valid(t *testing.T) {
	cfg := &Config{OTPCleanupInterval: "30m"}
	d := cfg.ParsedOTPCleanupInterval()
	if d != 30*time.Minute {
		t.Errorf("expected 30m, got %v", d)
	}
}

func TestParsedOTPCleanupInterval_Invalid(t *testing.T) {
	cfg := &Config{OTPCleanupInterval: "invalid"}
	d := cfg.ParsedOTPCleanupInterval()
	if d != time.Hour {
		t.Errorf("expected fallback 1h, got %v", d)
	}
}

func TestParsedOTPCleanupInterval_Various(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Duration
	}{
		{"1h", time.Hour},
		{"24h", 24 * time.Hour},
		{"5m", 5 * time.Minute},
		{"10s", 10 * time.Second},
	}

	for _, tc := range cases {
		cfg := &Config{OTPCleanupInterval: tc.input}
		got := cfg.ParsedOTPCleanupInterval()
		if got != tc.expected {
			t.Errorf("input=%q: expected %v, got %v", tc.input, tc.expected, got)
		}
	}
}
