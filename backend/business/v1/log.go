package businessv1

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func ReadLogsByLogID(ctx context.Context, logID string) ([]modelsv1.LogEntry, error) {
	var allLogs []modelsv1.LogEntry
	logDir := utils.GetLogDirectory()
	if logDir == "" {
		return nil, fmt.Errorf("log directory not found")
	}
	
	// Get date range to search
	dates := utils.GetLast3Days()
	
	for _, date := range dates {
		logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))
		
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			fmt.Printf("Log file not found: %s\n", logFile)
			continue // Skip if log file doesn't exist
		}

		logs, _, err := utils.ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
			return entry.LogID == logID
		})

		if err != nil {
			fmt.Printf("Error reading log file: %s\n", logFile)
			continue // Skip problematic files
		}

		allLogs = append(allLogs, logs...)
	}

	// Sort by timestamp
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Timestamp < allLogs[j].Timestamp
	})

	return allLogs, nil
}

func ReadLogsByDate(date string) ([]modelsv1.LogEntry, int, error) {
	// Construct log file path (adjust based on your log file naming convention)
	logDir := utils.GetLogDirectory()
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil, 0, fmt.Errorf("log file not found: %s", logFile)
	}

	return utils.ReadAndFilterLogs(logFile, nil)
}