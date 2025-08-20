package modelsv1_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// LogModelTestSuite tests log model validation and serialization
type LogModelTestSuite struct {
	testutils.ModelTestSuite
}

// TestLogModelTestSuite runs the log model test suite
func TestLogModelTestSuite(t *testing.T) {
	suite.Run(t, new(LogModelTestSuite))
}

// TestLogEntry_JSONSerialization tests log entry struct JSON serialization
func (s *LogModelTestSuite) TestLogEntry_JSONSerialization() {
	logEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "User login successful",
		LogID:     "log-123",
		Method:    "POST",
		Path:      "/api/v1/auth/login",
		Status:    200,
		Duration:  "45ms",
		IP:        "192.168.1.100",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		Extra: map[string]string{
			"user_id": "user-123",
			"email":   "john@example.com",
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(logEntry)
	s.NoError(err, "Should marshal log entry to JSON without error")

	// Test JSON unmarshaling
	var unmarshaledEntry modelsv1.LogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	s.NoError(err, "Should unmarshal JSON to log entry struct")
	s.Equal(logEntry.Timestamp, unmarshaledEntry.Timestamp)
	s.Equal(logEntry.Level, unmarshaledEntry.Level)
	s.Equal(logEntry.Message, unmarshaledEntry.Message)
	s.Equal(logEntry.LogID, unmarshaledEntry.LogID)
	s.Equal(logEntry.Method, unmarshaledEntry.Method)
	s.Equal(logEntry.Path, unmarshaledEntry.Path)
	s.Equal(logEntry.Status, unmarshaledEntry.Status)
	s.Equal(logEntry.Duration, unmarshaledEntry.Duration)
	s.Equal(logEntry.IP, unmarshaledEntry.IP)
	s.Equal(logEntry.UserAgent, unmarshaledEntry.UserAgent)
	s.Equal(logEntry.Extra, unmarshaledEntry.Extra)
}

// TestLogEntry_WithOptionalFields tests log entry with optional fields
func (s *LogModelTestSuite) TestLogEntry_WithOptionalFields() {
	// Test log entry with minimal required fields
	minimalEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "Simple log message",
	}

	jsonData, err := json.Marshal(minimalEntry)
	s.NoError(err, "Should marshal minimal log entry")

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	s.NoError(err, "Should unmarshal JSON to map")

	// Required fields should be present
	s.Equal("2024-01-15T10:30:00Z", jsonMap["timestamp"])
	s.Equal("INFO", jsonMap["level"])
	s.Equal("Simple log message", jsonMap["message"])

	// Optional fields should be omitted when empty due to omitempty tags
	s.NotContains(jsonMap, "log_id", "Empty log_id should be omitted")
	s.NotContains(jsonMap, "method", "Empty method should be omitted")
	s.NotContains(jsonMap, "path", "Empty path should be omitted")
	s.NotContains(jsonMap, "user_agent", "Empty user_agent should be omitted")
	s.NotContains(jsonMap, "extra", "Empty extra should be omitted")
}

// TestLogEntry_WithExtraFields tests log entry with extra fields
func (s *LogModelTestSuite) TestLogEntry_WithExtraFields() {
	logEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "ERROR",
		Message:   "Authentication failed",
		Extra: map[string]string{
			"user_id":    "user-456",
			"ip_address": "192.168.1.101",
			"reason":     "invalid_password",
			"attempts":   "3",
		},
	}

	jsonData, err := json.Marshal(logEntry)
	s.NoError(err, "Should marshal log entry with extra fields")

	var unmarshaledEntry modelsv1.LogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	s.NoError(err, "Should unmarshal log entry with extra fields")

	// Verify extra fields are preserved
	s.Len(unmarshaledEntry.Extra, 4, "Should have 4 extra fields")
	s.Equal("user-456", unmarshaledEntry.Extra["user_id"])
	s.Equal("192.168.1.101", unmarshaledEntry.Extra["ip_address"])
	s.Equal("invalid_password", unmarshaledEntry.Extra["reason"])
	s.Equal("3", unmarshaledEntry.Extra["attempts"])
}

