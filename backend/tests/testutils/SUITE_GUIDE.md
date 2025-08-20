# Base Test Suite Structure Guide

This guide explains how to use the base test suite structure provided by the PostEaze testing framework.

## Overview

The testing framework provides several specialized test suite types that extend the base functionality:

- **BaseTestSuite**: Core functionality for all test suites
- **APITestSuite**: Specialized for API handler testing
- **DatabaseTestSuite**: Enhanced database testing with transactions
- **BusinessLogicTestSuite**: For business logic layer testing
- **ModelTestSuite**: For model validation and serialization testing

## BaseTestSuite

The foundation for all test suites, providing common setup and teardown functionality.

### Features

- Automatic database setup (with graceful fallback if unavailable)
- Gin router initialization for HTTP testing
- Test fixture management
- Context management
- Common cleanup operations

### Usage

```go
type MyTestSuite struct {
    testutils.BaseTestSuite
}

func (s *MyTestSuite) TestSomething() {
    // Your test code here
    s.NotNil(s.Router, "Router should be available")
    s.NotNil(s.DB, "Database should be available (if configured)")
}

func TestMyTestSuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}
```

### Available Methods

- `LoadTestFixtures()`: Loads predefined test data
- `LoadCustomFixtures(fixtures...)`: Loads custom test data
- `GetTestUser(index)`: Gets predefined test user
- `GetTestTeam(index)`: Gets predefined test team
- `GetTestToken(index)`: Gets predefined test token
- `CreateTestUser(overrides...)`: Creates custom test user
- `CreateTestTeam(ownerID, overrides...)`: Creates custom test team
- `CreateTestToken(userID, overrides...)`: Creates custom test token

## APITestSuite

Specialized for testing API handlers with HTTP request/response utilities.

### Features

- Common API middleware setup
- Authenticated request creation
- Response assertion helpers
- JSON parsing utilities

### Usage

```go
type AuthAPITestSuite struct {
    testutils.APITestSuite
}

func (s *AuthAPITestSuite) TestLoginEndpoint() {
    // Create authenticated request
    ctx, recorder, err := s.CreateAuthenticatedRequest(
        "POST", "/api/v1/auth/login", 
        map[string]string{"email": "test@example.com"}, 
        "user123", "admin",
    )
    s.NoError(err)
    
    // Test your handler
    // handler(ctx)
    
    // Assert response
    s.AssertSuccessAPIResponse(recorder, nil)
}
```

### Available Methods

- `CreateAuthenticatedRequest(method, url, body, userID, userType)`: Creates authenticated HTTP request
- `AssertAPIResponse(recorder, expectedStatus, expectedMessage)`: General API response assertion
- `AssertSuccessAPIResponse(recorder, expectedData)`: Assert successful response
- `AssertErrorAPIResponse(recorder, expectedStatus, expectedMessage)`: Assert error response

## DatabaseTestSuite

Enhanced database testing with transaction support and database package integration.

### Features

- Automatic database package initialization
- Transaction management
- Enhanced cleanup operations
- Database-specific utilities

### Usage

```go
type UserRepositoryTestSuite struct {
    testutils.DatabaseTestSuite
}

func (s *UserRepositoryTestSuite) TestCreateUser() {
    // Begin transaction for test isolation
    tx := s.BeginTransaction()
    defer s.RollbackTransaction(tx)
    
    // Your database test code here
    // repo := NewUserRepository(s.DB)
    // user, err := repo.Create(ctx, userData)
    
    s.NoError(err)
    s.NotNil(user)
}
```

### Available Methods

- `BeginTransaction()`: Starts a test transaction
- `RollbackTransaction(tx)`: Rolls back a test transaction

## BusinessLogicTestSuite

Specialized for testing business logic layer with pre-loaded fixtures.

### Features

- Automatic fixture loading
- Mock context creation
- Business logic testing utilities

### Usage

```go
type AuthBusinessLogicTestSuite struct {
    testutils.BusinessLogicTestSuite
}

func (s *AuthBusinessLogicTestSuite) TestUserSignup() {
    ctx := s.CreateMockContext()
    
    // Your business logic test code here
    // service := NewAuthService(mockRepo)
    // result, err := service.Signup(ctx, signupData)
    
    s.NoError(err)
    s.NotNil(result)
}
```

### Available Methods

- `CreateMockContext()`: Creates a mock context for business logic testing

## ModelTestSuite

Specialized for testing model validation and serialization.

### Features

