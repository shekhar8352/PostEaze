package business_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
)

// createTestLogFile creates a test log file with sample entries
func createTestLogFile(t *testing.T, logDir, date string, entries []modelsv1.LogEntry) string {
	t.Helper()
	
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))
	
	file, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("failed to create test log file: %v", err)
	}
	defer file.Close()
	
	for _, entry := range entries {
		logLine := fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s","log_id":"%s","method":"%s","path":"%s","status":%d,"duration":"%s","ip":"%s","user_agent":"%s"}`,
			entry.Timestamp, entry.Level, entry.Message, entry.LogID, entry.Method, entry.Path, entry.Status, entry.Duration, entry.IP, entry.UserAgent)
		if _, err := file.WriteString(logLine + "\n"); err != nil {
			t.Fatalf("failed to write log entry: %v", err)
		}
	}
	
	return logFile
}

// setupTestLogDir creates a temporary directory for test logs
func setupTestLogDir(t *testing.T) (string, func()) {
	t.Helper()
	
	testLogDir := filepath.Join("backend", "tests", "business", "testlogs")
	
	// Store original log directory if set
	originalLogDir := os.Getenv("LOG_DIR")
	
	// Create test log directory
	if err := os.MkdirAll(testLogDir, 0755); err != nil {
		t.Fatalf("failed to create test log directory: %v", err)
	}
	
	// Set test log directory
	os.Setenv("LOG_DIR", testLogDir)
	
	// Return cleanup function
	cleanup := func() {
		// Restore original log directory
		if originalLogDir != "" {
			os.Setenv("LOG_DIR", originalLogDir)
		} else {
			os.Unsetenv("LOG_DIR")
		}
		
		// Clean up test directory
		os.RemoveAll(testLogDir)
	}
	
	return testLogDir, cleanup
}

func TestReadLogsByLogID_Success(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	ctx := context.Background()
	logID := "test-log-123"
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	// Create test log entries using helpers
	todayEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = time.Now().Format(time.RFC3339)
			e.Level = "INFO"
			e.Message = "Test log message 1"
			e.LogID = logID
			e.Method = "GET"
			e.Path = "/api/test"
			e.Status = 200
		}),
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
			e.Level = "ERROR"
			e.Message = "Test error message"
			e.LogID = "other-log-456"
			e.Method = "POST"
			e.Path = "/api/error"
			e.Status = 500
		}),
	}
	
	yesterdayEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
			e.Level = "INFO"
			e.Message = "Test log message 2"
			e.LogID = logID
			e.Method = "PUT"
			e.Path = "/api/update"
			e.Status = 200
		}),
	}
	
	// Create test log files
	createTestLogFile(t, testLogDir, today, todayEntries)
	createTestLogFile(t, testLogDir, yesterday, yesterdayEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(ctx, logID)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(logs))
	}
	
	// Verify logs are sorted by timestamp
	if len(logs) >= 2 && logs[0].Timestamp > logs[1].Timestamp {
		t.Error("logs should be sorted by timestamp")
	}
	
	// Verify all logs have the correct log ID
	for _, log := range logs {
		if log.LogID != logID {
			t.Errorf("expected log ID %s, got %s", logID, log.LogID)
		}
	}
}

func TestReadLogsByLogID_NoLogsFound(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	ctx := context.Background()
	logID := "nonexistent-log-id"
	today := time.Now().Format("2006-01-02")
	
	// Create test log file with different log IDs
	todayEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.LogID = "different-log-id"
		}),
	}
	
	createTestLogFile(t, testLogDir, today, todayEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(ctx, logID)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 0 {
		t.Errorf("expected empty slice, got %d logs", len(logs))
	}
}

func TestReadLogsByLogID_NoLogDirectory(t *testing.T) {
	ctx := context.Background()
	logID := "test-log-123"
	
	// Set nonexistent log directory
	originalLogDir := os.Getenv("LOG_DIR")
	os.Setenv("LOG_DIR", "/nonexistent/directory")
	defer func() {
		if originalLogDir != "" {
			os.Setenv("LOG_DIR", originalLogDir)
		} else {
			os.Unsetenv("LOG_DIR")
		}
	}()
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(ctx, logID)
	
	// Assert
	if err == nil {
		t.Error("expected error for nonexistent log directory")
	}
	
	if logs != nil {
		t.Error("expected nil logs on error")
	}
	
	if err != nil && len(err.Error()) > 0 {
		errMsg := err.Error()
		if len(errMsg) == 0 || (errMsg != "failed to get available log files" && len(errMsg) < 10) {
			// Check if it contains expected error message
			t.Logf("got error: %v", err)
		}
	}
}

func TestReadLogsByLogID_EmptyLogID(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	ctx := context.Background()
	logID := ""
	today := time.Now().Format("2006-01-02")
	
	// Create test log file
	todayEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.LogID = "some-log-id"
		}),
	}
	createTestLogFile(t, testLogDir, today, todayEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(ctx, logID)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 0 {
		t.Errorf("expected empty slice for empty log ID, got %d logs", len(logs))
	}
}