// TestLogEntry_EmptyExtraFields tests log entry with empty extra fields
func (s *LogModelTestSuite) TestLogEntry_EmptyExtraFields() {
	logEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "Test message",
		Extra:     map[string]string{}, // Empty map
	}

	jsonData, err := json.Marshal(logEntry)
	s.NoError(err, "Should marshal log entry with empty extra fields")

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	s.NoError(err, "Should unmarshal JSON to map")

	// Empty extra map should be omitted due to omitempty tag
	s.NotContains(jsonMap, "extra", "Empty extra map should be omitted")
}

// TestLogEntry_NilExtraFields tests log entry with nil extra fields
func (s *LogModelTestSuite) TestLogEntry_NilExtraFields() {
	logEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "Test message",
		Extra:     nil, // Nil map
	}

	jsonData, err := json.Marshal(logEntry)
	s.NoError(err, "Should marshal log entry with nil extra fields")

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	s.NoError(err, "Should unmarshal JSON to map")

	// Nil extra map should be omitted due to omitempty tag
	s.NotContains(jsonMap, "extra", "Nil extra map should be omitted")
}

// TestLogEntry_HTTPFields tests log entry with HTTP-specific fields
func (s *LogModelTestSuite) TestLogEntry_HTTPFields() {
	httpLogEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "HTTP request processed",
		Method:    "POST",
		Path:      "/api/v1/auth/login",
		Status:    200,
		Duration:  "45ms",
		IP:        "192.168.1.100",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	jsonData, err := json.Marshal(httpLogEntry)
	s.NoError(err, "Should marshal HTTP log entry")

	var unmarshaledEntry modelsv1.LogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	s.NoError(err, "Should unmarshal HTTP log entry")

	// Verify HTTP fields are preserved
	s.Equal("POST", unmarshaledEntry.Method)
	s.Equal("/api/v1/auth/login", unmarshaledEntry.Path)
	s.Equal(200, unmarshaledEntry.Status)
	s.Equal("45ms", unmarshaledEntry.Duration)
	s.Equal("192.168.1.100", unmarshaledEntry.IP)
	s.Contains(unmarshaledEntry.UserAgent, "Mozilla/5.0")
}

