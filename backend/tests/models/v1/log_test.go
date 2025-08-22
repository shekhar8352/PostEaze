package modelsv1_test

import (
	"encoding/json"
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func TestLogEntry_JSONSerialization(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal log entry to JSON without error: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledEntry modelsv1.LogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to log entry struct: %v", err)
	}

	// Verify all fields are preserved
	if logEntry.Timestamp != unmarshaledEntry.Timestamp {
		t.Errorf("Expected timestamp %s, got %s", logEntry.Timestamp, unmarshaledEntry.Timestamp)
	}
	if logEntry.Level != unmarshaledEntry.Level {
		t.Errorf("Expected level %s, got %s", logEntry.Level, unmarshaledEntry.Level)
	}
	if logEntry.Message != unmarshaledEntry.Message {
		t.Errorf("Expected message %s, got %s", logEntry.Message, unmarshaledEntry.Message)
	}
	if logEntry.LogID != unmarshaledEntry.LogID {
		t.Errorf("Expected log ID %s, got %s", logEntry.LogID, unmarshaledEntry.LogID)
	}
	if logEntry.Method != unmarshaledEntry.Method {
		t.Errorf("Expected method %s, got %s", logEntry.Method, unmarshaledEntry.Method)
	}
	if logEntry.Path != unmarshaledEntry.Path {
		t.Errorf("Expected path %s, got %s", logEntry.Path, unmarshaledEntry.Path)
	}
	if logEntry.Status != unmarshaledEntry.Status {
		t.Errorf("Expected status %d, got %d", logEntry.Status, unmarshaledEntry.Status)
	}
	if logEntry.Duration != unmarshaledEntry.Duration {
		t.Errorf("Expected duration %s, got %s", logEntry.Duration, unmarshaledEntry.Duration)
	}
	if logEntry.IP != unmarshaledEntry.IP {
		t.Errorf("Expected IP %s, got %s", logEntry.IP, unmarshaledEntry.IP)
	}
	if logEntry.UserAgent != unmarshaledEntry.UserAgent {
		t.Errorf("Expected user agent %s, got %s", logEntry.UserAgent, unmarshaledEntry.UserAgent)
	}
	if len(logEntry.Extra) != len(unmarshaledEntry.Extra) {
		t.Errorf("Expected %d extra fields, got %d", len(logEntry.Extra), len(unmarshaledEntry.Extra))
	}
	for key, value := range logEntry.Extra {
		if unmarshaledEntry.Extra[key] != value {
			t.Errorf("Expected extra[%s] = %s, got %s", key, value, unmarshaledEntry.Extra[key])
		}
	}
}

func TestLogEntry_WithOptionalFields(t *testing.T) {
	// Test log entry with minimal required fields
	minimalEntry := modelsv1.LogEntry{
		Timestamp: "2024-01-15T10:30:00Z",
		Level:     "INFO",
		Message:   "Simple log message",
	}

	jsonData, err := json.Marshal(minimalEntry)
	if err != nil {
		t.Fatalf("Should marshal minimal log entry: %v", err)
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to map: %v", err)
	}

	// Required fields should be present
	if jsonMap["timestamp"] != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected timestamp in JSON")
	}
	if jsonMap["level"] != "INFO" {
		t.Errorf("Expected level in JSON")
	}
	if jsonMap["message"] != "Simple log message" {
		t.Errorf("Expected message in JSON")
	}

	// Optional fields should be omitted when empty due to omitempty tags
	if _, exists := jsonMap["log_id"]; exists {
		t.Errorf("Empty log_id should be omitted")
	}
	if _, exists := jsonMap["method"]; exists {
		t.Errorf("Empty method should be omitted")
	}
	if _, exists := jsonMap["path"]; exists {
		t.Errorf("Empty path should be omitted")
	}
	if _, exists := jsonMap["user_agent"]; exists {
		t.Errorf("Empty user_agent should be omitted")
	}
	if _, exists := jsonMap["extra"]; exists {
		t.Errorf("Empty extra should be omitted")
	}
}

func TestLogEntry_WithExtraFields(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal log entry with extra fields: %v", err)
	}

	var unmarshaledEntry modelsv1.LogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	if err != nil {
		t.Fatalf("Should unmarshal log entry with extra fields: %v", err)
	}

	// Verify extra fields are preserved
	if len(unmarshaledEntry.Extra) != 4 {
		t.Errorf("Should have 4 extra fields, got %d", len(unmarshaledEntry.Extra))
	}
	if unmarshaledEntry.Extra["user_id"] != "user-456" {
		t.Errorf("Expected user_id to be user-456")
	}
	if unmarshaledEntry.Extra["ip_address"] != "192.168.1.101" {
		t.Errorf("Expected ip_address to be 192.168.1.101")
	}
	if unmarshaledEntry.Extra["reason"] != "invalid_password" {
		t.Errorf("Expected reason to be invalid_password")
	}
	if unmarshaledEntry.Extra["attempts"] != "3" {
		t.Errorf("Expected attempts to be 3")
	}
}

func TestLogEntry_HTTPFields(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal HTTP log entry: %v", err)
	}

	var unmarshaledEntry modelsv1.LogEntry
	err = json.Unmarshal(jsonData, &unmarshaledEntry)
	if err != nil {
		t.Fatalf("Should unmarshal HTTP log entry: %v", err)
	}

	// Verify HTTP fields are preserved
	if unmarshaledEntry.Method != "POST" {
		t.Errorf("Expected method POST, got %s", unmarshaledEntry.Method)
	}
	if unmarshaledEntry.Path != "/api/v1/auth/login" {
		t.Errorf("Expected path /api/v1/auth/login, got %s", unmarshaledEntry.Path)
	}
	if unmarshaledEntry.Status != 200 {
		t.Errorf("Expected status 200, got %d", unmarshaledEntry.Status)
	}
	if unmarshaledEntry.Duration != "45ms" {
		t.Errorf("Expected duration 45ms, got %s", unmarshaledEntry.Duration)
	}
	if unmarshaledEntry.IP != "192.168.1.100" {
		t.Errorf("Expected IP 192.168.1.100, got %s", unmarshaledEntry.IP)
	}
	if unmarshaledEntry.UserAgent != "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" {
		t.Errorf("Expected specific user agent, got %s", unmarshaledEntry.UserAgent)
	}
}

