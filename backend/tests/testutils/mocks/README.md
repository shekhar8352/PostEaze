# Mock Framework Documentation

This directory contains comprehensive mock implementations for testing the PostEaze backend application. The mock framework provides utilities for mocking database operations, HTTP clients, and configuration management.

## Overview

The mock framework consists of three main components:

1. **Database Mocks** (`database.go`) - Mock implementations for database operations
2. **HTTP Client Mocks** (`http.go`) - Mock implementations for HTTP client operations
3. **Configuration Mocks** (`config.go`) - Mock implementations for configuration management

## Database Mocks

### Basic Usage

```go
import "github.com/shekhar8352/PostEaze/tests/testutils/mocks"

func TestDatabaseOperation(t *testing.T) {
    // Create a mock database
    mockDB := mocks.NewMockDatabase()
    
    // Setup expectation
    mockDB.On("QueryRaw", mock.Anything, mock.Anything, 1).Return(nil)
    
    // Use the mock in your test
    err := mockDB.QueryRaw(context.Background(), nil, 1)
    
    // Assert
    assert.NoError(t, err)
    mockDB.AssertExpectations(t)
}
```

### Using Database Mock Manager

```go
func TestWithDatabaseManager(t *testing.T) {
    manager := mocks.NewMockDatabaseManager()
    
    // Setup successful query
    manager.SetupSuccessfulQuery(nil, 1)
    
    // Setup failed query
    manager.SetupFailedQuery(nil, 2, database.ErrNoRecords)
    
    // Use the mock database
    mockDB := manager.GetMockDB()
    
    // Test successful operation
    err := mockDB.QueryRaw(context.Background(), nil, 1)
    assert.NoError(t, err)
    
    // Test failed operation
    err = mockDB.QueryRaw(context.Background(), nil, 2)
    assert.Error(t, err)
}
```

### Repository Helper

```go
func TestRepositoryOperations(t *testing.T) {
    helper := mocks.NewMockRepositoryHelper()
    
    // Setup user creation
    helper.SetupUserCreation("user-123")
    
    // Setup user lookup
    helper.SetupUserLookup(true, nil) // user found
    
    // Use in your repository tests
    dbManager := helper.GetDatabaseManager()
    // ... test repository operations
}
```

## HTTP Client Mocks

### Basic Usage

