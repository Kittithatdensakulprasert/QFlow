package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(3, 10*time.Second)

	assert.True(t, rl.Allow("ip1"))
	assert.True(t, rl.Allow("ip1"))
	assert.True(t, rl.Allow("ip1"))
	assert.False(t, rl.Allow("ip1")) // 4th should be blocked
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rl := NewRateLimiter(2, 10*time.Second)

	assert.True(t, rl.Allow("ip1"))
	assert.True(t, rl.Allow("ip1"))
	assert.False(t, rl.Allow("ip1"))

	// Different IP should have its own counter
	assert.True(t, rl.Allow("ip2"))
	assert.True(t, rl.Allow("ip2"))
	assert.False(t, rl.Allow("ip2"))
}

func TestRateLimiter_WindowReset(t *testing.T) {
	rl := NewRateLimiter(2, 50*time.Millisecond)

	assert.True(t, rl.Allow("ip1"))
	assert.True(t, rl.Allow("ip1"))
	assert.False(t, rl.Allow("ip1"))

	time.Sleep(60 * time.Millisecond)

	// Window should have reset
	assert.True(t, rl.Allow("ip1"))
}

func TestRateLimiter_Middleware_Allow(t *testing.T) {
	rl := NewRateLimiter(3, 10*time.Second)

	r := gin.New()
	r.POST("/test", rl.Middleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimiter_Middleware_Block(t *testing.T) {
	rl := NewRateLimiter(2, 10*time.Second)

	r := gin.New()
	r.POST("/test", rl.Middleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 3rd request should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
