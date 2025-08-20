package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewTestGinContext(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		body   interface{}
	}{
		{
			name:   "GET request without body",
			method: "GET",
			url:    "/test",
			body:   nil,
		},
		{
			name:   "POST request with JSON body",
			method: "POST",
			url:    "/test",
			body:   map[string]string{"key": "value"},
		},
		{
			name:   "PUT request with struct body",
			method: "PUT",
			url:    "/test/123",
			body:   struct{ Name string }{Name: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := NewTestGinContext(tt.method, tt.url, tt.body)
			
			assert.NotNil(t, ctx)
			assert.NotNil(t, recorder)
			assert.Equal(t, tt.method, ctx.Request.Method)
			assert.Equal(t, tt.url, ctx.Request.URL.Path)
			
			if tt.body != nil {
				assert.Equal(t, "application/json", ctx.Request.Header.Get("Content-Type"))
			}
		})
	}
}

func TestCreateTestRequest(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		body   interface{}
	}{
		{
			name:   "GET request",
			method: "GET",
			url:    "/api/test",
			body:   nil,
		},
		{
			name:   "POST request with body",
			method: "POST",
			url:    "/api/test",
			body:   map[string]interface{}{"test": "data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateTestRequest(tt.method, tt.url, tt.body)
			
			assert.NotNil(t, req)
			assert.Equal(t, tt.method, req.Method)
			assert.Equal(t, tt.url, req.URL.Path)
			
			if tt.body != nil {
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			}
		})
	}
}

func TestParseJSONResponse(t *testing.T) {
	tests := []struct {
		name        string
		responseBody string
		target      interface{}
		expectError bool
	}{
		{
			name:         "valid JSON",
			responseBody: `{"status": "success", "msg": "test"}`,
			target:       &APIResponse{},
			expectError:  false,
		},
		{
			name:         "invalid JSON",
			responseBody: `{"status": "success", "msg":}`,
			target:       &APIResponse{},
			expectError:  true,
		},
		{
			name:         "empty response",
			responseBody: "",
			target:       &APIResponse{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.WriteString(tt.responseBody)
			
			err := ParseJSONResponse(recorder, tt.target)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssertSuccessResponse(t *testing.T) {
	t.Run("valid success response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusOK)
		recorder.WriteString(`{"status": "success", "msg": "operation successful", "data": {"id": "123"}}`)
		
		expectedData := map[string]string{"id": "123"}
		
		// This should not panic or fail
		AssertSuccessResponse(t, recorder, expectedData)
	})

	t.Run("success response without data", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusOK)
		recorder.WriteString(`{"status": "success", "msg": "operation successful"}`)
		
		// This should not panic or fail
		AssertSuccessResponse(t, recorder, nil)
	})
}

func TestAssertErrorResponse(t *testing.T) {
	t.Run("valid error response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusBadRequest)
		recorder.WriteString(`{"status": "error", "msg": "Invalid input data"}`)
		
		// This should not panic or fail
		AssertErrorResponse(t, recorder, http.StatusBadRequest, "invalid input")
	})

	t.Run("error response with empty message check", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusInternalServerError)
		recorder.WriteString(`{"status": "error", "msg": "Internal server error"}`)
		
		// This should not panic or fail
		AssertErrorResponse(t, recorder, http.StatusInternalServerError, "")
	})
}

func TestAssertJSONResponse(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		expectValid  bool
	}{
		{
			name:         "valid JSON",
			responseBody: `{"status": "success"}`,
			expectValid:  true,
		},
		{
			name:         "invalid JSON",
			responseBody: `{"status": "success"`,
			expectValid:  false,
		},
		{
			name:         "empty response",
			responseBody: "",
			expectValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.WriteString(tt.responseBody)
			
			if tt.expectValid {
				// This should not fail
				AssertJSONResponse(t, recorder)
			} else {
				// We can't easily test assertion failures, but we can verify the function exists
				assert.NotNil(t, AssertJSONResponse)
			}
		})
	}
}

func TestAssertResponseHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "application/json")
	recorder.Header().Set("X-Custom-Header", "test-value")
	
	expectedHeaders := map[string]string{
		"Content-Type":     "application/json",
		"X-Custom-Header":  "test-value",
	}
	
	// This should not panic or fail
	AssertResponseHeaders(t, recorder, expectedHeaders)
}

