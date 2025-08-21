package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// SearchLogs searches for logs containing a keyword (removed pagination)
func SearchLogs(query, date, level string) ([]modelsv1.LogEntry, int, error) {
	var logFiles []string
	logDir := GetLogDirectory()

	if date != "" {
		// Search in specific date file
		logFiles = []string{filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))}
	} else {
		// Search in recent log files (last 7 days)
		for i := 0; i < 7; i++ {
			date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			logFiles = append(logFiles, filepath.Join(logDir, fmt.Sprintf("app-%s.log", date)))
		}
	}

	var allLogs []modelsv1.LogEntry
	queryLower := strings.ToLower(query)

	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			continue
		}

		// Configure options for optimized search operations
		options := ReadLogsOptions{
			MaxResults:             5000, // Reasonable limit for search operations
			ChunkSize:              800,  // Medium chunks for search processing
			EnableEarlyTermination: true, // Enable early termination for search performance
		}

		logs, _, err := ReadAndFilterLogsWithOptions(logFile, func(entry modelsv1.LogEntry) bool {
			// Check if query matches message, path, or other fields
			if strings.Contains(strings.ToLower(entry.Message), queryLower) ||
				strings.Contains(strings.ToLower(entry.Path), queryLower) ||
				strings.Contains(strings.ToLower(entry.LogID), queryLower) {
				
				if level != "" && !strings.EqualFold(entry.Level, level) {
					return false
				}
				return true
			}
			return false
		}, options)

		if err != nil {
			continue
		}

		allLogs = append(allLogs, logs...)
	}

	// Sort by timestamp (newest first)
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].Timestamp > allLogs[j].Timestamp
	})

	return allLogs, len(allLogs), nil
}

// ReadLogsByDate reads logs from file system for a specific date
func ReadLogsByDate(date, level string) ([]modelsv1.LogEntry, int, error) {
	logDir := GetLogDirectory()
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))

	// Configure options for date-based log reading
	options := ReadLogsOptions{
		MaxResults:             0,    // No limit for date-based searches
		ChunkSize:              1000, // Larger chunks for date processing
		EnableEarlyTermination: false, // Don't enable early termination for date searches
	}

	return ReadAndFilterLogsWithOptions(logFile, func(entry modelsv1.LogEntry) bool {
		if level != "" && !strings.EqualFold(entry.Level, level) {
			return false
		}
		return true
	}, options)
}

// ReadAndFilterLogs reads and filters logs from a file with enhanced error handling and performance optimizations
func ReadAndFilterLogs(filename string, filter func(modelsv1.LogEntry) bool) ([]modelsv1.LogEntry, int, error) {
	return ReadAndFilterLogsWithOptions(filename, filter, ReadLogsOptions{})
}

// ReadLogsOptions provides configuration options for log reading operations
type ReadLogsOptions struct {
	// MaxResults limits the number of results returned (0 = no limit)
	MaxResults int
	// ChunkSize defines the number of lines to process in each chunk (0 = use default)
	ChunkSize int
	// EnableEarlyTermination allows stopping when MaxResults is reached
	EnableEarlyTermination bool
}

// ReadAndFilterLogsWithOptions reads and filters logs from a file with performance optimizations
func ReadAndFilterLogsWithOptions(filename string, filter func(modelsv1.LogEntry) bool, options ReadLogsOptions) ([]modelsv1.LogEntry, int, error) {
	// Check if file exists first
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, 0, NewLogFileNotFoundError(filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, NewLogFileReadError(filename, fmt.Errorf("failed to open file: %w", err))
	}
	// Ensure proper resource cleanup with explicit defer
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the close error but don't override the main error
			// Note: We can't use Logger here as it might cause circular dependency
		}
	}()

	// Set default chunk size if not specified
	chunkSize := options.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 1000 // Default chunk size for processing
	}

	var allLogs []modelsv1.LogEntry
	var skippedLines int
	scanner := bufio.NewScanner(file)

	// Increase buffer size for large log lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// Process logs in chunks for better memory management
	var currentChunk []string
	lineCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		currentChunk = append(currentChunk, line)
		lineCount++

		// Process chunk when it reaches the specified size
		if len(currentChunk) >= chunkSize {
			processedLogs, skipped := processLogChunk(currentChunk, filter)
			allLogs = append(allLogs, processedLogs...)
			skippedLines += skipped

			// Clear the chunk to free memory
			currentChunk = currentChunk[:0]

			// Early termination if we have enough results and it's enabled
			if options.EnableEarlyTermination && options.MaxResults > 0 && len(allLogs) >= options.MaxResults {
				// Truncate to exact max results
				if len(allLogs) > options.MaxResults {
					allLogs = allLogs[:options.MaxResults]
				}
				break
			}
		}
	}

	// Process remaining lines in the final chunk
	if len(currentChunk) > 0 {
		processedLogs, skipped := processLogChunk(currentChunk, filter)
		allLogs = append(allLogs, processedLogs...)
		skippedLines += skipped

		// Apply max results limit if early termination is enabled
		if options.EnableEarlyTermination && options.MaxResults > 0 && len(allLogs) > options.MaxResults {
			allLogs = allLogs[:options.MaxResults]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, NewLogFileReadError(filename, fmt.Errorf("error reading file: %w", err))
	}

	// Log information about skipped lines if any (but don't fail the operation)
	if skippedLines > 0 {
		// Note: We can't use Logger here as it might cause circular dependency
		// The caller should handle logging if needed
	}

	return allLogs, len(allLogs), nil
}

