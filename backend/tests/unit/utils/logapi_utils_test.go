package utils

import (
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LogParsingTestSuite defines the test suite for log parsing functionality
type LogParsingTestSuite struct {
	suite.Suite
}

// TestParseLogLine_ValidJSONEntries tests parsing of valid JSON log entries
func (suite *LogParsingTestSuite) TestParseLogLine_ValidJSONEntries() {
	tests := []struct {
		name     string
		logLine  string
		expected modelsv1.LogEntry
	}{
		{
			name:    "Complete JSON log entry with all fields",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","log_id":"req-123","level":"INFO","message":"Request processed successfully","file":"main.go","line":42,"function":"handleRequest"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				LogID:     "req-123",
				Level:     "INFO",
				Message:   "Request processed successfully",
				File:      "main.go",
				Line:      42,
				Function:  "handleRequest",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with minimal fields",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","level":"ERROR","message":"Database connection failed"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "ERROR",
				Message:   "Database connection failed",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with HTTP request details in message",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","log_id":"req-456","level":"INFO","message":"GET /api/v1/users | Status: 200 | Duration: 45ms | IP: 192.168.1.100 | User-Agent: Mozilla/5.0"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				LogID:     "req-456",
				Level:     "INFO",
				Message:   "GET /api/v1/users | Status: 200 | Duration: 45ms | IP: 192.168.1.100 | User-Agent: Mozilla/5.0",
				Method:    "GET",
				Path:      "/api/v1/users",
				Status:    200,
				Duration:  "45ms",
				IP:        "192.168.1.100",
				UserAgent: "Mozilla/5.0",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with POST request",
			logLine: `{"timestamp":"2025-08-20T10:31:00Z","log_id":"req-789","level":"INFO","message":"POST /api/v1/posts | Status: 201 | Duration: 120ms | IP: 10.0.0.1"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:31:00Z",
				LogID:     "req-789",
				Level:     "INFO",
				Message:   "POST /api/v1/posts | Status: 201 | Duration: 120ms | IP: 10.0.0.1",
				Method:    "POST",
				Path:      "/api/v1/posts",
				Status:    201,
				Duration:  "120ms",
				IP:        "10.0.0.1",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with complex User-Agent",
			logLine: `{"timestamp":"2025-08-20T10:32:00Z","level":"DEBUG","message":"PUT /api/v1/settings | Status: 204 | User-Agent: PostEaze-Client/1.0 (Windows NT 10.0; Win64; x64)"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:32:00Z",
				Level:     "DEBUG",
				Message:   "PUT /api/v1/settings | Status: 204 | User-Agent: PostEaze-Client/1.0 (Windows NT 10.0; Win64; x64)",
				Method:    "PUT",
				Path:      "/api/v1/settings",
				Status:    204,
				UserAgent: "PostEaze-Client/1.0 (Windows NT 10.0; Win64; x64)",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with all HTTP methods",
			logLine: `{"timestamp":"2025-08-20T10:33:00Z","level":"INFO","message":"DELETE /api/v1/posts/123 | Status: 404 | Duration: 15ms"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:33:00Z",
				Level:     "INFO",
				Message:   "DELETE /api/v1/posts/123 | Status: 404 | Duration: 15ms",
				Method:    "DELETE",
				Path:      "/api/v1/posts/123",
				Status:    404,
				Duration:  "15ms",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with PATCH method",
			logLine: `{"timestamp":"2025-08-20T10:34:00Z","level":"INFO","message":"PATCH /api/v1/users/456 | Status: 200"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:34:00Z",
				Level:     "INFO",
				Message:   "PATCH /api/v1/users/456 | Status: 200",
				Method:    "PATCH",
				Path:      "/api/v1/users/456",
				Status:    200,
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with HEAD method",
			logLine: `{"timestamp":"2025-08-20T10:35:00Z","level":"INFO","message":"HEAD /api/v1/health | Status: 200 | Duration: 5ms"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:35:00Z",
				Level:     "INFO",
				Message:   "HEAD /api/v1/health | Status: 200 | Duration: 5ms",
				Method:    "HEAD",
				Path:      "/api/v1/health",
				Status:    200,
				Duration:  "5ms",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON log entry with OPTIONS method",
			logLine: `{"timestamp":"2025-08-20T10:36:00Z","level":"INFO","message":"OPTIONS /api/v1/cors | Status: 204"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:36:00Z",
				Level:     "INFO",
				Message:   "OPTIONS /api/v1/cors | Status: 204",
				Method:    "OPTIONS",
				Path:      "/api/v1/cors",
				Status:    204,
				Extra:     make(map[string]string),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := utils.ParseLogLine(tt.logLine)
			
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.expected.Timestamp, result.Timestamp)
			assert.Equal(suite.T(), tt.expected.LogID, result.LogID)
			assert.Equal(suite.T(), tt.expected.Level, result.Level)
			assert.Equal(suite.T(), tt.expected.Message, result.Message)
			assert.Equal(suite.T(), tt.expected.File, result.File)
			assert.Equal(suite.T(), tt.expected.Line, result.Line)
			assert.Equal(suite.T(), tt.expected.Function, result.Function)
			assert.Equal(suite.T(), tt.expected.Method, result.Method)
			assert.Equal(suite.T(), tt.expected.Path, result.Path)
			assert.Equal(suite.T(), tt.expected.Status, result.Status)
			assert.Equal(suite.T(), tt.expected.Duration, result.Duration)
			assert.Equal(suite.T(), tt.expected.IP, result.IP)
			assert.Equal(suite.T(), tt.expected.UserAgent, result.UserAgent)
			assert.NotNil(suite.T(), result.Extra)
		})
	}
}

// TestParseLogLine_MalformedJSONEntries tests handling of malformed JSON entries
func (suite *LogParsingTestSuite) TestParseLogLine_MalformedJSONEntries() {
	tests := []struct {
		name        string
		logLine     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid JSON - missing closing brace",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test message"`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Invalid JSON - missing quotes",
			logLine:     `{timestamp:"2025-08-20T10:30:45Z",level:"INFO",message:"Test message"}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Invalid JSON - extra comma",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test message",}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Invalid JSON - malformed string",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test message with unescaped "quotes"}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Empty string",
			logLine:     ``,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Non-JSON string",
			logLine:     `This is not JSON at all`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Invalid JSON - wrong data type for line field",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test","line":"not_a_number"}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "Partial JSON object",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z"`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := utils.ParseLogLine(tt.logLine)
			
			if tt.expectError {
				assert.Error(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tt.errorMsg)
				assert.Equal(suite.T(), modelsv1.LogEntry{}, result)
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

// TestParseLogLine_HTTPRequestDetailsExtraction tests extraction of HTTP request details from message field
func (suite *LogParsingTestSuite) TestParseLogLine_HTTPRequestDetailsExtraction() {
	tests := []struct {
		name            string
		logLine         string
		expectedMethod  string
		expectedPath    string
		expectedStatus  int
		expectedDuration string
		expectedIP      string
		expectedUA      string
	}{
		{
			name:            "Complete HTTP request details",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/users | Status: 200 | Duration: 45ms | IP: 192.168.1.100 | User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)"}`,
			expectedMethod:  "GET",
			expectedPath:    "/api/v1/users",
			expectedStatus:  200,
			expectedDuration: "45ms",
			expectedIP:      "192.168.1.100",
			expectedUA:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		{
			name:            "HTTP request with query parameters",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/posts?page=1&limit=10 | Status: 200 | Duration: 30ms"}`,
			expectedMethod:  "GET",
			expectedPath:    "/api/v1/posts?page=1&limit=10",
			expectedStatus:  200,
			expectedDuration: "30ms",
		},
		{
			name:            "POST request with creation status",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"POST /api/v1/posts | Status: 201 | Duration: 150ms | IP: 10.0.0.1"}`,
			expectedMethod:  "POST",
			expectedPath:    "/api/v1/posts",
			expectedStatus:  201,
			expectedDuration: "150ms",
			expectedIP:      "10.0.0.1",
		},
		{
			name:            "Error status code",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/nonexistent | Status: 404 | Duration: 5ms"}`,
			expectedMethod:  "GET",
			expectedPath:    "/api/v1/nonexistent",
			expectedStatus:  404,
			expectedDuration: "5ms",
		},
		{
			name:            "Server error status",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"POST /api/v1/process | Status: 500 | Duration: 1000ms"}`,
			expectedMethod:  "POST",
			expectedPath:    "/api/v1/process",
			expectedStatus:  500,
			expectedDuration: "1000ms",
		},
		{
			name:            "Complex path with ID",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"PUT /api/v1/users/12345/profile | Status: 200"}`,
			expectedMethod:  "PUT",
			expectedPath:    "/api/v1/users/12345/profile",
			expectedStatus:  200,
		},
		{
			name:            "Duration in different formats",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/slow | Status: 200 | Duration: 2.5s"}`,
			expectedMethod:  "GET",
			expectedPath:    "/api/v1/slow",
			expectedStatus:  200,
			expectedDuration: "2.5s",
		},
		{
			name:            "Complex User-Agent with special characters",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/data | Status: 200 | User-Agent: PostEaze-Client/1.0 (Linux; Android 11; SM-G991B)"}`,
			expectedMethod:  "GET",
			expectedPath:    "/api/v1/data",
			expectedStatus:  200,
			expectedUA:      "PostEaze-Client/1.0 (Linux; Android 11; SM-G991B)",
		},
		{
			name:           "Message without HTTP details",
			logLine:        `{"timestamp":"2025-08-20T10:30:45Z","message":"Database connection established"}`,
			expectedMethod: "",
			expectedPath:   "",
			expectedStatus: 0,
		},
		{
			name:            "Partial HTTP details - only method and path",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"DELETE /api/v1/temp/cleanup"}`,
			expectedMethod:  "DELETE",
			expectedPath:    "/api/v1/temp/cleanup",
			expectedStatus:  0,
		},
		{
			name:            "IPv6 address",
			logLine:         `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/test | Status: 200 | IP: 2001:db8::1"}`,
			expectedMethod:  "GET",
			expectedPath:    "/api/v1/test",
			expectedStatus:  200,
			expectedIP:      "2001:db8::1",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := utils.ParseLogLine(tt.logLine)
			
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.expectedMethod, result.Method)
			assert.Equal(suite.T(), tt.expectedPath, result.Path)
			assert.Equal(suite.T(), tt.expectedStatus, result.Status)
			assert.Equal(suite.T(), tt.expectedDuration, result.Duration)
			assert.Equal(suite.T(), tt.expectedIP, result.IP)
			assert.Equal(suite.T(), tt.expectedUA, result.UserAgent)
		})
	}
}

