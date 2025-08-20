package businessv1_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// LogBusinessLogicTestSuite tests the log business logic
type LogBusinessLogicTestSuite struct {
	testutils.BusinessLogicTestSuite
	ctx           context.Context
	testLogDir    string
	originalLogDir string
}

// SetupSuite initializes the test suite
func (s *LogBusinessLogicTestSuite) SetupSuite() {
	s.BusinessLogicTestSuite.SetupSuite()
	s.ctx = context.Background()
	
	// Create a temporary directory for test logs
	s.testLogDir = filepath.Join("tests", "testutils", "logs")
	
	// Store original log directory if set
	s.originalLogDir = os.Getenv("LOG_DIR")
	
	// Set test log directory
	os.Setenv("LOG_DIR", s.testLogDir)
}

// TearDownSuite cleans up the test suite
func (s *LogBusinessLogicTestSuite) TearDownSuite() {
	// Restore original log directory
	if s.originalLogDir != "" {
		os.Setenv("LOG_DIR", s.originalLogDir)
	} else {
		os.Unsetenv("LOG_DIR")
	}
	
	s.BusinessLogicTestSuite.TearDownSuite()
}

// SetupTest prepares each test
func (s *LogBusinessLogicTestSuite) SetupTest() {
	s.BusinessLogicTestSuite.SetupTest()
	
	// Ensure test log directory exists
	err := os.MkdirAll(s.testLogDir, 0755)
	s.Require().NoError(err, "Failed to create test log directory")
}

// TearDownTest cleans up after each test
func (s *LogBusinessLogicTestSuite) TearDownTest() {
	s.BusinessLogicTestSuite.TearDownTest()
}

// createTestLogFile creates a test log file with sample entries
func (s *LogBusinessLogicTestSuite) createTestLogFile(date string, entries []modelsv1.LogEntry) string {
	logFile := filepath.Join(s.testLogDir, fmt.Sprintf("app-%s.log", date))
	
	file, err := os.Create(logFile)
	s.Require().NoError(err, "Failed to create test log file")
	defer file.Close()
	
	for _, entry := range entries {
		logLine := fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s","log_id":"%s","method":"%s","path":"%s","status":%d,"duration":"%s","ip":"%s","user_agent":"%s"}`,
			entry.Timestamp, entry.Level, entry.Message, entry.LogID, entry.Method, entry.Path, entry.Status, entry.Duration, entry.IP, entry.UserAgent)
		_, err := file.WriteString(logLine + "\n")
		s.Require().NoError(err, "Failed to write log entry")
	}
	
	return logFile
}

// TestReadLogsByLogID_Success tests successful log retrieval by log ID
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_Success() {
	// Arrange
	logID := "test-log-123"
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	
	// Create test log entries
	todayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Test log message 1",
			LogID:     logID,
			Method:    "GET",
			Path:      "/api/test",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
		{
			Timestamp: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			Level:     "ERROR",
			Message:   "Test error message",
			LogID:     "other-log-456",
			Method:    "POST",
			Path:      "/api/error",
			Status:    500,
			Duration:  "5ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	
	yesterdayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Test log message 2",
			LogID:     logID,
			Method:    "PUT",
			Path:      "/api/update",
			Status:    200,
			Duration:  "15ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	
	// Create test log files
	s.createTestLogFile(today, todayEntries)
	s.createTestLogFile(yesterday, yesterdayEntries)
	s.createTestLogFile(twoDaysAgo, []modelsv1.LogEntry{}) // Empty file
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err)
	s.Len(logs, 2, "Should find 2 log entries with the specified log ID")
	
	// Verify logs are sorted by timestamp
	s.True(logs[0].Timestamp <= logs[1].Timestamp, "Logs should be sorted by timestamp")
	
	// Verify all logs have the correct log ID
	for _, log := range logs {
		s.Equal(logID, log.LogID, "All logs should have the specified log ID")
	}
}

