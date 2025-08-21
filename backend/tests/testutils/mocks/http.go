package mocks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stretchr/testify/mock"
)

// HTTPClient interface for mocking standard http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MockHTTPClient is a mock implementation of HTTPClient
type MockHTTPClient struct {
	mock.Mock
}

// Do mocks the Do method of http.Client
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{}
}

// MockHTTPResponse represents a mock HTTP response
type MockHTTPResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
	Error      error
}

// HTTPClientMockManager provides utilities for managing HTTP client mocks
type HTTPClientMockManager struct {
	mockClient *MockHTTPClient
	responses  map[string]*MockHTTPResponse
}

// NewHTTPClientMockManager creates a new HTTP client mock manager
func NewHTTPClientMockManager() *HTTPClientMockManager {
	return &HTTPClientMockManager{
		mockClient: NewMockHTTPClient(),
		responses:  make(map[string]*MockHTTPResponse),
	}
}

// GetMockClient returns the mock HTTP client
func (m *HTTPClientMockManager) GetMockClient() *MockHTTPClient {
	return m.mockClient
}

// SetupResponse configures a mock response for a specific endpoint
func (m *HTTPClientMockManager) SetupResponse(method, url string, response *MockHTTPResponse) {
	key := fmt.Sprintf("%s:%s", method, url)
	m.responses[key] = response
	
	if response.Error != nil {
		m.mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == method && req.URL.String() == url
		})).Return((*http.Response)(nil), response.Error)
	} else {
		httpResponse := &http.Response{
			StatusCode: response.StatusCode,
			Body:       io.NopCloser(strings.NewReader(response.Body)),
			Header:     make(http.Header),
		}
		
		for k, v := range response.Headers {
			httpResponse.Header.Set(k, v)
		}
		
		m.mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == method && req.URL.String() == url
		})).Return(httpResponse, nil)
	}
}

// SetupJSONResponse configures a mock JSON response for a specific endpoint
func (m *HTTPClientMockManager) SetupJSONResponse(method, url string, statusCode int, responseData interface{}) error {
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return err
	}
	
	response := &MockHTTPResponse{
		StatusCode: statusCode,
		Body:       string(jsonData),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	
	m.SetupResponse(method, url, response)
	return nil
}

// SetupErrorResponse configures a mock error response for a specific endpoint
func (m *HTTPClientMockManager) SetupErrorResponse(method, url string, err error) {
	response := &MockHTTPResponse{
		Error: err,
	}
	m.SetupResponse(method, url, response)
}

// SetupSuccessResponse configures a mock success response for a specific endpoint
func (m *HTTPClientMockManager) SetupSuccessResponse(method, url string, body string) {
	response := &MockHTTPResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
	m.SetupResponse(method, url, response)
}

// VerifyRequestCalled verifies that a request was made to a specific endpoint
func (m *HTTPClientMockManager) VerifyRequestCalled(method, url string) bool {
	for _, call := range m.mockClient.Calls {
		if call.Method == "Do" && len(call.Arguments) > 0 {
			if req, ok := call.Arguments[0].(*http.Request); ok {
				if req.Method == method && req.URL.String() == url {
					return true
				}
			}
		}
	}
	return false
}

