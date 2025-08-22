package utils

import (
	"os"
	"testing"

	"github.com/shekhar8352/PostEaze/utils/env"
)

func TestInitEnv(t *testing.T) {
	// Store original environment variables that we'll modify
	originalEnvVars := make(map[string]string)
	testVars := []string{"TEST_INIT_VAR", "ANOTHER_VAR", "TEST_DOTENV_VAR", "ANOTHER_DOTENV_VAR"}

	for _, varName := range testVars {
		if val := os.Getenv(varName); val != "" {
			originalEnvVars[varName] = val
		}
	}

	// Cleanup function
	defer func() {
		for _, varName := range testVars {
			if originalVal, exists := originalEnvVars[varName]; exists {
				os.Setenv(varName, originalVal)
			} else {
				os.Unsetenv(varName)
			}
		}
	}()

	t.Run("loads environment variables", func(t *testing.T) {
		// Set some environment variables before initialization
		os.Setenv("TEST_INIT_VAR", "init_value")
		os.Setenv("ANOTHER_VAR", "another_value")

		// Initialize environment
		env.InitEnv()

		// Test that variables are loaded
		result := env.ApplyEnvironmentToString("Value is ${TEST_INIT_VAR}")
		if result != "Value is init_value" {
			t.Errorf("ApplyEnvironmentToString() = %v, want 'Value is init_value'", result)
		}

		result = env.ApplyEnvironmentToString("Another is ${ANOTHER_VAR}")
		if result != "Another is another_value" {
			t.Errorf("ApplyEnvironmentToString() = %v, want 'Another is another_value'", result)
		}
	})

	t.Run("handles empty environment", func(t *testing.T) {
		// Clear some environment variables
		os.Unsetenv("TEST_ENV_VAR")
		os.Unsetenv("DB_HOST")

		// Initialize environment
		env.InitEnv()

		// Should handle missing variables gracefully
		result := env.ApplyEnvironmentToString("Value is ${TEST_ENV_VAR}")
		if result != "Value is " {
			t.Errorf("ApplyEnvironmentToString() = %v, want 'Value is '", result)
		}
	})

	t.Run("loads dotenv file", func(t *testing.T) {
		// Change to a temporary directory and create a .env file
		originalDir, _ := os.Getwd()
		tempDir, err := os.MkdirTemp("", "env_test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		os.Chdir(tempDir)
		defer os.Chdir(originalDir)

		// Create .env file in temp directory
		envContent := `TEST_DOTENV_VAR=dotenv_value
ANOTHER_DOTENV_VAR=another_dotenv_value`

		err = os.WriteFile(".env", []byte(envContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write .env file: %v", err)
		}

		// Initialize environment (should load .env file)
		env.InitEnv()

		// Test that .env variables are loaded
		result := env.ApplyEnvironmentToString("Value is ${TEST_DOTENV_VAR}")
		if result != "Value is dotenv_value" {
			t.Errorf("ApplyEnvironmentToString() = %v, want 'Value is dotenv_value'", result)
		}
	})

	t.Run("handles invalid dotenv file", func(t *testing.T) {
		originalDir, _ := os.Getwd()
		tempDir, err := os.MkdirTemp("", "env_test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		os.Chdir(tempDir)
		defer os.Chdir(originalDir)

		// Create invalid .env file
		invalidEnvContent := `INVALID LINE WITHOUT EQUALS
=INVALID_EMPTY_KEY
VALID_VAR=valid_value`

		err = os.WriteFile(".env", []byte(invalidEnvContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write .env file: %v", err)
		}

		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("InitEnv() panicked: %v", r)
			}
		}()

		env.InitEnv()
	})
}

func TestApplyEnvironmentToString(t *testing.T) {
	// Store and restore environment variables
	testVars := []string{"SINGLE_VAR", "VAR1", "VAR2", "VAR3", "REPEAT_VAR", "START_VAR", "END_VAR", "ONLY_VAR", "SPECIAL_TEST_VAR", "PROTOCOL", "DOMAIN", "VERSION", "DB_USER", "DB_PASS", "DB_HOST_TEST", "DB_PORT", "DB_NAME", "OUTER_VAR", "INNER_VAR", "CASE_VAR", "case_var", "EQUALS_VAR", "MULTI_EQUALS", "WHITESPACE_VAR", "NEWLINE_VAR", "EMPTY_TEST_VAR"}
	originalEnvVars := make(map[string]string)

	for _, varName := range testVars {
		if val := os.Getenv(varName); val != "" {
			originalEnvVars[varName] = val
		}
	}

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
			name: "single variable",
			setup: func() {
				os.Setenv("SINGLE_VAR", "single_value")
				env.InitEnv()
			},
			input:    "The value is ${SINGLE_VAR}",
			expected: "The value is single_value",
		},
		{
			name: "multiple variables",
			setup: func() {
				os.Setenv("VAR1", "value1")
				os.Setenv("VAR2", "value2")
				os.Setenv("VAR3", "value3")
				env.InitEnv()
			},
			input:    "${VAR1} and ${VAR2} and ${VAR3}",
			expected: "value1 and value2 and value3",
		},
		{
			name: "no variables",
			setup: func() {
				env.InitEnv()
			},
			input:    "This string has no environment variables",
			expected: "This string has no environment variables",
		},
		{
			name: "empty string",
			setup: func() {
				env.InitEnv()
			},
			input:    "",
			expected: "",
		},
		{
			name: "non-existent variable",
			setup: func() {
				env.InitEnv()
			},
			input:    "Value is ${NON_EXISTENT_VAR}",
			expected: "Value is ",
		},
		{
			name: "empty variable",
			setup: func() {
				os.Setenv("EMPTY_TEST_VAR", "")
				env.InitEnv()
			},
			input:    "Value is ${EMPTY_TEST_VAR}",
			expected: "Value is ",
		},
		{
			name: "same variable multiple times",
			setup: func() {
				os.Setenv("REPEAT_VAR", "repeated")
				env.InitEnv()
			},
			input:    "${REPEAT_VAR} and ${REPEAT_VAR} again",
			expected: "repeated and repeated again",
		},
		{
			name: "variable at start",
			setup: func() {
				os.Setenv("START_VAR", "start")
				env.InitEnv()
			},
			input:    "${START_VAR} is at the beginning",
			expected: "start is at the beginning",
		},
		{
			name: "variable at end",
			setup: func() {
				os.Setenv("END_VAR", "end")
				env.InitEnv()
			},
			input:    "This is at the ${END_VAR}",
			expected: "This is at the end",
		},
		{
			name: "only variable",
			setup: func() {
				os.Setenv("ONLY_VAR", "only")
				env.InitEnv()
			},
			input:    "${ONLY_VAR}",
			expected: "only",
		},
		{
			name: "special characters in value",
			setup: func() {
				os.Setenv("SPECIAL_TEST_VAR", "value with spaces and !@#$%^&*()")
				env.InitEnv()
			},
			input:    "Special: ${SPECIAL_TEST_VAR}",
			expected: "Special: value with spaces and !@#$%^&*()",
		},
		{
			name: "URL template",
			setup: func() {
				os.Setenv("PROTOCOL", "https")
				os.Setenv("DOMAIN", "api.example.com")
				os.Setenv("VERSION", "v1")
				env.InitEnv()
			},
			input:    "${PROTOCOL}://${DOMAIN}/${VERSION}/users",
			expected: "https://api.example.com/v1/users",
		},
		{
			name: "database connection string",
			setup: func() {
				os.Setenv("DB_USER", "admin")
				os.Setenv("DB_PASS", "secret")
				os.Setenv("DB_HOST_TEST", "localhost")
				os.Setenv("DB_PORT", "5432")
				os.Setenv("DB_NAME", "myapp")
				env.InitEnv()
			},
			input:    "postgres://${DB_USER}:${DB_PASS}@${DB_HOST_TEST}:${DB_PORT}/${DB_NAME}",
			expected: "postgres://admin:secret@localhost:5432/myapp",
		},
		{
			name: "nested replacements",
			setup: func() {
				os.Setenv("OUTER_VAR", "${INNER_VAR}")
				os.Setenv("INNER_VAR", "inner_value")
				env.InitEnv()
			},
			input:    "Value is ${OUTER_VAR}",
			expected: "Value is ${INNER_VAR}", // Should not resolve nested variables
		},
		{
			name: "case sensitive",
			setup: func() {
				os.Setenv("CASE_VAR", "lowercase")
				os.Setenv("case_var", "uppercase")
				env.InitEnv()
			},
			input:    "${CASE_VAR} vs ${case_var}",
			expected: "lowercase vs uppercase",
		},
		{
			name: "variable with equals",
			setup: func() {
				os.Setenv("EQUALS_VAR", "key=value")
				env.InitEnv()
			},
			input:    "Config: ${EQUALS_VAR}",
			expected: "Config: key=value",
		},
		{
			name: "multiple equals",
			setup: func() {
				os.Setenv("MULTI_EQUALS", "key=value=another=value")
				env.InitEnv()
			},
			input:    "${MULTI_EQUALS}",
			expected: "key=value=another=value",
		},
		{
			name: "whitespace in values",
			setup: func() {
				os.Setenv("WHITESPACE_VAR", "  value with leading and trailing spaces  ")
				env.InitEnv()
			},
			input:    "${WHITESPACE_VAR}",
			expected: "  value with leading and trailing spaces  ",
		},
		{
			name: "newlines in values",
			setup: func() {
				os.Setenv("NEWLINE_VAR", "line1\nline2\nline3")
				env.InitEnv()
			},
			input:    "${NEWLINE_VAR}",
			expected: "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variables before each test
			for _, varName := range testVars {
				os.Unsetenv(varName)
			}

			tt.setup()
			result := env.ApplyEnvironmentToString(tt.input)
			if result != tt.expected {
				t.Errorf("ApplyEnvironmentToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}