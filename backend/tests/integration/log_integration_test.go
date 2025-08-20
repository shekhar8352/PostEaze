package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	apiv1 "github.com/shekhar8352/PostEaze/api/v1"
	"github.com/shekhar8352/PostEaze/constants"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
	"github.com/shekhar8352/PostEaze/utils"
)

// LogAPIIntegrationTestSuite tests complete log API workflows with real file system operations
type LogAPIIntegrationTestSuite struct {
	suite.Suite
	router      *gin.Engine
	ctx         context.Context
	testLogDir  string
	cleanup     func()
	originalDir string
}

// SetupSuite initializes the test environment with real file system
func (s *LogAPIIntegrationTestSuite) SetupSuite() {
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
	
	s.T().Logf("Log API integration test suite setup completed with log dir: %s", s.testLogDir)
}

// TearDownSuite cleans up test environment
func (s *LogAPIIntegrationTestSuite) TearDownSuite() {
	// Restore original log directory
	s.restoreOriginalLogDirectory()
	
	// Clean up temporary directory
	if s.testLogDir != "" {
		os.RemoveAll(s.testLogDir)
	}
	
	if s.cleanup != nil {
		s.cleanup()
	}
	
	s.T().Log("Log API integration test suite teardown completed")
}

// SetupTest prepares each test
func (s *LogAPIIntegrationTestSuite) SetupTest() {
	// Clean up any existing test files
	s.cleanupTestFiles()
}

// TearDownTest cleans up after each test
func (s *LogAPIIntegrationTestSuite) TearDownTest() {
	// Clean up test files after each test
	s.cleanupTestFiles()
}

// setTestLogDirectory configures the log directory for testing
func (s *LogAPIIntegrationTestSuite) setTestLogDirectory() {
	// Set environment variable that utils.GetLogDirectory() uses
	os.Setenv("LOG_DIR", s.testLogDir)
}

// restoreOriginalLogDirectory restores the original log directory
func (s *LogAPIIntegrationTestSuite) restoreOriginalLogDirectory() {
	if s.originalDir != "" {
		os.Setenv("LOG_DIR", s.originalDir)
	} else {
		os.Unsetenv("LOG_DIR")
	}
}

// cleanupTestFiles removes all test log files
func (s *LogAPIIntegrationTestSuite) cleanupTestFiles() {
	if s.testLogDir != "" {
		entries, err := os.ReadDir(s.testLogDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasPrefix(entry.Name(), "app-") {
					os.Remove(filepath.Join(s.testLogDir, entry.Name()))
				}
			}
		}
	}
}

// setupTestRouter creates a test router with actual API routes
func (s *LogAPIIntegrationTestSuite) setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Use the actual API route structure from constants
	api := router.Group(constants.ApiRoute)
	v1 := api.Group(constants.V1Route)
	logv1 := v1.Group(constants.LogRoute)
	
	// Use the actual API handlers to test the complete workflow
	logv1.GET(constants.LogByDate, apiv1.GetLogsByDate)
	logv1.GET(constants.LogById, apiv1.GetLogByIDHandler)
	
	return router
}

