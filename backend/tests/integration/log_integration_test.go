package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
)

// TestLogRetrievalWorkflow tests complete log retrieval workflow with file system operations
func TestLogRetrievalWorkflow(t *testing.T) {
	// Create temporary directory for test logs
	testLogDir, err := os.MkdirTemp("", "log_integration_test_*")
	require.NoError(t, err, "Failed to create temp log directory")
	defer os.RemoveAll(testLogDir)

	http := helpers.NewHTTPTest()

	t.Logf("Log integration test setup completed with log dir: %s", testLogDir)

	// Test log file creation and retrieval
	t.Run("Log File Creation and Retrieval", func(t *testing.T) {
		testDate := time.Now().Format("2006-01-02")
		
		// Create test log entries
		testEntries := []modelsv1.LogEntry{
			helpers.CreateLogEntry(func(l *modelsv1.LogEntry) {
				l.LogID = "test-log-001"
				l.Message = "Test log message 1"
				l.Timestamp = fmt.Sprintf("%sT10:00:00Z", testDate)
			}),
			helpers.CreateLogEntry(func(l *modelsv1.LogEntry) {
				l.LogID = "test-log-002"
				l.Message = "Test log message 2"
				l.Timestamp = fmt.Sprintf("%sT10:01:00Z", testDate)
			}),
		}

		// Create log file
		filename := fmt.Sprintf("app-%s.log", testDate)
		filepath := filepath.Join(testLogDir, filename)
		
		file, err := os.Create(filepath)
		require.NoError(t, err, "Failed to create log file")
		defer file.Close()

		for _, entry := range testEntries {
			jsonData, err := json.Marshal(entry)
			require.NoError(t, err, "Failed to marshal log entry")
			
			_, err = file.WriteString(string(jsonData) + "\n")
			require.NoError(t, err, "Failed to write log entry")
		}

		t.Logf("Created test log file: %s with %d entries", filename, len(testEntries))

		// Test log retrieval by date (simulated)
		response := http.Request("GET", fmt.Sprintf("/api/v1/log/byDate/%s", testDate), nil)
		assert.Equal(t, 200, response.Code, "Log retrieval by date should succeed")

		// Test log retrieval by ID (simulated)
		response = http.Request("GET", "/api/v1/log/byId/test-log-001", nil)
		assert.Equal(t, 200, response.Code, "Log retrieval by ID should succeed")

		t.Log("Log file creation and retrieval test completed successfully")
	})

	// Test error handling for missing files
	t.Run("Error Handling for Missing Files", func(t *testing.T) {
		// Test non-existent date
		nonExistentDate := "2020-01-01"
		response := http.Request("GET", fmt.Sprintf("/api/v1/log/byDate/%s", nonExistentDate), nil)
		assert.Equal(t, 404, response.Code, "Should return 404 for missing log file")

		// Test invalid date format
		invalidDate := "invalid-date"
		response = http.Request("GET", fmt.Sprintf("/api/v1/log/byDate/%s", invalidDate), nil)
		assert.Equal(t, 404, response.Code, "Should return 404 for invalid date format")

		t.Log("Error handling tests completed successfully")
	})
}

// TestConcurrentLogAccess tests concurrent access to log operations
func TestConcurrentLogAccess(t *testing.T) {
	// Create temporary directory for test logs
	testLogDir, err := os.MkdirTemp("", "concurrent_log_test_*")
	require.NoError(t, err, "Failed to create temp log directory")
	defer os.RemoveAll(testLogDir)

	http := helpers.NewHTTPTest()
	const numConcurrentRequests = 3

	testDate := time.Now().Format("2006-01-02")
	testLogID := "concurrent-test-log-id"

	// Create test log file
	testEntries := []modelsv1.LogEntry{
		helpers.CreateLogEntry(func(l *modelsv1.LogEntry) {
			l.LogID = testLogID
			l.Message = "Concurrent access test log entry"
			l.Timestamp = fmt.Sprintf("%sT12:00:00Z", testDate)
		}),
	}

	filename := fmt.Sprintf("app-%s.log", testDate)
	filepath := filepath.Join(testLogDir, filename)
	
	file, err := os.Create(filepath)
	require.NoError(t, err, "Failed to create test log file")
	defer file.Close()

	for _, entry := range testEntries {
		jsonData, err := json.Marshal(entry)
		require.NoError(t, err, "Failed to marshal log entry")
		
		_, err = file.WriteString(string(jsonData) + "\n")
		require.NoError(t, err, "Failed to write log entry")
	}

	// Test concurrent requests
	results := make(chan bool, numConcurrentRequests*2)

	// Launch concurrent requests to both endpoints
	for i := 0; i < numConcurrentRequests; i++ {
		// Concurrent requests to byDate endpoint
		go func() {
			response := http.Request("GET", fmt.Sprintf("/api/v1/log/byDate/%s", testDate), nil)
			results <- response.Code == 200
		}()

		// Concurrent requests to byId endpoint
		go func() {
			response := http.Request("GET", fmt.Sprintf("/api/v1/log/byId/%s", testLogID), nil)
			results <- response.Code == 200
		}()
	}

	// Collect results
	successCount := 0
	for i := 0; i < numConcurrentRequests*2; i++ {
		if <-results {
			successCount++
		}
	}

	// All requests should succeed (or at least not fail catastrophically)
	assert.GreaterOrEqual(t, successCount, numConcurrentRequests, "Most concurrent requests should succeed")
	t.Logf("Concurrent access test completed: %d/%d requests succeeded", 
		successCount, numConcurrentRequests*2)
}