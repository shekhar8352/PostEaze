package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
	"github.com/shekhar8352/PostEaze/utils"
)

// LogIntegrationTestSuite tests log retrieval with real file system operations
type LogIntegrationTestSuite struct {
	suite.Suite
	router      *gin.Engine
	ctx         context.Context
	testLogDir  string
	cleanup     func()
	originalDir string
}

// SetupSuite initializes the test environment with real file system
func (s *LogIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	gin.SetMode(gin.TestMode)
	
	// Create temporary directory for test logs
	var err error
	s.testLogDir, err = os.MkdirTemp("", "postease_test_logs_*")
	require.NoError(s.T(), err, "Failed to create temp log directory")
	
	// Store original log directory and set test directory
	s.originalDir = utils.GetLogDirectory()
	s.setTestLogDirectory()
	
	// Setup test router
	s.router = s.setupTestRouter()
	
	// Create test log files
	s.createTestLogFiles()
	
	s.T().Logf("Log integration test suite setup completed with log dir: %s", s.testLogDir)
}

// TearDownSuite cleans up test environment
func (s *LogIntegrationTestSuite) TearDownSuite() {
	// Restore original log directory
	s.restoreOriginalLogDirectory()
	
	// Clean up temporary directory
	if s.testLogDir != "" {
		os.RemoveAll(s.testLogDir)
	}
	
	if s.cleanup != nil {
		s.cleanup()
	}
	
	s.T().Log("Log integration test suite teardown completed")
}

// SetupTest prepares each test
func (s *LogIntegrationTestSuite) SetupTest() {
	// Ensure test log files exist for each test
	s.createTestLogFiles()
}

// TearDownTest cleans up after each test
func (s *LogIntegrationTestSuite) TearDownTest() {
	// Clean up any test-specific log files if needed
}

// setTestLogDirectory configures the log directory for testing
func (s *LogIntegrationTestSuite) setTestLogDirectory() {
	// Set environment variable that utils.GetLogDirectory() uses
	os.Setenv("LOG_DIR", s.testLogDir)
}

// restoreOriginalLogDirectory restores the original log directory
func (s *LogIntegrationTestSuite) restoreOriginalLogDirectory() {
	if s.originalDir != "" {
		os.Setenv("LOG_DIR", s.originalDir)
	} else {
		os.Unsetenv("LOG_DIR")
	}
}

// setupTestRouter creates a test router with log routes
func (s *LogIntegrationTestSuite) setupTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	api := router.Group("/api")
	v1 := api.Group("/v1")
	logv1 := v1.Group("/logs")
	
	// Add log routes that call actual business logic
	logv1.GET("/date/:date", s.getLogsByDateHandler)
	logv1.GET("/id/:log_id", s.getLogByIDHandler)
	
	return router
}

// getLogsByDateHandler handles log retrieval by date using real business logic
func (s *LogIntegrationTestSuite) getLogsByDateHandler(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Date is required"})
		return
	}
	
	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid date format"})
		return
	}
	
	// Call actual business logic
	logs, total, err := businessv1.ReadLogsByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Read logs by date successfully",
		"data":   gin.H{"logs": logs, "total": total},
	})
}

// getLogByIDHandler handles log retrieval by ID using real business logic
func (s *LogIntegrationTestSuite) getLogByIDHandler(c *gin.Context) {
	logID := c.Param("log_id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Log ID is required"})
		return
	}
	
	// Call actual business logic
	logs, err := businessv1.ReadLogsByLogID(c.Request.Context(), logID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Read log by ID successfully",
		"data":   logs,
	})
}

// createTestLogFiles creates test log files with sample data
func (s *LogIntegrationTestSuite) createTestLogFiles() {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	
	// Create log files for different dates
	s.createLogFile(today, s.generateLogEntries(today, 10))
	s.createLogFile(yesterday, s.generateLogEntries(yesterday, 15))
	s.createLogFile(twoDaysAgo, s.generateLogEntries(twoDaysAgo, 8))
}

// createLogFile creates a log file with specified content
func (s *LogIntegrationTestSuite) createLogFile(date string, entries []modelsv1.LogEntry) {
	filename := fmt.Sprintf("app-%s.log", date)
	filepath := filepath.Join(s.testLogDir, filename)
	
	file, err := os.Create(filepath)
	require.NoError(s.T(), err, "Failed to create log file: %s", filename)
	defer file.Close()
	
	for _, entry := range entries {
		jsonData, err := json.Marshal(entry)
		require.NoError(s.T(), err, "Failed to marshal log entry")
		
		_, err = file.WriteString(string(jsonData) + "\n")
		require.NoError(s.T(), err, "Failed to write log entry")
	}
	
	s.T().Logf("Created test log file: %s with %d entries", filename, len(entries))
}

