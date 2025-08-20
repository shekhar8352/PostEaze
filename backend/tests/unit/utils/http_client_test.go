package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpclient "github.com/shekhar8352/PostEaze/utils/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// HTTPClientTestSuite defines the test suite for HTTP client utilities
type HTTPClientTestSuite struct {
	suite.Suite
	testServer *httptest.Server
}

// SetupSuite runs before all tests in the suite
func (suite *HTTPClientTestSuite) SetupSuite() {
	// Create a test HTTP server
	mux := http.NewServeMux()
	
	// Success endpoint
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "success",
			"method":  r.Method,
			"headers": r.Header,
		}
		
		// Read body if present
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			if len(body) > 0 {
				response["body"] = string(body)
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	
	// Error endpoint
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "bad request",
		})
	})
	
	// Slow endpoint for timeout testing
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "slow response",
		})
	})
	
	// Echo endpoint that returns request details
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		response := map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"headers":     r.Header,
			"body":        string(body),
			"query":       r.URL.Query(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	
	suite.testServer = httptest.NewServer(mux)
}

// TearDownSuite runs after all tests in the suite
func (suite *HTTPClientTestSuite) TearDownSuite() {
	if suite.testServer != nil {
		suite.testServer.Close()
	}
}

// TestNewRequestConfig tests the request configuration creation
func (suite *HTTPClientTestSuite) TestNewRequestConfig_ValidConfig() {
	configMap := map[string]interface{}{
		"method":           "POST",
		"url":              suite.testServer.URL + "/success",
		"timeoutinmillis":  5000,
		"retrycount":       3,
		"headers": map[string]interface{}{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
		},
	}
	
	config := httpclient.NewRequestConfig("test-config", configMap)
	
	assert.NotNil(suite.T(), config)
	// Note: The fields are private, so we can't directly test them
	// We'll test the functionality through the HTTP client usage
}

func (suite *HTTPClientTestSuite) TestNewRequestConfig_EmptyConfig() {
	config := httpclient.NewRequestConfig("empty-config", nil)
	
	assert.NotNil(suite.T(), config)
}

func (suite *HTTPClientTestSuite) TestNewRequestConfig_PartialConfig() {
	configMap := map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/success",
	}
	
	config := httpclient.NewRequestConfig("partial-config", configMap)
	
	assert.NotNil(suite.T(), config)
}

func (suite *HTTPClientTestSuite) TestNewRequestConfig_InvalidTypes() {
	configMap := map[string]interface{}{
		"method":          "GET",
		"url":             suite.testServer.URL + "/success",
		"timeoutinmillis": "invalid", // Should be int
		"retrycount":      "invalid", // Should be int
	}
	
	config := httpclient.NewRequestConfig("invalid-config", configMap)
	
	assert.NotNil(suite.T(), config)
	// The config should still be created, but with default values for invalid fields
}

// TestInitHttp tests HTTP client initialization
func (suite *HTTPClientTestSuite) TestInitHttp_SingleConfig() {
	config := httpclient.NewRequestConfig("test-single", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/success",
	})
	
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPClientTestSuite) TestInitHttp_MultipleConfigs() {
	config1 := httpclient.NewRequestConfig("test-config-1", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/success",
	})
	
	config2 := httpclient.NewRequestConfig("test-config-2", map[string]interface{}{
		"method": "POST",
		"url":    suite.testServer.URL + "/echo",
	})
	
	httpclient.InitHttp(config1, config2)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPClientTestSuite) TestInitHttp_NilConfig() {
	httpclient.InitHttp(nil)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPClientTestSuite) TestInitHttp_EmptyName() {
	config := httpclient.NewRequestConfig("", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/success",
	})
	
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

// TestGetClient tests the client getter
func (suite *HTTPClientTestSuite) TestGetClient() {
	client := httpclient.GetClient()
	
	assert.NotNil(suite.T(), client)
	
	// Test that multiple calls return the same instance
	client2 := httpclient.GetClient()
	assert.Equal(suite.T(), client, client2)
}

// TestRequestBuilder tests the request builder functionality
func (suite *HTTPClientTestSuite) TestRequestBuilder_SetContext() {
	// Initialize client first
	config := httpclient.NewRequestConfig("context-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/success",
	})
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	ctx := context.WithValue(context.Background(), "test", "value")
	
	// Note: The request builder methods are not directly testable due to private fields
	// We would need to test through the actual HTTP execution
	assert.NotNil(suite.T(), client)
	assert.NotNil(suite.T(), ctx)
}