// VerifyRequestCalledWithBody verifies that a request was made with specific body content
func (m *HTTPClientMockManager) VerifyRequestCalledWithBody(method, url, expectedBody string) bool {
	for _, call := range m.mockClient.Calls {
		if call.Method == "Do" && len(call.Arguments) > 0 {
			if req, ok := call.Arguments[0].(*http.Request); ok {
				if req.Method == method && req.URL.String() == url {
					if req.Body != nil {
						bodyBytes, err := io.ReadAll(req.Body)
						if err == nil && string(bodyBytes) == expectedBody {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// VerifyRequestCalledWithHeaders verifies that a request was made with specific headers
func (m *HTTPClientMockManager) VerifyRequestCalledWithHeaders(method, url string, expectedHeaders map[string]string) bool {
	for _, call := range m.mockClient.Calls {
		if call.Method == "Do" && len(call.Arguments) > 0 {
			if req, ok := call.Arguments[0].(*http.Request); ok {
				if req.Method == method && req.URL.String() == url {
					for key, expectedValue := range expectedHeaders {
						if req.Header.Get(key) != expectedValue {
							return false
						}
					}
					return true
				}
			}
		}
	}
	return false
}

// GetRequestCount returns the number of requests made to a specific endpoint
func (m *HTTPClientMockManager) GetRequestCount(method, url string) int {
	count := 0
	for _, call := range m.mockClient.Calls {
		if call.Method == "Do" && len(call.Arguments) > 0 {
			if req, ok := call.Arguments[0].(*http.Request); ok {
				if req.Method == method && req.URL.String() == url {
					count++
				}
			}
		}
	}
	return count
}

// AssertExpectations asserts that all expectations were met
func (m *HTTPClientMockManager) AssertExpectations(t mock.TestingT) {
	m.mockClient.AssertExpectations(t)
}

// Reset clears all expectations and call history
func (m *HTTPClientMockManager) Reset() {
	m.mockClient.ExpectedCalls = nil
	m.mockClient.Calls = nil
	m.responses = make(map[string]*MockHTTPResponse)
}

// HTTPRequestBuilder provides utilities for building HTTP requests for testing
type HTTPRequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    string
	query   map[string]string
}

// NewHTTPRequestBuilder creates a new HTTP request builder
func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		headers: make(map[string]string),
		query:   make(map[string]string),
	}
}

// Method sets the HTTP method
func (b *HTTPRequestBuilder) Method(method string) *HTTPRequestBuilder {
	b.method = method
	return b
}

// URL sets the request URL
func (b *HTTPRequestBuilder) URL(url string) *HTTPRequestBuilder {
	b.url = url
	return b
}

// Header adds a header to the request
func (b *HTTPRequestBuilder) Header(key, value string) *HTTPRequestBuilder {
	b.headers[key] = value
	return b
}

// Headers adds multiple headers to the request
func (b *HTTPRequestBuilder) Headers(headers map[string]string) *HTTPRequestBuilder {
	for k, v := range headers {
		b.headers[k] = v
	}
	return b
}

// Body sets the request body
func (b *HTTPRequestBuilder) Body(body string) *HTTPRequestBuilder {
	b.body = body
	return b
}

// JSONBody sets the request body as JSON
func (b *HTTPRequestBuilder) JSONBody(data interface{}) *HTTPRequestBuilder {
	jsonData, _ := json.Marshal(data)
	b.body = string(jsonData)
	b.headers["Content-Type"] = "application/json"
	return b
}

// Query adds a query parameter
func (b *HTTPRequestBuilder) Query(key, value string) *HTTPRequestBuilder {
	b.query[key] = value
	return b
}

// Build creates the HTTP request
func (b *HTTPRequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	var bodyReader io.Reader
	if b.body != "" {
		bodyReader = strings.NewReader(b.body)
	}
	
	req, err := http.NewRequestWithContext(ctx, b.method, b.url, bodyReader)
	if err != nil {
		return nil, err
	}
	
	// Add headers
	for k, v := range b.headers {
		req.Header.Set(k, v)
	}
	
	// Add query parameters
	if len(b.query) > 0 {
		q := req.URL.Query()
		for k, v := range b.query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	
	return req, nil
}

// HTTPResponseBuilder provides utilities for building HTTP responses for testing
type HTTPResponseBuilder struct {
	statusCode int
	body       string
	headers    map[string]string
}

// NewHTTPResponseBuilder creates a new HTTP response builder
func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{
		statusCode: http.StatusOK,
		headers:    make(map[string]string),
	}
}

// StatusCode sets the response status code
func (b *HTTPResponseBuilder) StatusCode(code int) *HTTPResponseBuilder {
	b.statusCode = code
	return b
}

// Body sets the response body
func (b *HTTPResponseBuilder) Body(body string) *HTTPResponseBuilder {
	b.body = body
	return b
}

// JSONBody sets the response body as JSON
func (b *HTTPResponseBuilder) JSONBody(data interface{}) *HTTPResponseBuilder {
	jsonData, _ := json.Marshal(data)
	b.body = string(jsonData)
	b.headers["Content-Type"] = "application/json"
	return b
}

// Header adds a header to the response
func (b *HTTPResponseBuilder) Header(key, value string) *HTTPResponseBuilder {
	b.headers[key] = value
	return b
}

// Build creates the HTTP response
func (b *HTTPResponseBuilder) Build() *http.Response {
	response := &http.Response{
		StatusCode: b.statusCode,
		Body:       io.NopCloser(strings.NewReader(b.body)),
		Header:     make(http.Header),
	}
	
	for k, v := range b.headers {
		response.Header.Set(k, v)
	}
	
	return response
}

// Common HTTP errors for testing
var (
	ErrMockHTTPTimeout      = errors.New("mock: request timeout")
	ErrMockHTTPConnection   = errors.New("mock: connection failed")
	ErrMockUnauthorized     = errors.New("mock: unauthorized")
	ErrMockNotFound         = errors.New("mock: not found")
	ErrMockServerError      = errors.New("mock: internal server error")
	ErrMockBadRequest       = errors.New("mock: bad request")
)

// HTTPMockScenario represents a complete HTTP mock scenario
type HTTPMockScenario struct {
	Name        string
	Method      string
	URL         string
	RequestBody string
	Headers     map[string]string
	Response    *MockHTTPResponse
}

// HTTPScenarioManager manages multiple HTTP mock scenarios
type HTTPScenarioManager struct {
	manager   *HTTPClientMockManager
	scenarios []HTTPMockScenario
}

// NewHTTPScenarioManager creates a new HTTP scenario manager
func NewHTTPScenarioManager() *HTTPScenarioManager {
	return &HTTPScenarioManager{
		manager:   NewHTTPClientMockManager(),
		scenarios: make([]HTTPMockScenario, 0),
	}
}

// AddScenario adds a new HTTP mock scenario
func (s *HTTPScenarioManager) AddScenario(scenario HTTPMockScenario) {
	s.scenarios = append(s.scenarios, scenario)
	s.manager.SetupResponse(scenario.Method, scenario.URL, scenario.Response)
}

// GetManager returns the underlying HTTP client mock manager
func (s *HTTPScenarioManager) GetManager() *HTTPClientMockManager {
	return s.manager
}

// VerifyAllScenarios verifies that all scenarios were called
func (s *HTTPScenarioManager) VerifyAllScenarios() []string {
	var failures []string
	for _, scenario := range s.scenarios {
		if !s.manager.VerifyRequestCalled(scenario.Method, scenario.URL) {
			failures = append(failures, fmt.Sprintf("Scenario '%s' was not called: %s %s", 
				scenario.Name, scenario.Method, scenario.URL))
		}
	}
	return failures
}

// RequestMatcher provides utilities for matching HTTP requests
type RequestMatcher struct{}

// NewRequestMatcher creates a new request matcher
func NewRequestMatcher() *RequestMatcher {
	return &RequestMatcher{}
}

// MatchMethod creates a matcher for HTTP method
func (m *RequestMatcher) MatchMethod(expectedMethod string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		return req.Method == expectedMethod
	}
}

// MatchURL creates a matcher for URL
func (m *RequestMatcher) MatchURL(expectedURL string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		return req.URL.String() == expectedURL
	}
}

// MatchURLPattern creates a matcher for URL pattern
func (m *RequestMatcher) MatchURLPattern(pattern string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		return strings.Contains(req.URL.String(), pattern)
	}
}

// MatchHeader creates a matcher for specific header
func (m *RequestMatcher) MatchHeader(key, value string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		return req.Header.Get(key) == value
	}
}

// MatchBody creates a matcher for request body
func (m *RequestMatcher) MatchBody(expectedBody string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		if req.Body == nil {
			return expectedBody == ""
		}
		
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return false
		}
		
		// Reset body for potential re-reading
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		
		return string(bodyBytes) == expectedBody
	}
}

// MatchJSONBody creates a matcher for JSON request body
func (m *RequestMatcher) MatchJSONBody(expectedData interface{}) func(*http.Request) bool {
	return func(req *http.Request) bool {
		if req.Body == nil {
			return false
		}
		
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return false
		}
		
		// Reset body for potential re-reading
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		
		expectedJSON, err := json.Marshal(expectedData)
		if err != nil {
			return false
		}
		
		return string(bodyBytes) == string(expectedJSON)
	}
}

// MatchQuery creates a matcher for query parameters
func (m *RequestMatcher) MatchQuery(key, value string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		return req.URL.Query().Get(key) == value
	}
}

// CombineMatchers combines multiple matchers with AND logic
func (m *RequestMatcher) CombineMatchers(matchers ...func(*http.Request) bool) func(*http.Request) bool {
	return func(req *http.Request) bool {
		for _, matcher := range matchers {
			if !matcher(req) {
				return false
			}
		}
		return true
	}
}