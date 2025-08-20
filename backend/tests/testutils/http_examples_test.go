package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ExampleAPIHandler demonstrates a typical API handler for testing
func ExampleAPIHandler(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid JSON"})
		return
	}
	
	// Simulate some business logic
	if name, exists := body["name"]; !exists || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Name is required"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "User created successfully",
		"data":   gin.H{"id": "123", "name": body["name"]},
	})
}

// ExampleAuthenticatedHandler demonstrates an authenticated API handler
func ExampleAuthenticatedHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Authorization header required"})
		return
	}
	
	// Simulate token validation
	if authHeader != "Bearer valid-token" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Invalid token"})
		return
	}
	
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "User retrieved successfully",
		"data":   gin.H{"id": userID, "name": "John Doe"},
	})
}

// TestExampleAPIHandlerUsage demonstrates how to use HTTP utilities to test API handlers
func TestExampleAPIHandlerUsage(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		// Arrange
		requestBody := map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		}
		
		ctx, recorder := NewTestGinContext("POST", "/api/users", requestBody)
		
		// Act
		ExampleAPIHandler(ctx)
		
		// Assert
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertJSONResponse(t, recorder)
		
		expectedData := map[string]interface{}{
			"id":   "123",
			"name": "John Doe",
		}
		AssertSuccessResponse(t, recorder, expectedData)
		AssertContentType(t, recorder, "application/json")
	})
	
	t.Run("missing name field", func(t *testing.T) {
		// Arrange
		requestBody := map[string]interface{}{
			"email": "john@example.com",
		}
		
		ctx, recorder := NewTestGinContext("POST", "/api/users", requestBody)
		
		// Act
		ExampleAPIHandler(ctx)
		
		// Assert
		AssertStatusCode(t, recorder, http.StatusBadRequest)
		AssertErrorResponse(t, recorder, http.StatusBadRequest, "name is required")
	})
	
	t.Run("invalid JSON", func(t *testing.T) {
		// Arrange - Create context with invalid JSON
		ctx, recorder := NewTestGinContext("POST", "/api/users", nil)
		// Simulate invalid JSON by setting a malformed body
		ctx.Request.Header.Set("Content-Type", "application/json")
		
		// Act
		ExampleAPIHandler(ctx)
		
		// Assert
		AssertStatusCode(t, recorder, http.StatusBadRequest)
		AssertErrorResponse(t, recorder, http.StatusBadRequest, "invalid json")
	})
}

// TestExampleAuthenticatedHandlerUsage demonstrates testing authenticated endpoints
func TestExampleAuthenticatedHandlerUsage(t *testing.T) {
	t.Run("successful authenticated request", func(t *testing.T) {
		// Arrange
		ctx, recorder := NewTestGinContext("GET", "/api/users/123", nil)
		SetAuthorizationHeader(ctx, "valid-token")
		SetURLParam(ctx, "id", "123")
		
		// Act
		ExampleAuthenticatedHandler(ctx)
		
		// Assert
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertSuccessResponse(t, recorder, map[string]interface{}{
			"id":   "123",
			"name": "John Doe",
		})
	})
	
	t.Run("missing authorization header", func(t *testing.T) {
		// Arrange
		ctx, recorder := NewTestGinContext("GET", "/api/users/123", nil)
		SetURLParam(ctx, "id", "123")
		
		// Act
		ExampleAuthenticatedHandler(ctx)
		
		// Assert
		AssertStatusCode(t, recorder, http.StatusUnauthorized)
		AssertErrorResponse(t, recorder, http.StatusUnauthorized, "authorization header required")
	})
	
	t.Run("invalid token", func(t *testing.T) {
		// Arrange
		ctx, recorder := NewTestGinContext("GET", "/api/users/123", nil)
		SetAuthorizationHeader(ctx, "invalid-token")
		SetURLParam(ctx, "id", "123")
		
		// Act
		ExampleAuthenticatedHandler(ctx)
		
		// Assert
		AssertStatusCode(t, recorder, http.StatusUnauthorized)
		AssertErrorResponse(t, recorder, http.StatusUnauthorized, "invalid token")
	})
}

