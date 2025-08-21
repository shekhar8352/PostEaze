# PostEaze Backend Testing Framework

This directory contains the comprehensive testing framework for the PostEaze backend application. The framework provides utilities, fixtures, and patterns for testing API handlers, business logic, models, and utilities.

## Directory Structure

```
tests/
├── testutils/           # Reusable testing utilities
│   ├── database.go      # Database testing utilities
│   ├── database_mock.go # Mock database utilities
│   ├── http.go          # HTTP testing utilities
│   ├── auth.go          # Authentication testing utilities
│   ├── fixtures.go      # Test data fixtures
│   ├── suite.go         # Base test suite structure
│   └── *_test.go        # Tests for testing utilities
├── config/              # Test configuration files
│   ├── test.env         # Test environment variables
│   └── test_config.json # Test configuration
├── integration/         # Integration tests (future)
├── unit/               # Unit tests organized by package
│   ├── api/            # API handler tests
│   ├── business/       # Business logic tests
│   ├── models/         # Model tests
│   └── utils/          # Utility function tests
└── README.md           # This file
```

## Testing Utilities

### Database Testing (`testutils/database.go`)

Provides utilities for database testing with proper isolation and cleanup:

- `SetupTestDB(ctx)` - Sets up a test database connection
- `CleanupTestData(ctx, db)` - Cleans up test data between tests
- `LoadFixtures(ctx, db, fixtures...)` - Loads test fixtures into database
- `CreateTestTables(ctx, db)` - Creates necessary tables for testing
- `BeginTestTransaction(ctx, db)` - Starts a test transaction
- `RollbackTestTransaction(tx)` - Rolls back a test transaction

### HTTP Testing (`testutils/http.go`)

Provides utilities for HTTP handler testing with request/response handling:

- `NewTestGinContext(method, url, body)` - Creates test Gin context
- `ParseJSONResponse(recorder, target)` - Parses JSON responses
- `AssertSuccessResponse(t, recorder, expectedData)` - Asserts successful API responses
- `AssertErrorResponse(t, recorder, expectedStatus, expectedMessage)` - Asserts error responses
- `SetAuthorizationHeader(ctx, token)` - Sets Authorization header with Bearer token

### Authentication Testing (`testutils/auth.go`)

Provides comprehensive utilities for testing authentication and authorization:

#### JWT Token Generation
- `GenerateTestJWT(userID, role)` - Generates valid test JWT tokens
- `GenerateTestRefreshToken(userID)` - Generates test refresh tokens
- `GenerateExpiredTestJWT(userID, role)` - Generates expired tokens for testing
- `GenerateInvalidTestJWT(userID, role)` - Generates invalid tokens for testing

#### Context Creation
- `CreateAuthenticatedContext(method, url, userID, role, body)` - Creates authenticated test context
- `CreateUnauthenticatedContext(method, url, body)` - Creates unauthenticated test context

#### Mock Middleware
- `MockAuthMiddleware()` - Mock authentication middleware for testing
- `MockRequireRole(allowedRoles...)` - Mock role-based authorization middleware

#### Test Data Creation
- `CreateTestUser(userType)` - Creates test user data
- `CreateTestUserWithID(userID, userType)` - Creates test user with specific ID
- `CreateTestTeam(ownerID)` - Creates test team data
- `CreateTestTeamMember(teamID, userID, role)` - Creates test team member data

#### Test Scenarios
- `SetupTestAuthScenarios(userID, role)` - Creates common auth test scenarios
- `SetupTestJWTSecrets()` - Sets up test JWT secrets
- `CleanupTestJWTSecrets()` - Cleans up test JWT secrets

#### Example Usage
```go
// Generate test JWT token
token, err := GenerateTestJWT("user-123", "admin")

// Create authenticated context
ctx, recorder, err := CreateAuthenticatedContext("GET", "/protected", "user-123", "admin", nil)

// Use mock middleware in router
router := gin.New()
router.Use(MockAuthMiddleware())
router.Use(MockRequireRole("admin", "editor"))
router.GET("/admin", handler)

// Set up common auth scenarios
scenarios, err := SetupTestAuthScenarios("user-123", "admin")
// Use scenarios.ValidToken, scenarios.ExpiredToken, scenarios.InvalidToken
```

