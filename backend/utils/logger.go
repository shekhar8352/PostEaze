package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// LogEntry represents a single log entry with all necessary metadata
type LogEntry struct {
	LogID     string    `json:"log_id"`
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Function  string    `json:"function,omitempty"`
}

// MarshalJSON provides custom JSON formatting for better readability
func (le LogEntry) MarshalJSON() ([]byte, error) {
	type Alias LogEntry
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Timestamp: le.Timestamp.Format(time.RFC3339),
		Alias:     (*Alias)(&le),
	})
}

// ContextKey type for context keys to avoid collisions
type ContextKey string

const LogIDKey ContextKey = "log_id"

// StructuredLogger handles structured logging with rotation and context support
type StructuredLogger struct {
	mu         sync.RWMutex
	logFile    *os.File
	currentDay string
	config     LoggerConfig
}

// LoggerConfig holds configuration options for the logger
type LoggerConfig struct {
	LogsDir        string
	FilePrefix     string
	IncludeConsole bool
	MinLevel       LogLevel
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		LogsDir:        "logs",
		FilePrefix:     "app",
		IncludeConsole: true,
		MinLevel:       LevelInfo,
	}
}

// Global logger instance
var Logger = NewLogger(DefaultConfig())

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config LoggerConfig) *StructuredLogger {
	l := &StructuredLogger{
		config: config,
	}
	l.rotateLogFile()
	return l
}

// AddLogID adds a unique log ID to the context
func AddLogID(ctx context.Context) context.Context {
	return context.WithValue(ctx, LogIDKey, uuid.New().String())
}

// GetLogID retrieves the log ID from context, returns empty string if not found
func GetLogID(ctx context.Context) string {
	if val, ok := ctx.Value(LogIDKey).(string); ok {
		return val
	}
	return ""
}

// rotateLogFile creates a new log file if the day has changed
func (l *StructuredLogger) rotateLogFile() {
	today := time.Now().Format("2006-01-02")
	
	// Check if rotation is needed
	if today == l.currentDay && l.logFile != nil {
		return
	}

	// Close existing file
	if l.logFile != nil {
		l.logFile.Close()
	}

	// Create logs directory
	if err := os.MkdirAll(l.config.LogsDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %v", err))
	}

	// Create new log file
	filename := fmt.Sprintf("%s-%s.log", l.config.FilePrefix, today)
	logPath := filepath.Join(l.config.LogsDir, filename)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	l.logFile = file
	l.currentDay = today
}

// shouldLog determines if a log level should be written based on minimum level
func (l *StructuredLogger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}
	
	minLevel, exists := levels[l.config.MinLevel]
	if !exists {
		minLevel = 1 // default to INFO
	}
	
	currentLevel, exists := levels[level]
	if !exists {
		return true // log unknown levels
	}
	
	return currentLevel >= minLevel
}

// writeLog writes the log entry to file and optionally console
func (l *StructuredLogger) writeLog(entry LogEntry) {
	if !l.shouldLog(entry.Level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Rotate log file if needed
	l.rotateLogFile()

	// Marshal to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		// Fallback to basic format if JSON marshaling fails
		data = []byte(fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"JSON marshal error: %v"}`,
			entry.Timestamp.Format(time.RFC3339), entry.Level, err))
	}

	// Write to file
	if l.logFile != nil {
		l.logFile.Write(append(data, '\n'))
	}

	// Write to console if configured
	if l.config.IncludeConsole {
		fmt.Println(string(data))
	}
}

// getCallerInfo extracts file, line, and function information from the call stack
func getCallerInfo(depth int) (string, int, string) {
	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "unknown", 0, "unknown"
	}
	
	filename := filepath.Base(file)
	funcName := "unknown"
	
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = fn.Name()
		// Extract just the function name without package path
		if lastSlash := strings.LastIndex(funcName, "/"); lastSlash >= 0 {
			funcName = funcName[lastSlash+1:]
		}
		if lastDot := strings.LastIndex(funcName, "."); lastDot >= 0 {
			funcName = funcName[lastDot+1:]
		}
	}
	
	return filename, line, funcName
}

// createLogEntry creates a structured log entry with caller information
func (l *StructuredLogger) createLogEntry(ctx context.Context, level LogLevel, msg string, args ...interface{}) LogEntry {
	// Format message if args provided
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	// Get caller information (skip 3 frames: createLogEntry -> specific method -> user code)
	file, line, function := getCallerInfo(3)

	return LogEntry{
		LogID:     GetLogID(ctx),
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		File:      file,
		Line:      line,
		Function:  function,
	}
}

// Context-aware logging methods
func (l *StructuredLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	entry := l.createLogEntry(ctx, LevelDebug, msg, args...)
	l.writeLog(entry)
}

func (l *StructuredLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	entry := l.createLogEntry(ctx, LevelInfo, msg, args...)
	l.writeLog(entry)
}

func (l *StructuredLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	entry := l.createLogEntry(ctx, LevelWarn, msg, args...)
	l.writeLog(entry)
}

func (l *StructuredLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	entry := l.createLogEntry(ctx, LevelError, msg, args...)
	l.writeLog(entry)
}

// Close gracefully closes the logger and any open file handles
func (l *StructuredLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// SetLevel dynamically changes the minimum log level
func (l *StructuredLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.MinLevel = level
}