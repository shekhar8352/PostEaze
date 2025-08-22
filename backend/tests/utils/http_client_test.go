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
)

func TestNewRequestConfig(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer testServer.Close()

	tests := []struct {
		name      string
		configMap map[string]interface{}
	}{
		{
			name: "valid config",
			configMap: map[string]interface{}{
				"method":          "POST",
				"url":             testServer.URL + "/success",
				"timeoutinmillis": 5000,
				"retrycount":      3,
				"headers": map[string]interface{}{
					"Content-Type":  "application/json",
					"Authorization": "Bearer token123",
				},
			},
		},
		{
			name:      "empty config",
			configMap: nil,
		},
		{
			name: "partial config",
			configMap: map[string]interface{}{
				"method": "GET",
				"url":    testServer.URL + "/success",
			},
		},
		{
			name: "invalid types",
			configMap: map[string]interface{}{
				"method":          "GET",
				"url":             testServer.URL + "/success",
				"timeoutinmillis": "invalid", // Should be int
				"retrycount":      "invalid", // Should be int
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := httpclient.NewRequestConfig("test-config", tt.configMap)
			if config == nil {
				t.Error("NewRequestConfig() should not return nil")
			}
		})
	}
}

func TestInitHttp(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer testServer.Close()

	t.Run("single config", func(t *testing.T) {
		config := httpclient.NewRequestConfig("test-single", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/success",
		})

		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client after InitHttp")
		}
	})

	t.Run("multiple configs", func(t *testing.T) {
		config1 := httpclient.NewRequestConfig("test-config-1", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/success",
		})

		config2 := httpclient.NewRequestConfig("test-config-2", map[string]interface{}{
			"method": "POST",
			"url":    testServer.URL + "/echo",
		})

		httpclient.InitHttp(config1, config2)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client after InitHttp with multiple configs")
		}
	})

	t.Run("nil config", func(t *testing.T) {
		httpclient.InitHttp(nil)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client even with nil config")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		config := httpclient.NewRequestConfig("", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/success",
		})

		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client with empty name config")
		}
	})
}

func TestGetClient(t *testing.T) {
	client := httpclient.GetClient()
	if client == nil {
		t.Error("GetClient() should not return nil")
	}

	// Test that multiple calls return the same instance
	client2 := httpclient.GetClient()
	if client != client2 {
		t.Error("GetClient() should return the same instance on multiple calls")
	}
}

func TestHTTPExecution(t *testing.T) {
	// Create test server with multiple endpoints
	testServer := httptest.NewServer(http.NewServeMux())
	defer testServer.Close()

	mux := testServer.Config.Handler.(*http.ServeMux)

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
			"method":  r.Method,
			"url":     r.URL.String(),
			"headers": r.Header,
			"body":    string(body),
			"query":   r.URL.Query(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	t.Run("call success", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/success", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		err = httpclient.Call(req, client)
		if err != nil {
			t.Errorf("Call() error = %v", err)
		}
	})

	t.Run("call error", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/error", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		err = httpclient.Call(req, client)
		if err == nil {
			t.Error("Call() should return error for error endpoint")
		}
		if err != nil && err.Error() == "" {
			t.Error("Call() should return meaningful error message")
		}
	})

	t.Run("call timeout", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/slow", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 100 * time.Millisecond}

		err = httpclient.Call(req, client)
		if err == nil {
			t.Error("Call() should return error for timeout")
		}
	})

	t.Run("call and get response success", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/success", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		response, err := httpclient.CallAndGetResponse(req, client)
		if err != nil {
			t.Errorf("CallAndGetResponse() error = %v", err)
		}
		if len(response) == 0 {
			t.Error("CallAndGetResponse() should return non-empty response")
		}

		// Parse response to verify it's valid JSON
		var responseData map[string]interface{}
		err = json.Unmarshal(response, &responseData)
		if err != nil {
			t.Errorf("CallAndGetResponse() returned invalid JSON: %v", err)
		}
		if responseData["message"] != "success" {
			t.Errorf("CallAndGetResponse() message = %v, want success", responseData["message"])
		}
	})

	t.Run("call and get response error", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/error", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		response, err := httpclient.CallAndGetResponse(req, client)
		if err == nil {
			t.Error("CallAndGetResponse() should return error for error endpoint")
		}
		if response != nil {
			t.Error("CallAndGetResponse() should return nil response on error")
		}
	})

	t.Run("call and bind success", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/success", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		var response map[string]interface{}
		err = httpclient.CallAndBind(req, &response, client)
		if err != nil {
			t.Errorf("CallAndBind() error = %v", err)
		}
		if response["message"] != "success" {
			t.Errorf("CallAndBind() message = %v, want success", response["message"])
		}
		if response["method"] != "GET" {
			t.Errorf("CallAndBind() method = %v, want GET", response["method"])
		}
	})

	t.Run("call and bind invalid JSON", func(t *testing.T) {
		// Create a server that returns invalid JSON
		invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json response"))
		}))
		defer invalidServer.Close()

		req, err := http.NewRequest("GET", invalidServer.URL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		var response map[string]interface{}
		err = httpclient.CallAndBind(req, &response, client)
		if err == nil {
			t.Error("CallAndBind() should return error for invalid JSON")
		}
		if err != nil && err.Error() == "" {
			t.Error("CallAndBind() should return meaningful error message")
		}
	})

	t.Run("call and bind wrong type", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/success", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		// Try to bind to wrong type
		var response string
		err = httpclient.CallAndBind(req, &response, client)
		if err == nil {
			t.Error("CallAndBind() should return error for wrong type")
		}
	})
}