func (suite *HTTPClientTestSuite) TestRequestBuilder_SetMethod() {
	config := httpclient.NewRequestConfig("method-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
	
	// The actual method setting would be tested through HTTP execution
}

func (suite *HTTPClientTestSuite) TestRequestBuilder_SetURL() {
	config := httpclient.NewRequestConfig("url-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/success",
	})
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPClientTestSuite) TestRequestBuilder_SetQueryParams() {
	config := httpclient.NewRequestConfig("query-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
	
	// Query params would be tested through actual HTTP execution
}

func (suite *HTTPClientTestSuite) TestRequestBuilder_SetHeaders() {
	config := httpclient.NewRequestConfig("header-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
		"headers": map[string]interface{}{
			"X-Test-Header": "test-value",
		},
	})
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPClientTestSuite) TestRequestBuilder_SetBody() {
	config := httpclient.NewRequestConfig("body-test", map[string]interface{}{
		"method": "POST",
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
	
	// Body setting would be tested through actual HTTP execution
}

// TestHTTPExecution tests the HTTP execution functions
func (suite *HTTPClientTestSuite) TestCall_Success() {
	// Create a simple HTTP request
	req, err := http.NewRequest("GET", suite.testServer.URL+"/success", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	err = httpclient.Call(req, client)
	
	assert.NoError(suite.T(), err)
}

func (suite *HTTPClientTestSuite) TestCall_Error() {
	// Create a request to error endpoint
	req, err := http.NewRequest("GET", suite.testServer.URL+"/error", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	err = httpclient.Call(req, client)
	
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "non-success HTTP status")
}

func (suite *HTTPClientTestSuite) TestCall_Timeout() {
	// Create a request to slow endpoint with short timeout
	req, err := http.NewRequest("GET", suite.testServer.URL+"/slow", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 100 * time.Millisecond}
	
	err = httpclient.Call(req, client)
	
	assert.Error(suite.T(), err)
	// Error should be related to timeout or context cancellation
}

func (suite *HTTPClientTestSuite) TestCallAndGetResponse_Success() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/success", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	response, err := httpclient.CallAndGetResponse(req, client)
	
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response)
	
	// Parse response to verify it's valid JSON
	var responseData map[string]interface{}
	err = json.Unmarshal(response, &responseData)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", responseData["message"])
}

func (suite *HTTPClientTestSuite) TestCallAndGetResponse_Error() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/error", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	response, err := httpclient.CallAndGetResponse(req, client)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
}

func (suite *HTTPClientTestSuite) TestCallAndBind_Success() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/success", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["message"])
	assert.Equal(suite.T(), "GET", response["method"])
}

func (suite *HTTPClientTestSuite) TestCallAndBind_InvalidJSON() {
	// Create a server that returns invalid JSON
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json response"))
	}))
	defer invalidServer.Close()
	
	req, err := http.NewRequest("GET", invalidServer.URL, nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to unmarshal response")
}

func (suite *HTTPClientTestSuite) TestCallAndBind_WrongType() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/success", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	// Try to bind to wrong type
	var response string
	err = httpclient.CallAndBind(req, &response, client)
	
	assert.Error(suite.T(), err)
}

// TestHTTPRequestWithBody tests HTTP requests with body content
func (suite *HTTPClientTestSuite) TestHTTPRequestWithBody() {
	requestBody := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}
	
	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(suite.T(), err)
	
	req, err := http.NewRequest("POST", suite.testServer.URL+"/echo", bytes.NewReader(bodyBytes))
	require.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	
	client := http.Client{Timeout: 5 * time.Second}
	
	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "POST", response["method"])
	assert.Contains(suite.T(), response["body"], "test")
	assert.Contains(suite.T(), response["body"], "123")
}

// TestHTTPRequestWithHeaders tests HTTP requests with custom headers
func (suite *HTTPClientTestSuite) TestHTTPRequestWithHeaders() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/echo", nil)
	require.NoError(suite.T(), err)
	
	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("Authorization", "Bearer token123")
	
	client := http.Client{Timeout: 5 * time.Second}
	
	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	
	assert.NoError(suite.T(), err)
	
	headers := response["headers"].(map[string]interface{})
	assert.Contains(suite.T(), headers, "X-Custom-Header")
	assert.Contains(suite.T(), headers, "Authorization")
}

// TestHTTPRequestWithQueryParams tests HTTP requests with query parameters
func (suite *HTTPClientTestSuite) TestHTTPRequestWithQueryParams() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/echo?param1=value1&param2=value2", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	
	assert.NoError(suite.T(), err)
	
	query := response["query"].(map[string]interface{})
	assert.Contains(suite.T(), query, "param1")
	assert.Contains(suite.T(), query, "param2")
}

// TestHTTPErrorHandling tests various error scenarios
func (suite *HTTPClientTestSuite) TestHTTPErrorHandling_InvalidURL() {
	req, err := http.NewRequest("GET", "invalid-url", nil)
	
	if err == nil {
		client := http.Client{Timeout: 5 * time.Second}
		err = httpclient.Call(req, client)
		assert.Error(suite.T(), err)
	} else {
		assert.Error(suite.T(), err)
	}
}

func (suite *HTTPClientTestSuite) TestHTTPErrorHandling_ConnectionRefused() {
	// Use a port that's likely not in use
	req, err := http.NewRequest("GET", "http://localhost:99999/test", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 1 * time.Second}
	
	err = httpclient.Call(req, client)
	
	assert.Error(suite.T(), err)
}

// Run the test suite
func TestHTTPClientTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientTestSuite))
}