// processLogChunk processes a chunk of log lines and returns filtered results
func processLogChunk(lines []string, filter func(modelsv1.LogEntry) bool) ([]modelsv1.LogEntry, int) {
	var logs []modelsv1.LogEntry
	var skippedLines int

	for _, line := range lines {
		entry, err := ParseLogLine(line)
		if err != nil {
			// Skip malformed JSON entries and continue processing gracefully
			skippedLines++
			continue
		}
		
		// Apply filter if provided, otherwise include all entries
		if filter == nil || (entry.Timestamp != "" && filter(entry)) {
			logs = append(logs, entry)
		}
	}

	return logs, skippedLines
}

// ParseLogLine parses a single JSON log line into LogEntry struct
func ParseLogLine(line string) (modelsv1.LogEntry, error) {
	// First try to parse as JSON
	var jsonEntry struct {
		Timestamp string `json:"timestamp"`
		LogID     string `json:"log_id"`
		Level     string `json:"level"`
		Message   string `json:"message"`
		File      string `json:"file"`
		Line      int    `json:"line"`
		Function  string `json:"function"`
	}

	if err := json.Unmarshal([]byte(line), &jsonEntry); err != nil {
		return modelsv1.LogEntry{}, fmt.Errorf("failed to parse JSON log line: %w", err)
	}

	// Create LogEntry with all JSON fields properly mapped
	entry := modelsv1.LogEntry{
		Timestamp: jsonEntry.Timestamp,
		LogID:     jsonEntry.LogID,
		Level:     jsonEntry.Level,
		Message:   jsonEntry.Message,
		File:      jsonEntry.File,
		Line:      jsonEntry.Line,
		Function:  jsonEntry.Function,
		Extra:     make(map[string]string),
	}

	// Parse HTTP request details from message field for backward compatibility
	message := jsonEntry.Message
	
	// Extract HTTP method and path
	if httpMatch := regexp.MustCompile(`(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+([^\s\|]+)`).FindStringSubmatch(message); len(httpMatch) > 2 {
		entry.Method = httpMatch[1]
		entry.Path = httpMatch[2]
	}

	// Extract status code
	if statusMatch := regexp.MustCompile(`Status:\s*(\d+)`).FindStringSubmatch(message); len(statusMatch) > 1 {
		if status, err := strconv.Atoi(statusMatch[1]); err == nil {
			entry.Status = status
		}
	}

	// Extract duration
	if durationMatch := regexp.MustCompile(`Duration:\s*([^\s\|]+)`).FindStringSubmatch(message); len(durationMatch) > 1 {
		entry.Duration = durationMatch[1]
	}

	// Extract IP
	if ipMatch := regexp.MustCompile(`IP:\s*([^\s\|]+)`).FindStringSubmatch(message); len(ipMatch) > 1 {
		entry.IP = ipMatch[1]
	}

	// Extract User-Agent
	if uaMatch := regexp.MustCompile(`User-Agent:\s*([^|\r\n]+)`).FindStringSubmatch(message); len(uaMatch) > 1 {
		entry.UserAgent = strings.TrimSpace(uaMatch[1])
	}

	return entry, nil
}

// GetLogDirectory returns the log directory path
func GetLogDirectory() string {
	// Adjust this path based on your log file location
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "./logs" // Default log directory
	}
	return logDir
}

// GetDateRange returns a slice of dates between start and end dates
func GetDateRange(startDate, endDate string) []string {
	var dates []string
	
	if startDate == "" && endDate == "" {
		// Default to last 7 days
		for i := 6; i >= 0; i-- {
			date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			dates = append(dates, date)
		}
		return dates
	}

	start, err1 := time.Parse("2006-01-02", startDate)
	end, err2 := time.Parse("2006-01-02", endDate)

	if err1 != nil || err2 != nil {
		// Default to today if parsing fails
		return []string{time.Now().Format("2006-01-02")}
	}

	if end.Before(start) {
		start, end = end, start // Swap if end is before start
	}

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates
}

// GetLast3Days returns a slice of date strings for the last 3 days (including today)
func GetLast3Days() []string {
	var dates []string
	
	// Get last 3 days including today (0, 1, 2 days ago)
	for i := 2; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates = append(dates, date)
	}
	
	return dates
}

// GetAvailableLogFiles scans the log directory and returns all available log files
// sorted by date for efficient processing (newest first)
func GetAvailableLogFiles() ([]string, error) {
	logDir := GetLogDirectory()
	
	// Check if log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return nil, NewLogFileReadError(logDir, fmt.Errorf("log directory does not exist"))
	}
	
	// Read all files in the log directory
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return nil, NewLogFileReadError(logDir, fmt.Errorf("failed to read log directory: %w", err))
	}
	
	var logFiles []string
	logPattern := regexp.MustCompile(`^app-(\d{4}-\d{2}-\d{2})\.log$`)
	
	// Filter for log files matching the expected pattern
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		filename := entry.Name()
		if logPattern.MatchString(filename) {
			fullPath := filepath.Join(logDir, filename)
			// Verify the file is readable before adding to the list
			if _, err := os.Stat(fullPath); err == nil {
				logFiles = append(logFiles, fullPath)
			}
		}
	}
	
	// Sort log files by date (newest first) for efficient processing
	sort.Slice(logFiles, func(i, j int) bool {
		// Extract dates from filenames for comparison
		dateI := extractDateFromLogFile(logFiles[i])
		dateJ := extractDateFromLogFile(logFiles[j])
		return dateI > dateJ // Newest first
	})
	
	return logFiles, nil
}

// extractDateFromLogFile extracts the date string from a log file path
func extractDateFromLogFile(filePath string) string {
	filename := filepath.Base(filePath)
	logPattern := regexp.MustCompile(`^app-(\d{4}-\d{2}-\d{2})\.log$`)
	matches := logPattern.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}