// generateLogEntries generates sample log entries for testing
func (s *LogIntegrationTestSuite) generateLogEntries(date string, count int) []modelsv1.LogEntry {
	entries := make([]modelsv1.LogEntry, count)
	
	for i := 0; i < count; i++ {
		timestamp := fmt.Sprintf("%sT%02d:%02d:%02d.000Z", date, i%24, i%60, i%60)
		logID := fmt.Sprintf("log-%s-%03d", date, i+1)
		
		entries[i] = modelsv1.LogEntry{
			Timestamp: timestamp,
			Level:     s.getRandomLogLevel(i),
			Message:   fmt.Sprintf("Test log message %d for date %s", i+1, date),
			LogID:     logID,
			Method:    s.getRandomMethod(i),
			Path:      fmt.Sprintf("/api/v1/test/%d", i+1),
			Status:    s.getRandomStatus(i),
			Duration:  fmt.Sprintf("%dms", (i+1)*10),
			IP:        fmt.Sprintf("192.168.1.%d", (i%254)+1),
			UserAgent: "PostEaze-Test-Client/1.0",
			Extra: map[string]string{
				"request_id": fmt.Sprintf("req-%s-%03d", date, i+1),
				"user_id":    fmt.Sprintf("user-%03d", i+1),
			},
		}
	}
	
	return entries
}

// getRandomLogLevel returns a random log level for testing
func (s *LogIntegrationTestSuite) getRandomLogLevel(index int) string {
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR"}
	return levels[index%len(levels)]
}

// getRandomMethod returns a random HTTP method for testing
func (s *LogIntegrationTestSuite) getRandomMethod(index int) string {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	return methods[index%len(methods)]
}

// getRandomStatus returns a random HTTP status code for testing
func (s *LogIntegrationTestSuite) getRandomStatus(index int) int {
	statuses := []int{200, 201, 400, 404, 500}
	return statuses[index%len(statuses)]
}

// TestLogRetrievalByDate tests log retrieval by date with real file system operations
func (s *LogIntegrationTestSuite) TestLogRetrievalByDate() {
	today := time.Now().Format("2006-01-02")
	
	// Test retrieving logs for today
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", today), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Log retrieval by date should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse log response")
	
	assert.Equal(s.T(), "success", result["status"])
	assert.Contains(s.T(), result, "data")
	
	data := result["data"].(map[string]interface{})
	assert.Contains(s.T(), data, "logs")
	assert.Contains(s.T(), data, "total")
	
	logs := data["logs"].([]interface{})
	total := int(data["total"].(float64))
	
	assert.Equal(s.T(), len(logs), total, "Total count should match log entries count")
	assert.Greater(s.T(), len(logs), 0, "Should retrieve some log entries")
	
	// Verify log entry structure
	if len(logs) > 0 {
		firstLog := logs[0].(map[string]interface{})
		assert.Contains(s.T(), firstLog, "timestamp")
		assert.Contains(s.T(), firstLog, "level")
		assert.Contains(s.T(), firstLog, "message")
		assert.Contains(s.T(), firstLog, "log_id")
	}
	
	s.T().Logf("Successfully retrieved %d log entries for date %s", len(logs), today)
}

// TestLogRetrievalByID tests log retrieval by log ID with real file system operations
func (s *LogIntegrationTestSuite) TestLogRetrievalByID() {
	today := time.Now().Format("2006-01-02")
	testLogID := fmt.Sprintf("log-%s-001", today)
	
	// Test retrieving logs by specific log ID
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/id/%s", testLogID), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Log retrieval by ID should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse log response")
	
	assert.Equal(s.T(), "success", result["status"])
	assert.Contains(s.T(), result, "data")
	
	logs := result["data"].([]interface{})
	assert.GreaterOrEqual(s.T(), len(logs), 1, "Should retrieve at least one log entry")
	
	// Verify that retrieved logs have the correct log ID
	for _, logInterface := range logs {
		logEntry := logInterface.(map[string]interface{})
		assert.Equal(s.T(), testLogID, logEntry["log_id"], "Retrieved log should have correct log ID")
	}
	
	s.T().Logf("Successfully retrieved %d log entries for log ID %s", len(logs), testLogID)
}

// TestLogRetrievalWithNonExistentDate tests log retrieval for non-existent date
func (s *LogIntegrationTestSuite) TestLogRetrievalWithNonExistentDate() {
	nonExistentDate := "2020-01-01" // Date that doesn't have log files
	
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", nonExistentDate), nil)
	assert.Equal(s.T(), http.StatusInternalServerError, response.Code, "Should return error for non-existent date")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse error response")
	
	assert.Equal(s.T(), "error", result["status"])
	assert.Contains(s.T(), result["msg"], "log file not found")
	
	s.T().Logf("Correctly handled non-existent date: %s", nonExistentDate)
}

// TestLogRetrievalWithInvalidDate tests log retrieval with invalid date format
func (s *LogIntegrationTestSuite) TestLogRetrievalWithInvalidDate() {
	invalidDate := "invalid-date-format"
	
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", invalidDate), nil)
	assert.Equal(s.T(), http.StatusBadRequest, response.Code, "Should return bad request for invalid date")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse error response")
	
	assert.Equal(s.T(), "error", result["status"])
	assert.Contains(s.T(), result["msg"], "Invalid date format")
	
	s.T().Log("Correctly handled invalid date format")
}

