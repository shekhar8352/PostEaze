package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Validation errors
var (
	ErrInvalidLogIDFormat = errors.New("invalid log ID format")
	ErrLogIDTooLong       = errors.New("log ID too long")
	ErrInvalidDateRange   = errors.New("date is out of valid range")
)

// ValidateLogID validates the log ID parameter
func ValidateLogID(logID string) error {
	// Check if empty
	if strings.TrimSpace(logID) == "" {
		return ErrEmptyLogID
	}

	// Check length (reasonable limit to prevent abuse)
	if len(logID) > 100 {
		return ErrLogIDTooLong
	}

	// Check for valid characters (alphanumeric, hyphens, underscores, dots)
	// This prevents path traversal and other security issues
	validLogIDPattern := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	if !validLogIDPattern.MatchString(logID) {
		return ErrInvalidLogIDFormat
	}

	return nil
}

// ValidateDate validates the date parameter
func ValidateDate(dateStr string) error {
	// Check if empty
	if strings.TrimSpace(dateStr) == "" {
		return ErrEmptyDate
	}

	// Parse date in YYYY-MM-DD format
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ErrInvalidDateFormat
	}

	// Check if date is within reasonable range (not too far in past or future)
	now := time.Now()
	minDate := now.AddDate(-5, 0, 0) // 5 years ago
	maxDate := now.AddDate(0, 0, 1)  // 1 day in future

	if parsedDate.Before(minDate) || parsedDate.After(maxDate) {
		return ErrInvalidDateRange
	}

	return nil
}

// GetValidationErrorMessage returns user-friendly error messages for validation errors
func GetValidationErrorMessage(err error) string {
	switch {
	case errors.Is(err, ErrEmptyLogID):
		return "Log ID is required"
	case errors.Is(err, ErrInvalidLogIDFormat):
		return "Invalid log ID format. Only alphanumeric characters, hyphens, underscores, and dots are allowed"
	case errors.Is(err, ErrLogIDTooLong):
		return "Log ID is too long. Maximum length is 100 characters"
	case errors.Is(err, ErrEmptyDate):
		return "Date is required"
	case errors.Is(err, ErrInvalidDateFormat):
		return "Invalid date format. Expected YYYY-MM-DD"
	case errors.Is(err, ErrInvalidDateRange):
		return "Date is out of valid range. Please provide a date within the last 5 years"
	default:
		return "Invalid input"
	}
}