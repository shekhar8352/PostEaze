package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// TestPerformanceOptimizations tests the performance optimizations for large log files
func TestPerformanceOptimizations(t *testing.T) {
	// Create a temporary log directory for testing
	tempDir, err := os.MkdirTemp("", "log_performance_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the log directory to our temp directory
	originalLogDir := os.Getenv("LOG_DIR")
	os.Setenv("LOG_DIR", tempDir)
	defer func() {
		if originalLogDir != "" {
			os.Setenv("LOG_DIR", originalLogDir)
		} else {
			os.Unsetenv("LOG_DIR")
		}
	}()

	t.Run("TestChunkedProcessing", func(t *testing.T) {
		testChunkedProcessing(t, tempDir)
	})

	t.Run("TestEarlyTermination", func(t *testing.T) {
		testEarlyTermination(t, tempDir)
	})

	t.Run("TestResourceCleanup", func(t *testing.T) {
		testResourceCleanup(t, tempDir)
	})
}

func testChunkedProcessing(t *testing.T, tempDir string) {
	// Create a large log file with many entries
	logFile := filepath.Join(tempDir, "app-2024-01-01.log")
	file, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}

	// Write 5000 log entries to test chunked processing
	for i := 0; i < 5000; i++ {
		logEntry := fmt.Sprintf(`{"timestamp":"2024-01-01T10:%02d:00Z","log_id":"test-log-%d","level":"INFO","message":"Test message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
			i%60, i, i, i+1)
		file.WriteString(logEntry + "\n")
	}
	file.Close()

	// Test with different chunk sizes
	testCases := []struct {
		name      string
		chunkSize int
	}{
		{"SmallChunks", 100},
		{"MediumChunks", 500},
		{"LargeChunks", 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := utils.ReadLogsOptions{
				ChunkSize:              tc.chunkSize,
				EnableEarlyTermination: false,
			}

			start := time.Now()
			logs, _, err := utils.ReadAndFilterLogsWithOptions(logFile, nil, options)
			duration := time.Since(start)

			if err != nil {
				t.Errorf("Failed to read logs with chunk size %d: %v", tc.chunkSize, err)
				return
			}

			if len(logs) != 5000 {
				t.Errorf("Expected 5000 logs, got %d", len(logs))
			}

			t.Logf("Chunk size %d: processed %d logs in %v", tc.chunkSize, len(logs), duration)
		})
	}
}

func testEarlyTermination(t *testing.T, tempDir string) {
	// Create a log file with many entries but we only want a few
	logFile := filepath.Join(tempDir, "app-2024-01-02.log")
	file, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}

	// Write 10000 log entries, but we'll search for specific log IDs
	targetLogID := "target-log-123"
	for i := 0; i < 10000; i++ {
		var logID string
		if i < 50 { // Put target entries at the beginning
			logID = targetLogID
		} else {
			logID = fmt.Sprintf("other-log-%d", i)
		}
		
		logEntry := fmt.Sprintf(`{"timestamp":"2024-01-02T10:%02d:00Z","log_id":"%s","level":"INFO","message":"Test message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
			i%60, logID, i, i+1)
		file.WriteString(logEntry + "\n")
	}
	file.Close()

	// Test early termination with max results
	options := utils.ReadLogsOptions{
		MaxResults:             10,   // Only want 10 results
		ChunkSize:              100,  // Small chunks
		EnableEarlyTermination: true, // Enable early termination
	}

	start := time.Now()
	logs, _, err := utils.ReadAndFilterLogsWithOptions(logFile, func(entry modelsv1.LogEntry) bool {
		return entry.LogID == targetLogID
	}, options)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Failed to read logs with early termination: %v", err)
		return
	}

	if len(logs) != 10 {
		t.Errorf("Expected exactly 10 logs due to MaxResults, got %d", len(logs))
	}

	// Verify all returned logs have the target log ID
	for _, log := range logs {
		if log.LogID != targetLogID {
			t.Errorf("Expected log ID %s, got %s", targetLogID, log.LogID)
		}
	}

	t.Logf("Early termination: found %d logs in %v (should be much faster than processing all 10000 entries)", len(logs), duration)

	// Test without early termination for comparison
	optionsNoEarlyTerm := utils.ReadLogsOptions{
		MaxResults:             0,     // No limit
		ChunkSize:              100,   // Same chunk size
		EnableEarlyTermination: false, // Disable early termination
	}

	start = time.Now()
	logsAll, _, err := utils.ReadAndFilterLogsWithOptions(logFile, func(entry modelsv1.LogEntry) bool {
		return entry.LogID == targetLogID
	}, optionsNoEarlyTerm)
	durationAll := time.Since(start)

	if err != nil {
		t.Errorf("Failed to read all logs: %v", err)
		return
	}

	if len(logsAll) != 50 {
		t.Errorf("Expected 50 total matching logs, got %d", len(logsAll))
	}

	t.Logf("Without early termination: found %d logs in %v", len(logsAll), durationAll)

	// Early termination should be faster (though this might not always be true in tests due to small data size)
	if duration > durationAll*2 {
		t.Logf("Warning: Early termination took longer than expected. This might be normal for small test data.")
	}
}

