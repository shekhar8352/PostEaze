package utils

import (
	"testing"
	"time"
)

func TestValidateLogID(t *testing.T) {
	tests := []struct {
		name    string
		logID   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid log ID with alphanumeric",
			logID:   "abc123",
			wantErr: false,
		},
		{
			name:    "valid log ID with hyphens",
			logID:   "log-id-123",
			wantErr: false,
		},
		{
			name:    "valid log ID with underscores",
			logID:   "log_id_123",
			wantErr: false,
		},
		{
			name:    "valid log ID with dots",
			logID:   "log.id.123",
			wantErr: false,
		},
		{
			name:    "empty log ID",
			logID:   "",
			wantErr: true,
			errType: ErrEmptyLogID,
		},
		{
			name:    "whitespace only log ID",
			logID:   "   ",
			wantErr: true,
			errType: ErrEmptyLogID,
		},
		{
			name:    "log ID with invalid characters",
			logID:   "log@id#123",
			wantErr: true,
			errType: ErrInvalidLogIDFormat,
		},
		{
			name:    "log ID with spaces",
			logID:   "log id 123",
			wantErr: true,
			errType: ErrInvalidLogIDFormat,
		},
		{
			name:    "log ID with path traversal attempt",
			logID:   "../../../etc/passwd",
			wantErr: true,
			errType: ErrInvalidLogIDFormat,
		},
		{
			name:    "log ID too long",
			logID:   string(make([]byte, 101)), // 101 characters
			wantErr: true,
			errType: ErrLogIDTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogID(tt.logID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil && err != tt.errType {
				t.Errorf("ValidateLogID() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
		errType error
	}{
		{
			name:    "valid date",
			date:    "2024-01-15",
			wantErr: false,
		},
		{
			name:    "valid date today",
			date:    time.Now().Format("2006-01-02"),
			wantErr: false,
		},
		{
			name:    "empty date",
			date:    "",
			wantErr: true,
			errType: ErrEmptyDate,
		},
		{
			name:    "whitespace only date",
			date:    "   ",
			wantErr: true,
			errType: ErrEmptyDate,
		},
		{
			name:    "invalid date format - wrong separator",
			date:    "2024/01/15",
			wantErr: true,
			errType: ErrInvalidDateFormat,
		},
		{
			name:    "invalid date format - missing day",
			date:    "2024-01",
			wantErr: true,
			errType: ErrInvalidDateFormat,
		},
		{
			name:    "invalid date format - wrong order",
			date:    "01-15-2024",
			wantErr: true,
			errType: ErrInvalidDateFormat,
		},
		{
			name:    "invalid date - February 30th",
			date:    "2024-02-30",
			wantErr: true,
			errType: ErrInvalidDateFormat,
		},
		{
			name:    "date too far in past",
			date:    "2010-01-01",
			wantErr: true,
			errType: ErrInvalidDateRange,
		},
		{
			name:    "date too far in future",
			date:    time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
			wantErr: true,
			errType: ErrInvalidDateRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDate(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil && err != tt.errType {
				t.Errorf("ValidateDate() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestGetValidationErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "empty log ID error",
			err:      ErrEmptyLogID,
			expected: "Log ID is required",
		},
		{
			name:     "invalid log ID format error",
			err:      ErrInvalidLogIDFormat,
			expected: "Invalid log ID format. Only alphanumeric characters, hyphens, underscores, and dots are allowed",
		},
		{
			name:     "log ID too long error",
			err:      ErrLogIDTooLong,
			expected: "Log ID is too long. Maximum length is 100 characters",
		},
		{
			name:     "empty date error",
			err:      ErrEmptyDate,
			expected: "Date is required",
		},
		{
			name:     "invalid date format error",
			err:      ErrInvalidDateFormat,
			expected: "Invalid date format. Expected YYYY-MM-DD",
		},
		{
			name:     "invalid date range error",
			err:      ErrInvalidDateRange,
			expected: "Date is out of valid range. Please provide a date within the last 5 years",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetValidationErrorMessage(tt.err)
			if result != tt.expected {
				t.Errorf("GetValidationErrorMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}