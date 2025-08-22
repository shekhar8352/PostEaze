package utils

import (
	"errors"
	"fmt"
)

// Custom error types for log API operations
var (
	ErrLogFileNotFound    = errors.New("log file not found")
	ErrLogDirectoryNotFound = errors.New("log directory not found")
	ErrInvalidDateFormat  = errors.New("invalid date format")
	ErrEmptyLogID         = errors.New("log ID is required")
	ErrEmptyDate          = errors.New("date is required")
)

// LogFileError represents an error related to log file operations
type LogFileError struct {
	FilePath string
	Err      error
	Type     string
}

func (e *LogFileError) Error() string {
	return fmt.Sprintf("log file error (%s): %s - %v", e.Type, e.FilePath, e.Err)
}

func (e *LogFileError) Unwrap() error {
	return e.Err
}

// NewLogFileNotFoundError creates a new log file not found error
func NewLogFileNotFoundError(filePath string) *LogFileError {
	return &LogFileError{
		FilePath: filePath,
		Err:      ErrLogFileNotFound,
		Type:     ErrorTypeNotFound,
	}
}

// NewLogFileReadError creates a new log file read error
func NewLogFileReadError(filePath string, err error) *LogFileError {
	return &LogFileError{
		FilePath: filePath,
		Err:      err,
		Type:     ErrorTypeInternal,
	}
}

// IsLogFileNotFoundError checks if the error is a log file not found error
func IsLogFileNotFoundError(err error) bool {
	var logFileErr *LogFileError
	if errors.As(err, &logFileErr) {
		return logFileErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsLogFileReadError checks if the error is a log file read error
func IsLogFileReadError(err error) bool {
	var logFileErr *LogFileError
	if errors.As(err, &logFileErr) {
		return logFileErr.Type == ErrorTypeInternal
	}
	return false
}