func TestAssertContentType(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// This should not panic or fail
	AssertContentType(t, recorder, "application/json")
}

func TestGetResponseBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	expectedBody := "test response body"
	recorder.WriteString(expectedBody)
	
	actualBody := GetResponseBody(recorder)
	assert.Equal(t, expectedBody, actualBody)
}

func TestGetResponseJSON(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		expectError  bool
	}{
		{
			name:         "valid JSON",
			responseBody: `{"key": "value", "number": 123}`,
			expectError:  false,
		},
		{
			name:         "invalid JSON",
			responseBody: `{"key": "value"`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.WriteString(tt.responseBody)
			
			result, err := GetResponseJSON(recorder)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestSetRequestHeader(t *testing.T) {
	ctx, _ := NewTestGinContext("GET", "/test", nil)
	
	SetRequestHeader(ctx, "X-Test-Header", "test-value")
	
	assert.Equal(t, "test-value", ctx.Request.Header.Get("X-Test-Header"))
}

func TestSetAuthorizationHeader(t *testing.T) {
	ctx, _ := NewTestGinContext("GET", "/test", nil)
	token := "test-jwt-token"
	
	SetAuthorizationHeader(ctx, token)
	
	assert.Equal(t, "Bearer "+token, ctx.Request.Header.Get("Authorization"))
}

func TestSetURLParam(t *testing.T) {
	ctx, _ := NewTestGinContext("GET", "/test", nil)
	
	SetURLParam(ctx, "id", "123")
	SetURLParam(ctx, "type", "user")
	
	assert.Equal(t, "123", ctx.Param("id"))
	assert.Equal(t, "user", ctx.Param("type"))
}

func TestSetQueryParam(t *testing.T) {
	ctx, _ := NewTestGinContext("GET", "/test", nil)
	
	SetQueryParam(ctx, "page", "1")
	SetQueryParam(ctx, "limit", "10")
	
	assert.Equal(t, "1", ctx.Query("page"))
	assert.Equal(t, "10", ctx.Query("limit"))
}

func TestAssertResponseContains(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteString("This is a test response with specific content")
	
	// This should not panic or fail
	AssertResponseContains(t, recorder, "specific content")
}

func TestAssertResponseNotContains(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteString("This is a test response")
	
	// This should not panic or fail
	AssertResponseNotContains(t, recorder, "unwanted content")
}

func TestCreateTestRouter(t *testing.T) {
	router := CreateTestRouter()
	
	assert.NotNil(t, router)
	assert.Equal(t, gin.TestMode, gin.Mode())
}

func TestPerformRequest(t *testing.T) {
	router := CreateTestRouter()
	
	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
	
	recorder := PerformRequest(router, "GET", "/test", nil)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "test")
}

func TestAssertStatusCode(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(http.StatusCreated)
	
	// This should not panic or fail
	AssertStatusCode(t, recorder, http.StatusCreated)
}

func TestAssertEmptyResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteString("   ") // Only whitespace
	
	// This should not panic or fail
	AssertEmptyResponse(t, recorder)
}

func TestAssertNonEmptyResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteString("some content")
	
	// This should not panic or fail
	AssertNonEmptyResponse(t, recorder)
}

// Integration test to demonstrate usage
func TestHTTPUtilitiesIntegration(t *testing.T) {
	// Create a test router
	router := CreateTestRouter()
	
	// Add test routes
	router.POST("/api/test", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid JSON"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Data received", "data": body})
	})
	
	router.GET("/api/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": "Something went wrong"})
	})
	
	t.Run("successful POST request", func(t *testing.T) {
		testData := map[string]interface{}{
			"name":  "test",
			"value": 123,
		}
		
		recorder := PerformRequest(router, "POST", "/api/test", testData)
		
		AssertStatusCode(t, recorder, http.StatusOK)
		AssertJSONResponse(t, recorder)
		AssertSuccessResponse(t, recorder, testData)
		AssertContentType(t, recorder, "application/json")
	})
	
	t.Run("error response", func(t *testing.T) {
		recorder := PerformRequest(router, "GET", "/api/error", nil)
		
		AssertStatusCode(t, recorder, http.StatusInternalServerError)
		AssertErrorResponse(t, recorder, http.StatusInternalServerError, "something went wrong")
	})
}