func TestHTTPRequestWithBody(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		response := map[string]interface{}{
			"method":  r.Method,
			"body":    string(body),
			"headers": r.Header,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer testServer.Close()

	requestBody := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", testServer.URL, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 5 * time.Second}

	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	if err != nil {
		t.Errorf("CallAndBind() error = %v", err)
	}
	if response["method"] != "POST" {
		t.Errorf("Response method = %v, want POST", response["method"])
	}
	if !bytes.Contains([]byte(response["body"].(string)), []byte("test")) {
		t.Error("Response body should contain 'test'")
	}
	if !bytes.Contains([]byte(response["body"].(string)), []byte("123")) {
		t.Error("Response body should contain '123'")
	}
}

func TestHTTPRequestWithHeaders(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"headers": r.Header,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer testServer.Close()

	req, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("Authorization", "Bearer token123")

	client := http.Client{Timeout: 5 * time.Second}

	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	if err != nil {
		t.Errorf("CallAndBind() error = %v", err)
	}

	headers := response["headers"].(map[string]interface{})
	if _, exists := headers["X-Custom-Header"]; !exists {
		t.Error("Response should contain X-Custom-Header")
	}
	if _, exists := headers["Authorization"]; !exists {
		t.Error("Response should contain Authorization header")
	}
}

func TestHTTPRequestWithQueryParams(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"query": r.URL.Query(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer testServer.Close()

	req, err := http.NewRequest("GET", testServer.URL+"?param1=value1&param2=value2", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}

	var response map[string]interface{}
	err = httpclient.CallAndBind(req, &response, client)
	if err != nil {
		t.Errorf("CallAndBind() error = %v", err)
	}

	query := response["query"].(map[string]interface{})
	if _, exists := query["param1"]; !exists {
		t.Error("Response should contain param1")
	}
	if _, exists := query["param2"]; !exists {
		t.Error("Response should contain param2")
	}
}

func TestHTTPErrorHandling(t *testing.T) {
	t.Run("invalid URL", func(t *testing.T) {
		req, err := http.NewRequest("GET", "invalid-url", nil)

		if err == nil {
			client := http.Client{Timeout: 5 * time.Second}
			err = httpclient.Call(req, client)
			if err == nil {
				t.Error("Call() should return error for invalid URL")
			}
		} else {
			if err == nil {
				t.Error("NewRequest() should return error for invalid URL")
			}
		}
	})

	t.Run("connection refused", func(t *testing.T) {
		// Use a port that's likely not in use
		req, err := http.NewRequest("GET", "http://localhost:99999/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 1 * time.Second}

		err = httpclient.Call(req, client)
		if err == nil {
			t.Error("Call() should return error for connection refused")
		}
	})
}

func TestRequestBuilder(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer testServer.Close()

	t.Run("set context", func(t *testing.T) {
		config := httpclient.NewRequestConfig("context-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/success",
		})
		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		ctx := context.WithValue(context.Background(), "test", "value")

		if client == nil {
			t.Error("GetClient() should return client")
		}
		if ctx == nil {
			t.Error("Context should not be nil")
		}
	})

	t.Run("set method", func(t *testing.T) {
		config := httpclient.NewRequestConfig("method-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
		})
		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client")
		}
	})

	t.Run("set URL", func(t *testing.T) {
		config := httpclient.NewRequestConfig("url-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/success",
		})
		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client")
		}
	})

	t.Run("set headers", func(t *testing.T) {
		config := httpclient.NewRequestConfig("header-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
			"headers": map[string]interface{}{
				"X-Test-Header": "test-value",
			},
		})
		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client")
		}
	})

	t.Run("set body", func(t *testing.T) {
		config := httpclient.NewRequestConfig("body-test", map[string]interface{}{
			"method": "POST",
			"url":    testServer.URL + "/echo",
		})
		httpclient.InitHttp(config)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client")
		}
	})
}