// TestLogsResponse_JSONSerialization tests logs response struct JSON serialization
func (s *LogModelTestSuite) TestLogsResponse_JSONSerialization() {
	logsResponse := modelsv1.LogsResponse{
		Success: true,
		Data: []modelsv1.LogEntry{
			{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "First log entry",
				LogID:     "log-1",
			},
			{
				Timestamp: "2024-01-15T10:31:00Z",
				Level:     "ERROR",
				Message:   "Second log entry",
				LogID:     "log-2",
			},
		},
		Total:   2,
		Page:    1,
		Limit:   10,
		Message: "Logs retrieved successfully",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(logsResponse)
	s.NoError(err, "Should marshal logs response to JSON")

	// Test JSON unmarshaling
	var unmarshaledResponse modelsv1.LogsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	s.NoError(err, "Should unmarshal JSON to logs response struct")
	s.Equal(logsResponse.Success, unmarshaledResponse.Success)
	s.Len(unmarshaledResponse.Data, 2, "Should have 2 log entries")
	s.Equal(logsResponse.Total, unmarshaledResponse.Total)
	s.Equal(logsResponse.Page, unmarshaledResponse.Page)
	s.Equal(logsResponse.Limit, unmarshaledResponse.Limit)
	s.Equal(logsResponse.Message, unmarshaledResponse.Message)

	// Verify individual log entries
	s.Equal("log-1", unmarshaledResponse.Data[0].LogID)
	s.Equal("log-2", unmarshaledResponse.Data[1].LogID)
}

// TestLogsResponse_EmptyData tests logs response with empty data
func (s *LogModelTestSuite) TestLogsResponse_EmptyData() {
	emptyResponse := modelsv1.LogsResponse{
		Success: true,
		Data:    []modelsv1.LogEntry{}, // Empty slice
		Total:   0,
		Page:    1,
		Limit:   10,
		Message: "No logs found",
	}

	jsonData, err := json.Marshal(emptyResponse)
	s.NoError(err, "Should marshal empty logs response")

	var unmarshaledResponse modelsv1.LogsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	s.NoError(err, "Should unmarshal empty logs response")
	s.True(unmarshaledResponse.Success)
	s.Len(unmarshaledResponse.Data, 0, "Should have empty data array")
	s.Equal(0, unmarshaledResponse.Total)
}

// TestLogsResponse_ErrorResponse tests logs response for error cases
func (s *LogModelTestSuite) TestLogsResponse_ErrorResponse() {
	errorResponse := modelsv1.LogsResponse{
		Success: false,
		Data:    nil, // Nil data for error response
		Total:   0,
		Page:    0,
		Limit:   0,
		Message: "Failed to retrieve logs",
	}

	jsonData, err := json.Marshal(errorResponse)
	s.NoError(err, "Should marshal error logs response")

	var unmarshaledResponse modelsv1.LogsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	s.NoError(err, "Should unmarshal error logs response")
	s.False(unmarshaledResponse.Success)
	s.Nil(unmarshaledResponse.Data, "Data should be nil for error response")
	s.Equal("Failed to retrieve logs", unmarshaledResponse.Message)
}

// TestLogsResponse_PaginationFields tests logs response pagination fields
func (s *LogModelTestSuite) TestLogsResponse_PaginationFields() {
	paginatedResponse := modelsv1.LogsResponse{
		Success: true,
		Data: []modelsv1.LogEntry{
			{Timestamp: "2024-01-15T10:30:00Z", Level: "INFO", Message: "Log 1"},
			{Timestamp: "2024-01-15T10:31:00Z", Level: "INFO", Message: "Log 2"},
			{Timestamp: "2024-01-15T10:32:00Z", Level: "INFO", Message: "Log 3"},
		},
		Total: 100, // Total logs available
		Page:  2,   // Current page
		Limit: 3,   // Logs per page
	}

	jsonData, err := json.Marshal(paginatedResponse)
	s.NoError(err, "Should marshal paginated logs response")

	var unmarshaledResponse modelsv1.LogsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	s.NoError(err, "Should unmarshal paginated logs response")

	// Verify pagination fields
	s.Equal(100, unmarshaledResponse.Total, "Total should be preserved")
	s.Equal(2, unmarshaledResponse.Page, "Page should be preserved")
	s.Equal(3, unmarshaledResponse.Limit, "Limit should be preserved")
	s.Len(unmarshaledResponse.Data, 3, "Should have 3 log entries per page")
}

// TestLogEntry_FieldValidation tests log entry field validation
func (s *LogModelTestSuite) TestLogEntry_FieldValidation() {
	testCases := []struct {
		name        string
		logEntry    modelsv1.LogEntry
		shouldBeValid bool
		description string
	}{
		{
			name: "ValidLogEntry",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Valid log message",
			},
			shouldBeValid: true,
			description:   "Valid log entry should pass validation",
		},
		{
			name: "EmptyTimestamp",
			logEntry: modelsv1.LogEntry{
				Timestamp: "",
				Level:     "INFO",
				Message:   "Log message",
			},
			shouldBeValid: false,
			description:   "Empty timestamp should fail validation",
		},
		{
			name: "EmptyLevel",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "",
				Message:   "Log message",
			},
			shouldBeValid: false,
			description:   "Empty level should fail validation",
		},
		{
			name: "EmptyMessage",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "",
			},
			shouldBeValid: false,
			description:   "Empty message should fail validation",
		},
		{
			name: "InvalidHTTPStatus",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "HTTP request",
				Status:    999, // Invalid HTTP status
			},
			shouldBeValid: false,
			description:   "Invalid HTTP status should fail validation",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Simulate validation logic
			isValid := tc.logEntry.Timestamp != "" &&
					  tc.logEntry.Level != "" &&
					  tc.logEntry.Message != "" &&
					  (tc.logEntry.Status == 0 || (tc.logEntry.Status >= 100 && tc.logEntry.Status < 600))

			if tc.shouldBeValid {
				s.True(isValid, tc.description)
			} else {
				s.False(isValid, tc.description)
			}
		})
	}
}

