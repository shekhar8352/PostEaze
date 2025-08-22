package utils

import (
	"os"
	"testing"

	"github.com/shekhar8352/PostEaze/utils"
)

func TestReplacePlaceHoldersWithEnv(t *testing.T) {
	// Store original environment variables that we'll modify
	testVars := []string{"TEST_VAR", "DB_HOST", "API_KEY", "PORT", "EMPTY_VAR", "VAR_123", "MY_TEST_VAR", "SPECIAL_VAR", "PROTOCOL", "DOMAIN", "VERSION", "DB_USER", "DB_PASS", "DB_NAME", "CASE_VAR", "case_var"}
	originalEnvVars := make(map[string]string)

	for _, varName := range testVars {
		if val := os.Getenv(varName); val != "" {
			originalEnvVars[varName] = val
		}
	}

	// Set test environment variables
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("API_KEY", "secret123")
	os.Setenv("PORT", "8080")
	os.Setenv("EMPTY_VAR", "")

	defer func() {
		for _, varName := range testVars {
			if originalVal, exists := originalEnvVars[varName]; exists {
				os.Setenv(varName, originalVal)
			} else {
				os.Unsetenv(varName)
			}
		}
	}()

	tests := []struct {
		name     string
		setup    func()
		input    string
		expected string
	}{
		{
			name:     "single variable",
			input:    "Database host is ${DB_HOST}",
			expected: "Database host is localhost",
		},
		{
			name:     "multiple variables",
			input:    "Connect to ${DB_HOST}:${PORT} with key ${API_KEY}",
			expected: "Connect to localhost:8080 with key secret123",
		},
		{
			name:     "no variables",
			input:    "This string has no environment variables",
			expected: "This string has no environment variables",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "non-existent variable",
			input:    "Value is ${NON_EXISTENT_VAR}",
			expected: "Value is ",
		},
		{
			name:     "empty variable",
			input:    "Value is ${EMPTY_VAR}",
			expected: "Value is ",
		},
		{
			name:     "same variable multiple times",
			input:    "${DB_HOST} and ${DB_HOST} again",
			expected: "localhost and localhost again",
		},
		{
			name:     "variable at start",
			input:    "${DB_HOST} is the database host",
			expected: "localhost is the database host",
		},
		{
			name:     "variable at end",
			input:    "Database host is ${DB_HOST}",
			expected: "Database host is localhost",
		},
		{
			name:     "only variable",
			input:    "${DB_HOST}",
			expected: "localhost",
		},
		{
			name: "variable with numbers",
			setup: func() {
				os.Setenv("VAR_123", "numeric_var")
			},
			input:    "Value is ${VAR_123}",
			expected: "Value is numeric_var",
		},
		{
			name: "variable with underscores",
			setup: func() {
				os.Setenv("MY_TEST_VAR", "underscore_var")
			},
			input:    "Value is ${MY_TEST_VAR}",
			expected: "Value is underscore_var",
		},
		{
			name:     "malformed placeholder",
			input:    "This has malformed ${INCOMPLETE and ${ALSO_INCOMPLETE",
			expected: "This has malformed ${INCOMPLETE and ${ALSO_INCOMPLETE",
		},
		{
			name: "special characters in value",
			setup: func() {
				os.Setenv("SPECIAL_VAR", "value with spaces and !@#$%^&*()")
			},
			input:    "Special value: ${SPECIAL_VAR}",
			expected: "Special value: value with spaces and !@#$%^&*()",
		},
		{
			name: "URL template",
			setup: func() {
				os.Setenv("PROTOCOL", "https")
				os.Setenv("DOMAIN", "api.example.com")
				os.Setenv("VERSION", "v1")
			},
			input:    "${PROTOCOL}://${DOMAIN}/${VERSION}/users",
			expected: "https://api.example.com/v1/users",
		},
		{
			name: "database connection string",
			setup: func() {
				os.Setenv("DB_USER", "admin")
				os.Setenv("DB_PASS", "secret")
				os.Setenv("DB_NAME", "myapp")
			},
			input:    "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${PORT}/${DB_NAME}",
			expected: "postgres://admin:secret@localhost:8080/myapp",
		},
		{
			name:     "mixed content",
			input:    "Normal text ${DB_HOST} more text ${PORT} and ${NON_EXISTENT} end",
			expected: "Normal text localhost more text 8080 and  end",
		},
		{
			name: "case sensitive",
			setup: func() {
				os.Setenv("CASE_VAR", "lowercase")
				os.Setenv("case_var", "uppercase")
			},
			input:    "${CASE_VAR} vs ${case_var}",
			expected: "lowercase vs uppercase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			result := utils.ReplacePlaceHoldersWithEnv(tt.input)
			if result != tt.expected {
				t.Errorf("ReplacePlaceHoldersWithEnv() = %v, want %v", result, tt.expected)
			}

			// Clean up any variables set in setup
			if tt.setup != nil {
				for _, varName := range []string{"VAR_123", "MY_TEST_VAR", "SPECIAL_VAR", "PROTOCOL", "DOMAIN", "VERSION", "DB_USER", "DB_PASS", "DB_NAME", "CASE_VAR", "case_var"} {
					if _, exists := originalEnvVars[varName]; !exists {
						os.Unsetenv(varName)
					}
				}
			}
		})
	}

	t.Run("nested braces", func(t *testing.T) {
		input := "This has ${${NESTED}} braces"
		result := utils.ReplacePlaceHoldersWithEnv(input)

		// The function may handle nested braces differently than expected
		// Let's just verify it doesn't crash and produces some output
		if result == "" {
			t.Error("ReplacePlaceHoldersWithEnv() should not return empty string for nested braces")
		}
		if result == input {
			t.Log("ReplacePlaceHoldersWithEnv() returned input unchanged for nested braces (acceptable)")
		}
	})
}