# PostEaze Backend Testing Best Practices Guide

This guide provides comprehensive best practices, patterns, and conventions for writing effective tests in the PostEaze backend application.

## Table of Contents

1. [Testing Philosophy](#testing-philosophy)
2. [Test Organization](#test-organization)
3. [Naming Conventions](#naming-conventions)
4. [Test Structure Patterns](#test-structure-patterns)
5. [Data Management](#data-management)
6. [Mocking Guidelines](#mocking-guidelines)
7. [Assertion Best Practices](#assertion-best-practices)
8. [Performance Considerations](#performance-considerations)
9. [Error Testing](#error-testing)
10. [Integration Testing](#integration-testing)
11. [Code Coverage](#code-coverage)
12. [Common Anti-Patterns](#common-anti-patterns)

## Testing Philosophy

### Test Pyramid
Follow the test pyramid principle:
- **Unit Tests (70%)**: Fast, isolated tests for individual components
- **Integration Tests (20%)**: Test component interactions
- **End-to-End Tests (10%)**: Test complete user workflows

### Testing Principles
1. **Fast**: Tests should run quickly to provide rapid feedback
2. **Independent**: Tests should not depend on each other
3. **Repeatable**: Tests should produce consistent results
4. **Self-Validating**: Tests should have clear pass/fail outcomes
5. **Timely**: Tests should be written alongside or before the code

## Test Organization

### Directory Structure
```
tests/
├── unit/                    # Unit tests (fast, isolated)
│   ├── api/v1/             # API handler tests
│   ├── business/v1/        # Business logic tests
│   ├── models/v1/          # Model tests
│   └── utils/              # Utility tests
├── integration/            # Integration tests (slower, with dependencies)
│   ├── auth_integration_test.go
│   ├── log_integration_test.go
│   └── database_integration_test.go
├── testutils/              # Test utilities and helpers
├── config/                 # Test configuration
└── examples/               # Example test implementations
```

### Package Organization
- Group tests by the package they're testing
- Use consistent naming: `package_test` for external testing
- Keep test files close to the code they test

## Naming Conventions

### Test Files
```go
// Good
auth_handler_test.go
user_service_test.go
jwt_utils_test.go

// Avoid
test_auth.go
auth_tests.go
```

### Test Functions
Use descriptive names that explain the scenario:

```go
// Good
func TestUserService_CreateUser_ValidInput_ReturnsUser(t *testing.T)
func TestAuthHandler_Login_InvalidCredentials_ReturnsUnauthorized(t *testing.T)
func TestJWTUtils_GenerateToken_ExpiredSecret_ReturnsError(t *testing.T)

// Avoid
func TestCreateUser(t *testing.T)
func TestLogin(t *testing.T)
func TestJWT(t *testing.T)
```

### Test Suites
```go
// Good
type AuthHandlerTestSuite struct {
    testutils.APITestSuite
}

type UserServiceTestSuite struct {
    testutils.BusinessLogicTestSuite
}

// Avoid
type AuthTests struct {
    suite.Suite
}
```

### Variables and Constants
```go
// Good
const (
    testUserEmail    = "test@example.com"
    testUserPassword = "password123"
    validJWTToken    = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
)

var (
    testUserID   = "user-123"
    testTeamID   = "team-456"
    testTimeout  = 5 * time.Second
)

// Avoid
const email = "test@example.com"
var id = "123"
```

## Test Structure Patterns

### AAA Pattern (Arrange-Act-Assert)
Structure all tests using the Arrange-Act-Assert pattern:

```go
func TestUserService_CreateUser_ValidInput_ReturnsUser(t *testing.T) {
    // Arrange
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
    
    // Act
    result, err := service.CreateUser(context.Background(), userData)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedUser.ID, result.ID)
    assert.Equal(t, expectedUser.Name, result.Name)
    assert.Equal(t, expectedUser.Email, result.Email)
    
    mockRepo.AssertExpectations(t)
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
            name:        "missing domain",
            email:       "user@",
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

### Subtests for Related Scenarios
Use subtests to group related test scenarios:

```go
func TestAuthHandler_Login(t *testing.T) {
    t.Run("successful login", func(t *testing.T) {
        // Test successful login scenario
    })
    
    t.Run("invalid credentials", func(t *testing.T) {
        // Test invalid credentials scenario
    })
    
    t.Run("missing email", func(t *testing.T) {
        // Test missing email scenario
    })
    
    t.Run("missing password", func(t *testing.T) {
        // Test missing password scenario
    })
}
```

## Data Management

### Test Fixtures
Use consistent test fixtures for predictable test data:

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

### Database Test Data
Use transactions for database test isolation:

```go
func TestUserRepository_Create(t *testing.T) {
    // Setup
    db := testutils.SetupTestDB(t)
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

### Cleanup Strategies
Always clean up test data:

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    tempFile := createTempFile(t)
    defer os.Remove(tempFile) // Cleanup file
    
    mockServer := httptest.NewServer(handler)
    defer mockServer.Close() // Cleanup server
    
    // Test code...
}
```

## Mocking Guidelines

### Interface-Based Mocking
Always mock interfaces, not concrete types:

```go
// Good: Mock the interface
type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
}

func TestUserService(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    service := NewUserService(mockRepo)
    // Test with mock...
}

// Avoid: Mocking concrete types
type ConcreteUserRepository struct {
    db *sql.DB
}
```

### Mock Expectations
Set up clear, specific mock expectations:

```go
// Good: Specific expectations
mockRepo.On("GetByID", mock.Anything, "user-123").Return(expectedUser, nil)
mockRepo.On("GetByID", mock.Anything, "user-404").Return(nil, ErrUserNotFound)

// Avoid: Overly broad expectations
mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(expectedUser, nil)
```

### Mock Verification
Always verify mock expectations:

```go
func TestWithMocks(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    defer mockRepo.AssertExpectations(t) // Verify all expectations were met
    
    // Setup expectations
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)
    
    // Test code...
    
    // Additional specific verifications if needed
    mockRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*models.User"))
}
```

## Assertion Best Practices

### Use Specific Assertions
Choose the most specific assertion for your needs:

```go
// Good: Specific assertions
assert.Equal(t, expectedValue, actualValue)
assert.True(t, condition)
assert.NotNil(t, pointer)
assert.Contains(t, slice, element)
assert.Len(t, slice, expectedLength)

// Avoid: Generic assertions
assert.True(t, actualValue == expectedValue)
assert.True(t, pointer != nil)
assert.True(t, len(slice) == expectedLength)
```

### Error Assertions
Test both error presence and error content:

```go
// Good: Test error content
err := someFunction()
assert.Error(t, err)
assert.Contains(t, err.Error(), "expected error message")
assert.IsType(t, &ValidationError{}, err)

// Also good: Test no error
err := someFunction()
assert.NoError(t, err)

// Avoid: Only testing error presence
assert.NotNil(t, err)
```

### Custom Assertions
Create custom assertions for complex validations:

```go
func assertValidUser(t *testing.T, user *models.User) {
    t.Helper() // Mark as helper function
    
    assert.NotNil(t, user)
    assert.NotEmpty(t, user.ID)
    assert.NotEmpty(t, user.Name)
    assert.NotEmpty(t, user.Email)
    assert.Contains(t, []models.UserType{models.UserTypeIndividual, models.UserTypeTeam}, user.UserType)
    assert.False(t, user.CreatedAt.IsZero())
}

// Usage
user, err := service.CreateUser(ctx, userData)
assert.NoError(t, err)
assertValidUser(t, user)
```

## Performance Considerations

### Fast Test Execution
Keep tests fast by avoiding unnecessary operations:

```go
// Good: Use in-memory database for unit tests
func TestUserRepository_Unit(t *testing.T) {
    db := testutils.SetupInMemoryDB(t)
    defer db.Close()
    // Fast test...
}

// Good: Use mocks for external dependencies
func TestUserService_Unit(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    mockHTTP := mocks.NewMockHTTPClient()
    // Fast test with mocks...
}

// Avoid: Real database connections in unit tests
func TestUserService_Slow(t *testing.T) {
    db := connectToRealDatabase() // Slow!
    // Slow test...
}
```

### Parallel Test Execution
Make tests safe for parallel execution:

```go
func TestParallelSafe(t *testing.T) {
    t.Parallel() // Mark as safe for parallel execution
    
    // Use unique test data
    userID := fmt.Sprintf("user-%d", time.Now().UnixNano())
    
    // Avoid shared state
    // Test code...
}
```

### Resource Management
Properly manage test resources:

```go
func TestWithResources(t *testing.T) {
    // Use context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Use resource pools
    db := testutils.GetDBFromPool(t)
    defer testutils.ReturnDBToPool(db)
    
    // Test code...
}
```

## Error Testing

### Test All Error Paths
Ensure all error conditions are tested:

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
        service := NewUserService(nil) // No repo needed for validation
        user, err := service.CreateUser(context.Background(), &models.User{}) // Invalid user
        
        assert.Error(t, err)
        assert.Nil(t, user)
        assert.IsType(t, &ValidationError{}, err)
    })
}
```

### Error Type Testing
Test specific error types and codes:

```go
func TestErrorTypes(t *testing.T) {
    err := someFunction()
    
    // Test error type
    var validationErr *ValidationError
    assert.True(t, errors.As(err, &validationErr))
    
    // Test error code
    assert.Equal(t, ErrorCodeInvalidInput, validationErr.Code)
    
    // Test error fields
    assert.Contains(t, validationErr.Fields, "email")
}
```

## Integration Testing

### Database Integration
Use real database connections for integration tests:

```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    db := testutils.SetupIntegrationDB(t)
    defer testutils.CleanupIntegrationDB(t, db)
    
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

### HTTP Integration
Test complete HTTP workflows:

```go
func TestAuthFlow_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup test server
    server := testutils.SetupTestServer(t)
    defer server.Close()
    
    client := &http.Client{Timeout: 10 * time.Second}
    
    // Test signup
    signupResp := testutils.PostJSON(t, client, server.URL+"/auth/signup", signupData)
    assert.Equal(t, http.StatusCreated, signupResp.StatusCode)
    
    // Test login
    loginResp := testutils.PostJSON(t, client, server.URL+"/auth/login", loginData)
    assert.Equal(t, http.StatusOK, loginResp.StatusCode)
    
    // Extract token and test protected endpoint
    token := extractTokenFromResponse(t, loginResp)
    protectedResp := testutils.GetWithAuth(t, client, server.URL+"/protected", token)
    assert.Equal(t, http.StatusOK, protectedResp.StatusCode)
}
```

## Code Coverage

### Coverage Goals
Set appropriate coverage targets:
- **Unit Tests**: 80-90% line coverage
- **Integration Tests**: Focus on critical paths
- **Overall**: 80% minimum, 90% target

### Coverage Analysis
Use coverage tools effectively:

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Set coverage threshold
go test -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//' | awk '{if($1<80) exit 1}'
```

### Coverage Best Practices
- Focus on critical business logic
- Don't chase 100% coverage at the expense of test quality
- Use coverage to identify untested code paths
- Exclude generated code and vendor dependencies

## Common Anti-Patterns

### Avoid These Patterns

#### 1. Testing Implementation Details
```go
// Bad: Testing internal implementation
func TestUserService_Bad(t *testing.T) {
    service := NewUserService(mockRepo)
    
    // Don't test internal method calls
    assert.True(t, service.validateEmailCalled)
    assert.Equal(t, 1, service.hashPasswordCallCount)
}

// Good: Test behavior and outcomes
func TestUserService_Good(t *testing.T) {
    service := NewUserService(mockRepo)
    
    user, err := service.CreateUser(ctx, userData)
    
    // Test the outcome
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotEmpty(t, user.ID)
}
```

#### 2. Shared Test State
```go
// Bad: Shared mutable state
var globalTestUser *models.User

func TestA(t *testing.T) {
    globalTestUser = createTestUser()
    globalTestUser.Name = "Modified"
    // Test affects other tests
}

// Good: Isolated test data
func TestA(t *testing.T) {
    user := createTestUser()
    user.Name = "Modified"
    // Test is isolated
}
```

#### 3. Overly Complex Tests
```go
// Bad: Complex test doing too much
func TestComplexScenario(t *testing.T) {
    // 50 lines of setup
    // Multiple operations
    // Multiple assertions testing different things
}

// Good: Focused, single-purpose tests
func TestUserCreation(t *testing.T) {
    // Test only user creation
}

func TestUserValidation(t *testing.T) {
    // Test only user validation
}
```

#### 4. Ignoring Test Failures
```go
// Bad: Ignoring or commenting out failing tests
// func TestSomething(t *testing.T) {
//     // This test is failing, will fix later
// }

// Good: Fix or properly skip tests
func TestSomething(t *testing.T) {
    t.Skip("TODO: Fix this test - issue #123")
    // Or fix the test
}
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
- [ ] Documentation is updated if needed

## Continuous Improvement

### Regular Review
- Review test failures and flaky tests
- Refactor tests when code changes
- Update test utilities and helpers
- Share testing knowledge with the team

### Metrics to Track
- Test execution time
- Test coverage percentage
- Test failure rate
- Flaky test frequency

### Learning Resources
- Go testing documentation
- Testify framework documentation
- Testing best practices articles
- Team code review feedback

This guide should be treated as a living document that evolves with the project and team needs.