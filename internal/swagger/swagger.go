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
	Schemas         map[string]openAPISchema         `json:"schemas"`
}

type openAPISecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

type openAPISchema struct {
	Ref                  string                   `json:"$ref,omitempty"`
	Type                 string                   `json:"type,omitempty"`
	Format               string                   `json:"format,omitempty"`
	Required             []string                 `json:"required,omitempty"`
	Enum                 []string                 `json:"enum,omitempty"`
	Items                *openAPISchema           `json:"items,omitempty"`
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
			RequestBody: requestBodyForRoute(route.Method, path),
			Responses:   responsesForRoute(route.Method, path),
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
			Schemas: componentSchemas(),
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
				Schema:      schemaForParameter(name),
			})
		}
	}
	return params
}

func schemaForParameter(name string) openAPISchema {
	normalized := strings.ToLower(name)
	switch {
	case strings.Contains(normalized, "number"):
		return openAPISchema{Type: "integer", Format: "int32"}
	case strings.HasSuffix(normalized, "id") || normalized == "id":
		return openAPISchema{Type: "integer", Format: "uint"}
	default:
		return openAPISchema{Type: "string"}
	}
}

func requestBodyForRoute(method, path string) *openAPIRequestBody {
	if routeHasNoBody(method, path) {
		return nil
	}

	if schema, ok := requestSchemaForRoute(method, path); ok {
		return jsonRequestBody(schema)
	}

	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return jsonRequestBody(freeFormObjectSchema())
	default:
		return nil
	}
}

func requestSchemaForRoute(method, path string) (openAPISchema, bool) {
	schemas := map[string]openAPISchema{
		"POST /api/auth/request-otp":     schemaRef("RequestOTPRequest"),
		"POST /api/auth/verify-otp":      schemaRef("VerifyOTPRequest"),
		"POST /api/auth/register":        schemaRef("RegisterRequest"),
		"PUT /api/auth/me":               schemaRef("UpdateProfileRequest"),
		"POST /api/categories":           schemaRef("CategoryRequest"),
		"PUT /api/categories/{id}":       schemaRef("CategoryRequest"),
		"POST /api/providers":            schemaRef("CreateProviderRequest"),
		"POST /api/providers/{id}/zones": schemaRef("CreateZoneRequest"),
		"POST /api/queues/book":          schemaRef("BookQueueRequest"),
		"POST /api/notifications/send":   schemaRef("SendNotificationRequest"),
	}

	schema, ok := schemas[method+" "+path]
	if !ok {
		return openAPISchema{}, false
	}
	return schema, true
}

func routeHasNoBody(method, path string) bool {
	noBodyRoutes := map[string]struct{}{
		"PATCH /api/zones/{id}/toggle":           {},
		"PATCH /api/queues/{id}/cancel":          {},
		"PATCH /api/manage/queues/{id}/call":     {},
		"PATCH /api/manage/queues/{id}/complete": {},
		"PATCH /api/manage/queues/{id}/skip":     {},
		"PATCH /api/notifications/{id}/read":     {},
		"DELETE /api/categories/{id}":            {},
		"DELETE /api/notifications/{id}":         {},
	}
	_, ok := noBodyRoutes[method+" "+path]
	return ok
}

func responsesForRoute(method, path string) map[string]openAPIResponse {
	successStatus := "200"
	if method == http.MethodPost {
		successStatus = "201"
	}
	if method == http.MethodDelete {
		successStatus = "204"
	}

	responses := map[string]openAPIResponse{
		successStatus: jsonResponse("Successful response", responseSchemaForRoute(method, path)),
		"400":         jsonResponse("Bad request", schemaRef("ErrorResponse")),
		"500":         jsonResponse("Internal server error", schemaRef("ErrorResponse")),
	}
	if method == http.MethodDelete {
		responses[successStatus] = openAPIResponse{Description: "Deleted successfully"}
	}
	return responses
}

func responseSchemaForRoute(method, path string) openAPISchema {
	schemas := map[string]openAPISchema{
		"GET /api/categories":                arraySchema(schemaRef("Category")),
		"GET /api/categories/{id}":           schemaRef("Category"),
		"POST /api/categories":               schemaRef("Category"),
		"PUT /api/categories/{id}":           schemaRef("Category"),
		"GET /api/providers":                 arraySchema(schemaRef("Provider")),
		"POST /api/providers":                schemaRef("Provider"),
		"GET /api/providers/{id}/zones":      arraySchema(schemaRef("Zone")),
		"POST /api/providers/{id}/zones":     schemaRef("Zone"),
		"PATCH /api/zones/{id}/toggle":       schemaRef("Zone"),
		"POST /api/queues/book":              schemaRef("Queue"),
		"GET /api/queues/{queueNumber}":      schemaRef("Queue"),
		"GET /api/queues/history":            arraySchema(schemaRef("Queue")),
		"GET /api/notifications":             arraySchema(schemaRef("Notification")),
		"POST /api/notifications/send":       schemaRef("Notification"),
		"POST /api/auth/request-otp":         schemaRef("MessageResponse"),
		"POST /api/auth/verify-otp":          schemaRef("AuthTokenResponse"),
		"POST /api/auth/register":            schemaRef("User"),
		"GET /api/auth/me":                   schemaRef("User"),
		"PUT /api/auth/me":                   schemaRef("User"),
		"PATCH /api/queues/{id}/cancel":      schemaRef("MessageResponse"),
		"PATCH /api/notifications/{id}/read": schemaRef("MessageResponse"),
	}

	if schema, ok := schemas[method+" "+path]; ok {
		return schema
	}
	return freeFormObjectSchema()
}