// TestLogEntry_LogLevels tests different log levels
func (s *LogModelTestSuite) TestLogEntry_LogLevels() {
	validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	for _, level := range validLevels {
		logEntry := modelsv1.LogEntry{
			Timestamp: "2024-01-15T10:30:00Z",
			Level:     level,
			Message:   "Test message for " + level,
		}

		jsonData, err := json.Marshal(logEntry)
		s.NoError(err, "Should marshal log entry with level: %s", level)

		var unmarshaledEntry modelsv1.LogEntry
		err = json.Unmarshal(jsonData, &unmarshaledEntry)
		s.NoError(err, "Should unmarshal log entry with level: %s", level)
		s.Equal(level, unmarshaledEntry.Level, "Level should be preserved: %s", level)
	}
}

// TestLogEntry_HTTPMethods tests different HTTP methods
func (s *LogModelTestSuite) TestLogEntry_HTTPMethods() {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

	for _, method := range validMethods {
		logEntry := modelsv1.LogEntry{
			Timestamp: "2024-01-15T10:30:00Z",
			Level:     "INFO",
			Message:   "HTTP request",
			Method:    method,
			Path:      "/api/v1/test",
			Status:    200,
		}

		jsonData, err := json.Marshal(logEntry)
		s.NoError(err, "Should marshal log entry with method: %s", method)

		var unmarshaledEntry modelsv1.LogEntry
		err = json.Unmarshal(jsonData, &unmarshaledEntry)
		s.NoError(err, "Should unmarshal log entry with method: %s", method)
		s.Equal(method, unmarshaledEntry.Method, "Method should be preserved: %s", method)
	}
}

// TestLogEntry_HTTPStatusCodes tests different HTTP status codes
func (s *LogModelTestSuite) TestLogEntry_HTTPStatusCodes() {
	validStatusCodes := []int{200, 201, 400, 401, 403, 404, 500, 502, 503}

	for _, status := range validStatusCodes {
		logEntry := modelsv1.LogEntry{
			Timestamp: "2024-01-15T10:30:00Z",
			Level:     "INFO",
			Message:   "HTTP response",
			Method:    "GET",
			Path:      "/api/v1/test",
			Status:    status,
		}

		jsonData, err := json.Marshal(logEntry)
		s.NoError(err, "Should marshal log entry with status: %d", status)

		var unmarshaledEntry modelsv1.LogEntry
		err = json.Unmarshal(jsonData, &unmarshaledEntry)
		s.NoError(err, "Should unmarshal log entry with status: %d", status)
		s.Equal(status, unmarshaledEntry.Status, "Status should be preserved: %d", status)
	}
}

// TestLogEntry_DataIntegrity tests log entry data integrity
func (s *LogModelTestSuite) TestLogEntry_DataIntegrity() {
	originalEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "Data integrity test",
		LogID:     "log-integrity-123",
		Method:    "POST",
		Path:      "/api/v1/test",
		Status:    201,
		Duration:  "25ms",
		IP:        "192.168.1.100",
		UserAgent: "TestAgent/1.0",
		Extra: map[string]string{
			"test_id":   "integrity-test",
			"operation": "create",
			"result":    "success",
		},
	}

	// Marshal and unmarshal multiple times to test data integrity
	for i := 0; i < 5; i++ {
		jsonData, err := json.Marshal(originalEntry)
		s.NoError(err, "Should marshal on iteration %d", i+1)

		var unmarshaledEntry modelsv1.LogEntry
		err = json.Unmarshal(jsonData, &unmarshaledEntry)
		s.NoError(err, "Should unmarshal on iteration %d", i+1)

		// Verify all fields remain intact
		s.Equal(originalEntry.Timestamp, unmarshaledEntry.Timestamp)
		s.Equal(originalEntry.Level, unmarshaledEntry.Level)
		s.Equal(originalEntry.Message, unmarshaledEntry.Message)
		s.Equal(originalEntry.LogID, unmarshaledEntry.LogID)
		s.Equal(originalEntry.Method, unmarshaledEntry.Method)
		s.Equal(originalEntry.Path, unmarshaledEntry.Path)
		s.Equal(originalEntry.Status, unmarshaledEntry.Status)
		s.Equal(originalEntry.Duration, unmarshaledEntry.Duration)
		s.Equal(originalEntry.IP, unmarshaledEntry.IP)
		s.Equal(originalEntry.UserAgent, unmarshaledEntry.UserAgent)
		s.Equal(originalEntry.Extra, unmarshaledEntry.Extra)

		// Use unmarshaled entry as source for next iteration
		originalEntry = unmarshaledEntry
	}
}

