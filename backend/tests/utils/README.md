# Utility Tests

This directory contains simplified unit tests for utility functions in the PostEaze backend.

## Structure

- `auth_utils_test.go` - Tests for authentication utilities (JWT, password hashing)
- `config_utils_test.go` - Tests for configuration management utilities
- `database_utils_test.go` - Tests for database utilities and interfaces
- `env_replacement_test.go` - Tests for environment variable replacement functions
- `env_utils_test.go` - Tests for environment initialization and management
- `general_utils_test.go` - Tests for general HTTP response utilities
- `http_client_test.go` - Tests for HTTP client utilities
- `http_request_utils_test.go` - Tests for HTTP request building utilities
- `legacy_jwt_test.go` - Tests for legacy JWT functions
- `logapi_utils_test.go` - Tests for log parsing utilities

## Testing Approach

These tests follow a simplified approach compared to the previous test suite:

- **Standard Go Testing**: Uses standard `testing` package without complex test suites
- **Table-Driven Tests**: Uses table-driven test patterns for multiple scenarios
- **Minimal Setup**: Reduces complex setup and teardown in favor of simple test functions
- **Clear Naming**: Self-documenting test names that describe the scenario being tested
- **Focused Tests**: Each test focuses on a specific function or behavior

## Running Tests

Run all utility tests:
```bash
go test ./backend/tests/utils/...
```

Run specific test file:
```bash
go test ./backend/tests/utils/auth_utils_test.go
```

Run with verbose output:
```bash
go test -v ./backend/tests/utils/...
```

## Test Patterns

### Basic Function Test
```go
func TestFunctionName(t *testing.T) {
    input := "test input"
    expected := "expected output"
    
    result := utils.FunctionName(input)
    
    if result != expected {
        t.Errorf("FunctionName() = %v, want %v", result, expected)
    }
}
```

### Table-Driven Test
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "result", false},
        {"invalid input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := utils.FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Setup and Cleanup
```go
func TestFunctionWithSetup(t *testing.T) {
    // Store original state
    original := os.Getenv("TEST_VAR")
    defer func() {
        if original != "" {
            os.Setenv("TEST_VAR", original)
        } else {
            os.Unsetenv("TEST_VAR")
        }
    }()
    
    // Set test state
    os.Setenv("TEST_VAR", "test_value")
    
    // Run test
    result := utils.FunctionName()
    if result != "expected" {
        t.Errorf("FunctionName() = %v, want expected", result)
    }
}
```

## Guidelines

1. **Keep tests simple** - Avoid complex abstractions and helper functions
2. **Test one thing** - Each test should focus on a single behavior
3. **Use descriptive names** - Test names should clearly describe what is being tested
4. **Handle cleanup** - Always restore original state after tests that modify global state
5. **Test error cases** - Include tests for both success and failure scenarios
6. **Use standard patterns** - Follow Go testing conventions and idioms