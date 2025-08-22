# PostEaze Backend Testing Guide

This comprehensive guide covers everything you need to know about writing and running tests in the PostEaze backend application using our simplified testing framework.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Directory Structure](#directory-structure)
3. [Writing Tests](#writing-tests)
4. [Running Tests](#running-tests)
5. [Test Utilities](#test-utilities)
6. [Best Practices](#best-practices)
7. [Coverage and Performance](#coverage-and-performance)
8. [Troubleshooting](#troubleshooting)

## Quick Start

### Most Common Commands (Fixed and Working)

```bash
# Run all tests
./scripts/test.sh

# Run unit tests with verbose output (for debugging)
./scripts/test.sh --type unit --verbose

# Run specific package tests
./scripts/test.sh --type unit --package utils

# Run with coverage
./scripts/test.sh --coverage --threshold 80

# Run integration tests
./scripts/test.sh --type integration

# Run benchmark tests
./scripts/test.sh --type benchmark
```

### Alternative Methods

```bash
# Using the simplified test runner directly
cd backend && go run tests/test.go

# Using Go directly (for individual packages)
cd backend && go test ./tests/utils/
```

### Run Specific Test Types
```bash
# Unit tests only
./scripts/test.sh --type unit

# Integration tests only
./scripts/test.sh --type integration

# Benchmark tests only
./scripts/test.sh --type benchmark

# With coverage
./scripts/test.sh --coverage

# Specific package tests
./scripts/test.sh --type unit --package utils
./scripts/test.sh --type unit --package api

# With verbose output for debugging
./scripts/test.sh --type unit --verbose
```

### Windows Users
```powershell
# PowerShell script
.\scripts\test.ps1 -Type all

# With coverage
.\scripts\test.ps1 -Coverage

# Specific suite
.\scripts\test.ps1 -Type unit -Suite api
```

## Directory Structure

```
tests/
├── api/                   # API handler tests (unit tests)
├── business/              # Business logic tests (unit tests)
├── models/                # Model tests (unit tests)
├── utils/                 # Utility function tests (unit tests)
├── integration/           # Integration tests
│   ├── auth_integration_test.go
│   ├── log_integration_test.go
│   ├── database_integration_test.go
│   └── end_to_end_integration_test.go
├── helpers/               # Simple test helpers and utilities
├── examples/              # Example test implementations
├── benchmarks/            # Performance benchmarks
├── test.go                # Simplified test runner
├── .env.test              # Test environment variables
└── TESTING_GUIDE.md       # This file
```

### Test Types

- **Unit Tests (70%)**: Fast, isolated tests in `api/`, `business/`, `models/`, `utils/`
- **Integration Tests (20%)**: Component interaction tests in `integration/`
- **End-to-End Tests (10%)**: Complete workflow tests

## Writing Tests

### Naming Conventions

#### Test Files
```go
// Good
auth_handler_test.go
user_service_test.go
jwt_utils_test.go

// Avoid
test_auth.go
auth_tests.go
```

#### Test Functions
```go
// Good - Descriptive names that explain the scenario
func TestUserService_CreateUser_ValidInput_ReturnsUser(t *testing.T)
func TestAuthHandler_Login_InvalidCredentials_ReturnsUnauthorized(t *testing.T)
func TestJWTUtils_GenerateToken_ExpiredSecret_ReturnsError(t *testing.T)

// Avoid - Generic names
func TestCreateUser(t *testing.T)
func TestLogin(t *testing.T)
```

### Test Structure Pattern

Use the **Arrange-Act-Assert (AAA)** pattern for all tests:

```go
func TestUserService_CreateUser_ValidInput_ReturnsUser(t *testing.T) {
    // Arrange - Setup test data and mocks
    mockRepo := mocks.NewMockUserRepository()
    service := NewUserService(mockRepo)
    
    userData := &models.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    expectedUser := &models.User{
        ID:    "user-123",
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    mockRepo.On("Create", mock.Anything, userData).Return(expectedUser, nil)
    
    // Act - Execute the function under test
    result, err := service.CreateUser(context.Background(), userData)
    
    // Assert - Verify the results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedUser.ID, result.ID)
    assert.Equal(t, expectedUser.Name, result.Name)
    assert.Equal(t, expectedUser.Email, result.Email)
    
    mockRepo.AssertExpectations(t)
}
```

### Example API Handler Test

```go
func TestAuthHandler_Signup_ValidInput_Success(t *testing.T) {
    // Arrange
    signupData := modelsv1.SignupParams{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
        UserType: modelsv1.UserTypeIndividual,
    }
    
    // Act
    ctx, recorder := helpers.NewTestGinContext("POST", "/auth/signup", signupData)
    handler.Signup(ctx)
    
    // Assert
    assert.Equal(t, http.StatusCreated, recorder.Code)
    
    var response map[string]interface{}
    err := helpers.ParseJSONResponse(recorder, &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
}
```

### Example Business Logic Test

```go
func TestAuthService_Login_ValidCredentials_Success(t *testing.T) {
    // Arrange
    mockDB, err := helpers.SetupMockDB()
    require.NoError(t, err)
    defer mockDB.Close()
    
    // Set up mock expectations
    mockDB.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
        WithArgs("test@example.com").
        WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
            AddRow("user-1", "$2a$10$hashedpassword"))
    
    service := business.NewAuthService(mockDB.DB)
    
    // Act
    result, err := service.Login("test@example.com", "password")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "user-1", result.UserID)
    
    err = mockDB.ExpectationsWereMet()
    assert.NoError(t, err)
}
```

### Table-Driven Tests

Use table-driven tests for multiple similar scenarios:

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        expectValid bool
        expectError string
    }{
        {
            name:        "valid email",
            email:       "user@example.com",
            expectValid: true,
        },
        {
            name:        "missing @ symbol",
            email:       "userexample.com",
            expectValid: false,
            expectError: "invalid email format",
        },
        {
            name:        "empty email",
            email:       "",
            expectValid: false,
            expectError: "email is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            valid, err := ValidateEmail(tt.email)
            
            // Assert
            assert.Equal(t, tt.expectValid, valid)
            if tt.expectError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectError)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Running Tests

### Test Execution Scripts

#### Main Test Script (`scripts/test.sh`)

```bash
./scripts/test.sh [OPTIONS]
```

**Options:**
- `--type TYPE`: Test type (unit|integration|benchmark|all)
- `--package PKG`: Specific package to test (api|business|models|utils)
- `--verbose`: Enable verbose output for debugging
- `--coverage`: Generate coverage report
- `--threshold NUM`: Coverage threshold percentage (0-100)

**Examples:**
```bash
# Run all tests with coverage
./scripts/test.sh --coverage

# Run unit tests for utils package with verbose output
./scripts/test.sh --type unit --package utils --verbose

# Run integration tests
./scripts/test.sh --type integration

# Run with coverage threshold
./scripts/test.sh --coverage --threshold 80

# Run benchmark tests
./scripts/test.sh --type benchmark --verbose
```

#### Direct Go Commands (Alternative)

```bash
# Run tests for specific packages (from backend directory)
cd backend
go test ./tests/utils/      # Utils package tests
go test ./tests/business/   # Business logic tests
go test ./tests/integration/ # Integration tests

# Run with coverage
go test -cover ./tests/utils/

# Run with verbose output
go test -v ./tests/utils/

# Run specific test
go test -v -run TestHashPassword ./tests/utils/
```

#### Simplified Test Runner

```bash
# Run all tests
cd backend && go run tests/test.go

# Run specific test type
cd backend && go run tests/test.go -type=unit

# Run with coverage
cd backend && go run tests/test.go -coverage

# Run specific package
cd backend && go run tests/test.go -type=unit -package=utils

# Run with verbose output
cd backend && go run tests/test.go -type=unit -verbose

# Run with coverage threshold
cd backend && go run tests/test.go -coverage -coverage-threshold=80
```

## Test Utilities

### Simple Test Helpers (`helpers/`)

The simplified testing framework provides essential utilities:

#### Database Helpers (`helpers/database.go`)
- `SetupTestDB(t)` - Simple database setup
- `SetupMockDB()` - Mock database for unit tests
- `SetupInMemoryDB(t)` - In-memory database for fast tests
- Basic fixture loading utilities
- Transaction management for test isolation

#### HTTP Helpers (`helpers/http.go`)
- `NewTestGinContext(method, url, body)` - Create test Gin context
- `ParseJSONResponse(recorder, target)` - Parse JSON responses
- `CreateAuthenticatedContext(method, url, userID, role, body)` - Create authenticated context
- Response assertion helpers

#### Test Data (`helpers/fixtures.go`)
- Simple test data creation functions
- Basic fixture management
- Minimal test data structures

### Test Data Management

#### Test Fixtures
```go
// Good: Use fixture functions
func createTestUser() *models.User {
    return &models.User{
        ID:       "user-123",
        Name:     "Test User",
        Email:    "test@example.com",
        UserType: models.UserTypeIndividual,
        CreatedAt: time.Now(),
    }
}

func createTestUserWithOverrides(overrides func(*models.User)) *models.User {
    user := createTestUser()
    if overrides != nil {
        overrides(user)
    }
    return user
}

// Usage
user := createTestUserWithOverrides(func(u *models.User) {
    u.Email = "custom@example.com"
    u.UserType = models.UserTypeTeam
})
```

#### Database Test Data
```go
func TestUserRepository_Create(t *testing.T) {
    // Setup
    db := helpers.SetupTestDB(t)
    defer db.Close()
    
    tx, err := db.Begin()
    require.NoError(t, err)
    defer tx.Rollback() // Always rollback in tests
    
    repo := NewUserRepository(tx)
    
    // Test
    user := createTestUser()
    result, err := repo.Create(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.NotEmpty(t, result.ID)
}
```

## Best Practices

### Test Organization
- Use simple, focused test functions
- Use descriptive test names that explain the scenario
- Keep tests focused on a single behavior
- Use the Arrange-Act-Assert pattern
- Group tests by the package they're testing

### Mocking Guidelines
- Mock external dependencies (database, HTTP clients, etc.)
- Use interfaces to enable easy mocking
- Keep mocks simple and focused
- Always verify mock expectations
- Avoid complex mock frameworks

```go
// Good: Mock the interface
type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
}

func TestUserService(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    defer mockRepo.AssertExpectations(t) // Always verify expectations
    
    service := NewUserService(mockRepo)
    // Test with mock...
}
```

### Assertion Best Practices
```go
// Good: Use specific assertions
assert.Equal(t, expectedValue, actualValue)
assert.True(t, condition)
assert.NotNil(t, pointer)
assert.Contains(t, slice, element)
assert.Len(t, slice, expectedLength)

// Good: Test error content
err := someFunction()
assert.Error(t, err)
assert.Contains(t, err.Error(), "expected error message")

// Good: Test no error
err := someFunction()
assert.NoError(t, err)
```

### Performance Considerations
- Use in-memory databases for unit tests
- Use mocks for external dependencies in unit tests
- Run tests in parallel when safe (`t.Parallel()`)
- Avoid unnecessary setup/teardown
- Use table-driven tests for multiple scenarios

### Error Testing
Always test error conditions:

```go
func TestUserService_CreateUser_ErrorScenarios(t *testing.T) {
    t.Run("database error", func(t *testing.T) {
        mockRepo := mocks.NewMockUserRepository()
        mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
        
        service := NewUserService(mockRepo)
        user, err := service.CreateUser(context.Background(), &models.User{})
        
        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Contains(t, err.Error(), "database error")
    })
    
    t.Run("validation error", func(t *testing.T) {
        service := NewUserService(nil)
        user, err := service.CreateUser(context.Background(), &models.User{}) // Invalid user
        
        assert.Error(t, err)
        assert.Nil(t, user)
        assert.IsType(t, &ValidationError{}, err)
    })
}
```

## Coverage and Performance

### Coverage Goals
- **Unit Tests**: 80-90% line coverage
- **Integration Tests**: Focus on critical paths
- **Overall**: 80% minimum, 90% target

### Coverage Generation
```bash
# Generate coverage report
./scripts/test.sh --coverage

# With specific threshold
./scripts/test.sh --coverage --threshold 90

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Benchmark Tests
```bash
# Run all benchmarks
./scripts/test.sh --benchmark

# Benchmarks with memory profiling
make perf-test

# CPU and memory profiling
make perf-profile
```

### Integration Testing
```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    db := helpers.SetupIntegrationDB(t)
    defer helpers.CleanupIntegrationDB(t, db)
    
    repo := NewUserRepository(db)
    
    // Test with real database
    user, err := repo.Create(context.Background(), createTestUser())
    assert.NoError(t, err)
    assert.NotNil(t, user)
    
    // Verify persistence
    retrieved, err := repo.GetByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.ID, retrieved.ID)
}
```

## Troubleshooting

### Common Issues

#### Test Runner Issues

**Problem: "package tests/utils/integration_test.go is not in std" error**
```bash
# This error occurs when the test runner tries to run individual files instead of packages
# Solution: Use the fixed test runner that runs tests at package level
./scripts/test.sh --type unit --package utils
```

**Problem: "cannot find main module" error**
```bash
# This happens when running from wrong directory
# Solution: Always run from the backend directory
cd backend
./scripts/test.sh --type unit
```

**Problem: Tests not found for package**
```bash
# Check if test files exist in the package directory
ls -la tests/utils/*_test.go

# Run with verbose to see what's happening
./scripts/test.sh --type unit --package utils --verbose
```

#### Database Connection Issues
```bash
# For integration tests, ensure database is available
# Most unit tests use in-memory databases and don't require external DB
```

#### Coverage Threshold Failures
```bash
# Check current coverage
go tool cover -func=coverage.out | tail -1

# Run with lower threshold
./scripts/test.sh --coverage --threshold 70
```

#### Race Condition Detection
```bash
# Run with race detection
./scripts/test.sh --race

# Fix race conditions in code
go run -race main.go
```

#### Test Timeout Issues
```bash
# Run with shorter timeout
go test -timeout 5m ./...

# Run in short mode
./scripts/test.sh --short
```

### Debug Mode
```bash
# Verbose test output
./scripts/test.sh --verbose

# Go test verbose mode
go test -v ./tests/...

# Debug specific test
go test -v -run TestSpecificFunction ./tests/api/
```

### Environment Issues
```bash
# Check test environment
env | grep -E "(MODE|DATABASE|JWT)"

# Set test environment
export MODE=test
export DATABASE_DRIVER=sqlite3
export DATABASE_URL=":memory:"
```

## Current Test Status

### Test Suite Status (as of latest fixes)

✅ **Test Runner**: Fixed and working correctly
- Fixed package-level test execution
- Fixed working directory issues
- Proper JSON output parsing
- Coverage reporting functional

✅ **Unit Tests**: 
- **Utils Package**: 292 tests found, 283 passing, 9 failing
- **API Package**: No test files found (needs implementation)
- **Business Package**: Limited test files
- **Models Package**: No test files found (needs implementation)

✅ **Integration Tests**: Framework working (no tests found)
✅ **Benchmark Tests**: Framework working (no benchmarks found)

### Known Issues
- 9 failing tests in utils package (legitimate test failures, not framework issues)
- Missing test implementations in api and models packages
- Some test files may need environment setup

### Next Steps
1. Fix the 9 failing tests in utils package
2. Implement missing tests in api and models packages
3. Add integration tests for critical workflows
4. Add benchmark tests for performance-critical functions

## Configuration

### Test Environment (`.env.test`)
```bash
# Database Configuration
TEST_DATABASE_URL=:memory:
TEST_DATABASE_DRIVER=sqlite3

# JWT Configuration
TEST_JWT_SECRET=test-secret-key-for-simplified-testing-framework
TEST_JWT_EXPIRY=1h

# Logging Configuration
TEST_LOG_LEVEL=error
TEST_LOG_OUTPUT=stdout

# Test Execution Configuration
TEST_TIMEOUT=30s
TEST_PARALLEL=true
TEST_COVERAGE_THRESHOLD=80
```

## Testing Checklist

Before submitting code, ensure:

- [ ] All tests pass locally
- [ ] New code has appropriate test coverage
- [ ] Tests follow naming conventions
- [ ] Tests are independent and can run in any order
- [ ] Mocks are properly set up and verified
- [ ] Error cases are tested
- [ ] Tests are fast and focused
- [ ] Integration tests cover critical workflows

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Coverage Tools](https://golang.org/doc/code.html#Testing)
- [PostEaze Backend README](../README.md)

---

This guide should be treated as a living document that evolves with the project and team needs. For questions or suggestions, please reach out to the development team.