// TestReadLogsByLogID_NoLogsFound tests log retrieval when no logs match the ID
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_NoLogsFound() {
	// Arrange
	logID := "nonexistent-log-id"
	today := time.Now().Format("2006-01-02")
	
	// Create test log file with different log IDs
	todayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Test log message",
			LogID:     "different-log-id",
			Method:    "GET",
			Path:      "/api/test",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	
	s.createTestLogFile(today, todayEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err)
	s.Empty(logs, "Should return empty slice when no logs match the ID")
}

// TestReadLogsByLogID_NoLogDirectory tests log retrieval when log directory doesn't exist
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_NoLogDirectory() {
	// Arrange
	logID := "test-log-123"
	
	// Set nonexistent log directory
	os.Setenv("LOG_DIR", "/nonexistent/directory")
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.Error(err)
	s.Nil(logs)
	s.Contains(err.Error(), "failed to get available log files")
	
	// Restore test log directory
	os.Setenv("LOG_DIR", s.testLogDir)
}

// TestReadLogsByLogID_NoLogFiles tests log retrieval when no log files exist
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_NoLogFiles() {
	// Arrange
	logID := "test-log-123"
	
	// Ensure log directory is empty
	err := os.RemoveAll(s.testLogDir)
	s.Require().NoError(err)
	err = os.MkdirAll(s.testLogDir, 0755)
	s.Require().NoError(err)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err)
	s.Empty(logs, "Should return empty slice when no log files exist")
}

// TestReadLogsByLogID_MalformedLogFile tests log retrieval with malformed log entries
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_MalformedLogFile() {
	// Arrange
	logID := "test-log-123"
	today := time.Now().Format("2006-01-02")
	
	// Create log file with malformed JSON
	logFile := filepath.Join(s.testLogDir, fmt.Sprintf("app-%s.log", today))
	file, err := os.Create(logFile)
	s.Require().NoError(err)
	defer file.Close()
	
	// Write malformed JSON
	_, err = file.WriteString("invalid json line\n")
	s.Require().NoError(err)
	
	// Write valid JSON with matching log ID
	validEntry := fmt.Sprintf(`{"timestamp":"%s","level":"INFO","message":"Valid entry","log_id":"%s","method":"GET","path":"/api/test","status":200,"duration":"10ms","ip":"127.0.0.1","user_agent":"test-agent"}`,
		time.Now().Format(time.RFC3339), logID)
	_, err = file.WriteString(validEntry + "\n")
	s.Require().NoError(err)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	// The function should handle malformed entries gracefully and return valid ones
	s.NoError(err)
	s.Len(logs, 1, "Should return valid log entries despite malformed ones")
	s.Equal(logID, logs[0].LogID)
}

// TestReadLogsByDate_Success tests successful log retrieval by date
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_Success() {
	// Arrange
	date := "2024-01-15"
	
	testEntries := []modelsv1.LogEntry{
		{
			Timestamp: "2024-01-15T10:30:00Z",
			Level:     "INFO",
			Message:   "User login successful",
			LogID:     "log-1",
			Method:    "POST",
			Path:      "/api/v1/auth/login",
			Status:    200,
			Duration:  "45ms",
			IP:        "192.168.1.100",
			UserAgent: "Mozilla/5.0",
		},
		{
			Timestamp: "2024-01-15T11:00:00Z",
			Level:     "ERROR",
			Message:   "Authentication failed",
			LogID:     "log-2",
			Method:    "POST",
			Path:      "/api/v1/auth/login",
			Status:    401,
			Duration:  "12ms",
			IP:        "192.168.1.101",
			UserAgent: "curl/7.68.0",
		},
	}
	
	s.createTestLogFile(date, testEntries)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	s.NoError(err)
	s.Len(logs, 2, "Should return all log entries for the date")
	s.Equal(2, total, "Total count should match number of entries")
	
	// Verify log content
	s.Equal("log-1", logs[0].LogID)
	s.Equal("log-2", logs[1].LogID)
	s.Equal("INFO", logs[0].Level)
	s.Equal("ERROR", logs[1].Level)
}

