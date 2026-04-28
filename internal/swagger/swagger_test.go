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
	require.Nil(t, requestBodyForRoute(http.MethodGet, "/api/categories"))
	require.NotNil(t, requestBodyForRoute(http.MethodPost, "/api/providers"))
	require.NotNil(t, requestBodyForRoute(http.MethodPut, "/api/categories/{id}"))
	require.Nil(t, requestBodyForRoute(http.MethodPatch, "/api/zones/{id}/toggle"))
	require.NotNil(t, requestBodyForRoute(http.MethodPatch, "/api/example"))
}

func TestKnownRouteUsesDetailedRequestSchema(t *testing.T) {
	body := requestBodyForRoute(http.MethodPost, "/api/providers")

	require.NotNil(t, body)
	require.Equal(t, "#/components/schemas/CreateProviderRequest", body.Content["application/json"].Schema.Ref)
}

func TestPathParameterSchemaDetection(t *testing.T) {
	idSchema := schemaForParameter("id")
	require.Equal(t, "integer", idSchema.Type)
	require.Equal(t, "uint", idSchema.Format)

	zoneIDSchema := schemaForParameter("zoneId")
	require.Equal(t, "integer", zoneIDSchema.Type)
	require.Equal(t, "uint", zoneIDSchema.Format)

	queueNumberSchema := schemaForParameter("queueNumber")
	require.Equal(t, "integer", queueNumberSchema.Type)
	require.Equal(t, "int32", queueNumberSchema.Format)

	slugSchema := schemaForParameter("slug")
	require.Equal(t, "string", slugSchema.Type)
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
	require.Contains(t, docRecorder.Body.String(), "components")
	require.Contains(t, docRecorder.Body.String(), "ErrorResponse")

	uiRecorder := httptest.NewRecorder()
	uiRequest := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	r.ServeHTTP(uiRecorder, uiRequest)
	require.Equal(t, http.StatusOK, uiRecorder.Code)
	require.Contains(t, uiRecorder.Body.String(), "swagger-ui")
}
