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
	var skippedFiles []string
	
	// Get all available log files instead of just last 3 days
	logFiles, err := utils.GetAvailableLogFiles()
	if err != nil {
		return nil, utils.NewLogFileReadError("", fmt.Errorf("failed to get available log files: %w", err))
	}
	
	// Configure options for optimized log ID search
	options := utils.ReadLogsOptions{
		MaxResults:             1000, // Reasonable limit for log ID searches
		ChunkSize:              500,  // Smaller chunks for better responsiveness
		EnableEarlyTermination: true, // Enable early termination for performance
	}
	
	// Search through all available log files (already sorted by date, newest first)
	for _, logFile := range logFiles {
		// Check if file exists before processing
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			skippedFiles = append(skippedFiles, logFile)
			continue // Skip if log file doesn't exist
		}

		// Use the optimized utility function with performance options
		logs, _, err := utils.ReadAndFilterLogsWithOptions(logFile, func(entry modelsv1.LogEntry) bool {
			return entry.LogID == logID
		}, options)

		if err != nil {
			// Log the error but continue processing other files gracefully
			utils.Logger.Warn(ctx, "Failed to read log file %s: %v", logFile, err)
			skippedFiles = append(skippedFiles, logFile)
			continue // Skip problematic files and continue processing others
		}

		allLogs = append(allLogs, logs...)
		
		// Early termination if we've found enough results across all files
		if len(allLogs) >= options.MaxResults {
			utils.Logger.Info(ctx, "Early termination: found %d log entries for log ID %s", len(allLogs), logID)
			break
		}
	}

	// Log information about skipped files if any
	if len(skippedFiles) > 0 {
		utils.Logger.Info(ctx, "Skipped %d log files due to errors or missing files", len(skippedFiles))
	}

	// Sort by timestamp using consistent sorting logic
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Timestamp < allLogs[j].Timestamp
	})

	return allLogs, nil
}

func ReadLogsByDate(date string) ([]modelsv1.LogEntry, int, error) {
	// Use utility function to get log directory for consistency
	logDir := utils.GetLogDirectory()
	if logDir == "" {
		return nil, 0, utils.NewLogFileReadError("", fmt.Errorf("log directory not configured"))
	}
	
	// Construct log file path using consistent naming convention
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil, 0, utils.NewLogFileNotFoundError(logFile)
	}

	// Configure options for optimized date-based log reading
	options := utils.ReadLogsOptions{
		MaxResults:             0,    // No limit for date-based searches (return all logs for the date)
		ChunkSize:              1000, // Larger chunks for date-based processing
		EnableEarlyTermination: false, // Don't enable early termination for date searches
	}

	// Use the optimized utility function with performance options
	logs, total, err := utils.ReadAndFilterLogsWithOptions(logFile, nil, options)
	if err != nil {
		return nil, 0, utils.NewLogFileReadError(logFile, err)
	}
	
	return logs, total, nil
}