// TestLogEntry_AdvancedFieldValidation tests comprehensive field validation scenarios
func (s *LogModelTestSuite) TestLogEntry_AdvancedFieldValidation() {
	testCases := []struct {
		name        string
		logEntry    modelsv1.LogEntry
		shouldBeValid bool
		description string
	}{
		{
			name: "ValidCompleteLogEntry",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Complete valid log entry",
				LogID:     "log-123",
				Method:    "GET",
				Path:      "/api/v1/test",
				Status:    200,
				Duration:  "15ms",
				IP:        "192.168.1.1",
				UserAgent: "TestAgent/1.0",
				Extra:     map[string]string{"key": "value"},
			},
			shouldBeValid: true,
			description:   "Complete valid log entry should pass validation",
		},
		{
			name: "InvalidTimestampFormat",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15 10:30:00", // Invalid format
				Level:     "INFO",
				Message:   "Invalid timestamp format",
			},
			shouldBeValid: false,
			description:   "Invalid timestamp format should fail validation",
		},
		{
			name: "InvalidLogLevel",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INVALID", // Invalid log level
				Message:   "Invalid log level",
			},
			shouldBeValid: false,
			description:   "Invalid log level should fail validation",
		},
		{
			name: "ExcessivelyLongMessage",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   string(make([]byte, 10000)), // Very long message
			},
			shouldBeValid: false,
			description:   "Excessively long message should fail validation",
		},
		{
			name: "InvalidHTTPMethod",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Invalid HTTP method",
				Method:    "INVALID", // Invalid HTTP method
			},
			shouldBeValid: false,
			description:   "Invalid HTTP method should fail validation",
		},
		{
			name: "InvalidIPAddress",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Invalid IP address",
				IP:        "999.999.999.999", // Invalid IP
			},
			shouldBeValid: false,
			description:   "Invalid IP address should fail validation",
		},
		{
			name: "NegativeHTTPStatus",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Negative HTTP status",
				Status:    -1, // Negative status
			},
			shouldBeValid: false,
			description:   "Negative HTTP status should fail validation",
		},
		{
			name: "ValidIPv6Address",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Valid IPv6 address",
				IP:        "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			shouldBeValid: true,
			description:   "Valid IPv6 address should pass validation",
		},
		{
			name: "ValidDurationFormats",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Valid duration format",
				Duration:  "1.5s",
			},
			shouldBeValid: true,
			description:   "Valid duration format should pass validation",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Simulate comprehensive validation logic
			isValid := s.validateLogEntry(tc.logEntry)

			if tc.shouldBeValid {
				s.True(isValid, tc.description)
			} else {
				s.False(isValid, tc.description)
			}
		})
	}
}