// TestLogRetrievalWithNonExistentLogID tests log retrieval for non-existent log ID
func (s *LogIntegrationTestSuite) TestLogRetrievalWithNonExistentLogID() {
	nonExistentLogID := "non-existent-log-id"
	
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/id/%s", nonExistentLogID), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Should succeed but return empty results")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse response")
	
	assert.Equal(s.T(), "success", result["status"])
	logs := result["data"].([]interface{})
	assert.Equal(s.T(), 0, len(logs), "Should return empty log array for non-existent log ID")
	
	s.T().Log("Correctly handled non-existent log ID")
}

// TestLogRetrievalAcrossMultipleDates tests log retrieval that spans multiple dates
func (s *LogIntegrationTestSuite) TestLogRetrievalAcrossMultipleDates() {
	// Create a log ID that appears in multiple days
	commonLogID := "common-log-id"
	
	// Add the common log ID to multiple date files
	dates := []string{
		time.Now().Format("2006-01-02"),
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
	}
	
	for _, date := range dates {
		entry := modelsv1.LogEntry{
			Timestamp: fmt.Sprintf("%sT12:00:00.000Z", date),
			Level:     "INFO",
			Message:   fmt.Sprintf("Common log entry for %s", date),
			LogID:     commonLogID,
			Method:    "GET",
			Path:      "/api/v1/common",
			Status:    200,
			Duration:  "50ms",
			IP:        "192.168.1.100",
			UserAgent: "PostEaze-Test-Client/1.0",
		}
		
		filename := fmt.Sprintf("app-%s.log", date)
		filepath := filepath.Join(s.testLogDir, filename)
		
		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
		require.NoError(s.T(), err, "Failed to open log file for appending")
		
		jsonData, err := json.Marshal(entry)
		require.NoError(s.T(), err, "Failed to marshal log entry")
		
		_, err = file.WriteString(string(jsonData) + "\n")
		require.NoError(s.T(), err, "Failed to write log entry")
		
		file.Close()
	}
	
	// Test retrieving the common log ID
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/id/%s", commonLogID), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Log retrieval across dates should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse response")
	
	assert.Equal(s.T(), "success", result["status"])
	logs := result["data"].([]interface{})
	assert.Equal(s.T(), len(dates), len(logs), "Should retrieve logs from all dates")
	
	s.T().Logf("Successfully retrieved %d log entries across multiple dates", len(logs))
}

// TestConcurrentLogRetrieval tests concurrent log retrieval operations
func (s *LogIntegrationTestSuite) TestConcurrentLogRetrieval() {
	const numConcurrentRequests = 10
	today := time.Now().Format("2006-01-02")
	
	// Create channels to collect results
	responses := make(chan *httptest.ResponseRecorder, numConcurrentRequests)
	
	// Launch concurrent log retrieval requests
	for i := 0; i < numConcurrentRequests; i++ {
		go func(index int) {
			path := fmt.Sprintf("/api/v1/logs/date/%s", today)
			response := s.performRequest("GET", path, nil)
			responses <- response
		}(i)
	}
	
	// Collect and verify responses
	successCount := 0
	for i := 0; i < numConcurrentRequests; i++ {
		response := <-responses
		if response.Code == http.StatusOK {
			successCount++
		}
	}
	
	assert.Equal(s.T(), numConcurrentRequests, successCount, "All concurrent requests should succeed")
	s.T().Logf("Concurrent log retrieval test completed - %d/%d requests succeeded", successCount, numConcurrentRequests)
}

// TestLogFileSystemOperations tests various file system operations
func (s *LogIntegrationTestSuite) TestLogFileSystemOperations() {
	// Test with empty log directory
	emptyLogDir, err := os.MkdirTemp("", "empty_logs_*")
	require.NoError(s.T(), err, "Failed to create empty log directory")
	defer os.RemoveAll(emptyLogDir)
	
	// Temporarily change log directory to empty one
	originalDir := s.testLogDir
	s.testLogDir = emptyLogDir
	s.setTestLogDirectory()
	
	// Test log retrieval with empty directory
	today := time.Now().Format("2006-01-02")
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", today), nil)
	assert.Equal(s.T(), http.StatusInternalServerError, response.Code, "Should fail with empty log directory")
	
	// Restore original log directory
	s.testLogDir = originalDir
	s.setTestLogDirectory()
	
	s.T().Log("File system operations test completed")
}

// performRequest is a helper method to perform HTTP requests on the test router
func (s *LogIntegrationTestSuite) performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	req := testutils.CreateTestRequest(method, path, body)
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, req)
	return recorder
}

// TestLogIntegrationSuite runs the log integration test suite
func TestLogIntegrationSuite(t *testing.T) {
	suite.Run(t, new(LogIntegrationTestSuite))
}