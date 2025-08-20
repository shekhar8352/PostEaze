# Utility Function Tests

This directory contains comprehensive unit tests for all utility functions in the PostEaze backend application.

## Test Files

### Authentication Utilities
- **`auth_utils_test.go`** - Tests for modern JWT utilities (access/refresh tokens)
  - Token generation and parsing
  - Token expiration handling
  - User ID extraction
  - Error handling for invalid tokens

- **`legacy_jwt_test.go`** - Tests for legacy JWT utilities
  - Legacy JWT generation and parsing
  - Backward compatibility testing
  - Error handling for malformed tokens

### General Utilities
- **`general_utils_test.go`** - Tests for HTTP response utilities
  - `SendError()` function testing
  - `SendSuccess()` function testing
  - Response format consistency
  - Various data type handling

### Environment Variable Utilities
- **`env_replacement_test.go`** - Tests for environment variable replacement
  - `ReplacePlaceHoldersWithEnv()` function testing
  - Single and multiple variable replacement
  - Edge cases and error handling
  - Special characters and nested scenarios

- **`env_utils_test.go`** - Tests for environment initialization utilities
  - `InitEnv()` function testing
  - `.env` file loading
  - `ApplyEnvironmentToString()` function testing
  - Environment variable handling edge cases

### Configuration Utilities
- **`config_utils_test.go`** - Tests for configuration management
  - Development configuration initialization
  - Release configuration initialization
  - Error handling for invalid configurations
  - Multiple configuration file handling

### HTTP Client Utilities
- **`http_client_test.go`** - Tests for HTTP client utilities
  - Request configuration creation
  - HTTP client initialization
  - Request execution and response handling
  - Error handling and timeout scenarios

## Running Tests

### Run all utility tests:
```bash
go test ./tests/unit/utils/... -v
```

### Run specific test files:
```bash
go test ./tests/unit/utils/auth_utils_test.go -v
go test ./tests/unit/utils/general_utils_test.go -v
go test ./tests/unit/utils/env_replacement_test.go -v
```

### Run with coverage:
```bash
go test ./tests/unit/utils/... -v -cover
```

## Test Coverage

The tests cover the following utility functions:

### Password and Authentication
- `HashPassword()` - Password hashing with bcrypt
- `CheckPasswordHash()` - Password verification
- `GenerateAccessToken()` - JWT access token generation
- `GenerateRefreshToken()` - JWT refresh token generation
- `ParseToken()` - JWT token parsing and validation
- `GetUserIDFromToken()` - User ID extraction from tokens
- `GetRefreshTokenExpiry()` - Refresh token expiry calculation

### Legacy JWT Functions
- `GenerateJWT()` - Legacy JWT generation
- `ParseJWT()` - Legacy JWT parsing

### HTTP Response Utilities
- `SendError()` - Error response formatting
- `SendSuccess()` - Success response formatting

### Environment Variable Handling
- `ReplacePlaceHoldersWithEnv()` - Environment variable replacement
- `InitEnv()` - Environment initialization
- `ApplyEnvironmentToString()` - Environment variable application

### Configuration Management
- `InitDev()` - Development configuration initialization
- `InitRelease()` - Release configuration initialization
- `Get()` - Configuration client getter

### HTTP Client Utilities
- `NewRequestConfig()` - Request configuration creation
- `InitHttp()` - HTTP client initialization
- `GetClient()` - HTTP client getter
- `Call()` - HTTP request execution
- `CallAndGetResponse()` - HTTP request with response
- `CallAndBind()` - HTTP request with JSON binding

## Test Patterns

All tests follow these patterns:

1. **Test Suite Structure** - Using testify/suite for organized test setup and teardown
2. **Environment Isolation** - Tests save and restore environment variables
3. **Comprehensive Coverage** - Tests cover happy paths, edge cases, and error conditions
4. **Clear Naming** - Test names clearly indicate what is being tested
5. **Proper Assertions** - Using testify assertions for clear error messages

## Dependencies

The tests use the following testing libraries:
- `github.com/stretchr/testify/suite` - Test suite organization
- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/require` - Required assertions
- Standard Go `testing` package

## Notes

- Some tests may behave differently in different environments due to OS-specific behavior
- Environment variable tests properly clean up after themselves
- HTTP client tests use a test server for realistic testing
- JWT tests use test secrets to avoid affecting production configurations