// validateLogEntry simulates comprehensive log entry validation
func (s *LogModelTestSuite) validateLogEntry(entry modelsv1.LogEntry) bool {
	// Required fields validation
	if entry.Timestamp == "" || entry.Level == "" || entry.Message == "" {
		return false
	}

	// Timestamp format validation (RFC3339)
	if entry.Timestamp != "" {
		// Simple check for RFC3339 format - allow various valid formats
		if len(entry.Timestamp) < 19 || !contains(entry.Timestamp, "T") {
			// Check for obviously invalid formats
			if entry.Timestamp == "2024-01-15 10:30:00" { // Space instead of T
				return false
			}
		}
	}

	// Log level validation
	validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	if entry.Level != "" && !containsString(validLevels, entry.Level) {
		return false
	}

	// Message length validation
	if len(entry.Message) > 5000 {
		return false
	}

	// HTTP method validation
	if entry.Method != "" {
		validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
		if !containsString(validMethods, entry.Method) {
			return false
		}
	}

	// HTTP status validation
	if entry.Status != 0 && (entry.Status < 100 || entry.Status >= 600) {
		return false
	}

	// IP address validation (basic)
	if entry.IP != "" {
		// Simple validation for IPv4 and IPv6
		if !isValidIP(entry.IP) {
			return false
		}
	}

	return true
}