func jsonRequestBody(schema openAPISchema) *openAPIRequestBody {
	return &openAPIRequestBody{
		Required: true,
		Content: map[string]openAPIMediaSchema{
			"application/json": {Schema: schema},
		},
	}
}

func jsonResponse(description string, schema openAPISchema) openAPIResponse {
	return openAPIResponse{
		Description: description,
		Content: map[string]openAPIMediaSchema{
			"application/json": {Schema: schema},
		},
	}
}

func schemaRef(name string) openAPISchema {
	return openAPISchema{Ref: "#/components/schemas/" + name}
}

func arraySchema(item openAPISchema) openAPISchema {
	return openAPISchema{Type: "array", Items: &item}
}

func freeFormObjectSchema() openAPISchema {
	return openAPISchema{
		Type:                 "object",
		AdditionalProperties: &openAPISchema{},
	}
}

func componentSchemas() map[string]openAPISchema {
	return map[string]openAPISchema{
		"RequestOTPRequest": objectSchema([]string{"phone"}, map[string]openAPISchema{
			"phone": {Type: "string"},
		}),
		"VerifyOTPRequest": objectSchema([]string{"phone", "code"}, map[string]openAPISchema{
			"phone": {Type: "string"},
			"code":  {Type: "string"},
		}),
		"RegisterRequest": objectSchema([]string{"phone", "name"}, map[string]openAPISchema{
			"phone": {Type: "string"},
			"name":  {Type: "string"},
			"role":  {Type: "string", Enum: []string{"user", "provider", "admin"}},
		}),
		"UpdateProfileRequest": objectSchema(nil, map[string]openAPISchema{
			"name": {Type: "string"},
		}),
		"CategoryRequest": objectSchema([]string{"name"}, map[string]openAPISchema{
			"name": {Type: "string"},
		}),
		"CreateProviderRequest": objectSchema([]string{"name"}, map[string]openAPISchema{
			"name":        {Type: "string"},
			"category_id": {Type: "integer", Format: "uint"},
		}),
		"CreateZoneRequest": objectSchema([]string{"name"}, map[string]openAPISchema{
			"name": {Type: "string"},
		}),
		"BookQueueRequest": objectSchema([]string{"zone_id"}, map[string]openAPISchema{
			"zone_id": {Type: "integer", Format: "uint"},
			"user_id": {Type: "integer", Format: "uint"},
		}),
		"SendNotificationRequest": objectSchema([]string{"user_id", "message"}, map[string]openAPISchema{
			"user_id": {Type: "integer", Format: "uint"},
			"message": {Type: "string"},
		}),
		"Category": objectSchema(nil, timestampedProperties(map[string]openAPISchema{
			"id":   {Type: "integer", Format: "uint"},
			"name": {Type: "string"},
		})),
		"Provider": objectSchema(nil, timestampedProperties(map[string]openAPISchema{
			"id":          {Type: "integer", Format: "uint"},
			"name":        {Type: "string"},
			"category_id": {Type: "integer", Format: "uint"},
			"category":    schemaRef("Category"),
			"zones":       arraySchema(schemaRef("Zone")),
		})),
		"Zone": objectSchema(nil, timestampedProperties(map[string]openAPISchema{
			"id":          {Type: "integer", Format: "uint"},
			"provider_id": {Type: "integer", Format: "uint"},
			"name":        {Type: "string"},
			"is_open":     {Type: "boolean"},
			"queue_count": {Type: "integer", Format: "int32"},
		})),
		"Queue": objectSchema(nil, timestampedProperties(map[string]openAPISchema{
			"id":           {Type: "integer", Format: "uint"},
			"queue_number": {Type: "integer", Format: "int32"},
			"zone_id":      {Type: "integer", Format: "uint"},
			"user_id":      {Type: "integer", Format: "uint"},
			"status":       {Type: "string", Enum: []string{"waiting", "called", "completed", "skipped", "cancelled"}},
		})),
		"Notification": objectSchema(nil, timestampedProperties(map[string]openAPISchema{
			"id":      {Type: "integer", Format: "uint"},
			"user_id": {Type: "integer", Format: "uint"},
			"message": {Type: "string"},
			"is_read": {Type: "boolean"},
		})),
		"User": objectSchema(nil, timestampedProperties(map[string]openAPISchema{
			"id":    {Type: "integer", Format: "uint"},
			"phone": {Type: "string"},
			"name":  {Type: "string"},
			"role":  {Type: "string", Enum: []string{"user", "provider", "admin"}},
		})),
		"AuthTokenResponse": objectSchema(nil, map[string]openAPISchema{
			"token": {Type: "string"},
			"user":  schemaRef("User"),
		}),
		"MessageResponse": objectSchema(nil, map[string]openAPISchema{
			"message": {Type: "string"},
		}),
		"ErrorResponse": objectSchema(nil, map[string]openAPISchema{
			"status":  {Type: "integer", Format: "int32"},
			"error":   {Type: "string"},
			"message": {Type: "string"},
			"path":    {Type: "string"},
		}),
	}
}

func objectSchema(required []string, properties map[string]openAPISchema) openAPISchema {
	return openAPISchema{
		Type:       "object",
		Required:   required,
		Properties: properties,
	}
}

func timestampedProperties(properties map[string]openAPISchema) map[string]openAPISchema {
	properties["created_at"] = openAPISchema{Type: "string", Format: "date-time"}
	properties["updated_at"] = openAPISchema{Type: "string", Format: "date-time"}
	return properties
}