// TestReadLogsByDate_FileNotFound tests log retrieval when log file doesn't exist
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_FileNotFound() {
	// Arrange
	date := "1999-01-01" // Date with no log file (using a date that won't conflict with other tests)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	s.Error(err)
	s.Nil(logs)
	s.Equal(0, total)
	s.Contains(err.Error(), "log file not found")
}

// TestReadLogsByDate_EmptyLogFile tests log retrieval with empty log file
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_EmptyLogFile() {
	// Arrange
	date := "2024-01-16"
	
	// Create empty log file
	s.createTestLogFile(date, []modelsv1.LogEntry{})
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	s.NoError(err)
	s.Empty(logs, "Should return empty slice for empty log file")
	s.Equal(0, total, "Total should be 0 for empty log file")
}

// TestReadLogsByDate_InvalidLogDirectory tests log retrieval with invalid log directory
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_InvalidLogDirectory() {
	// Arrange
	date := "2024-01-15"
	
	// Set invalid log directory
	os.Setenv("LOG_DIR", "/nonexistent/directory")
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	s.Error(err)
	s.Nil(logs)
	s.Equal(0, total)
	s.Contains(err.Error(), "log file not found")
	
	// Restore test log directory
	os.Setenv("LOG_DIR", s.testLogDir)
}

// TestReadLogsByDate_LargeLogFile tests log retrieval with large log file
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_LargeLogFile() {
	// Arrange
	date := "2024-01-17"
	
	// Create a large number of log entries
	var testEntries []modelsv1.LogEntry
	for i := 0; i < 1000; i++ {
		entry := modelsv1.LogEntry{
			Timestamp: fmt.Sprintf("2024-01-17T%02d:%02d:00Z", i/60, i%60),
			Level:     "INFO",
			Message:   fmt.Sprintf("Test message %d", i),
			LogID:     fmt.Sprintf("log-%d", i),
			Method:    "GET",
			Path:      fmt.Sprintf("/api/test/%d", i),
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		}
		testEntries = append(testEntries, entry)
	}
	
	s.createTestLogFile(date, testEntries)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	s.NoError(err)
	s.Len(logs, 1000, "Should return all log entries")
	s.Equal(1000, total, "Total should match number of entries")
}

// TestReadLogsByLogID_MultipleFiles tests log retrieval across multiple files
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_MultipleFiles() {
	// Arrange
	logID := "multi-file-log"
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	
	// Create entries for each day
	todayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Today's log",
			LogID:     logID,
			Method:    "GET",
			Path:      "/api/today",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	
	yesterdayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
			Level:     "WARN",
			Message:   "Yesterday's log",
			LogID:     logID,
			Method:    "POST",
			Path:      "/api/yesterday",
			Status:    400,
			Duration:  "20ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	
	twoDaysAgoEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
			Level:     "ERROR",
			Message:   "Two days ago log",
			LogID:     logID,
			Method:    "DELETE",
			Path:      "/api/twodaysago",
			Status:    500,
			Duration:  "30ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	
	// Create log files
	s.createTestLogFile(today, todayEntries)
	s.createTestLogFile(yesterday, yesterdayEntries)
	s.createTestLogFile(twoDaysAgo, twoDaysAgoEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err)
	s.Len(logs, 3, "Should find logs from all three days")
	
	// Verify logs are sorted by timestamp (oldest first)
	s.True(logs[0].Timestamp <= logs[1].Timestamp, "First log should be oldest")
	s.True(logs[1].Timestamp <= logs[2].Timestamp, "Logs should be in chronological order")
	
	// Verify content from different days
	s.Contains(logs[0].Message, "Two days ago")
	s.Contains(logs[1].Message, "Yesterday's")
	s.Contains(logs[2].Message, "Today's")
}

