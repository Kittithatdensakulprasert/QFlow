package swagger

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type openAPIDocument struct {
	OpenAPI    string                 `json:"openapi"`
	Info       openAPIInfo            `json:"info"`
	Servers    []openAPIServer        `json:"servers"`
	Paths      map[string]openAPIPath `json:"paths"`
	Components openAPIComponents      `json:"components"`
}

type openAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type openAPIServer struct {
	URL string `json:"url"`
}

type openAPIPath map[string]openAPIOperation

type openAPIOperation struct {
	Tags        []string                   `json:"tags,omitempty"`
	Summary     string                     `json:"summary"`
	OperationID string                     `json:"operationId"`
	Parameters  []openAPIParameter         `json:"parameters,omitempty"`
	RequestBody *openAPIRequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]openAPIResponse `json:"responses"`
}

type openAPIParameter struct {
	Name        string        `json:"name"`
	In          string        `json:"in"`
	Required    bool          `json:"required"`
	Description string        `json:"description,omitempty"`
	Schema      openAPISchema `json:"schema"`
}

type openAPIRequestBody struct {
	Required bool                          `json:"required"`
	Content  map[string]openAPIMediaSchema `json:"content"`
}

type openAPIMediaSchema struct {
	Schema openAPISchema `json:"schema"`
}

type openAPIResponse struct {
	Description string                        `json:"description"`
	Content     map[string]openAPIMediaSchema `json:"content,omitempty"`
}

type openAPIComponents struct {
	SecuritySchemes map[string]openAPISecurityScheme `json:"securitySchemes"`
}

type openAPISecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

type openAPISchema struct {
	Type                 string                   `json:"type,omitempty"`
	Format               string                   `json:"format,omitempty"`
	AdditionalProperties *openAPISchema           `json:"additionalProperties,omitempty"`
	Properties           map[string]openAPISchema `json:"properties,omitempty"`
}

func Register(r *gin.Engine) {
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, buildDocument(r, c.Request.Host))
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi.json"),
		ginSwagger.DeepLinking(true),
		ginSwagger.PersistAuthorization(true),
	))
}

func buildDocument(r *gin.Engine, host string) openAPIDocument {
	paths := make(map[string]openAPIPath)
	routes := r.Routes()

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	for _, route := range routes {
		if strings.HasPrefix(route.Path, "/swagger") || route.Path == "/openapi.json" {
			continue
		}

		path := ginPathToOpenAPIPath(route.Path)
		if _, ok := paths[path]; !ok {
			paths[path] = make(openAPIPath)
		}

		paths[path][strings.ToLower(route.Method)] = openAPIOperation{
			Tags:        []string{tagForPath(path)},
			Summary:     summaryForRoute(route.Method, path),
			OperationID: operationID(route.Method, path),
			Parameters:  parametersForPath(path),
			RequestBody: requestBodyForMethod(route.Method),
			Responses: map[string]openAPIResponse{
				"200": responseWithJSON("Successful response"),
				"400": responseWithJSON("Bad request"),
				"500": responseWithJSON("Internal server error"),
			},
		}
	}

	return openAPIDocument{
		OpenAPI: "3.0.3",
		Info: openAPIInfo{
			Title:       "QFlow API",
			Description: "Runtime-generated Swagger documentation from registered Gin routes.",
			Version:     "1.0.0",
		},
		Servers: []openAPIServer{
			{URL: "http://" + host},
		},
		Paths: paths,
		Components: openAPIComponents{
			SecuritySchemes: map[string]openAPISecurityScheme{
				"BearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		},
	}
}

func ginPathToOpenAPIPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + strings.TrimPrefix(part, ":") + "}"
		}
	}
	return strings.Join(parts, "/")
}

func tagForPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "api" {
		return title(parts[1])
	}
	return "Default"
}

func title(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func summaryForRoute(method, path string) string {
	return method + " " + path
}

func operationID(method, path string) string {
	replacer := strings.NewReplacer("/", "_", "{", "", "}", "", "-", "_")
	return strings.Trim(replacer.Replace(strings.ToLower(method)+path), "_")
}

func parametersForPath(path string) []openAPIParameter {
	var params []openAPIParameter
	for _, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			name := strings.Trim(part, "{}")
			params = append(params, openAPIParameter{
				Name:        name,
				In:          "path",
				Required:    true,
				Description: "Path parameter: " + name,
				Schema:      openAPISchema{Type: schemaTypeForParameter(name)},
			})
		}
	}
	return params
}

func schemaTypeForParameter(name string) string {
	if strings.Contains(strings.ToLower(name), "number") {
		return "integer"
	}
	return "integer"
}

func requestBodyForMethod(method string) *openAPIRequestBody {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return &openAPIRequestBody{
			Required: true,
			Content: map[string]openAPIMediaSchema{
				"application/json": {
					Schema: openAPISchema{
						Type:                 "object",
						AdditionalProperties: &openAPISchema{},
					},
				},
			},
		}
	default:
		return nil
	}
}

func responseWithJSON(description string) openAPIResponse {
	return openAPIResponse{
		Description: description,
		Content: map[string]openAPIMediaSchema{
			"application/json": {
				Schema: openAPISchema{
					Type:                 "object",
					AdditionalProperties: &openAPISchema{},
				},
			},
		},
	}
}