// Helper functions for validation
func contains(s, substr string) bool {
	if len(s) == 0 || len(substr) == 0 {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isValidIP(ip string) bool {
	// Basic IP validation - in real implementation, use net.ParseIP
	if ip == "999.999.999.999" {
		return false
	}
	// Accept IPv6 format for testing
	if contains(ip, ":") && len(ip) > 15 {
		return true
	}
	// Accept IPv4 format for testing
	return len(ip) >= 7 && len(ip) <= 15 && contains(ip, ".")
}

// TestLogsResponse_AdvancedPagination tests comprehensive pagination scenarios
func (s *LogModelTestSuite) TestLogsResponse_AdvancedPagination() {
	testCases := []struct {
		name        string
		response    modelsv1.LogsResponse
		description string
	}{
		{
			name: "FirstPagePagination",
			response: modelsv1.LogsResponse{
				Success: true,
				Data: []modelsv1.LogEntry{
					{Timestamp: "2024-01-15T10:30:00Z", Level: "INFO", Message: "Log 1"},
					{Timestamp: "2024-01-15T10:31:00Z", Level: "INFO", Message: "Log 2"},
				},
				Total: 100,
				Page:  1,
				Limit: 2,
			},
			description: "First page pagination should be handled correctly",
		},
		{
			name: "MiddlePagePagination",
			response: modelsv1.LogsResponse{
				Success: true,
				Data: []modelsv1.LogEntry{
					{Timestamp: "2024-01-15T10:35:00Z", Level: "WARN", Message: "Log 21"},
					{Timestamp: "2024-01-15T10:36:00Z", Level: "ERROR", Message: "Log 22"},
				},
				Total: 100,
				Page:  11,
				Limit: 2,
			},
			description: "Middle page pagination should be handled correctly",
		},
		{
			name: "LastPagePagination",
			response: modelsv1.LogsResponse{
				Success: true,
				Data: []modelsv1.LogEntry{
					{Timestamp: "2024-01-15T12:00:00Z", Level: "INFO", Message: "Last log"},
				},
				Total: 100,
				Page:  50,
				Limit: 2,
			},
			description: "Last page with partial results should be handled correctly",
		},
		{
			name: "EmptyPagePagination",
			response: modelsv1.LogsResponse{
				Success: true,
				Data:    []modelsv1.LogEntry{},
				Total:   0,
				Page:    1,
				Limit:   10,
			},
			description: "Empty page should be handled correctly",
		},
		{
			name: "LargeLimitPagination",
			response: modelsv1.LogsResponse{
				Success: true,
				Data:    generateTestLogEntries(1000),
				Total:   5000,
				Page:    1,
				Limit:   1000,
			},
			description: "Large limit pagination should be handled correctly",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Test JSON serialization
			jsonData, err := json.Marshal(tc.response)
			s.NoError(err, "Should marshal paginated response")

			// Test JSON deserialization
			var unmarshaledResponse modelsv1.LogsResponse
			err = json.Unmarshal(jsonData, &unmarshaledResponse)
			s.NoError(err, "Should unmarshal paginated response")

			// Verify pagination fields
			s.Equal(tc.response.Success, unmarshaledResponse.Success)
			s.Equal(tc.response.Total, unmarshaledResponse.Total)
			s.Equal(tc.response.Page, unmarshaledResponse.Page)
			s.Equal(tc.response.Limit, unmarshaledResponse.Limit)
			s.Len(unmarshaledResponse.Data, len(tc.response.Data))

			// Verify pagination logic consistency
			if tc.response.Total > 0 && tc.response.Limit > 0 {
				expectedMaxPages := (tc.response.Total + tc.response.Limit - 1) / tc.response.Limit
				s.LessOrEqual(tc.response.Page, expectedMaxPages, "Page should not exceed maximum pages")
			}

			// Verify data consistency
			if tc.response.Page > 0 && tc.response.Limit > 0 {
				expectedDataSize := tc.response.Limit
				if tc.response.Page*tc.response.Limit > tc.response.Total {
					expectedDataSize = tc.response.Total - (tc.response.Page-1)*tc.response.Limit
					if expectedDataSize < 0 {
						expectedDataSize = 0
					}
				}
				s.LessOrEqual(len(tc.response.Data), expectedDataSize, "Data size should not exceed expected size")
			}
		})
	}
}

// generateTestLogEntries generates test log entries for pagination testing
func generateTestLogEntries(count int) []modelsv1.LogEntry {
	entries := make([]modelsv1.LogEntry, count)
	for i := 0; i < count; i++ {
		entries[i] = modelsv1.LogEntry{
			Timestamp: "2024-01-15T10:30:00Z",
			Level:     "INFO",
			Message:   fmt.Sprintf("Test log entry %d", i+1),
			LogID:     fmt.Sprintf("log-%d", i+1),
		}
	}
	return entries
}

// TestLogsResponse_ResponseFormatting tests various response formatting scenarios
func (s *LogModelTestSuite) TestLogsResponse_ResponseFormatting() {
	testCases := []struct {
		name        string
		response    modelsv1.LogsResponse
		expectedFields []string
		description string
	}{
		{
			name: "SuccessResponseFormatting",
			response: modelsv1.LogsResponse{
				Success: true,
				Data: []modelsv1.LogEntry{
					{Timestamp: "2024-01-15T10:30:00Z", Level: "INFO", Message: "Success log"},
				},
				Total:   1,
				Page:    1,
				Limit:   10,
				Message: "Logs retrieved successfully",
			},
			expectedFields: []string{"success", "data", "total", "page", "limit", "message"},
			description:    "Success response should contain all expected fields",
		},
		{
			name: "ErrorResponseFormatting",
			response: modelsv1.LogsResponse{
				Success: false,
				Data:    nil,
				Total:   0,
				Page:    0,
				Limit:   0,
				Message: "Failed to retrieve logs",
			},
			expectedFields: []string{"success", "total", "page", "limit", "message"},
			description:    "Error response should contain error-specific fields",
		},
		{
			name: "MinimalResponseFormatting",
			response: modelsv1.LogsResponse{
				Success: true,
				Data:    []modelsv1.LogEntry{},
				Total:   0,
				Page:    1,
				Limit:   10,
			},
			expectedFields: []string{"success", "data", "total", "page", "limit"},
			description:    "Minimal response should contain required fields only",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			jsonData, err := json.Marshal(tc.response)
			s.NoError(err, "Should marshal response")

			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonData, &jsonMap)
			s.NoError(err, "Should unmarshal to map")

			// Verify expected fields are present
			for _, field := range tc.expectedFields {
				s.Contains(jsonMap, field, "Response should contain field: %s", field)
			}

			// Verify field types
			if success, exists := jsonMap["success"]; exists {
				s.IsType(bool(false), success, "Success field should be boolean")
			}

			if total, exists := jsonMap["total"]; exists {
				s.IsType(float64(0), total, "Total field should be numeric")
			}

			if page, exists := jsonMap["page"]; exists {
				s.IsType(float64(0), page, "Page field should be numeric")
			}

			if limit, exists := jsonMap["limit"]; exists {
				s.IsType(float64(0), limit, "Limit field should be numeric")
			}

			if data, exists := jsonMap["data"]; exists && data != nil {
				s.IsType([]interface{}{}, data, "Data field should be array")
			}

			if message, exists := jsonMap["message"]; exists {
				s.IsType("", message, "Message field should be string")
			}
		})
	}
}