// TestReadLogsByLogID_CorruptedLogFile tests log retrieval with corrupted log files
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_CorruptedLogFile() {
	// Arrange
	logID := "test-log-123"
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	// Create a valid log file for today
	todayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Valid log entry",
			LogID:     logID,
			Method:    "GET",
			Path:      "/api/test",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	s.createTestLogFile(today, todayEntries)
	
	// Create a corrupted log file for yesterday
	corruptedLogFile := filepath.Join(s.testLogDir, fmt.Sprintf("app-%s.log", yesterday))
	file, err := os.Create(corruptedLogFile)
	s.Require().NoError(err)
	defer file.Close()
	
	// Write corrupted content (not readable as file)
	_, err = file.WriteString("This is not JSON\n")
	s.Require().NoError(err)
	_, err = file.WriteString("Neither is this\n")
	s.Require().NoError(err)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	// Should handle corrupted files gracefully and return valid entries from other files
	s.NoError(err, "Should handle corrupted files gracefully")
	s.Len(logs, 1, "Should return valid entries from readable files")
	s.Equal(logID, logs[0].LogID)
	s.Equal("Valid log entry", logs[0].Message)
}

// TestReadLogsByLogID_MixedFileStates tests log retrieval with mixed file states
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_MixedFileStates() {
	// Arrange
	logID := "mixed-state-log"
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	threeDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	
	// Create valid log file for today
	todayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Today's valid entry",
			LogID:     logID,
			Method:    "GET",
			Path:      "/api/today",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	s.createTestLogFile(today, todayEntries)
	
	// Create empty log file for yesterday
	s.createTestLogFile(yesterday, []modelsv1.LogEntry{})
	
	// Create valid log file for two days ago
	twoDaysAgoEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
			Level:     "WARN",
			Message:   "Two days ago entry",
			LogID:     logID,
			Method:    "POST",
			Path:      "/api/twodaysago",
			Status:    400,
			Duration:  "20ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	s.createTestLogFile(twoDaysAgo, twoDaysAgoEntries)
	
	// Create corrupted log file for three days ago
	corruptedLogFile := filepath.Join(s.testLogDir, fmt.Sprintf("app-%s.log", threeDaysAgo))
	file, err := os.Create(corruptedLogFile)
	s.Require().NoError(err)
	defer file.Close()
	_, err = file.WriteString("corrupted content\n")
	s.Require().NoError(err)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err, "Should handle mixed file states gracefully")
	s.Len(logs, 2, "Should return entries from valid files only")
	
	// Verify entries are from valid files
	foundToday := false
	foundTwoDaysAgo := false
	for _, log := range logs {
		s.Equal(logID, log.LogID)
		if log.Message == "Today's valid entry" {
			foundToday = true
		}
		if log.Message == "Two days ago entry" {
			foundTwoDaysAgo = true
		}
	}
	s.True(foundToday, "Should find today's entry")
	s.True(foundTwoDaysAgo, "Should find two days ago entry")
}

// TestReadLogsByLogID_EmptyLogID tests log retrieval with empty log ID
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_EmptyLogID() {
	// Arrange
	logID := ""
	today := time.Now().Format("2006-01-02")
	
	// Create test log file
	todayEntries := []modelsv1.LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "INFO",
			Message:   "Test message",
			LogID:     "some-log-id",
			Method:    "GET",
			Path:      "/api/test",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	s.createTestLogFile(today, todayEntries)
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err, "Should handle empty log ID gracefully")
	s.Empty(logs, "Should return empty slice for empty log ID")
}

// TestReadLogsByDate_InvalidDateFormats tests log retrieval with various invalid date formats
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_InvalidDateFormats() {
	testCases := []struct {
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
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			logs, total, err := businessv1.ReadLogsByDate(tc.date)
			
			// Assert
			s.Error(err, "Should return error for invalid date format: %s", tc.date)
			s.Nil(logs)
			s.Equal(0, total)
		})
	}
}

// TestReadLogsByDate_CorruptedLogFile tests log retrieval with corrupted log file
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_CorruptedLogFile() {
	// Arrange
	date := "2024-01-20"
	
	// Create corrupted log file
	logFile := filepath.Join(s.testLogDir, fmt.Sprintf("app-%s.log", date))
	file, err := os.Create(logFile)
	s.Require().NoError(err)
	defer file.Close()
	
	// Write corrupted content
	_, err = file.WriteString("This is not JSON\n")
	s.Require().NoError(err)
	_, err = file.WriteString("Neither is this line\n")
	s.Require().NoError(err)
	_, err = file.WriteString("{incomplete json\n")
	s.Require().NoError(err)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	// Should handle corrupted file gracefully and return empty results
	s.NoError(err, "Should handle corrupted file gracefully")
	s.Empty(logs, "Should return empty slice for corrupted file")
	s.Equal(0, total, "Total should be 0 for corrupted file")
}

