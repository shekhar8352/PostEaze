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
	
	// Set invalid log directory
	os.Setenv("LOG_DIR", "")
	
	// Act
	logs, err := businessv1.ReadLogsByLogID(s.ctx, logID)
	
	// Assert
	s.Error(err)
	s.Nil(logs)
	s.Contains(err.Error(), "log directory not found")
	
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
	date := "2024-01-01" // Date with no log file
	
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

// TestRunner runs all the log business logic tests
func TestLogBusinessLogic(t *testing.T) {
	suite.Run(t, new(LogBusinessLogicTestSuite))
}