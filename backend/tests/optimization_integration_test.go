package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// TestOptimizationIntegration tests that the optimized functions work correctly in integration
func TestOptimizationIntegration(t *testing.T) {
	// Create a temporary log directory for testing
	tempDir, err := os.MkdirTemp("", "optimization_integration_test")
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

	// Create test log files with known data
	testData := []struct {
		filename string
		logID    string
		count    int
	}{
		{"app-2024-01-01.log", "test-log-123", 10},
		{"app-2024-01-02.log", "test-log-456", 15},
		{"app-2024-01-03.log", "test-log-123", 5}, // Same log ID in different file
	}

	for _, td := range testData {
		logFile := filepath.Join(tempDir, td.filename)
		file, err := os.Create(logFile)
		if err != nil {
			t.Fatalf("Failed to create log file %s: %v", logFile, err)
		}

		// Write test log entries
		for i := 0; i < td.count; i++ {
			logEntry := fmt.Sprintf(`{"timestamp":"2024-01-01T10:%02d:00Z","log_id":"%s","level":"INFO","message":"Test message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
				i%60, td.logID, i, i+1)
			file.WriteString(logEntry + "\n")
		}
		
		// Add some other log entries with different log IDs
		for i := 0; i < 50; i++ {
			logEntry := fmt.Sprintf(`{"timestamp":"2024-01-01T10:%02d:00Z","log_id":"other-log-%d","level":"INFO","message":"Other message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
				i%60, i, i, i+1)
			file.WriteString(logEntry + "\n")
		}
		
		file.Close()
	}

	t.Run("TestBusinessLogicWithOptimizations", func(t *testing.T) {
		ctx := context.Background()
		
		// Test ReadLogsByLogID with optimizations
		logs, err := businessv1.ReadLogsByLogID(ctx, "test-log-123")
		if err != nil {
			t.Errorf("Failed to read logs by ID: %v", err)
			return
		}

		// Should find entries from both files with test-log-123
		expectedCount := 15 // 10 from first file + 5 from third file
		if len(logs) != expectedCount {
			t.Errorf("Expected %d logs for test-log-123, got %d", expectedCount, len(logs))
		}

		// Verify all returned logs have the correct log ID
		for _, log := range logs {
			if log.LogID != "test-log-123" {
				t.Errorf("Expected log ID test-log-123, got %s", log.LogID)
			}
		}

		t.Logf("Successfully found %d logs for log ID test-log-123", len(logs))
	})

	t.Run("TestUtilsWithOptimizations", func(t *testing.T) {
		// Test ReadLogsByDate with optimizations
		logs, total, err := businessv1.ReadLogsByDate("2024-01-01")
		if err != nil {
			t.Errorf("Failed to read logs by date: %v", err)
			return
		}

		// Should find all entries from the first file (10 + 50 = 60)
		expectedCount := 60
		if len(logs) != expectedCount || total != expectedCount {
			t.Errorf("Expected %d logs for date 2024-01-01, got %d (total: %d)", expectedCount, len(logs), total)
		}

		t.Logf("Successfully found %d logs for date 2024-01-01", len(logs))
	})

	t.Run("TestEarlyTerminationInPractice", func(t *testing.T) {
		// Create a large log file to test early termination
		largeLogFile := filepath.Join(tempDir, "app-2024-01-04.log")
		file, err := os.Create(largeLogFile)
		if err != nil {
			t.Fatalf("Failed to create large log file: %v", err)
		}

		// Write many entries with the target log ID at the beginning
		targetLogID := "early-term-test"
		for i := 0; i < 100; i++ {
			logEntry := fmt.Sprintf(`{"timestamp":"2024-01-04T10:%02d:00Z","log_id":"%s","level":"INFO","message":"Target message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
				i%60, targetLogID, i, i+1)
			file.WriteString(logEntry + "\n")
		}

		// Write many more entries with different log IDs
		for i := 0; i < 5000; i++ {
			logEntry := fmt.Sprintf(`{"timestamp":"2024-01-04T10:%02d:00Z","log_id":"other-log-%d","level":"INFO","message":"Other message %d","file":"test.go","line":%d,"function":"testFunc"}`, 
				i%60, i, i, i+1)
			file.WriteString(logEntry + "\n")
		}
		file.Close()

		// Test with early termination options
		options := utils.ReadLogsOptions{
			MaxResults:             20,   // Limit results
			ChunkSize:              50,   // Small chunks
			EnableEarlyTermination: true, // Enable early termination
		}

		logs, _, err := utils.ReadAndFilterLogsWithOptions(largeLogFile, func(entry modelsv1.LogEntry) bool {
			return entry.LogID == targetLogID
		}, options)

		if err != nil {
			t.Errorf("Failed to read logs with early termination: %v", err)
			return
		}

		if len(logs) != 20 {
			t.Errorf("Expected exactly 20 logs due to early termination, got %d", len(logs))
		}

		t.Logf("Early termination worked correctly: got %d logs as expected", len(logs))
	})
}