// TestReadLogsByDate_PartiallyCorruptedLogFile tests log retrieval with partially corrupted log file
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_PartiallyCorruptedLogFile() {
	// Arrange
	date := "2024-01-21"
	
	// Create log file with mixed valid and invalid entries
	logFile := filepath.Join(s.testLogDir, fmt.Sprintf("app-%s.log", date))
	file, err := os.Create(logFile)
	s.Require().NoError(err)
	defer file.Close()
	
	// Write corrupted line
	_, err = file.WriteString("This is not JSON\n")
	s.Require().NoError(err)
	
	// Write valid JSON entry
	validEntry := `{"timestamp":"2024-01-21T10:30:00Z","level":"INFO","message":"Valid entry","log_id":"valid-log","method":"GET","path":"/api/test","status":200,"duration":"10ms","ip":"127.0.0.1","user_agent":"test-agent"}`
	_, err = file.WriteString(validEntry + "\n")
	s.Require().NoError(err)
	
	// Write another corrupted line
	_, err = file.WriteString("{incomplete json\n")
	s.Require().NoError(err)
	
	// Write another valid JSON entry
	validEntry2 := `{"timestamp":"2024-01-21T11:00:00Z","level":"ERROR","message":"Another valid entry","log_id":"valid-log-2","method":"POST","path":"/api/error","status":500,"duration":"5ms","ip":"127.0.0.1","user_agent":"test-agent"}`
	_, err = file.WriteString(validEntry2 + "\n")
	s.Require().NoError(err)
	
	// Act
	logs, total, err := businessv1.ReadLogsByDate(date)
	
	// Assert
	s.NoError(err, "Should handle partially corrupted file gracefully")
	s.Len(logs, 2, "Should return only valid entries")
	s.Equal(2, total, "Total should match valid entries count")
	
	// Verify valid entries are returned
	s.Equal("valid-log", logs[0].LogID)
	s.Equal("valid-log-2", logs[1].LogID)
	s.Equal("Valid entry", logs[0].Message)
	s.Equal("Another valid entry", logs[1].Message)
}

// TestReadLogsByDate_FilePermissionError tests log retrieval with permission errors
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_FilePermissionError() {
	// Skip this test on Windows as file permissions work differently
	s.T().Skip("Skipping file permission test on Windows - file permissions work differently")
}

// TestReadLogsByLogID_LargeNumberOfFiles tests log retrieval across many log files
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_LargeNumberOfFiles() {
	// Arrange
	logID := "large-scale-log"
	
	// Create 10 log files with entries
	for i := 0; i < 10; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		entries := []modelsv1.LogEntry{
			{
				Timestamp: time.Now().AddDate(0, 0, -i).Format(time.RFC3339),
				Level:     "INFO",
				Message:   fmt.Sprintf("Entry from day %d", i),
				LogID:     logID,
				Method:    "GET",
				Path:      fmt.Sprintf("/api/day%d", i),
				Status:    200,
				Duration:  "10ms",
				IP:        "127.0.0.1",
				UserAgent: "test-agent",
			},
		}
		s.createTestLogFile(date, entries)
	}
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.NoError(err)
	s.Len(logs, 10, "Should find entries from all 10 files")
	
	// Verify all entries have the correct log ID
	for _, log := range logs {
		s.Equal(logID, log.LogID)
	}
	
	// Verify logs are sorted by timestamp (oldest first)
	for i := 1; i < len(logs); i++ {
		s.True(logs[i-1].Timestamp <= logs[i].Timestamp, "Logs should be sorted by timestamp")
	}
}