### Test Fixtures (`testutils/fixtures.go`)

Provides predefined and customizable test data:

#### Predefined Fixtures
- `TestUsers` - Array of test user data
- `TestTeams` - Array of test team data  
- `TestTokens` - Array of test refresh token data
- `TestTeamMembers` - Array of test team member data

#### Fixture Creation Functions
- `CreateUserFixture(overrides...)` - Creates customizable user fixture
- `CreateTeamFixture(ownerID, overrides...)` - Creates customizable team fixture
- `CreateRefreshTokenFixture(userID, overrides...)` - Creates customizable token fixture

#### Example Usage
```go
// Create a custom user fixture
user := CreateUserFixture(func(u *UserFixture) {
    u.Name = "Custom User"
    u.Email = "custom@example.com"
    u.UserType = string(modelsv1.UserTypeTeam)
})

// Load fixtures into database
err := LoadFixtures(ctx, db, user)
```

### Test Suites (`testutils/suite.go`)

Provides base test suite structures with common setup and teardown:

#### BaseTestSuite
Basic test suite with database and HTTP router setup:

```go
type MyTestSuite struct {
    testutils.BaseTestSuite
}

func (s *MyTestSuite) TestSomething() {
    // Your test code here
    s.LoadTestFixtures() // Load predefined fixtures
    
    user := s.GetTestUser(0) // Get first test user
    // Test logic...
}

func TestMyTestSuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}
```

#### DatabaseTestSuite
Extended test suite with database package integration:

```go
type MyDatabaseTestSuite struct {
    testutils.DatabaseTestSuite
}

func (s *MyDatabaseTestSuite) TestDatabaseOperation() {
    tx := s.BeginTransaction()
    defer s.RollbackTransaction(tx)
    
    // Test database operations within transaction
}
```

### Mock Database (`testutils/database_mock.go`)

Provides SQL mock utilities for unit testing without real database:

```go
mockDB, err := SetupMockDB()
defer mockDB.Close()

// Set up expectations
mockDB.ExpectQuery("SELECT (.+) FROM users").
    WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
        AddRow("1", "Test User"))

// Execute code under test
// ...

// Verify expectations
err = mockDB.ExpectationsWereMet()
```

## Configuration

### Test Environment (`config/test.env`)
Contains environment variables for testing:
- Database configuration
- JWT secrets
- Logging settings
- Test-specific flags

### Test Configuration (`config/test_config.json`)
JSON configuration for test settings:
- Database connection parameters
- JWT configuration
- API settings
- Test behavior flags

## Running Tests

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Specific Package Tests
```bash
go test ./tests/testutils -v
go test ./tests/unit/api/v1 -v
```

### Run Tests with Coverage
```bash
go test ./tests/... -v -cover
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Tests in Parallel
```bash
go test ./tests/... -v -parallel 4
```

## Writing Tests

### Test Naming Convention
- Test files: `*_test.go`
- Test functions: `Test<FunctionName>_<Scenario>`
- Test suites: `<Package>TestSuite`

### Test Structure Pattern
```go
func TestFunctionName_Scenario(t *testing.T) {
    // Arrange
    // Setup test data and mocks
    
    // Act
    // Execute the function under test
    
    // Assert
    // Verify the results
}
```

### Example API Handler Test
```go
type AuthHandlerTestSuite struct {
    testutils.BaseTestSuite
}

func (s *AuthHandlerTestSuite) TestSignup_ValidInput_Success() {
    // Arrange
    s.LoadTestFixtures()
    
    signupData := modelsv1.SignupParams{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
        UserType: modelsv1.UserTypeIndividual,
    }
    
    // Act
    ctx, recorder := testutils.NewTestGinContext("POST", "/auth/signup", signupData)
    handler.Signup(ctx)
    
    // Assert
    s.Equal(http.StatusCreated, recorder.Code)
    
    var response map[string]interface{}
    err := testutils.ParseJSONResponse(recorder, &response)
    s.NoError(err)
    s.Equal("success", response["status"])
}

