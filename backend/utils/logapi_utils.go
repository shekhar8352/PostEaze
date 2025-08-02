package utils

import (
	"bufio"
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

		logs, _, err := ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
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
		})

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

	return ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
		if level != "" && !strings.EqualFold(entry.Level, level) {
			return false
		}
		return true
	})
}

// ReadLogsByLogID reads all logs for a specific log ID (last 3 days)
func ReadLogsByLogID(logID string) ([]modelsv1.LogEntry, error) {
	var allLogs []modelsv1.LogEntry
	logDir := GetLogDirectory()

	// Get last 3 days (including today)
	dates := GetLast3Days()
	
	for _, date := range dates {
		logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", date))
		
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			continue // Skip if log file doesn't exist
		}

		logs, _, err := ReadAndFilterLogs(logFile, func(entry modelsv1.LogEntry) bool {
			return entry.LogID == logID
		})

		if err != nil {
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

// ReadAndFilterLogs reads and filters logs from a file
func ReadAndFilterLogs(filename string, filter func(modelsv1.LogEntry) bool) ([]modelsv1.LogEntry, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var allLogs []modelsv1.LogEntry
	scanner := bufio.NewScanner(file)

	// Increase buffer size for large log lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// Regex patterns for parsing different log formats
	logRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?)\s+(\w+)\s+(.+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry := ParseLogLine(line, logRegex)
		if entry.Timestamp != "" && filter(entry) {
			allLogs = append(allLogs, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return allLogs, len(allLogs), nil
}

// ParseLogLine parses a single log line into LogEntry struct
func ParseLogLine(line string, regex *regexp.Regexp) modelsv1.LogEntry {
	matches := regex.FindStringSubmatch(line)
	if len(matches) < 4 {
		return modelsv1.LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     "UNKNOWN",
			Message:   line,
			Extra:     make(map[string]string),
		}
	}

	entry := modelsv1.LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
		Extra:     make(map[string]string),
	}

	// Parse additional fields from message
	message := matches[3]
	
	// Extract LogID - try multiple patterns
	if logIDMatch := regexp.MustCompile(`LogID:\s*([A-Za-z0-9\-_]+)`).FindStringSubmatch(message); len(logIDMatch) > 1 {
		entry.LogID = logIDMatch[1]
	} else if logIDMatch := regexp.MustCompile(`logid[:\s]+([A-Za-z0-9\-_]+)`).FindStringSubmatch(strings.ToLower(message)); len(logIDMatch) > 1 {
		entry.LogID = logIDMatch[1]
	} else if logIDMatch := regexp.MustCompile(`log_id[:\s]+([A-Za-z0-9\-_]+)`).FindStringSubmatch(strings.ToLower(message)); len(logIDMatch) > 1 {
		entry.LogID = logIDMatch[1]
	}

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

	return entry
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