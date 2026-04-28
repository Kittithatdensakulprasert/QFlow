package swagger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestBuildDocumentUsesRegisteredGinRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/categories/:id", func(c *gin.Context) {})
	r.POST("/api/providers", func(c *gin.Context) {})
	Register(r)

	doc := buildDocument(r, "localhost:3000")

	require.Contains(t, doc.Paths, "/api/categories/{id}")
	require.Contains(t, doc.Paths["/api/categories/{id}"], "get")
	require.Contains(t, doc.Paths, "/api/providers")
	require.Contains(t, doc.Paths["/api/providers"], "post")
	require.NotContains(t, doc.Paths, "/swagger/index.html")
	require.NotContains(t, doc.Paths, "/openapi.json")
}

func TestRequestBodyOnlyAddedForBodyMethods(t *testing.T) {
	require.Nil(t, requestBodyForMethod(http.MethodGet))
	require.NotNil(t, requestBodyForMethod(http.MethodPost))
	require.NotNil(t, requestBodyForMethod(http.MethodPut))
	require.NotNil(t, requestBodyForMethod(http.MethodPatch))
}

func TestRegisterServesSwaggerEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/categories", func(c *gin.Context) {})
	Register(r)

	docRecorder := httptest.NewRecorder()
	docRequest := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	r.ServeHTTP(docRecorder, docRequest)
	require.Equal(t, http.StatusOK, docRecorder.Code)
	require.Contains(t, docRecorder.Body.String(), "/api/categories")

	uiRecorder := httptest.NewRecorder()
	uiRequest := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	r.ServeHTTP(uiRecorder, uiRequest)
	require.Equal(t, http.StatusOK, uiRecorder.Code)
	require.Contains(t, uiRecorder.Body.String(), "swagger-ui")
}