// TestParseLogLine_EdgeCases tests edge cases and boundary conditions
func (suite *LogParsingTestSuite) TestParseLogLine_EdgeCases() {
	tests := []struct {
		name     string
		logLine  string
		expected modelsv1.LogEntry
	}{
		{
			name:    "JSON with null values",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","log_id":null,"level":"INFO","message":"Test message","file":null,"line":0,"function":null}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "INFO",
				Message:   "Test message",
				Line:      0,
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON with empty strings",
			logLine: `{"timestamp":"","log_id":"","level":"","message":"","file":"","line":0,"function":""}`,
			expected: modelsv1.LogEntry{
				Line:  0,
				Extra: make(map[string]string),
			},
		},
		{
			name:    "JSON with very long message",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"This is a very long message that contains a lot of text and should be handled properly by the parsing function without any issues or truncation problems that might occur during processing of large log entries"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "INFO",
				Message:   "This is a very long message that contains a lot of text and should be handled properly by the parsing function without any issues or truncation problems that might occur during processing of large log entries",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON with special characters in message",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","level":"ERROR","message":"Error: Failed to process request with special chars: @#$%^&*()[]{}|\\;':\",./<>?"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "ERROR",
				Message:   "Error: Failed to process request with special chars: @#$%^&*()[]{}|\\;':\",./<>?",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON with Unicode characters",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"User login: 用户登录成功 🎉"}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "INFO",
				Message:   "User login: 用户登录成功 🎉",
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON with negative line number",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","level":"DEBUG","message":"Debug info","line":-1}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "DEBUG",
				Message:   "Debug info",
				Line:      -1,
				Extra:     make(map[string]string),
			},
		},
		{
			name:    "JSON with very large line number",
			logLine: `{"timestamp":"2025-08-20T10:30:45Z","level":"TRACE","message":"Trace info","line":999999}`,
			expected: modelsv1.LogEntry{
				Timestamp: "2025-08-20T10:30:45Z",
				Level:     "TRACE",
				Message:   "Trace info",
				Line:      999999,
				Extra:     make(map[string]string),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := utils.ParseLogLine(tt.logLine)
			
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.expected.Timestamp, result.Timestamp)
			assert.Equal(suite.T(), tt.expected.LogID, result.LogID)
			assert.Equal(suite.T(), tt.expected.Level, result.Level)
			assert.Equal(suite.T(), tt.expected.Message, result.Message)
			assert.Equal(suite.T(), tt.expected.File, result.File)
			assert.Equal(suite.T(), tt.expected.Line, result.Line)
			assert.Equal(suite.T(), tt.expected.Function, result.Function)
			assert.NotNil(suite.T(), result.Extra)
		})
	}
}

// Run the test suite
func TestLogParsingTestSuite(t *testing.T) {
	suite.Run(t, new(LogParsingTestSuite))
}