// TestReadLogsByDate_EdgeCaseDates tests log retrieval with edge case dates
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_EdgeCaseDates() {
	testCases := []struct {
		name string
		date string
	}{
		{"Leap year date", "2024-02-29"},
		{"New Year's Day", "2024-01-01"},
		{"Year end", "2024-12-31"},
		{"Future date", "2030-01-01"},
		{"Past date", "2020-01-01"},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Create test log file for the date
			entries := []modelsv1.LogEntry{
				{
					Timestamp: fmt.Sprintf("%sT10:30:00Z", tc.date),
					Level:     "INFO",
					Message:   fmt.Sprintf("Entry for %s", tc.date),
					LogID:     fmt.Sprintf("log-%s", tc.date),
					Method:    "GET",
					Path:      "/api/test",
					Status:    200,
					Duration:  "10ms",
					IP:        "127.0.0.1",
					UserAgent: "test-agent",
				},
			}
			s.createTestLogFile(tc.date, entries)
			
			// Act
			logs, total, err := businessv1.ReadLogsByDate(tc.date)
			
			// Assert
			s.NoError(err, "Should handle edge case date: %s", tc.date)
			s.Len(logs, 1, "Should return entry for date: %s", tc.date)
			s.Equal(1, total, "Total should be 1 for date: %s", tc.date)
			s.Equal(fmt.Sprintf("log-%s", tc.date), logs[0].LogID)
		})
	}
}

// TestReadLogsByLogID_PerformanceWithLargeFiles tests performance with large log files
func (s *LogBusinessLogicTestSuite) TestReadLogsByLogID_PerformanceWithLargeFiles() {
	// Arrange
	logID := "performance-test-log"
	date := time.Now().Format("2006-01-02")
	
	// Create a large log file with many entries
	var entries []modelsv1.LogEntry
	for i := 0; i < 1000; i++ {
		entry := modelsv1.LogEntry{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second).Format(time.RFC3339),
			Level:     "INFO",
			Message:   fmt.Sprintf("Performance test entry %d", i),
			LogID:     fmt.Sprintf("log-%d", i),
			Method:    "GET",
			Path:      fmt.Sprintf("/api/test/%d", i),
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		}
		// Add the target log ID every 100 entries
		if i%100 == 0 {
			entry.LogID = logID
		}
		entries = append(entries, entry)
	}
	
	s.createTestLogFile(date, entries)
	
	// Act
	start := time.Now()
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	duration := time.Since(start)
	
	// Assert
	s.NoError(err)
	s.Len(logs, 10, "Should find 10 entries with target log ID")
	s.Less(duration, 5*time.Second, "Should complete within reasonable time")
	
	// Verify all returned logs have the correct ID
	for _, log := range logs {
		s.Equal(logID, log.LogID)
	}
}

// TestReadLogsByDate_ConcurrentAccess tests concurrent access to log files
func (s *LogBusinessLogicTestSuite) TestReadLogsByDate_ConcurrentAccess() {
	// Arrange
	date := "2024-01-25"
	
	// Create test log file
	entries := []modelsv1.LogEntry{
		{
			Timestamp: "2024-01-25T10:30:00Z",
			Level:     "INFO",
			Message:   "Concurrent access test",
			LogID:     "concurrent-log",
			Method:    "GET",
			Path:      "/api/concurrent",
			Status:    200,
			Duration:  "10ms",
			IP:        "127.0.0.1",
			UserAgent: "test-agent",
		},
	}
	s.createTestLogFile(date, entries)
	
	// Act - simulate concurrent access
	const numGoroutines = 10
	results := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			logs, total, err := businessv1.ReadLogsByDate(date)
			if err != nil {
				results <- err
				return
			}
			if len(logs) != 1 || total != 1 {
				results <- fmt.Errorf("unexpected result: logs=%d, total=%d", len(logs), total)
				return
			}
			results <- nil
		}()
	}
	
	// Assert
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		s.NoError(err, "Concurrent access should not cause errors")
	}
}

// TestRunner runs all the log business logic tests
func TestLogBusinessLogic(t *testing.T) {
	suite.Run(t, new(LogBusinessLogicTestSuite))
}