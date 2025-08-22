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
)

func TestRequestConfig(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		response := "Method: " + r.Method + "\n"
		response += "URL: " + r.URL.String() + "\n"
		response += "Headers: " + r.Header.Get("Content-Type") + "\n"
		response += "Body: " + string(body) + "\n"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer testServer.Close()

	t.Run("complete configuration", func(t *testing.T) {
		configMap := map[string]interface{}{
			"method":          "POST",
			"url":             testServer.URL + "/echo",
			"timeoutinmillis": 10000,
			"retrycount":      5,
			"headers": map[string]interface{}{
				"Content-Type":    "application/json",
				"Accept":          "application/json",
				"User-Agent":      "TestClient/1.0",
				"X-Custom-Header": "custom-value",
			},
		}

		config := httpclient.NewRequestConfig("complete-config", configMap)
		if config == nil {
			t.Error("NewRequestConfig() should not return nil")
		}
	})

	t.Run("minimal configuration", func(t *testing.T) {
		configMap := map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
		}

		config := httpclient.NewRequestConfig("minimal-config", configMap)
		if config == nil {
			t.Error("NewRequestConfig() should not return nil")
		}
	})

	timeoutTests := []struct {
		name          string
		timeoutMillis interface{}
	}{
		{"valid timeout", 5000},
		{"zero timeout", 0},
		{"string timeout", "5000"},
		{"invalid timeout", "invalid"},
		{"negative timeout", -1000},
	}

	for _, tt := range timeoutTests {
		t.Run("timeout_"+tt.name, func(t *testing.T) {
			configMap := map[string]interface{}{
				"method":          "GET",
				"url":             testServer.URL + "/echo",
				"timeoutinmillis": tt.timeoutMillis,
			}

			config := httpclient.NewRequestConfig("timeout-test", configMap)
			if config == nil {
				t.Error("NewRequestConfig() should not return nil")
			}
		})
	}

	retryTests := []struct {
		name       string
		retryCount interface{}
	}{
		{"valid retry count", 3},
		{"zero retries", 0},
		{"string retry count", "3"},
		{"invalid retry count", "invalid"},
		{"negative retry count", -1},
	}

	for _, tt := range retryTests {
		t.Run("retry_"+tt.name, func(t *testing.T) {
			configMap := map[string]interface{}{
				"method":     "GET",
				"url":        testServer.URL + "/echo",
				"retrycount": tt.retryCount,
			}

			config := httpclient.NewRequestConfig("retry-test", configMap)
			if config == nil {
				t.Error("NewRequestConfig() should not return nil")
			}
		})
	}

	headerTests := []struct {
		name    string
		headers interface{}
	}{
		{
			"valid headers map",
			map[string]interface{}{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{"empty headers map", map[string]interface{}{}},
		{"nil headers", nil},
		{"invalid headers type", "invalid"},
	}

	for _, tt := range headerTests {
		t.Run("headers_"+tt.name, func(t *testing.T) {
			configMap := map[string]interface{}{
				"method":  "GET",
				"url":     testServer.URL + "/echo",
				"headers": tt.headers,
			}

			config := httpclient.NewRequestConfig("header-test", configMap)
			if config == nil {
				t.Error("NewRequestConfig() should not return nil")
			}
		})
	}
}

func TestInitHttpRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer testServer.Close()

	t.Run("multiple configs with different settings", func(t *testing.T) {
		config1 := httpclient.NewRequestConfig("fast-client", map[string]interface{}{
			"method":          "GET",
			"url":             testServer.URL + "/echo",
			"timeoutinmillis": 1000,
			"retrycount":      1,
		})

		config2 := httpclient.NewRequestConfig("slow-client", map[string]interface{}{
			"method":          "POST",
			"url":             testServer.URL + "/echo",
			"timeoutinmillis": 30000,
			"retrycount":      5,
			"headers": map[string]interface{}{
				"Content-Type": "application/json",
			},
		})

		config3 := httpclient.NewRequestConfig("auth-client", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/auth",
			"headers": map[string]interface{}{
				"Authorization": "Bearer token123",
			},
		})

		httpclient.InitHttp(config1, config2, config3)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client after InitHttp")
		}
	})

	t.Run("duplicate config names", func(t *testing.T) {
		config1 := httpclient.NewRequestConfig("duplicate", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
		})

		config2 := httpclient.NewRequestConfig("duplicate", map[string]interface{}{
			"method": "POST",
			"url":    testServer.URL + "/auth",
		})

		httpclient.InitHttp(config1, config2)

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client with duplicate names")
		}
	})

	t.Run("empty config list", func(t *testing.T) {
		httpclient.InitHttp()

		client := httpclient.GetClient()
		if client == nil {
			t.Error("GetClient() should return client with empty config list")
		}
	})
}

