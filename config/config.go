package config

import (
	"os"
	"strings"
)

type Config struct {
	Port      string
	DSN       string
	JWTSecret string
	AppEnv    string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "3000"),
		DSN:       getEnv("DATABASE_URL", ""),
		JWTSecret: getEnv("JWT_SECRET", ""),
		AppEnv:    getEnv("APP_ENV", "production"),
	}
}

func (c *Config) ExposeOTPInResponse() bool {
	env := strings.ToLower(c.AppEnv)
	return env == "dev" || env == "development" || env == "test"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
