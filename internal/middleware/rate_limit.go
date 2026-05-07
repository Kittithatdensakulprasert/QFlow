package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

type RateLimiter struct {
	mu     sync.Mutex
	store  map[string]*rateLimitEntry
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		store:  make(map[string]*rateLimitEntry),
		limit:  limit,
		window: window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.store {
			if now.After(entry.windowEnd) {
				delete(rl.store, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.store[key]

	if !exists || now.After(entry.windowEnd) {
		rl.store[key] = &rateLimitEntry{
			count:     1,
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// OTPRequestLimiter limits OTP requests: 5 per 10 minutes per IP
var OTPRequestLimiter = NewRateLimiter(5, 10*time.Minute)

// OTPVerifyLimiter limits OTP verify attempts: 10 per 10 minutes per IP
var OTPVerifyLimiter = NewRateLimiter(10, 10*time.Minute)