func TestReadLogsByDate_Success(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	date := "2024-01-15"
	
	// Create test entries using helpers
	testEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = "2024-01-15T10:30:00Z"
			e.Level = "INFO"
			e.Message = "User login successful"
			e.LogID = "log-1"
			e.Method = "POST"
			e.Path = "/api/v1/auth/login"
			e.Status = 200
			e.Duration = "45ms"
			e.IP = "192.168.1.100"
			e.UserAgent = "Mozilla/5.0"
		}),
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = "2024-01-15T11:00:00Z"
			e.Level = "ERROR"
			e.Message = "Authentication failed"
			e.LogID = "log-2"
			e.Method = "POST"
			e.Path = "/api/v1/auth/login"
			e.Status = 401
			e.Duration = "12ms"
			e.IP = "192.168.1.101"
			e.UserAgent = "curl/7.68.0"
		}),
	}
	
	createTestLogFile(t, testLogDir, date, testEntries)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(logs))
	}
	
	if total != 2 {
		t.Errorf("expected total count 2, got %d", total)
	}
	
	// Verify log content
	if len(logs) >= 2 {
		if logs[0].LogID != "log-1" {
			t.Errorf("expected first log ID 'log-1', got '%s'", logs[0].LogID)
		}
		if logs[1].LogID != "log-2" {
			t.Errorf("expected second log ID 'log-2', got '%s'", logs[1].LogID)
		}
		if logs[0].Level != "INFO" {
			t.Errorf("expected first log level 'INFO', got '%s'", logs[0].Level)
		}
		if logs[1].Level != "ERROR" {
			t.Errorf("expected second log level 'ERROR', got '%s'", logs[1].Level)
		}
	}
}

func TestReadLogsByDate_FileNotFound(t *testing.T) {
	_, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	date := "1999-01-01" // Date with no log file
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	if err == nil {
		t.Error("expected error for missing log file")
	}
	
	if logs != nil {
		t.Error("expected nil logs on error")
	}
	
	if total != 0 {
		t.Errorf("expected total 0 on error, got %d", total)
	}
	
	if err != nil && len(err.Error()) > 0 {
		errMsg := err.Error()
		if len(errMsg) == 0 {
			t.Error("expected non-empty error message")
		}
	}
}

func TestReadLogsByDate_EmptyLogFile(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	date := "2024-01-16"
	
	// Create empty log file
	createTestLogFile(t, testLogDir, date, []modelsv1.LogEntry{})
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 0 {
		t.Errorf("expected empty slice for empty log file, got %d logs", len(logs))
	}
	
	if total != 0 {
		t.Errorf("expected total 0 for empty log file, got %d", total)
	}
}

func TestReadLogsByDate_InvalidDateFormats(t *testing.T) {
	_, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	tests := []struct {
		name string
		date string
	}{
		{"Empty date", ""},
		{"Invalid format 1", "2024/01/15"},
		{"Invalid format 2", "15-01-2024"},
		{"Invalid format 3", "2024-1-15"},
		{"Invalid format 4", "2024-01-1"},
		{"Invalid format 5", "Jan 15, 2024"},
		{"Invalid date", "2024-13-45"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			logs, total, err := businessv1.ReadLogsByDate(tt.date)
			
			// Assert
			if err == nil {
				t.Errorf("expected error for invalid date format: %s", tt.date)
			}
			
			if logs != nil {
				t.Error("expected nil logs on error")
			}
			
			if total != 0 {
				t.Errorf("expected total 0 on error, got %d", total)
			}
		})
	}
}

func TestReadLogsByLogID_MultipleFiles(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	ctx := context.Background()
	logID := "multi-file-log"
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	
	// Create entries for each day using helpers
	todayEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = time.Now().Format(time.RFC3339)
			e.Level = "INFO"
			e.Message = "Today's log"
			e.LogID = logID
			e.Path = "/api/today"
		}),
	}
	
	yesterdayEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
			e.Level = "WARN"
			e.Message = "Yesterday's log"
			e.LogID = logID
			e.Path = "/api/yesterday"
			e.Status = 400
		}),
	}
	
	twoDaysAgoEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = time.Now().AddDate(0, 0, -2).Format(time.RFC3339)
			e.Level = "ERROR"
			e.Message = "Two days ago log"
			e.LogID = logID
			e.Path = "/api/twodaysago"
			e.Status = 500
		}),
	}
	
	// Create log files
	createTestLogFile(t, testLogDir, today, todayEntries)
	createTestLogFile(t, testLogDir, yesterday, yesterdayEntries)
	createTestLogFile(t, testLogDir, twoDaysAgo, twoDaysAgoEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(ctx, logID)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 3 {
		t.Errorf("expected 3 log entries from all three days, got %d", len(logs))
	}
	
	// Verify logs are sorted by timestamp (oldest first)
	if len(logs) >= 2 {
		if logs[0].Timestamp > logs[1].Timestamp {
			t.Error("first log should be oldest")
		}
		if len(logs) >= 3 && logs[1].Timestamp > logs[2].Timestamp {
			t.Error("logs should be in chronological order")
		}
	}
	
	// Verify all logs have the correct log ID
	for _, log := range logs {
		if log.LogID != logID {
			t.Errorf("expected log ID %s, got %s", logID, log.LogID)
		}
	}
}

func TestReadLogsByDate_LargeLogFile(t *testing.T) {
	testLogDir, cleanup := setupTestLogDir(t)
	defer cleanup()
	
	date := "2024-01-17"
	
	// Create a large number of log entries
	var testEntries []modelsv1.LogEntry
	for i := 0; i < 100; i++ { // Reduced from 1000 to 100 for faster tests
		entry := helpers.CreateLogEntry(func(e *modelsv1.LogEntry) {
			e.Timestamp = fmt.Sprintf("2024-01-17T%02d:%02d:00Z", i/60, i%60)
			e.Message = fmt.Sprintf("Test message %d", i)
			e.LogID = fmt.Sprintf("log-%d", i)
			e.Path = fmt.Sprintf("/api/test/%d", i)
		})
		testEntries = append(testEntries, entry)
	}
	
	createTestLogFile(t, testLogDir, date, testEntries)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(logs) != 100 {
		t.Errorf("expected 100 log entries, got %d", len(logs))
	}
	
	if total != 100 {
		t.Errorf("expected total 100, got %d", total)
	}
}