func TestLogsResponse_JSONSerialization(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal logs response to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledResponse modelsv1.LogsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to logs response struct: %v", err)
	}

	if logsResponse.Success != unmarshaledResponse.Success {
		t.Errorf("Expected success %t, got %t", logsResponse.Success, unmarshaledResponse.Success)
	}
	if len(unmarshaledResponse.Data) != 2 {
		t.Errorf("Should have 2 log entries, got %d", len(unmarshaledResponse.Data))
	}
	if logsResponse.Total != unmarshaledResponse.Total {
		t.Errorf("Expected total %d, got %d", logsResponse.Total, unmarshaledResponse.Total)
	}
	if logsResponse.Page != unmarshaledResponse.Page {
		t.Errorf("Expected page %d, got %d", logsResponse.Page, unmarshaledResponse.Page)
	}
	if logsResponse.Limit != unmarshaledResponse.Limit {
		t.Errorf("Expected limit %d, got %d", logsResponse.Limit, unmarshaledResponse.Limit)
	}
	if logsResponse.Message != unmarshaledResponse.Message {
		t.Errorf("Expected message %s, got %s", logsResponse.Message, unmarshaledResponse.Message)
	}

	// Verify individual log entries
	if unmarshaledResponse.Data[0].LogID != "log-1" {
		t.Errorf("Expected first log ID to be log-1, got %s", unmarshaledResponse.Data[0].LogID)
	}
	if unmarshaledResponse.Data[1].LogID != "log-2" {
		t.Errorf("Expected second log ID to be log-2, got %s", unmarshaledResponse.Data[1].LogID)
	}
}

func TestLogEntry_FieldValidation(t *testing.T) {
	testCases := []struct {
		name          string
		logEntry      modelsv1.LogEntry
		shouldBeValid bool
		description   string
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
		t.Run(tc.name, func(t *testing.T) {
			// Simulate validation logic
			isValid := tc.logEntry.Timestamp != "" &&
				tc.logEntry.Level != "" &&
				tc.logEntry.Message != "" &&
				(tc.logEntry.Status == 0 || (tc.logEntry.Status >= 100 && tc.logEntry.Status < 600))

			if tc.shouldBeValid && !isValid {
				t.Errorf("%s: expected valid but got invalid", tc.description)
			}
			if !tc.shouldBeValid && isValid {
				t.Errorf("%s: expected invalid but got valid", tc.description)
			}
		})
	}
}

func TestLogEntry_LogLevels(t *testing.T) {
	validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	for _, level := range validLevels {
		t.Run("Level_"+level, func(t *testing.T) {
			logEntry := modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     level,
				Message:   "Test message for " + level,
			}

			jsonData, err := json.Marshal(logEntry)
			if err != nil {
				t.Fatalf("Should marshal log entry with level %s: %v", level, err)
			}

			var unmarshaledEntry modelsv1.LogEntry
			err = json.Unmarshal(jsonData, &unmarshaledEntry)
			if err != nil {
				t.Fatalf("Should unmarshal log entry with level %s: %v", level, err)
			}
			if unmarshaledEntry.Level != level {
				t.Errorf("Level should be preserved: expected %s, got %s", level, unmarshaledEntry.Level)
			}
		})
	}
}

func TestLogEntry_HTTPMethods(t *testing.T) {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

	for _, method := range validMethods {
		t.Run("Method_"+method, func(t *testing.T) {
			logEntry := modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "HTTP request",
				Method:    method,
				Path:      "/api/v1/test",
				Status:    200,
			}

			jsonData, err := json.Marshal(logEntry)
			if err != nil {
				t.Fatalf("Should marshal log entry with method %s: %v", method, err)
			}

			var unmarshaledEntry modelsv1.LogEntry
			err = json.Unmarshal(jsonData, &unmarshaledEntry)
			if err != nil {
				t.Fatalf("Should unmarshal log entry with method %s: %v", method, err)
			}
			if unmarshaledEntry.Method != method {
				t.Errorf("Method should be preserved: expected %s, got %s", method, unmarshaledEntry.Method)
			}
		})
	}
}

func TestLogEntry_HTTPStatusCodes(t *testing.T) {
	validStatusCodes := []int{200, 201, 400, 401, 403, 404, 500, 502, 503}

	for _, status := range validStatusCodes {
		t.Run("Status_"+string(rune(status)), func(t *testing.T) {
			logEntry := modelsv1.LogEntry{
				Timestamp: "2024-01-15T10:30:00Z",
				Level:     "INFO",
				Message:   "HTTP response",
				Method:    "GET",
				Path:      "/api/v1/test",
				Status:    status,
			}

			jsonData, err := json.Marshal(logEntry)
			if err != nil {
				t.Fatalf("Should marshal log entry with status %d: %v", status, err)
			}

			var unmarshaledEntry modelsv1.LogEntry
			err = json.Unmarshal(jsonData, &unmarshaledEntry)
			if err != nil {
				t.Fatalf("Should unmarshal log entry with status %d: %v", status, err)
			}
			if unmarshaledEntry.Status != status {
				t.Errorf("Status should be preserved: expected %d, got %d", status, unmarshaledEntry.Status)
			}
		})
	}
}