```go
func TestHTTPRequest(t *testing.T) {
    mockClient := mocks.NewMockHTTPClient()
    
    // Create expected response
    expectedResp := &http.Response{
        StatusCode: 200,
        Body:       io.NopCloser(strings.NewReader("success")),
    }
    
    // Setup expectation
    req, _ := http.NewRequest("GET", "http://example.com", nil)
    mockClient.On("Do", req).Return(expectedResp, nil)
    
    // Test
    resp, err := mockClient.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

### Using HTTP Client Manager

```go
func TestWithHTTPManager(t *testing.T) {
    manager := mocks.NewHTTPClientMockManager()
    
    // Setup JSON response
    responseData := map[string]string{"status": "success"}
    manager.SetupJSONResponse("GET", "http://api.example.com/data", 200, responseData)
    
    // Setup error response
    manager.SetupErrorResponse("POST", "http://api.example.com/error", errors.New("network error"))
    
    // Use the mock client
    mockClient := manager.GetMockClient()
    
    // Test successful request
    req, _ := http.NewRequest("GET", "http://api.example.com/data", nil)
    resp, err := mockClient.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verify request was made
    assert.True(t, manager.VerifyRequestCalled("GET", "http://api.example.com/data"))
}
```

### HTTP Scenario Manager

```go
func TestHTTPScenarios(t *testing.T) {
    scenarioManager := mocks.NewHTTPScenarioManager()
    
    // Add scenarios
    scenarioManager.AddScenario(mocks.HTTPMockScenario{
        Name:   "Get User",
        Method: "GET",
        URL:    "http://api.example.com/users/123",
        Response: &mocks.MockHTTPResponse{
            StatusCode: 200,
            Body:       `{"id": "123", "name": "John"}`,
            Headers:    map[string]string{"Content-Type": "application/json"},
        },
    })
    
    // Use the manager
    manager := scenarioManager.GetManager()
    mockClient := manager.GetMockClient()
    
    // Test the scenario
    req, _ := http.NewRequest("GET", "http://api.example.com/users/123", nil)
    resp, err := mockClient.Do(req)
    assert.NoError(t, err)
    
    // Verify all scenarios were called
    failures := scenarioManager.VerifyAllScenarios()
    assert.Empty(t, failures)
}
```

## Configuration Mocks

### Basic Usage

```go
func TestConfiguration(t *testing.T) {
    mockConfig := mocks.NewMockConfigClient()
    
    // Setup expectations
    mockConfig.On("GetString", "database.url").Return("localhost:5432", nil)
    mockConfig.On("GetInt", "database.maxConnections").Return(10, nil)
    
    // Test
    url, err := mockConfig.GetString("database.url")
    assert.NoError(t, err)
    assert.Equal(t, "localhost:5432", url)
    
    maxConn, err := mockConfig.GetInt("database.maxConnections")
    assert.NoError(t, err)
    assert.Equal(t, 10, maxConn)
}
```

### Using Configuration Manager

```go
func TestWithConfigManager(t *testing.T) {
    manager := mocks.NewConfigMockManager()
    
    // Set configuration values
    manager.SetValue("app.name", "PostEaze")
    manager.SetValue("app.port", 8080)
    manager.SetValue("app.debug", true)
    
    // Set database configuration
    manager.SetDatabaseConfig("postgres", "localhost:5432", 10, 5, time.Hour, time.Minute*30)
    
    // Set error for missing key
    manager.SetError("missing.key", mocks.ErrMockConfigNotFound)
    
    // Use the mock client
    mockClient := manager.GetMockClient()
    
    // Test configuration access
    name, err := mockClient.GetString("app.name")
    assert.NoError(t, err)
    assert.Equal(t, "PostEaze", name)
    
    // Test error case
    _, err = mockClient.GetString("missing.key")
    assert.Error(t, err)
    assert.Equal(t, mocks.ErrMockConfigNotFound, err)
}
```

### Environment Variable Mocking

```go
func TestEnvironmentVariables(t *testing.T) {
    envManager := mocks.NewEnvironmentMockManager()
    defer envManager.RestoreEnvironment() // Always restore after test
    
    // Set test environment variables
    envManager.SetEnv("DATABASE_URL", "test-db-url")
    envManager.SetEnv("JWT_SECRET", "test-secret")
    
    // Or setup complete test environment
    envManager.SetupTestEnvironment()
    
    // Test environment access
    dbURL := envManager.GetEnv("DATABASE_URL")
    assert.Equal(t, "test-db-url", dbURL)
    
    // Test with default
    timeout := envManager.GetEnvWithDefault("REQUEST_TIMEOUT", "30s")
    assert.Equal(t, "30s", timeout)
}
```

### Configuration Builder

```go
func TestConfigBuilder(t *testing.T) {
    // Build configuration using fluent interface
    manager := mocks.NewConfigBuilder().
        WithDatabaseConfig("postgres", "localhost:5432", 10, 5).
        WithLoggerConfig("info").
        WithAPIConfig(map[string]interface{}{
            "port":    8080,
            "timeout": "30s",
        }).
        WithCustomConfig("feature.enabled", true).
        Build()
    
    // Test the configuration
    driver, err := manager.GetMockClient().GetString("database.driverName")
    assert.NoError(t, err)
    assert.Equal(t, "postgres", driver)
    
    level, err := manager.GetMockClient().GetString("logger.level")
    assert.NoError(t, err)
    assert.Equal(t, "info", level)
}
```

### Configuration Test Helper

```go
func TestWithConfigHelper(t *testing.T) {
    helper := mocks.NewConfigTestHelper()
    defer helper.Cleanup() // Always cleanup after test
    
    // Setup default test configuration
    helper.SetupDefaultTestConfig()
    
    // Access configuration and environment managers
    configManager := helper.GetConfigManager()
    envManager := helper.GetEnvManager()
    
    // Test configuration
    driver, err := configManager.GetMockClient().GetString("database.driverName")
    assert.NoError(t, err)
    assert.Equal(t, "sqlite3", driver)
    
    // Test environment
    appEnv := envManager.GetEnv("APP_ENV")
    assert.Equal(t, "test", appEnv)
    
    // Verify configuration access
    accessResults := helper.VerifyConfigurationAccess("database.driverName", "logger.level")
    assert.True(t, accessResults["database.driverName"])
}
```

## Best Practices

1. **Always use defer for cleanup**: When using environment mocks, always defer the cleanup:
   ```go
   envManager := mocks.NewEnvironmentMockManager()
   defer envManager.RestoreEnvironment()
   ```

2. **Assert expectations**: Always assert that mock expectations were met:
   ```go
   mockDB.AssertExpectations(t)
   manager.AssertExpectations(t)
   ```

3. **Use builders for complex setups**: For complex configuration setups, use the builder pattern:
   ```go
   manager := mocks.NewConfigBuilder().
       WithDatabaseConfig(...).
       WithLoggerConfig(...).
       Build()
   ```

4. **Verify interactions**: Use verification methods to ensure mocks were called as expected:
   ```go
   assert.True(t, manager.VerifyRequestCalled("GET", "http://example.com"))
   assert.Equal(t, 2, manager.GetRequestCount("POST", "http://api.com"))
   ```

5. **Reset mocks between tests**: If reusing mocks, reset them between tests:
   ```go
   manager.Reset()
   ```

## Common Error Types

The mock framework provides common error types for testing:

- **Database Errors**: `ErrMockNoRecords`, `ErrMockNoRowsAffected`, `ErrMockConnection`, `ErrMockTransaction`, `ErrMockConstraint`
- **HTTP Errors**: `ErrMockHTTPTimeout`, `ErrMockHTTPConnection`, `ErrMockUnauthorized`, `ErrMockNotFound`, `ErrMockServerError`, `ErrMockBadRequest`
- **Configuration Errors**: `ErrMockConfigNotFound`, `ErrMockConfigInvalid`, `ErrMockConfigTypeMismatch`, `ErrMockConfigConnection`

## Running Tests

To run the mock tests:

```bash
# From the backend directory
go test -mod=mod -v ./tests/testutils/mocks/

# Or run all tests
go test -mod=mod -v ./tests/...
```

## Integration with Test Suites

The mocks can be easily integrated with testify test suites:

```go
type MyTestSuite struct {
    suite.Suite
    dbManager     *mocks.MockDatabaseManager
    httpManager   *mocks.HTTPClientMockManager
    configHelper  *mocks.ConfigTestHelper
}

func (s *MyTestSuite) SetupTest() {
    s.dbManager = mocks.NewMockDatabaseManager()
    s.httpManager = mocks.NewHTTPClientMockManager()
    s.configHelper = mocks.NewConfigTestHelper()
    s.configHelper.SetupDefaultTestConfig()
}

func (s *MyTestSuite) TearDownTest() {
    s.dbManager.Reset()
    s.httpManager.Reset()
    s.configHelper.Cleanup()
}

func TestMyTestSuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}
```

This mock framework provides comprehensive testing utilities that make it easy to write reliable, isolated unit tests for the PostEaze backend application.