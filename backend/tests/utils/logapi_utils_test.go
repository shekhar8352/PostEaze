package utils

import (
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func TestParseLogLine_ValidJSONEntries(t *testing.T) {
	tests := []struct {
		name     string
		logLine  string
		expected modelsv1.LogEntry
	}{
		{
			name:    "complete JSON log entry with all fields",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ParseLogLine(tt.logLine)
			if err != nil {
				t.Errorf("ParseLogLine() error = %v", err)
				return
			}

			if result.Timestamp != tt.expected.Timestamp {
				t.Errorf("ParseLogLine() timestamp = %v, want %v", result.Timestamp, tt.expected.Timestamp)
			}
			if result.LogID != tt.expected.LogID {
				t.Errorf("ParseLogLine() logID = %v, want %v", result.LogID, tt.expected.LogID)
			}
			if result.Level != tt.expected.Level {
				t.Errorf("ParseLogLine() level = %v, want %v", result.Level, tt.expected.Level)
			}
			if result.Message != tt.expected.Message {
				t.Errorf("ParseLogLine() message = %v, want %v", result.Message, tt.expected.Message)
			}
			if result.File != tt.expected.File {
				t.Errorf("ParseLogLine() file = %v, want %v", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("ParseLogLine() line = %v, want %v", result.Line, tt.expected.Line)
			}
			if result.Function != tt.expected.Function {
				t.Errorf("ParseLogLine() function = %v, want %v", result.Function, tt.expected.Function)
			}
			if result.Method != tt.expected.Method {
				t.Errorf("ParseLogLine() method = %v, want %v", result.Method, tt.expected.Method)
			}
			if result.Path != tt.expected.Path {
				t.Errorf("ParseLogLine() path = %v, want %v", result.Path, tt.expected.Path)
			}
			if result.Status != tt.expected.Status {
				t.Errorf("ParseLogLine() status = %v, want %v", result.Status, tt.expected.Status)
			}
			if result.Duration != tt.expected.Duration {
				t.Errorf("ParseLogLine() duration = %v, want %v", result.Duration, tt.expected.Duration)
			}
			if result.IP != tt.expected.IP {
				t.Errorf("ParseLogLine() ip = %v, want %v", result.IP, tt.expected.IP)
			}
			if result.UserAgent != tt.expected.UserAgent {
				t.Errorf("ParseLogLine() userAgent = %v, want %v", result.UserAgent, tt.expected.UserAgent)
			}
			if result.Extra == nil {
				t.Error("ParseLogLine() extra should not be nil")
			}
		})
	}
}

func TestParseLogLine_MalformedJSONEntries(t *testing.T) {
	tests := []struct {
		name        string
		logLine     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "invalid JSON - missing closing brace",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test message"`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "invalid JSON - missing quotes",
			logLine:     `{timestamp:"2025-08-20T10:30:45Z",level:"INFO",message:"Test message"}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "invalid JSON - extra comma",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test message",}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "invalid JSON - malformed string",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test message with unescaped "quotes"}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "empty string",
			logLine:     ``,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "non-JSON string",
			logLine:     `This is not JSON at all`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "invalid JSON - wrong data type for line field",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z","level":"INFO","message":"Test","line":"not_a_number"}`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
		{
			name:        "partial JSON object",
			logLine:     `{"timestamp":"2025-08-20T10:30:45Z"`,
			expectError: true,
			errorMsg:    "failed to parse JSON log line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ParseLogLine(tt.logLine)

			if tt.expectError {
				if err == nil {
					t.Error("ParseLogLine() should return error")
				}
				if err != nil && err.Error() == "" {
					t.Error("ParseLogLine() should return meaningful error message")
				}
				if result.Timestamp != "" || result.Level != "" || result.Message != "" {
					t.Error("ParseLogLine() should return empty LogEntry on error")
				}
			} else {
				if err != nil {
					t.Errorf("ParseLogLine() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestParseLogLine_HTTPRequestDetailsExtraction(t *testing.T) {
	tests := []struct {
		name             string
		logLine          string
		expectedMethod   string
		expectedPath     string
		expectedStatus   int
		expectedDuration string
		expectedIP       string
		expectedUA       string
	}{
		{
			name:             "complete HTTP request details",
			logLine:          `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/users | Status: 200 | Duration: 45ms | IP: 192.168.1.100 | User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)"}`,
			expectedMethod:   "GET",
			expectedPath:     "/api/v1/users",
			expectedStatus:   200,
			expectedDuration: "45ms",
			expectedIP:       "192.168.1.100",
			expectedUA:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		{
			name:             "HTTP request with query parameters",
			logLine:          `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/posts?page=1&limit=10 | Status: 200 | Duration: 30ms"}`,
			expectedMethod:   "GET",
			expectedPath:     "/api/v1/posts?page=1&limit=10",
			expectedStatus:   200,
			expectedDuration: "30ms",
		},
		{
			name:             "POST request with creation status",
			logLine:          `{"timestamp":"2025-08-20T10:30:45Z","message":"POST /api/v1/posts | Status: 201 | Duration: 150ms | IP: 10.0.0.1"}`,
			expectedMethod:   "POST",
			expectedPath:     "/api/v1/posts",
			expectedStatus:   201,
			expectedDuration: "150ms",
			expectedIP:       "10.0.0.1",
		},
		{
			name:             "error status code",
			logLine:          `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/nonexistent | Status: 404 | Duration: 5ms"}`,
			expectedMethod:   "GET",
			expectedPath:     "/api/v1/nonexistent",
			expectedStatus:   404,
			expectedDuration: "5ms",
		},
		{
			name:             "server error status",
			logLine:          `{"timestamp":"2025-08-20T10:30:45Z","message":"POST /api/v1/process | Status: 500 | Duration: 1000ms"}`,
			expectedMethod:   "POST",
			expectedPath:     "/api/v1/process",
			expectedStatus:   500,
			expectedDuration: "1000ms",
		},
		{
			name:           "message without HTTP details",
			logLine:        `{"timestamp":"2025-08-20T10:30:45Z","message":"Database connection established"}`,
			expectedMethod: "",
			expectedPath:   "",
			expectedStatus: 0,
		},
		{
			name:           "partial HTTP details - only method and path",
			logLine:        `{"timestamp":"2025-08-20T10:30:45Z","message":"DELETE /api/v1/temp/cleanup"}`,
			expectedMethod: "DELETE",
			expectedPath:   "/api/v1/temp/cleanup",
			expectedStatus: 0,
		},
		{
			name:           "IPv6 address",
			logLine:        `{"timestamp":"2025-08-20T10:30:45Z","message":"GET /api/v1/test | Status: 200 | IP: 2001:db8::1"}`,
			expectedMethod: "GET",
			expectedPath:   "/api/v1/test",
			expectedStatus: 200,
			expectedIP:     "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ParseLogLine(tt.logLine)
			if err != nil {
				t.Errorf("ParseLogLine() error = %v", err)
				return
			}

			if result.Method != tt.expectedMethod {
				t.Errorf("ParseLogLine() method = %v, want %v", result.Method, tt.expectedMethod)
			}
			if result.Path != tt.expectedPath {
				t.Errorf("ParseLogLine() path = %v, want %v", result.Path, tt.expectedPath)
			}
			if result.Status != tt.expectedStatus {
				t.Errorf("ParseLogLine() status = %v, want %v", result.Status, tt.expectedStatus)
			}
			if result.Duration != tt.expectedDuration {
				t.Errorf("ParseLogLine() duration = %v, want %v", result.Duration, tt.expectedDuration)
			}
			if result.IP != tt.expectedIP {
				t.Errorf("ParseLogLine() ip = %v, want %v", result.IP, tt.expectedIP)
			}
			if result.UserAgent != tt.expectedUA {
				t.Errorf("ParseLogLine() userAgent = %v, want %v", result.UserAgent, tt.expectedUA)
			}
		})
	}
}

func TestParseLogLine_EdgeCases(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ParseLogLine(tt.logLine)
			if err != nil {
				t.Errorf("ParseLogLine() error = %v", err)
				return
			}

			if result.Timestamp != tt.expected.Timestamp {
				t.Errorf("ParseLogLine() timestamp = %v, want %v", result.Timestamp, tt.expected.Timestamp)
			}
			if result.LogID != tt.expected.LogID {
				t.Errorf("ParseLogLine() logID = %v, want %v", result.LogID, tt.expected.LogID)
			}
			if result.Level != tt.expected.Level {
				t.Errorf("ParseLogLine() level = %v, want %v", result.Level, tt.expected.Level)
			}
			if result.Message != tt.expected.Message {
				t.Errorf("ParseLogLine() message = %v, want %v", result.Message, tt.expected.Message)
			}
			if result.File != tt.expected.File {
				t.Errorf("ParseLogLine() file = %v, want %v", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("ParseLogLine() line = %v, want %v", result.Line, tt.expected.Line)
			}
			if result.Function != tt.expected.Function {
				t.Errorf("ParseLogLine() function = %v, want %v", result.Function, tt.expected.Function)
			}
			if result.Extra == nil {
				t.Error("ParseLogLine() extra should not be nil")
			}
		})
	}
}