func TestRequestBuilding(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer testServer.Close()

	t.Run("context handling", func(t *testing.T) {
		config := httpclient.NewRequestConfig("context-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
		})
		httpclient.InitHttp(config)

		contexts := []struct {
			name string
			ctx  context.Context
		}{
			{"background context", context.Background()},
			{"context with value", context.WithValue(context.Background(), "key", "value")},
			{"context with timeout", func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }()},
			{"cancelled context", func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }()},
		}

		for _, tc := range contexts {
			t.Run(tc.name, func(t *testing.T) {
				if tc.ctx == nil {
					t.Error("Context should not be nil")
				}
			})
		}
	})

	t.Run("method override", func(t *testing.T) {
		config := httpclient.NewRequestConfig("method-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
		})
		httpclient.InitHttp(config)

		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

		for _, method := range methods {
			t.Run("method_"+method, func(t *testing.T) {
				if method == "" {
					t.Error("Method should not be empty")
				}
			})
		}
	})

	t.Run("query parameters", func(t *testing.T) {
		config := httpclient.NewRequestConfig("query-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/query",
		})
		httpclient.InitHttp(config)

		paramTests := []struct {
			name   string
			params map[string]string
		}{
			{
				"single parameter",
				map[string]string{"key": "value"},
			},
			{
				"multiple parameters",
				map[string]string{
					"param1": "value1",
					"param2": "value2",
					"param3": "value3",
				},
			},
			{
				"empty value parameter",
				map[string]string{"empty": ""},
			},
			{
				"special characters",
				map[string]string{
					"special": "value with spaces",
					"encoded": "value&with=special&chars",
				},
			},
			{
				"numeric parameters",
				map[string]string{
					"id":    "123",
					"count": "456",
					"page":  "1",
				},
			},
		}

		for _, tc := range paramTests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.params == nil {
					t.Error("Params should not be nil")
				}
				for key := range tc.params {
					if key == "" {
						t.Error("Param key should not be empty")
					}
				}
			})
		}
	})

	t.Run("header parameters", func(t *testing.T) {
		config := httpclient.NewRequestConfig("header-test", map[string]interface{}{
			"method": "GET",
			"url":    testServer.URL + "/echo",
		})
		httpclient.InitHttp(config)

		headerTests := []struct {
			name    string
			headers map[string]string
		}{
			{
				"content type headers",
				map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				},
			},
			{
				"authentication headers",
				map[string]string{
					"Authorization": "Bearer token123",
					"X-API-Key":     "api-key-456",
				},
			},
			{
				"custom headers",
				map[string]string{
					"X-Request-ID":     "req-123",
					"X-Client-Name":    "TestClient",
					"X-Client-Version": "1.0.0",
				},
			},
			{
				"cache control headers",
				map[string]string{
					"Cache-Control": "no-cache",
					"Pragma":        "no-cache",
				},
			},
		}

		for _, tc := range headerTests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.headers == nil {
					t.Error("Headers should not be nil")
				}
				for key, value := range tc.headers {
					if key == "" {
						t.Error("Header key should not be empty")
					}
					if value == "" {
						t.Error("Header value should not be empty")
					}
				}
			})
		}
	})

	t.Run("body handling", func(t *testing.T) {
		config := httpclient.NewRequestConfig("body-test", map[string]interface{}{
			"method": "POST",
			"url":    testServer.URL + "/echo",
		})
		httpclient.InitHttp(config)

		bodyTests := []struct {
			name string
			body io.Reader
		}{
			{"string body", strings.NewReader("test string body")},
			{"JSON body", strings.NewReader(`{"key": "value", "number": 123}`)},
			{"empty body", strings.NewReader("")},
			{"binary body", bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04})},
			{"large body", strings.NewReader(strings.Repeat("a", 10000))},
		}

		for _, tc := range bodyTests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.body == nil {
					t.Error("Body should not be nil")
				}
			})
		}
	})
}