func (s *AuthHandlerTestSuite) TestProtectedEndpoint_WithAuth_Success() {
    // Arrange
    userID := "user-123"
    role := "admin"
    
    // Act
    ctx, recorder, err := testutils.CreateAuthenticatedContext("GET", "/protected", userID, role, nil)
    s.NoError(err)
    
    handler.ProtectedEndpoint(ctx)
    
    // Assert
    s.Equal(http.StatusOK, recorder.Code)
    
    var response map[string]interface{}
    err = testutils.ParseJSONResponse(recorder, &response)
    s.NoError(err)
    s.Equal("success", response["status"])
    s.Equal(userID, response["user_id"])
}

func (s *AuthHandlerTestSuite) TestProtectedEndpoint_WithoutAuth_Unauthorized() {
    // Act
    ctx, recorder := testutils.CreateUnauthenticatedContext("GET", "/protected", nil)
    handler.ProtectedEndpoint(ctx)
    
    // Assert
    s.Equal(http.StatusUnauthorized, recorder.Code)
}
```

### Example Business Logic Test
```go
func TestAuthService_Login_ValidCredentials_Success(t *testing.T) {
    // Arrange
    mockDB, err := testutils.SetupMockDB()
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

## Best Practices

### Test Organization
- Group related tests in test suites
- Use descriptive test names that explain the scenario
- Keep tests focused on a single behavior
- Use the Arrange-Act-Assert pattern

### Test Data Management
- Use fixtures for consistent test data
- Clean up test data between tests
- Use transactions for database tests when possible
- Avoid hardcoded values in tests

### Mocking Guidelines
- Mock external dependencies (database, HTTP clients, etc.)
- Use interfaces to enable easy mocking
- Verify mock expectations in tests
- Keep mocks simple and focused

### Performance Considerations
- Use in-memory databases for fast tests
- Run tests in parallel when safe
- Avoid unnecessary setup/teardown
- Use table-driven tests for multiple scenarios

## Troubleshooting

### Common Issues

#### Database Connection Errors
If you see database connection errors, ensure:
- Test database is available (for integration tests)
- Database credentials are correct in test configuration
- Database driver is imported

#### Mock Expectation Failures
If mock expectations fail:
- Verify the exact SQL query matches expectation
- Check parameter values and types
- Ensure all expectations are set before execution

#### Test Isolation Issues
If tests interfere with each other:
- Ensure proper cleanup between tests
- Use transactions that are rolled back
- Check for shared state between tests

### Getting Help
- Check test logs for detailed error messages
- Use `-v` flag for verbose test output
- Review the test utilities source code for examples
- Consult the Go testing documentation

## Additional Documentation

For comprehensive testing guidance, see these additional resources:

### Core Documentation
- **[Testing Best Practices Guide](TESTING_BEST_PRACTICES.md)** - Comprehensive patterns, conventions, and best practices
- **[Test Execution Guide](TEST_EXECUTION_GUIDE.md)** - Detailed test execution and CI/CD setup
- **[Coverage and Execution Guide](COVERAGE_AND_EXECUTION_GUIDE.md)** - Coverage reporting and performance testing

### Specialized Guides
- **[HTTP Testing Guide](testutils/HTTP_TESTING_GUIDE.md)** - HTTP handler and API testing utilities
- **[Test Suite Guide](testutils/SUITE_GUIDE.md)** - Base test suite structures and patterns
- **[Mock Framework Guide](testutils/mocks/README.md)** - Comprehensive mocking utilities
- **[Configuration Mocking Guide](testutils/mocks/CONFIG_MOCKING_GUIDE.md)** - Configuration and environment mocking

### Example Implementations
- **[API Handler Examples](examples/api_handler_example_test.go)** - Complete API testing examples
- **[Business Logic Examples](examples/business_logic_example_test.go)** - Service layer testing patterns
- **[Integration Testing Examples](examples/integration_testing_example_test.go)** - End-to-end workflow testing
- **[Model Validation Examples](examples/model_validation_example_test.go)** - Model validation and serialization testing

## Future Enhancements

Planned improvements to the testing framework:
- Enhanced integration test utilities
- Advanced performance testing utilities
- Automated test data generation tools
- Real-time test coverage monitoring
- Advanced mock external service utilities
- Test result analytics and reporting