// TestRouterIntegrationExample demonstrates testing with a full router setup
func TestRouterIntegrationExample(t *testing.T) {
	// Setup router
	router := CreateTestRouter()
	router.POST("/api/users", ExampleAPIHandler)
	router.GET("/api/users/:id", ExampleAuthenticatedHandler)
	
	t.Run("POST /api/users - success", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"name":  "Jane Doe",
			"email": "jane@example.com",
		}
		
		recorder := PerformRequest(router, "POST", "/api/users", requestBody)
		
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertSuccessResponse(t, recorder, map[string]interface{}{
			"id":   "123",
			"name": "Jane Doe",
		})
	})
	
	t.Run("GET /api/users/:id - authenticated", func(t *testing.T) {
		req := CreateTestRequest("GET", "/api/users/456", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		
		recorder := PerformRequest(router, "GET", "/api/users/456", nil)
		recorder = performRequestWithHeaders(router, req)
		
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseContains(t, recorder, "John Doe")
	})
}

// Helper function for performing requests with custom headers
func performRequestWithHeaders(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

// TestAdvancedHTTPUtilities demonstrates advanced usage patterns
func TestAdvancedHTTPUtilities(t *testing.T) {
	t.Run("testing query parameters", func(t *testing.T) {
		ctx, recorder := NewTestGinContext("GET", "/api/users", nil)
		SetQueryParam(ctx, "page", "2")
		SetQueryParam(ctx, "limit", "10")
		SetQueryParam(ctx, "sort", "name")
		
		// Handler that uses query parameters
		handler := func(c *gin.Context) {
			page := c.Query("page")
			limit := c.Query("limit")
			sort := c.Query("sort")
			
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"msg":    "Users retrieved",
				"data": gin.H{
					"page":  page,
					"limit": limit,
					"sort":  sort,
					"users": []string{"user1", "user2"},
				},
			})
		}
		
		handler(ctx)
		
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseContains(t, recorder, "page")
		AssertResponseContains(t, recorder, "limit")
		AssertResponseContains(t, recorder, "sort")
	})
	
	t.Run("testing custom headers", func(t *testing.T) {
		ctx, recorder := NewTestGinContext("POST", "/api/data", map[string]string{"key": "value"})
		SetRequestHeader(ctx, "X-Client-Version", "1.0.0")
		SetRequestHeader(ctx, "X-Request-ID", "req-123")
		
		// Handler that uses custom headers
		handler := func(c *gin.Context) {
			clientVersion := c.GetHeader("X-Client-Version")
			requestID := c.GetHeader("X-Request-ID")
			
			c.Header("X-Response-ID", requestID)
			c.JSON(http.StatusOK, gin.H{
				"status":         "success",
				"msg":           "Data processed",
				"client_version": clientVersion,
				"request_id":     requestID,
			})
		}
		
		handler(ctx)
		
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertResponseHeaders(t, recorder, map[string]string{
			"X-Response-ID": "req-123",
		})
		AssertResponseContains(t, recorder, "1.0.0")
	})
	
	t.Run("testing response parsing", func(t *testing.T) {
		ctx, recorder := NewTestGinContext("GET", "/api/status", nil)
		
		// Handler that returns structured data
		handler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"msg":    "System status",
				"data": gin.H{
					"uptime":    "24h",
					"version":   "1.2.3",
					"healthy":   true,
					"services": []string{"api", "database", "cache"},
				},
			})
		}
		
		handler(ctx)
		
		// Parse and verify response structure
		responseJSON, err := GetResponseJSON(recorder)
		assert.NoError(t, err)
		assert.Equal(t, "success", responseJSON["status"])
		
		data, ok := responseJSON["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "1.2.3", data["version"])
		assert.Equal(t, true, data["healthy"])
		
		services, ok := data["services"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, services, 3)
	})
}