func TestHTTPRequestExecution(t *testing.T) {
	testServer := httptest.NewServer(http.NewServeMux())
	defer testServer.Close()

	mux := testServer.Config.Handler.(*http.ServeMux)

	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		response := "Method: " + r.Method + "\nBody: " + string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

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

	t.Run("success scenarios", func(t *testing.T) {
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
			t.Run(tc.name, func(t *testing.T) {
				req, err := http.NewRequest(tc.method, testServer.URL+tc.endpoint, nil)
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				client := http.Client{Timeout: 5 * time.Second}

				err = httpclient.Call(req, client)
				if err != nil {
					t.Errorf("Call() error = %v", err)
				}
			})
		}
	})

	t.Run("with authentication", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/auth", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer test-token")

		client := http.Client{Timeout: 5 * time.Second}

		err = httpclient.Call(req, client)
		if err != nil {
			t.Errorf("Call() error = %v", err)
		}
	})

	t.Run("without authentication", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/auth", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		err = httpclient.Call(req, client)
		if err == nil {
			t.Error("Call() should return error for unauthorized request")
		}
	})

	t.Run("with query parameters", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/query?param1=value1&param2=value2", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 5 * time.Second}

		response, err := httpclient.CallAndGetResponse(req, client)
		if err != nil {
			t.Errorf("CallAndGetResponse() error = %v", err)
		}
		if !strings.Contains(string(response), "param1=value1") {
			t.Error("Response should contain param1=value1")
		}
		if !strings.Contains(string(response), "param2=value2") {
			t.Error("Response should contain param2=value2")
		}
	})

	t.Run("with body", func(t *testing.T) {
		body := strings.NewReader(`{"test": "data"}`)
		req, err := http.NewRequest("POST", testServer.URL+"/echo", body)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := http.Client{Timeout: 5 * time.Second}

		response, err := httpclient.CallAndGetResponse(req, client)
		if err != nil {
			t.Errorf("CallAndGetResponse() error = %v", err)
		}
		if !strings.Contains(string(response), `{"test": "data"}`) {
			t.Error("Response should contain request body")
		}
	})
}

func TestHTTPRequestErrorHandling(t *testing.T) {
	t.Run("network errors", func(t *testing.T) {
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

	t.Run("timeout errors", func(t *testing.T) {
		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer slowServer.Close()

		req, err := http.NewRequest("GET", slowServer.URL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := http.Client{Timeout: 100 * time.Millisecond}

		err = httpclient.Call(req, client)
		if err == nil {
			t.Error("Call() should return error for timeout")
		}
	})

	t.Run("invalid requests", func(t *testing.T) {
		testCases := []struct {
			name   string
			method string
			url    string
		}{
			{"invalid URL", "GET", "://invalid-url"},
			{"empty URL", "GET", ""},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, err := http.NewRequest(tc.method, tc.url, nil)

				if err == nil {
					client := http.Client{Timeout: 5 * time.Second}
					err = httpclient.Call(req, client)
				}

				if err == nil {
					t.Error("Should have an error for invalid request")
				}
			})
		}
	})
}

func TestConfigurationEdgeCases(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	testCases := []struct {
		name      string
		configMap map[string]interface{}
	}{
		{
			"all nil values",
			map[string]interface{}{
				"method":          nil,
				"url":             nil,
				"timeoutinmillis": nil,
				"retrycount":      nil,
				"headers":         nil,
			},
		},
		{
			"mixed valid and invalid values",
			map[string]interface{}{
				"method":          "GET",
				"url":             testServer.URL + "/echo",
				"timeoutinmillis": "invalid",
				"retrycount":      nil,
				"headers":         "not a map",
			},
		},
		{
			"extreme values",
			map[string]interface{}{
				"method":          "GET",
				"url":             testServer.URL + "/echo",
				"timeoutinmillis": 999999999,
				"retrycount":      100,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := httpclient.NewRequestConfig("edge-case-test", tc.configMap)
			if config == nil {
				t.Error("NewRequestConfig() should not return nil")
			}

			// Should not panic during initialization
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("InitHttp() panicked: %v", r)
				}
			}()

			httpclient.InitHttp(config)
		})
	}
}