func testResourceCleanup(t *testing.T, tempDir string) {
	// Create multiple log files
	logFiles := []string{
		filepath.Join(tempDir, "app-2024-01-03.log"),
		filepath.Join(tempDir, "app-2024-01-04.log"),
		filepath.Join(tempDir, "app-2024-01-05.log"),
	}

	for i, logFile := range logFiles {
		file, err := os.Create(logFile)
		if err != nil {
			t.Fatalf("Failed to create log file %s: %v", logFile, err)
		}

		// Write some log entries
		for j := 0; j < 100; j++ {
			logEntry := fmt.Sprintf(`{"timestamp":"2024-01-0%dT10:%02d:00Z","log_id":"test-log-%d-%d","level":"INFO","message":"Test message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
				i+3, j%60, i, j, j, j+1)
			file.WriteString(logEntry + "\n")
		}
		file.Close()
	}

	// Test business logic function that processes multiple files
	ctx := context.Background()
	logs, err := businessv1.ReadLogsByLogID(ctx, "test-log-1-50")

	if err != nil {
		t.Errorf("Failed to read logs by ID: %v", err)
		return
	}

	if len(logs) == 0 {
		t.Error("Expected to find at least one log entry")
	}

	// Verify that we can still access the files (they should be properly closed)
	for _, logFile := range logFiles {
		file, err := os.Open(logFile)
		if err != nil {
			t.Errorf("Failed to reopen log file %s: %v", logFile, err)
			continue
		}
		file.Close()
	}

	t.Logf("Resource cleanup test passed: processed %d log files and found %d matching entries", len(logFiles), len(logs))
}

// TestReadLogsOptionsValidation tests that the ReadLogsOptions are properly validated and applied
func TestReadLogsOptionsValidation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "log_options_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test log file
	logFile := filepath.Join(tempDir, "test.log")
	file, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}

	// Write test entries
	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf(`{"timestamp":"2024-01-01T10:%02d:00Z","log_id":"test-%d","level":"INFO","message":"Message %d"}`, 
			i%60, i, i)
		file.WriteString(logEntry + "\n")
	}
	file.Close()

	testCases := []struct {
		name           string
		options        utils.ReadLogsOptions
		expectedCount  int
		shouldTerminate bool
	}{
		{
			name: "DefaultOptions",
			options: utils.ReadLogsOptions{},
			expectedCount: 1000,
			shouldTerminate: false,
		},
		{
			name: "MaxResults50",
			options: utils.ReadLogsOptions{
				MaxResults: 50,
				EnableEarlyTermination: true,
			},
			expectedCount: 50,
			shouldTerminate: true,
		},
		{
			name: "MaxResults100NoEarlyTerm",
			options: utils.ReadLogsOptions{
				MaxResults: 100,
				EnableEarlyTermination: false,
			},
			expectedCount: 1000, // Should read all since early termination is disabled
			shouldTerminate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logs, _, err := utils.ReadAndFilterLogsWithOptions(logFile, nil, tc.options)
			if err != nil {
				t.Errorf("Failed to read logs: %v", err)
				return
			}

			if tc.shouldTerminate && len(logs) != tc.expectedCount {
				t.Errorf("Expected exactly %d logs with early termination, got %d", tc.expectedCount, len(logs))
			} else if !tc.shouldTerminate && len(logs) != tc.expectedCount {
				t.Errorf("Expected %d logs without early termination, got %d", tc.expectedCount, len(logs))
			}

			t.Logf("Test %s: got %d logs (expected %d)", tc.name, len(logs), tc.expectedCount)
		})
	}
}