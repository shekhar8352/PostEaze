package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

// NewTestGinContext creates a new Gin context for testing with the specified method, URL, and body
func NewTestGinContext(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	
	var req *http.Request
	var err error
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal request body: %v", err))
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			panic(fmt.Sprintf("Failed to create request: %v", err))
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			panic(fmt.Sprintf("Failed to create request: %v", err))
		}
	}
	
	ctx.Request = req
	return ctx, recorder
}

// CreateTestRequest creates an HTTP request for testing
func CreateTestRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	var err error
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal request body: %v", err))
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			panic(fmt.Sprintf("Failed to create request: %v", err))
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			panic(fmt.Sprintf("Failed to create request: %v", err))
		}
	}
	
	return req
}

// ParseJSONResponse parses the JSON response from a ResponseRecorder into the target interface
func ParseJSONResponse(recorder *httptest.ResponseRecorder, target interface{}) error {
	if recorder.Body.Len() == 0 {
		return fmt.Errorf("response body is empty")
	}
	
	return json.Unmarshal(recorder.Body.Bytes(), target)
}

// AssertSuccessResponse asserts that the response is a successful API response with expected data
func AssertSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedData interface{}) {
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	
	var response APIResponse
	err := ParseJSONResponse(recorder, &response)
	require.NoError(t, err, "Failed to parse JSON response")
	
	assert.Equal(t, "success", response.Status, "Expected success status")
	assert.NotEmpty(t, response.Msg, "Expected non-empty message")
	
	if expectedData != nil {
		// Convert both to JSON strings for comparison to handle different types
		expectedJSON, err := json.Marshal(expectedData)
		require.NoError(t, err, "Failed to marshal expected data")
		
		actualJSON, err := json.Marshal(response.Data)
		require.NoError(t, err, "Failed to marshal actual data")
		
		assert.JSONEq(t, string(expectedJSON), string(actualJSON), "Response data does not match expected data")
	}
}

// AssertErrorResponse asserts that the response is an error response with expected status and message
func AssertErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	assert.Equal(t, expectedStatus, recorder.Code, "Expected status code %d", expectedStatus)
	
	var response APIResponse
	err := ParseJSONResponse(recorder, &response)
	require.NoError(t, err, "Failed to parse JSON response")
	
	assert.Equal(t, "error", response.Status, "Expected error status")
	
	if expectedMessage != "" {
		assert.Contains(t, strings.ToLower(response.Msg), strings.ToLower(expectedMessage), 
			"Expected error message to contain: %s", expectedMessage)
	}
}

// AssertJSONResponse asserts that the response contains valid JSON
func AssertJSONResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	assert.True(t, json.Valid(recorder.Body.Bytes()), "Response should contain valid JSON")
}

// AssertResponseHeaders asserts that the response contains expected headers
func AssertResponseHeaders(t *testing.T, recorder *httptest.ResponseRecorder, expectedHeaders map[string]string) {
	for key, expectedValue := range expectedHeaders {
		actualValue := recorder.Header().Get(key)
		assert.Equal(t, expectedValue, actualValue, "Expected header %s to be %s", key, expectedValue)
	}
}

// AssertContentType asserts that the response has the expected content type
func AssertContentType(t *testing.T, recorder *httptest.ResponseRecorder, expectedContentType string) {
	contentType := recorder.Header().Get("Content-Type")
	assert.Contains(t, contentType, expectedContentType, "Expected content type to contain %s", expectedContentType)
}

// GetResponseBody returns the response body as a string
func GetResponseBody(recorder *httptest.ResponseRecorder) string {
	return recorder.Body.String()
}

// GetResponseJSON parses the response body as JSON and returns it as a map
func GetResponseJSON(recorder *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := ParseJSONResponse(recorder, &result)
	return result, err
}

// SetRequestHeader sets a header on the request in the Gin context
func SetRequestHeader(ctx *gin.Context, key, value string) {
	ctx.Request.Header.Set(key, value)
}

// SetAuthorizationHeader sets the Authorization header with Bearer token
func SetAuthorizationHeader(ctx *gin.Context, token string) {
	ctx.Request.Header.Set("Authorization", "Bearer "+token)
}

// SetURLParam sets a URL parameter in the Gin context
func SetURLParam(ctx *gin.Context, key, value string) {
	ctx.Params = append(ctx.Params, gin.Param{Key: key, Value: value})
}

// SetQueryParam sets a query parameter in the Gin context
func SetQueryParam(ctx *gin.Context, key, value string) {
	q := ctx.Request.URL.Query()
	q.Set(key, value)
	ctx.Request.URL.RawQuery = q.Encode()
}

// AssertResponseContains asserts that the response body contains the expected string
func AssertResponseContains(t *testing.T, recorder *httptest.ResponseRecorder, expected string) {
	body := recorder.Body.String()
	assert.Contains(t, body, expected, "Expected response body to contain: %s", expected)
}

// AssertResponseNotContains asserts that the response body does not contain the specified string
func AssertResponseNotContains(t *testing.T, recorder *httptest.ResponseRecorder, notExpected string) {
	body := recorder.Body.String()
	assert.NotContains(t, body, notExpected, "Expected response body to not contain: %s", notExpected)
}

// CreateTestRouter creates a new Gin router for testing
func CreateTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// PerformRequest performs an HTTP request on the given router and returns the response recorder
func PerformRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	req := CreateTestRequest(method, path, body)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

// AssertStatusCode asserts that the response has the expected status code
func AssertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, recorder.Code, "Expected status code %d, got %d", expectedStatus, recorder.Code)
}

// AssertEmptyResponse asserts that the response body is empty or contains only whitespace
func AssertEmptyResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	body := strings.TrimSpace(recorder.Body.String())
	assert.Empty(t, body, "Expected empty response body")
}

// AssertNonEmptyResponse asserts that the response body is not empty
func AssertNonEmptyResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	body := strings.TrimSpace(recorder.Body.String())
	assert.NotEmpty(t, body, "Expected non-empty response body")
}