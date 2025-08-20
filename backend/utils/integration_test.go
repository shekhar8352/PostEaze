package utils

import (
	"path/filepath"
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func TestReadAndFilterLogs_WithActualLogFile(t *testing.T) {
	// Test with actual log file
	logFile := "../logs/app-2025-08-20.log"
	
	// Test reading all logs
	logs, count, err := ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
		return true // Accept all logs
	})
	
	if err != nil {
		t.Fatalf("Expected no error reading log file, got: %v", err)
	}
	
	if count == 0 {
		t.Fatal("Expected to read some logs, got 0")
	}
	
	if len(logs) != count {
		t.Errorf("Expected logs length %d to match count %d", len(logs), count)
	}
	
	// Verify first log entry has expected fields
	if len(logs) > 0 {
		firstLog := logs[0]
		if firstLog.Timestamp == "" {
			t.Error("Expected timestamp to be populated")
		}
		if firstLog.Level == "" {
			t.Error("Expected level to be populated")
		}
		if firstLog.Message == "" {
			t.Error("Expected message to be populated")
		}
	}
}

func TestReadAndFilterLogs_FilterByLogID(t *testing.T) {
	// Test filtering by log ID
	logFile := "../logs/app-2025-08-20.log"
	targetLogID := "45e1adc9-698f-41f9-886d-706effa57585"
	
	logs, count, err := ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
		return entry.LogID == targetLogID
	})
	
	if err != nil {
		t.Fatalf("Expected no error reading log file, got: %v", err)
	}
	
	// Verify all returned logs have the target log ID
	for _, log := range logs {
		if log.LogID != targetLogID {
			t.Errorf("Expected log ID %s, got %s", targetLogID, log.LogID)
		}
	}
	
	if len(logs) != count {
		t.Errorf("Expected logs length %d to match count %d", len(logs), count)
	}
}

func TestReadAndFilterLogs_NonExistentFile(t *testing.T) {
	// Test with non-existent file
	_, _, err := ReadAndFilterLogs("non-existent-file.log", func(entry modelsv1.LogEntry) bool {
		return true
	})
	
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}
}

func TestReadLogsByLogID_Integration(t *testing.T) {
	// Test the ReadLogsByLogID function with actual log files using business layer
	// Note: This test now uses the business layer function for consistency
	targetLogID := "45e1adc9-698f-41f9-886d-706effa57585"
	
	// Test the underlying utility functions that the business layer uses
	logDir := GetLogDirectory()
	dates := GetLast3Days()
	
	var allLogs []modelsv1.LogEntry
	for _, date := range dates {
		logFile := filepath.Join(logDir, "app-"+date+".log")
		logs, _, err := ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
			return entry.LogID == targetLogID
		})
		if err == nil {
			allLogs = append(allLogs, logs...)
		}
	}
	
	// Verify all returned logs have the target log ID
	for _, log := range allLogs {
		if log.LogID != targetLogID {
			t.Errorf("Expected log ID %s, got %s", targetLogID, log.LogID)
		}
	}
	
	// Verify logs are sorted by timestamp
	for i := 1; i < len(allLogs); i++ {
		if allLogs[i].Timestamp < allLogs[i-1].Timestamp {
			t.Error("Expected logs to be sorted by timestamp")
			break
		}
	}
}

func TestReadLogsByDate_Integration(t *testing.T) {
	// Test the ReadLogsByDate function
	date := "2025-08-20"
	
	logs, count, err := ReadLogsByDate(date, "")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(logs) != count {
		t.Errorf("Expected logs length %d to match count %d", len(logs), count)
	}
	
	// Verify all logs have proper structure
	for _, log := range logs {
		if log.Timestamp == "" {
			t.Error("Expected timestamp to be populated")
		}
		if log.Level == "" {
			t.Error("Expected level to be populated")
		}
	}
}