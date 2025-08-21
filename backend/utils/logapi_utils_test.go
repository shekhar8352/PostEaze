package utils

import (
	"testing"
)

func TestParseLogLine_JSONFormat(t *testing.T) {
	// Test with actual JSON log line format
	jsonLine := `{"timestamp":"2025-08-20T10:18:48+05:30","log_id":"45e1adc9-698f-41f9-886d-706effa57585","level":"INFO","message":"Started GET /api/health | IP: ::1 | User-Agent: PostmanRuntime/7.45.0","file":"log_middleware.go","line":26,"function":"func2"}`

	entry, err := ParseLogLine(jsonLine)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify all JSON fields are properly mapped
	if entry.Timestamp != "2025-08-20T10:18:48+05:30" {
		t.Errorf("Expected timestamp '2025-08-20T10:18:48+05:30', got '%s'", entry.Timestamp)
	}

	if entry.LogID != "45e1adc9-698f-41f9-886d-706effa57585" {
		t.Errorf("Expected log_id '45e1adc9-698f-41f9-886d-706effa57585', got '%s'", entry.LogID)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level 'INFO', got '%s'", entry.Level)
	}

	if entry.File != "log_middleware.go" {
		t.Errorf("Expected file 'log_middleware.go', got '%s'", entry.File)
	}

	if entry.Line != 26 {
		t.Errorf("Expected line 26, got %d", entry.Line)
	}

	if entry.Function != "func2" {
		t.Errorf("Expected function 'func2', got '%s'", entry.Function)
	}

	// Verify HTTP request details are extracted from message
	if entry.Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", entry.Method)
	}

	if entry.Path != "/api/health" {
		t.Errorf("Expected path '/api/health', got '%s'", entry.Path)
	}

	if entry.IP != "::1" {
		t.Errorf("Expected IP '::1', got '%s'", entry.IP)
	}

	if entry.UserAgent != "PostmanRuntime/7.45.0" {
		t.Errorf("Expected UserAgent 'PostmanRuntime/7.45.0', got '%s'", entry.UserAgent)
	}
}

func TestParseLogLine_WithStatusAndDuration(t *testing.T) {
	// Test with log line containing status and duration
	jsonLine := `{"timestamp":"2025-08-20T10:18:48+05:30","log_id":"45e1adc9-698f-41f9-886d-706effa57585","level":"INFO","message":"Completed GET /api/health | Status: 200 | Duration: 507.1µs | LogID: 45e1adc9-698f-41f9-886d-706effa57585","file":"log_middleware.go","line":36,"function":"func2"}`

	entry, err := ParseLogLine(jsonLine)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify HTTP request details are extracted
	if entry.Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", entry.Method)
	}

	if entry.Path != "/api/health" {
		t.Errorf("Expected path '/api/health', got '%s'", entry.Path)
	}

	if entry.Status != 200 {
		t.Errorf("Expected status 200, got %d", entry.Status)
	}

	if entry.Duration != "507.1µs" {
		t.Errorf("Expected duration '507.1µs', got '%s'", entry.Duration)
	}
}

func TestParseLogLine_MalformedJSON(t *testing.T) {
	// Test with malformed JSON
	malformedLine := `{"timestamp":"2025-08-20T10:18:48+05:30","log_id":"test","level":"INFO","message":"test message"`

	_, err := ParseLogLine(malformedLine)
	if err == nil {
		t.Error("Expected error for malformed JSON, got nil")
	}
}

func TestParseLogLine_EmptyFields(t *testing.T) {
	// Test with empty optional fields
	jsonLine := `{"timestamp":"2025-08-20T10:18:48+05:30","log_id":"","level":"ERROR","message":"Test error message","file":"","line":0,"function":""}`

	entry, err := ParseLogLine(jsonLine)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify required fields are set
	if entry.Timestamp != "2025-08-20T10:18:48+05:30" {
		t.Errorf("Expected timestamp '2025-08-20T10:18:48+05:30', got '%s'", entry.Timestamp)
	}

	if entry.Level != "ERROR" {
		t.Errorf("Expected level 'ERROR', got '%s'", entry.Level)
	}

	if entry.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got '%s'", entry.Message)
	}

	// Verify empty optional fields are handled correctly
	if entry.LogID != "" {
		t.Errorf("Expected empty log_id, got '%s'", entry.LogID)
	}

	if entry.File != "" {
		t.Errorf("Expected empty file, got '%s'", entry.File)
	}

	if entry.Line != 0 {
		t.Errorf("Expected line 0, got %d", entry.Line)
	}

	if entry.Function != "" {
		t.Errorf("Expected empty function, got '%s'", entry.Function)
	}
}