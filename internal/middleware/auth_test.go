package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"qflow/internal/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupJWT() *jwt.JWTManager {
	return jwt.NewJWTManager("testsecretkey")
}

func TestJWTAuth_NoHeader(t *testing.T) {
	jm := setupJWT()
	r := gin.New()
	r.GET("/", JWTAuth(jm), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_InvalidFormat(t *testing.T) {
	jm := setupJWT()
	r := gin.New()
	r.GET("/", JWTAuth(jm), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	cases := []string{
		"InvalidToken",
		"Token abc",
		"Bearer",
		"Bearer a b",
	}

	for _, h := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", h)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "header: %s", h)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	jm := setupJWT()
	r := gin.New()
	r.GET("/", JWTAuth(jm), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer not.valid.token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_ValidToken(t *testing.T) {
	jm := setupJWT()
	token, _ := jm.GenerateToken(42, "0812345678", "user")

	r := gin.New()
	r.GET("/", JWTAuth(jm), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_NoRoleInContext(t *testing.T) {
	r := gin.New()
	r.GET("/", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRole_WrongRole(t *testing.T) {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	}, RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_CorrectRole(t *testing.T) {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	}, RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_MultipleRoles(t *testing.T) {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Set("role", "provider")
		c.Next()
	}, RequireRole("admin", "provider"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_InvalidRoleType(t *testing.T) {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Set("role", 12345) // wrong type
		c.Next()
	}, RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