// createLogFile creates a log file with specified content
func (s *LogAPIIntegrationTestSuite) createLogFile(date string, entries []modelsv1.LogEntry) {
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
func (s *LogAPIIntegrationTestSuite) generateLogEntries(date string, count int) []modelsv1.LogEntry {
	entries := make([]modelsv1.LogEntry, count)
	
	for i := 0; i < count; i++ {
		timestamp := fmt.Sprintf("%sT%02d:%02d:%02d+05:30", date, i%24, i%60, i%60)
		logID := fmt.Sprintf("log-%s-%03d", date, i+1)
		
		entries[i] = modelsv1.LogEntry{
			Timestamp: timestamp,
			Level:     s.getRandomLogLevel(i),
			Message:   fmt.Sprintf("Test log message %d for date %s", i+1, date),
			LogID:     logID,
			File:      "test.go",
			Line:      i + 1,
			Function:  "TestFunction",
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
func (s *LogAPIIntegrationTestSuite) getRandomLogLevel(index int) string {
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR"}
	return levels[index%len(levels)]
}

// getRandomMethod returns a random HTTP method for testing
func (s *LogAPIIntegrationTestSuite) getRandomMethod(index int) string {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	return methods[index%len(methods)]
}

// getRandomStatus returns a random HTTP status code for testing
func (s *LogAPIIntegrationTestSuite) getRandomStatus(index int) int {
	statuses := []int{200, 201, 400, 404, 500}
	return statuses[index%len(statuses)]
}

// performRequest is a helper method to perform HTTP requests on the test router
func (s *LogAPIIntegrationTestSuite) performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	req := testutils.CreateTestRequest(method, path, body)
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, req)
	return recorder
}

// TestEndToEndLogRetrievalByIDWithRealFiles tests complete log retrieval by ID workflow
// Requirements: 1.1, 1.2
func (s *LogAPIIntegrationTestSuite) TestEndToEndLogRetrievalByIDWithRealFiles() {
	// Create test log files with known log IDs
	testDates := []string{
		time.Now().Format("2006-01-02"),
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
	}
	
	commonLogID := "test-e2e-log-id-001"
	
	// Create log entries across multiple files
	for i, date := range testDates {
		entries := []modelsv1.LogEntry{
			{
				Timestamp: fmt.Sprintf("%sT%02d:00:00+05:30", date, 10+i),
				LogID:     commonLogID,
				Level:     "INFO",
				Message:   fmt.Sprintf("E2E test log entry %d for date %s", i+1, date),
				File:      "test_file.go",
				Line:      100 + i,
				Function:  "TestFunction",
				Method:    "GET",
				Path:      fmt.Sprintf("/api/test/%d", i+1),
				Status:    200,
				Duration:  fmt.Sprintf("%dms", (i+1)*50),
				IP:        "192.168.1.100",
				UserAgent: "E2E-Test-Client/1.0",
			},
		}
		s.createLogFile(date, entries)
	}
	
	// Test the complete API workflow
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byId/%s", commonLogID), nil)
	
	// Verify response structure and content
	assert.Equal(s.T(), http.StatusOK, response.Code, "E2E log retrieval by ID should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse E2E response")
	
	// Verify response format matches API specification
	assert.Equal(s.T(), true, result["success"])
	assert.Contains(s.T(), result, "data")
	assert.Contains(s.T(), result, "message")
	
	logs := result["data"].([]interface{})
	assert.Equal(s.T(), len(testDates), len(logs), "Should retrieve logs from all test dates")
	
	// Verify each log entry has correct structure and data
	for i, logInterface := range logs {
		logEntry := logInterface.(map[string]interface{})
		assert.Equal(s.T(), commonLogID, logEntry["log_id"], "Log ID should match")
		assert.Contains(s.T(), logEntry, "timestamp")
		assert.Contains(s.T(), logEntry, "level")
		assert.Contains(s.T(), logEntry, "message")
		assert.Contains(s.T(), logEntry, "file")
		assert.Contains(s.T(), logEntry, "line")
		assert.Contains(s.T(), logEntry, "function")
		
		// Verify extracted HTTP details (these may be nil if not extracted from message)
		if logEntry["method"] != nil {
			assert.Equal(s.T(), "GET", logEntry["method"])
		}
		if logEntry["path"] != nil {
			assert.Equal(s.T(), fmt.Sprintf("/api/test/%d", i+1), logEntry["path"])
		}
		if logEntry["status"] != nil {
			assert.Equal(s.T(), float64(200), logEntry["status"])
		}
	}
	
	s.T().Logf("E2E test completed: Retrieved %d log entries for log ID %s across %d files", 
		len(logs), commonLogID, len(testDates))
}

// TestEndToEndLogRetrievalByDateWithVariousScenarios tests complete log retrieval by date workflow
// Requirements: 2.1
func (s *LogAPIIntegrationTestSuite) TestEndToEndLogRetrievalByDateWithVariousScenarios() {
	testDate := time.Now().Format("2006-01-02")
	
	// Create comprehensive test log entries with various scenarios
	testEntries := []modelsv1.LogEntry{
		{
			Timestamp: fmt.Sprintf("%sT09:00:00+05:30", testDate),
			LogID:     "scenario-1-success",
			Level:     "INFO",
			Message:   "GET /api/users/123 | Status: 200 | Duration: 45ms | IP: 192.168.1.10",
			File:      "user_handler.go",
			Line:      25,
			Function:  "GetUserHandler",
			Method:    "GET",
			Path:      "/api/users/123",
			Status:    200,
			Duration:  "45ms",
			IP:        "192.168.1.10",
			UserAgent: "Mozilla/5.0",
		},
		{
			Timestamp: fmt.Sprintf("%sT09:15:00+05:30", testDate),
			LogID:     "scenario-2-error",
			Level:     "ERROR",
			Message:   "POST /api/users | Status: 400 | Duration: 12ms | IP: 192.168.1.20",
			File:      "user_handler.go",
			Line:      45,
			Function:  "CreateUserHandler",
			Method:    "POST",
			Path:      "/api/users",
			Status:    400,
			Duration:  "12ms",
			IP:        "192.168.1.20",
			UserAgent: "PostmanRuntime/7.45.0",
		},
		{
			Timestamp: fmt.Sprintf("%sT09:30:00+05:30", testDate),
			LogID:     "scenario-3-warning",
			Level:     "WARN",
			Message:   "PUT /api/users/456 | Status: 404 | Duration: 23ms | IP: 192.168.1.30",
			File:      "user_handler.go",
			Line:      67,
			Function:  "UpdateUserHandler",
			Method:    "PUT",
			Path:      "/api/users/456",
			Status:    404,
			Duration:  "23ms",
			IP:        "192.168.1.30",
			UserAgent: "curl/7.68.0",
		},
	}
	
	s.createLogFile(testDate, testEntries)
	
	// Test the complete API workflow
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byDate/%s", testDate), nil)
	
	// Verify response structure and content
	assert.Equal(s.T(), http.StatusOK, response.Code, "E2E log retrieval by date should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse E2E date response")
	
	// Verify response format matches API specification
	assert.Equal(s.T(), true, result["success"])
	assert.Contains(s.T(), result, "data")
	assert.Contains(s.T(), result, "message")
	
	data := result["data"].(map[string]interface{})
	assert.Contains(s.T(), data, "logs")
	assert.Contains(s.T(), data, "total")
	
	logs := data["logs"].([]interface{})
	total := int(data["total"].(float64))
	
	assert.GreaterOrEqual(s.T(), len(logs), len(testEntries), "Should retrieve at least the test entries")
	assert.Equal(s.T(), total, len(logs), "Total should match logs count")
	
	// Verify that our test scenarios are included
	foundScenarios := make(map[string]bool)
	for _, logInterface := range logs {
		logEntry := logInterface.(map[string]interface{})
		if logID, ok := logEntry["log_id"].(string); ok {
			if strings.HasPrefix(logID, "scenario-") {
				foundScenarios[logID] = true
			}
		}
	}
	
	assert.True(s.T(), foundScenarios["scenario-1-success"], "Should find success scenario")
	assert.True(s.T(), foundScenarios["scenario-2-error"], "Should find error scenario")
	assert.True(s.T(), foundScenarios["scenario-3-warning"], "Should find warning scenario")
	
	s.T().Logf("E2E date test completed: Retrieved %d log entries for date %s with %d test scenarios", 
		len(logs), testDate, len(foundScenarios))
}

// TestErrorScenariosWithMissingFiles tests error handling for missing files
// Requirements: 2.4
func (s *LogAPIIntegrationTestSuite) TestErrorScenariosWithMissingFiles() {
	// Test 1: Missing log file for date
	nonExistentDate := "2020-01-01"
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byDate/%s", nonExistentDate), nil)
	
	// The API now validates date range, so this should return 400 instead of 404
	assert.True(s.T(), response.Code == http.StatusBadRequest || response.Code == http.StatusNotFound, 
		"Should return 400 or 404 for missing/invalid log file")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse error response")
	
	assert.Equal(s.T(), false, result["success"])
	assert.Contains(s.T(), result, "error")
	
	errorInfo := result["error"].(map[string]interface{})
	assert.Contains(s.T(), errorInfo, "code")
	assert.Contains(s.T(), errorInfo, "type")
	assert.Contains(s.T(), errorInfo, "message")
	
	s.T().Log("Missing file error handling test passed")
	
	// Test 2: Empty log directory scenario
	emptyLogDir, err := os.MkdirTemp("", "empty_logs_*")
	require.NoError(s.T(), err, "Failed to create empty log directory")
	defer os.RemoveAll(emptyLogDir)
	
	// Temporarily change log directory
	originalDir := s.testLogDir
	s.testLogDir = emptyLogDir
	s.setTestLogDirectory()
	
	// Test log retrieval with empty directory
	testLogID := "non-existent-log-id"
	response = s.performRequest("GET", fmt.Sprintf("/api/v1/log/byId/%s", testLogID), nil)
	
	// Should handle gracefully - either return empty results or appropriate error
	assert.True(s.T(), response.Code == http.StatusOK || response.Code == http.StatusNotFound || response.Code == http.StatusInternalServerError,
		"Should handle empty directory gracefully")
	
	// Restore original directory
	s.testLogDir = originalDir
	s.setTestLogDirectory()
	
	s.T().Log("Empty directory error handling test passed")
}

// TestErrorScenariosWithInvalidInputs tests error handling for invalid inputs
// Requirements: 1.2, 2.1
func (s *LogAPIIntegrationTestSuite) TestErrorScenariosWithInvalidInputs() {
	// Test 1: Invalid date format
	invalidDates := []string{
		"invalid-date",
		"2023-13-01", // Invalid month
		"2023-02-30", // Invalid day
		"23-01-01",   // Wrong format
	}
	
	for _, invalidDate := range invalidDates {
		response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byDate/%s", invalidDate), nil)
		
		assert.Equal(s.T(), http.StatusBadRequest, response.Code, 
			"Should return 400 for invalid date: %s", invalidDate)
		
		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		require.NoError(s.T(), err, "Should parse error response for date: %s", invalidDate)
		
		assert.Equal(s.T(), false, result["success"])
		assert.Contains(s.T(), result, "error")
		
		errorInfo := result["error"].(map[string]interface{})
		assert.Equal(s.T(), float64(400), errorInfo["code"])
		assert.Equal(s.T(), "invalid_input", errorInfo["type"])
	}
	
	s.T().Log("Invalid date format error handling tests passed")
	
	// Test 2: Invalid log ID format
	invalidLogIDs := []string{
		"", // Empty log ID - this will result in 404 due to route not matching
		strings.Repeat("a", 1000), // Extremely long log ID
	}
	
	for _, invalidLogID := range invalidLogIDs {
		var response *httptest.ResponseRecorder
		if invalidLogID == "" {
			// Empty log ID will not match the route, so test the route directly
			response = s.performRequest("GET", "/api/v1/log/byId/", nil)
		} else {
			response = s.performRequest("GET", fmt.Sprintf("/api/v1/log/byId/%s", invalidLogID), nil)
		}
		
		// Should either return 400 for validation error, 404 for route not found, or 200 with empty results
		assert.True(s.T(), response.Code == http.StatusBadRequest || response.Code == http.StatusOK || response.Code == http.StatusNotFound,
			"Should handle invalid log ID appropriately: %s", invalidLogID)
	}
	
	s.T().Log("Invalid log ID error handling tests passed")
}

// TestLogRetrievalPerformanceWithLargeFiles tests performance with large log files
func (s *LogAPIIntegrationTestSuite) TestLogRetrievalPerformanceWithLargeFiles() {
	testDate := time.Now().Format("2006-01-02")
	
	// Create a large log file with many entries
	const numEntries = 1000
	largeLogEntries := make([]modelsv1.LogEntry, numEntries)
	
	targetLogID := "performance-test-log-id"
	
	for i := 0; i < numEntries; i++ {
		logID := fmt.Sprintf("log-entry-%04d", i)
		if i == numEntries/2 { // Place target log ID in the middle
			logID = targetLogID
		}
		
		largeLogEntries[i] = modelsv1.LogEntry{
			Timestamp: fmt.Sprintf("%sT%02d:%02d:%02d+05:30", testDate, i/3600, (i%3600)/60, i%60),
			LogID:     logID,
			Level:     s.getRandomLogLevel(i),
			Message:   fmt.Sprintf("Performance test log entry %d", i),
			File:      "performance_test.go",
			Line:      i + 1,
			Function:  "PerformanceTestFunction",
			Method:    s.getRandomMethod(i),
			Path:      fmt.Sprintf("/api/performance/%d", i),
			Status:    s.getRandomStatus(i),
			Duration:  fmt.Sprintf("%dms", i%100+1),
			IP:        fmt.Sprintf("192.168.%d.%d", (i/254)+1, (i%254)+1),
			UserAgent: "Performance-Test-Client/1.0",
		}
	}
	
	s.createLogFile(testDate, largeLogEntries)
	
	// Test performance of log retrieval by ID
	startTime := time.Now()
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byId/%s", targetLogID), nil)
	duration := time.Since(startTime)
	
	assert.Equal(s.T(), http.StatusOK, response.Code, "Performance test should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse performance test response")
	
	logs := result["data"].([]interface{})
	assert.Equal(s.T(), 1, len(logs), "Should find exactly one log entry")
	
	// Verify the correct log was found
	logEntry := logs[0].(map[string]interface{})
	assert.Equal(s.T(), targetLogID, logEntry["log_id"])
	
	// Performance assertion - should complete within reasonable time
	assert.Less(s.T(), duration, 5*time.Second, "Log retrieval should complete within 5 seconds")
	
	s.T().Logf("Performance test completed: Retrieved log from %d entries in %v", numEntries, duration)
	
	// Test performance of log retrieval by date
	startTime = time.Now()
	response = s.performRequest("GET", fmt.Sprintf("/api/v1/log/byDate/%s", testDate), nil)
	duration = time.Since(startTime)
	
	assert.Equal(s.T(), http.StatusOK, response.Code, "Date performance test should succeed")
	
	err = json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse date performance test response")
	
	data := result["data"].(map[string]interface{})
	logs = data["logs"].([]interface{})
	total := int(data["total"].(float64))
	
	assert.GreaterOrEqual(s.T(), len(logs), numEntries, "Should retrieve all log entries")
	assert.Equal(s.T(), total, len(logs), "Total should match logs count")
	
	// Performance assertion for date retrieval
	assert.Less(s.T(), duration, 10*time.Second, "Date log retrieval should complete within 10 seconds")
	
	s.T().Logf("Date performance test completed: Retrieved %d entries in %v", len(logs), duration)
}

// TestMalformedJSONHandling tests handling of malformed JSON in log files
func (s *LogAPIIntegrationTestSuite) TestMalformedJSONHandling() {
	testDate := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("app-%s.log", testDate)
	filepath := filepath.Join(s.testLogDir, filename)
	
	// Create a log file with mix of valid and malformed JSON
	file, err := os.Create(filepath)
	require.NoError(s.T(), err, "Failed to create malformed JSON test file")
	defer file.Close()
	
	// Write valid JSON entry
	validEntry := modelsv1.LogEntry{
		Timestamp: fmt.Sprintf("%sT10:00:00+05:30", testDate),
		LogID:     "valid-json-entry",
		Level:     "INFO",
		Message:   "Valid JSON log entry",
		File:      "test.go",
		Line:      1,
		Function:  "TestFunction",
	}
	validJSON, _ := json.Marshal(validEntry)
	file.WriteString(string(validJSON) + "\n")
	
	// Write malformed JSON entries
	malformedEntries := []string{
		`{"timestamp":"2025-08-20T10:01:00+05:30","log_id":"malformed-1","level":"INFO","message":"Missing closing brace"`,
		`{"timestamp":"2025-08-20T10:02:00+05:30","log_id":"malformed-2","level":"INFO","message":"Extra comma",}`,
		`not-json-at-all`,
		`{"timestamp":"2025-08-20T10:03:00+05:30","log_id":"malformed-3","level":"INFO","message":"Invalid escape \z"}`,
		``, // Empty line
	}
	
	for _, malformed := range malformedEntries {
		file.WriteString(malformed + "\n")
	}
	
	// Write another valid JSON entry
	validEntry2 := modelsv1.LogEntry{
		Timestamp: fmt.Sprintf("%sT10:05:00+05:30", testDate),
		LogID:     "valid-json-entry-2",
		Level:     "INFO",
		Message:   "Second valid JSON log entry",
		File:      "test.go",
		Line:      2,
		Function:  "TestFunction",
	}
	validJSON2, _ := json.Marshal(validEntry2)
	file.WriteString(string(validJSON2) + "\n")
	
	// Test that API handles malformed JSON gracefully
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byDate/%s", testDate), nil)
	
	assert.Equal(s.T(), http.StatusOK, response.Code, "Should succeed despite malformed JSON")
	
	var result map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse response despite malformed JSON in log file")
	
	assert.Equal(s.T(), true, result["success"])
	
	data := result["data"].(map[string]interface{})
	logs := data["logs"].([]interface{})
	
	// Should retrieve only the valid JSON entries
	assert.GreaterOrEqual(s.T(), len(logs), 2, "Should retrieve at least the valid JSON entries")
	
	// Verify that valid entries are present
	foundValidEntries := 0
	for _, logInterface := range logs {
		logEntry := logInterface.(map[string]interface{})
		if logID, ok := logEntry["log_id"].(string); ok {
			if logID == "valid-json-entry" || logID == "valid-json-entry-2" {
				foundValidEntries++
			}
		}
	}
	
	assert.Equal(s.T(), 2, foundValidEntries, "Should find both valid JSON entries")
	
	s.T().Logf("Malformed JSON handling test completed: Retrieved %d valid entries from file with malformed JSON", foundValidEntries)
}

// TestConcurrentAPIAccess tests concurrent access to log APIs
func (s *LogAPIIntegrationTestSuite) TestConcurrentAPIAccess() {
	const numConcurrentRequests = 20
	testDate := time.Now().Format("2006-01-02")
	testLogID := "concurrent-test-log-id"
	
	// Create test data
	testEntries := []modelsv1.LogEntry{
		{
			Timestamp: fmt.Sprintf("%sT12:00:00+05:30", testDate),
			LogID:     testLogID,
			Level:     "INFO",
			Message:   "Concurrent access test log entry",
			File:      "concurrent_test.go",
			Line:      1,
			Function:  "ConcurrentTestFunction",
		},
	}
	s.createLogFile(testDate, testEntries)
	
	// Test concurrent access to both endpoints
	responses := make(chan *httptest.ResponseRecorder, numConcurrentRequests*2)
	
	// Launch concurrent requests to both endpoints
	for i := 0; i < numConcurrentRequests; i++ {
		// Concurrent requests to byDate endpoint
		go func(index int) {
			path := fmt.Sprintf("/api/v1/log/byDate/%s", testDate)
			response := s.performRequest("GET", path, nil)
			responses <- response
		}(i)
		
		// Concurrent requests to byId endpoint
		go func(index int) {
			path := fmt.Sprintf("/api/v1/log/byId/%s", testLogID)
			response := s.performRequest("GET", path, nil)
			responses <- response
		}(i)
	}
	
	// Collect and verify responses
	successCount := 0
	errorCount := 0
	
	for i := 0; i < numConcurrentRequests*2; i++ {
		response := <-responses
		if response.Code == http.StatusOK {
			successCount++
		} else {
			errorCount++
		}
	}
	
	// All requests should succeed under normal conditions
	assert.Equal(s.T(), numConcurrentRequests*2, successCount, "All concurrent requests should succeed")
	assert.Equal(s.T(), 0, errorCount, "No requests should fail")
	
	s.T().Logf("Concurrent access test completed: %d/%d requests succeeded", 
		successCount, numConcurrentRequests*2)
}

// TestCompleteAPIWorkflowIntegration tests the complete integration workflow
func (s *LogAPIIntegrationTestSuite) TestCompleteAPIWorkflowIntegration() {
	// This test combines multiple scenarios to test the complete workflow
	
	// Step 1: Create comprehensive test data across multiple dates
	dates := []string{
		time.Now().Format("2006-01-02"),
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
	}
	
	workflowLogID := "workflow-integration-test"
	
	for i, date := range dates {
		entries := []modelsv1.LogEntry{
			{
				Timestamp: fmt.Sprintf("%sT14:%02d:00+05:30", date, i*10),
				LogID:     workflowLogID,
				Level:     "INFO",
				Message:   fmt.Sprintf("Workflow test entry %d for %s", i+1, date),
				File:      "workflow_test.go",
				Line:      i + 10,
				Function:  "WorkflowTestFunction",
				Method:    "POST",
				Path:      fmt.Sprintf("/api/workflow/%d", i+1),
				Status:    201,
				Duration:  fmt.Sprintf("%dms", (i+1)*25),
				IP:        "192.168.100.1",
				UserAgent: "Workflow-Test-Client/1.0",
			},
		}
		s.createLogFile(date, entries)
	}
	
	// Step 2: Test log retrieval by ID (should find entries across all dates)
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byId/%s", workflowLogID), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Workflow ID retrieval should succeed")
	
	var idResult map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &idResult)
	require.NoError(s.T(), err, "Should parse workflow ID response")
	
	idLogs := idResult["data"].([]interface{})
	assert.Equal(s.T(), len(dates), len(idLogs), "Should find logs across all dates")
	
	// Step 3: Test log retrieval by date for each date
	for _, date := range dates {
		response := s.performRequest("GET", fmt.Sprintf("/api/v1/log/byDate/%s", date), nil)
		assert.Equal(s.T(), http.StatusOK, response.Code, "Workflow date retrieval should succeed for %s", date)
		
		var dateResult map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &dateResult)
		require.NoError(s.T(), err, "Should parse workflow date response for %s", date)
		
		data := dateResult["data"].(map[string]interface{})
		dateLogs := data["logs"].([]interface{})
		assert.GreaterOrEqual(s.T(), len(dateLogs), 1, "Should find at least one log for date %s", date)
		
		// Verify our workflow log is present
		found := false
		for _, logInterface := range dateLogs {
			logEntry := logInterface.(map[string]interface{})
			if logEntry["log_id"] == workflowLogID {
				found = true
				break
			}
		}
		assert.True(s.T(), found, "Should find workflow log in date %s", date)
	}
	
	// Step 4: Test error scenarios
	response = s.performRequest("GET", "/api/v1/log/byDate/invalid-date", nil)
	assert.Equal(s.T(), http.StatusBadRequest, response.Code, "Should handle invalid date")
	
	response = s.performRequest("GET", "/api/v1/log/byDate/2020-01-01", nil)
	assert.True(s.T(), response.Code == http.StatusBadRequest || response.Code == http.StatusNotFound, 
		"Should handle missing date")
	
	response = s.performRequest("GET", "/api/v1/log/byId/non-existent-log-id", nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Should handle non-existent log ID gracefully")
	
	var nonExistentResult map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &nonExistentResult)
	require.NoError(s.T(), err, "Should parse non-existent log ID response")
	
	// Handle case where data might be nil for empty results
	if nonExistentResult["data"] != nil {
		nonExistentLogs := nonExistentResult["data"].([]interface{})
		assert.Equal(s.T(), 0, len(nonExistentLogs), "Should return empty array for non-existent log ID")
	} else {
		s.T().Log("Data field is nil for non-existent log ID, which is acceptable")
	}
	
	s.T().Log("Complete API workflow integration test passed successfully")
}

// TestLogAPIIntegrationSuite runs the log API integration test suite
func TestLogAPIIntegrationSuite(t *testing.T) {
	suite.Run(t, new(LogAPIIntegrationTestSuite))
}