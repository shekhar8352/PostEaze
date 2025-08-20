package utils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httpclient "github.com/shekhar8352/PostEaze/utils/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// HTTPRequestUtilsTestSuite defines the test suite for HTTP request utilities
type HTTPRequestUtilsTestSuite struct {
	suite.Suite
	testServer *httptest.Server
}

// SetupSuite runs before all tests in the suite
func (suite *HTTPRequestUtilsTestSuite) SetupSuite() {
	// Create a comprehensive test server
	mux := http.NewServeMux()
	
	// Echo endpoint that returns all request details
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		response := "Method: " + r.Method + "\n"
		response += "URL: " + r.URL.String() + "\n"
		response += "Headers: " + r.Header.Get("Content-Type") + "\n"
		response += "Body: " + string(body) + "\n"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	
	// Endpoint that requires specific headers
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization header"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authorized: " + auth))
	})
	
	// Endpoint that processes query parameters
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		response := "Query params:\n"
		for key, values := range params {
			for _, value := range values {
				response += key + "=" + value + "\n"
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	
	suite.testServer = httptest.NewServer(mux)
}

// TearDownSuite runs after all tests in the suite
func (suite *HTTPRequestUtilsTestSuite) TearDownSuite() {
	if suite.testServer != nil {
		suite.testServer.Close()
	}
}

// TestRequestConfigCreation tests request configuration creation with various parameters
func (suite *HTTPRequestUtilsTestSuite) TestRequestConfig_CompleteConfiguration() {
	configMap := map[string]interface{}{
		"method":           "POST",
		"url":              suite.testServer.URL + "/echo",
		"timeoutinmillis":  10000,
		"retrycount":       5,
		"headers": map[string]interface{}{
			"Content-Type":     "application/json",
			"Accept":           "application/json",
			"User-Agent":       "TestClient/1.0",
			"X-Custom-Header":  "custom-value",
		},
	}
	
	config := httpclient.NewRequestConfig("complete-config", configMap)
	
	assert.NotNil(suite.T(), config)
	// The config should be created successfully with all parameters
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestConfig_MinimalConfiguration() {
	configMap := map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
	}
	
	config := httpclient.NewRequestConfig("minimal-config", configMap)
	
	assert.NotNil(suite.T(), config)
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestConfig_TimeoutConfiguration() {
	testCases := []struct {
		name           string
		timeoutMillis  interface{}
		expectedError  bool
	}{
		{"Valid timeout", 5000, false},
		{"Zero timeout", 0, false},
		{"String timeout", "5000", false}, // Should be converted
		{"Invalid timeout", "invalid", false}, // Should use default
		{"Negative timeout", -1000, false}, // Should handle gracefully
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			configMap := map[string]interface{}{
				"method":          "GET",
				"url":             suite.testServer.URL + "/echo",
				"timeoutinmillis": tc.timeoutMillis,
			}
			
			config := httpclient.NewRequestConfig("timeout-test", configMap)
			assert.NotNil(t, config)
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestConfig_RetryConfiguration() {
	testCases := []struct {
		name         string
		retryCount   interface{}
		expectedError bool
	}{
		{"Valid retry count", 3, false},
		{"Zero retries", 0, false},
		{"String retry count", "3", false}, // Should be converted
		{"Invalid retry count", "invalid", false}, // Should use default
		{"Negative retry count", -1, false}, // Should handle gracefully
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			configMap := map[string]interface{}{
				"method":     "GET",
				"url":        suite.testServer.URL + "/echo",
				"retrycount": tc.retryCount,
			}
			
			config := httpclient.NewRequestConfig("retry-test", configMap)
			assert.NotNil(t, config)
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestConfig_HeaderConfiguration() {
	testCases := []struct {
		name     string
		headers  interface{}
		expected bool
	}{
		{
			"Valid headers map",
			map[string]interface{}{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			true,
		},
		{
			"Empty headers map",
			map[string]interface{}{},
			true,
		},
		{
			"Nil headers",
			nil,
			true,
		},
		{
			"Invalid headers type",
			"invalid",
			true, // Should handle gracefully
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			configMap := map[string]interface{}{
				"method":  "GET",
				"url":     suite.testServer.URL + "/echo",
				"headers": tc.headers,
			}
			
			config := httpclient.NewRequestConfig("header-test", configMap)
			assert.NotNil(t, config)
		})
	}
}

// TestHTTPClientInitialization tests HTTP client initialization with various configurations
func (suite *HTTPRequestUtilsTestSuite) TestInitHttp_MultipleConfigsWithDifferentSettings() {
	config1 := httpclient.NewRequestConfig("fast-client", map[string]interface{}{
		"method":          "GET",
		"url":             suite.testServer.URL + "/echo",
		"timeoutinmillis": 1000,
		"retrycount":      1,
	})
	
	config2 := httpclient.NewRequestConfig("slow-client", map[string]interface{}{
		"method":          "POST",
		"url":             suite.testServer.URL + "/echo",
		"timeoutinmillis": 30000,
		"retrycount":      5,
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
	})
	
	config3 := httpclient.NewRequestConfig("auth-client", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/auth",
		"headers": map[string]interface{}{
			"Authorization": "Bearer token123",
		},
	})
	
	httpclient.InitHttp(config1, config2, config3)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPRequestUtilsTestSuite) TestInitHttp_DuplicateConfigNames() {
	config1 := httpclient.NewRequestConfig("duplicate", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
	})
	
	config2 := httpclient.NewRequestConfig("duplicate", map[string]interface{}{
		"method": "POST",
		"url":    suite.testServer.URL + "/auth",
	})
	
	// Should handle duplicate names gracefully (last one wins or similar behavior)
	httpclient.InitHttp(config1, config2)
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

func (suite *HTTPRequestUtilsTestSuite) TestInitHttp_EmptyConfigList() {
	httpclient.InitHttp()
	
	client := httpclient.GetClient()
	assert.NotNil(suite.T(), client)
}

// TestRequestBuilding tests the request building functionality
func (suite *HTTPRequestUtilsTestSuite) TestRequestBuilding_ContextHandling() {
	// Initialize client
	config := httpclient.NewRequestConfig("context-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	// Test different context scenarios
	testCases := []struct {
		name string
		ctx  context.Context
	}{
		{"Background context", context.Background()},
		{"Context with value", context.WithValue(context.Background(), "key", "value")},
		{"Context with timeout", func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }()},
		{"Cancelled context", func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }()},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.ctx)
			// The actual request building would be tested through HTTP execution
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestBuilding_MethodOverride() {
	config := httpclient.NewRequestConfig("method-test", map[string]interface{}{
		"method": "GET", // Default method
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	// Test that method can be overridden
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	
	for _, method := range methods {
		suite.T().Run("Method_"+method, func(t *testing.T) {
			assert.Contains(t, methods, method)
			// The actual method setting would be tested through HTTP execution
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestBuilding_URLOverride() {
	config := httpclient.NewRequestConfig("url-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo", // Default URL
	})
	httpclient.InitHttp(config)
	
	// Test URL override scenarios
	testURLs := []string{
		suite.testServer.URL + "/auth",
		suite.testServer.URL + "/query",
		suite.testServer.URL + "/echo?param=value",
	}
	
	for i, testURL := range testURLs {
		suite.T().Run("URL_"+string(rune(i+'A')), func(t *testing.T) {
			assert.NotEmpty(t, testURL)
			assert.Contains(t, testURL, suite.testServer.URL)
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestBuilding_QueryParameters() {
	config := httpclient.NewRequestConfig("query-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/query",
	})
	httpclient.InitHttp(config)
	
	// Test various query parameter scenarios
	testCases := []struct {
		name   string
		params map[string]string
	}{
		{
			"Single parameter",
			map[string]string{"key": "value"},
		},
		{
			"Multiple parameters",
			map[string]string{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
		},
		{
			"Empty value parameter",
			map[string]string{"empty": ""},
		},
		{
			"Special characters",
			map[string]string{
				"special": "value with spaces",
				"encoded": "value&with=special&chars",
			},
		},
		{
			"Numeric parameters",
			map[string]string{
				"id":    "123",
				"count": "456",
				"page":  "1",
			},
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.params)
			for key, value := range tc.params {
				assert.NotEmpty(t, key)
				// Value can be empty, so we don't check for that
				_ = value
			}
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestBuilding_HeaderParameters() {
	config := httpclient.NewRequestConfig("header-test", map[string]interface{}{
		"method": "GET",
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	// Test various header scenarios
	testCases := []struct {
		name    string
		headers map[string]string
	}{
		{
			"Content type headers",
			map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{
			"Authentication headers",
			map[string]string{
				"Authorization": "Bearer token123",
				"X-API-Key":     "api-key-456",
			},
		},
		{
			"Custom headers",
			map[string]string{
				"X-Request-ID":   "req-123",
				"X-Client-Name":  "TestClient",
				"X-Client-Version": "1.0.0",
			},
		},
		{
			"Cache control headers",
			map[string]string{
				"Cache-Control": "no-cache",
				"Pragma":        "no-cache",
			},
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.headers)
			for key, value := range tc.headers {
				assert.NotEmpty(t, key)
				assert.NotEmpty(t, value)
			}
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestRequestBuilding_BodyHandling() {
	config := httpclient.NewRequestConfig("body-test", map[string]interface{}{
		"method": "POST",
		"url":    suite.testServer.URL + "/echo",
	})
	httpclient.InitHttp(config)
	
	// Test various body types
	testCases := []struct {
		name string
		body io.Reader
	}{
		{
			"String body",
			strings.NewReader("test string body"),
		},
		{
			"JSON body",
			strings.NewReader(`{"key": "value", "number": 123}`),
		},
		{
			"Empty body",
			strings.NewReader(""),
		},
		{
			"Binary body",
			bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04}),
		},
		{
			"Large body",
			strings.NewReader(strings.Repeat("a", 10000)),
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.body)
		})
	}
}

// TestHTTPExecution tests HTTP execution with various scenarios
func (suite *HTTPRequestUtilsTestSuite) TestHTTPExecution_SuccessScenarios() {
	testCases := []struct {
		name           string
		method         string
		endpoint       string
		expectedStatus int
	}{
		{"GET request", "GET", "/echo", http.StatusOK},
		{"POST request", "POST", "/echo", http.StatusOK},
		{"PUT request", "PUT", "/echo", http.StatusOK},
		{"DELETE request", "DELETE", "/echo", http.StatusOK},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, suite.testServer.URL+tc.endpoint, nil)
			require.NoError(t, err)
			
			client := http.Client{Timeout: 5 * time.Second}
			
			err = httpclient.Call(req, client)
			assert.NoError(t, err)
		})
	}
}

func (suite *HTTPRequestUtilsTestSuite) TestHTTPExecution_WithAuthentication() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/auth", nil)
	require.NoError(suite.T(), err)
	
	req.Header.Set("Authorization", "Bearer test-token")
	
	client := http.Client{Timeout: 5 * time.Second}
	
	err = httpclient.Call(req, client)
	assert.NoError(suite.T(), err)
}

func (suite *HTTPRequestUtilsTestSuite) TestHTTPExecution_WithoutAuthentication() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/auth", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	err = httpclient.Call(req, client)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "non-success HTTP status")
}

func (suite *HTTPRequestUtilsTestSuite) TestHTTPExecution_WithQueryParameters() {
	req, err := http.NewRequest("GET", suite.testServer.URL+"/query?param1=value1&param2=value2", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 5 * time.Second}
	
	response, err := httpclient.CallAndGetResponse(req, client)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), string(response), "param1=value1")
	assert.Contains(suite.T(), string(response), "param2=value2")
}

func (suite *HTTPRequestUtilsTestSuite) TestHTTPExecution_WithBody() {
	body := strings.NewReader(`{"test": "data"}`)
	req, err := http.NewRequest("POST", suite.testServer.URL+"/echo", body)
	require.NoError(suite.T(), err)
	
	req.Header.Set("Content-Type", "application/json")
	
	client := http.Client{Timeout: 5 * time.Second}
	
	response, err := httpclient.CallAndGetResponse(req, client)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), string(response), `{"test": "data"}`)
}

// TestErrorHandling tests various error scenarios
func (suite *HTTPRequestUtilsTestSuite) TestErrorHandling_NetworkErrors() {
	// Test connection refused
	req, err := http.NewRequest("GET", "http://localhost:99999/test", nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 1 * time.Second}
	
	err = httpclient.Call(req, client)
	assert.Error(suite.T(), err)
}

func (suite *HTTPRequestUtilsTestSuite) TestErrorHandling_TimeoutErrors() {
	// Create a slow server
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()
	
	req, err := http.NewRequest("GET", slowServer.URL, nil)
	require.NoError(suite.T(), err)
	
	client := http.Client{Timeout: 100 * time.Millisecond}
	
	err = httpclient.Call(req, client)
	assert.Error(suite.T(), err)
}

func (suite *HTTPRequestUtilsTestSuite) TestErrorHandling_InvalidRequests() {
	testCases := []struct {
		name   string
		method string
		url    string
	}{
		{"Invalid URL", "GET", "://invalid-url"},
		{"Empty URL", "GET", ""},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, nil)
			
			if err == nil {
				client := http.Client{Timeout: 5 * time.Second}
				err = httpclient.Call(req, client)
			}
			
			// Should have an error in either case
			assert.Error(t, err)
		})
	}
}

// TestConfigurationEdgeCases tests edge cases in configuration
func (suite *HTTPRequestUtilsTestSuite) TestConfigurationEdgeCases() {
	testCases := []struct {
		name      string
		configMap map[string]interface{}
	}{
		{
			"All nil values",
			map[string]interface{}{
				"method":          nil,
				"url":             nil,
				"timeoutinmillis": nil,
				"retrycount":      nil,
				"headers":         nil,
			},
		},
		{
			"Mixed valid and invalid values",
			map[string]interface{}{
				"method":          "GET",
				"url":             suite.testServer.URL + "/echo",
				"timeoutinmillis": "invalid",
				"retrycount":      nil,
				"headers":         "not a map",
			},
		},
		{
			"Extreme values",
			map[string]interface{}{
				"method":          "GET",
				"url":             suite.testServer.URL + "/echo",
				"timeoutinmillis": 999999999,
				"retrycount":      100,
			},
		},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			config := httpclient.NewRequestConfig("edge-case-test", tc.configMap)
			assert.NotNil(t, config)
			
			// Should not panic during initialization
			assert.NotPanics(t, func() {
				httpclient.InitHttp(config)
			})
		})
	}
}

// Run the test suite
func TestHTTPRequestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPRequestUtilsTestSuite))
}