- Model validation testing utilities
- Serialization testing helpers
- Validation error assertions

### Usage

```go
type UserModelTestSuite struct {
    testutils.ModelTestSuite
}

func (s *UserModelTestSuite) TestUserValidation() {
    user := &User{Email: "invalid-email"}
    
    err := user.Validate()
    s.AssertValidationError(err, "email")
    
    user.Email = "valid@example.com"
    err = user.Validate()
    s.AssertNoValidationError(err)
}
```

### Available Methods

- `AssertValidationError(err, fieldName)`: Assert validation error occurred
- `AssertNoValidationError(err)`: Assert no validation error occurred

## TestSuiteHelper

Additional utility functions for advanced testing scenarios.

### Features

- Timeout-based test execution
- Eventually-true assertions
- Advanced testing patterns

### Usage

```go
func TestWithHelper(t *testing.T) {
    helper := testutils.NewTestSuiteHelper()
    
    // Run function with timeout
    helper.RunWithTimeout(t, 5*time.Second, func() {
        // Your test code that should complete within 5 seconds
    })
    
    // Assert condition becomes true eventually
    counter := 0
    helper.AssertEventuallyTrue(t, func() bool {
        counter++
        return counter >= 3
    }, 100*time.Millisecond, "Counter should reach 3")
}
```

### Available Methods

- `RunWithTimeout(t, timeout, fn)`: Run function with timeout
- `AssertEventuallyTrue(t, condition, timeout, message)`: Assert condition becomes true

## Best Practices

### 1. Choose the Right Suite Type

- Use `BaseTestSuite` for general testing
- Use `APITestSuite` for HTTP handler testing
- Use `DatabaseTestSuite` for repository/database testing
- Use `BusinessLogicTestSuite` for service layer testing
- Use `ModelTestSuite` for model validation testing

### 2. Test Isolation

- Each test method runs with fresh fixtures
- Database tests automatically clean up between tests
- Use transactions in `DatabaseTestSuite` for better isolation

### 3. Fixture Management

```go
func (s *MyTestSuite) TestWithCustomData() {
    // Create custom test data
    user := s.CreateTestUser(func(u *testutils.UserFixture) {
        u.Email = "custom@test.com"
        u.UserType = "admin"
    })
    
    // Use the custom user in your test
    s.Equal("custom@test.com", user.Email)
}
```

### 4. Error Handling

```go
func (s *MyTestSuite) TestErrorScenario() {
    // Test error conditions
    result, err := someFunction()
    
    s.Error(err, "Should return error for invalid input")
    s.Nil(result, "Result should be nil on error")
    s.Contains(err.Error(), "expected error message")
}
```

### 5. HTTP Testing

```go
func (s *APITestSuite) TestAPIEndpoint() {
    // Create request context
    ctx, recorder := testutils.NewTestGinContext("POST", "/api/endpoint", requestData)
    
    // Call your handler
    handler(ctx)
    
    // Assert response
    s.AssertSuccessAPIResponse(recorder, expectedData)
}
```

## Running Tests

### Run All Suite Tests
```bash
go test -v ./tests/testutils -run "TestSuite"
```

### Run Specific Suite Type
```bash
go test -v ./tests/testutils -run "TestAPITestSuite"
```

### Run with Coverage
```bash
go test -v -cover ./tests/testutils
```

## Common Patterns

### Setup and Teardown

```go
func (s *MyTestSuite) SetupTest() {
    // Called before each test method
    s.LoadTestFixtures()
}

func (s *MyTestSuite) TearDownTest() {
    // Called after each test method
    // Additional cleanup if needed
}
```

### Parameterized Tests

```go
func (s *MyTestSuite) TestMultipleScenarios() {
    testCases := []struct {
        name     string
        input    string
        expected string
    }{
        {"valid input", "test", "TEST"},
        {"empty input", "", ""},
    }
    
    for _, tc := range testCases {
        s.Run(tc.name, func() {
            result := processInput(tc.input)
            s.Equal(tc.expected, result)
        })
    }
}
```

### Database Transaction Testing

```go
func (s *DatabaseTestSuite) TestWithTransaction() {
    tx := s.BeginTransaction()
    defer s.RollbackTransaction(tx)
    
    // All database operations within this test
    // will be rolled back automatically
    repo := NewRepository(tx)
    result, err := repo.Create(testData)
    
    s.NoError(err)
    s.NotNil(result)
}
```

This base test suite structure provides a solid foundation for comprehensive testing across all layers of the PostEaze application.