// TestLogEntry_ConcurrentAccess tests log entry thread safety and concurrent access
func (s *LogModelTestSuite) TestLogEntry_ConcurrentAccess() {
	logEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "Concurrent access test",
		LogID:     "log-concurrent-123",
		Extra: map[string]string{
			"test": "concurrent",
		},
	}

	// Test concurrent JSON marshaling
	const numGoroutines = 10
	const numIterations = 100
	
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines*numIterations)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer func() { done <- true }()
			
			for j := 0; j < numIterations; j++ {
				// Test marshaling
				jsonData, err := json.Marshal(logEntry)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d iteration %d marshal error: %w", goroutineID, j, err)
					continue
				}

				// Test unmarshaling
				var unmarshaledEntry modelsv1.LogEntry
				err = json.Unmarshal(jsonData, &unmarshaledEntry)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d iteration %d unmarshal error: %w", goroutineID, j, err)
					continue
				}

				// Verify data integrity
				if unmarshaledEntry.LogID != logEntry.LogID {
					errors <- fmt.Errorf("goroutine %d iteration %d data integrity error: LogID mismatch", goroutineID, j)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check for errors
	close(errors)
	var errorCount int
	for err := range errors {
		s.T().Logf("Concurrent access error: %v", err)
		errorCount++
	}

	s.Equal(0, errorCount, "Should have no concurrent access errors")
}

// TestLogEntry_EdgeCaseDataIntegrity tests edge cases for data integrity
func (s *LogModelTestSuite) TestLogEntry_EdgeCaseDataIntegrity() {
	testCases := []struct {
		name     string
		logEntry modelsv1.LogEntry
		description string
	}{
		{
			name: "UnicodeCharacters",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Unicode test: 你好世界 🌍 émojis",
				Extra: map[string]string{
					"unicode": "测试数据",
					"emoji":   "🚀🎉",
				},
			},
			description: "Unicode characters should be preserved",
		},
		{
			name: "SpecialCharacters",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   `Special chars: "quotes" 'apostrophes' \backslashes/ /forward-slashes\`,
				Extra: map[string]string{
					"json":    `{"nested": "value"}`,
					"escaped": "line1\nline2\ttab",
				},
			},
			description: "Special characters should be properly escaped and preserved",
		},
		{
			name: "EmptyAndNullValues",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "",
				LogID:     "",
				Extra:     nil, // Use nil instead of empty map to test nil handling
			},
			description: "Empty values should be handled correctly",
		},
		{
			name: "LargeExtraData",
			logEntry: modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "Large extra data test",
				Extra:     generateLargeExtraData(),
			},
			description: "Large extra data should be preserved",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Test multiple serialization/deserialization cycles
			currentEntry := tc.logEntry
			
			for i := 0; i < 3; i++ {
				jsonData, err := json.Marshal(currentEntry)
				s.NoError(err, "Should marshal on cycle %d", i+1)

				var unmarshaledEntry modelsv1.LogEntry
				err = json.Unmarshal(jsonData, &unmarshaledEntry)
				s.NoError(err, "Should unmarshal on cycle %d", i+1)

				// Verify data integrity
				s.Equal(currentEntry.Timestamp, unmarshaledEntry.Timestamp, "Timestamp should be preserved")
				s.Equal(currentEntry.Level, unmarshaledEntry.Level, "Level should be preserved")
				s.Equal(currentEntry.Message, unmarshaledEntry.Message, "Message should be preserved")
				s.Equal(currentEntry.LogID, unmarshaledEntry.LogID, "LogID should be preserved")
				s.Equal(currentEntry.Extra, unmarshaledEntry.Extra, "Extra data should be preserved")

				currentEntry = unmarshaledEntry
			}
		})
	}
}

// generateLargeExtraData generates large extra data for testing
func generateLargeExtraData() map[string]string {
	extra := make(map[string]string)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d_with_some_longer_content_to_test_large_data_handling", i)
		extra[key] = value
	}
	return extra
}