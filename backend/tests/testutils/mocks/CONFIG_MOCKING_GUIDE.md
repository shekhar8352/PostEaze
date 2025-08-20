# Configuration Mocking Guide

This guide explains how to use the configuration mocking utilities in the PostEaze backend testing framework.

## Overview

The configuration mocking system provides comprehensive utilities for testing code that depends on configuration values and environment variables. It includes:

- **MockConfigClient**: Mock implementation of the configuration client interface
- **ConfigMockManager**: High-level manager for setting up configuration mocks
- **EnvironmentMockManager**: Manager for mocking environment variables
- **ConfigBuilder**: Fluent interface for building configuration mocks
- **ConfigTestHelper**: Comprehensive helper combining config and environment mocking
- **MockConfigProvider**: Provider for setting up common test configurations

## Quick Start

### Basic Configuration Mocking

```go
func TestMyFunction(t *testing.T) {
    // Create a configuration mock manager
    configManager := NewConfigMockManager()
    
    // Set up database configuration
    configManager.SetDatabaseConfig("sqlite3", ":memory:", 10, 5, time.Hour, time.Minute*30)
    
    // Get the mock client
    mockClient := configManager.GetMockClient()
    
    // Your test code here...
    // The mockClient can be used wherever the real config client is expected
}
```

### Environment Variable Mocking

```go
func TestWithEnvironment(t *testing.T) {
    envManager := NewEnvironmentMockManager()
    defer envManager.RestoreEnvironment() // Always restore after test
    
    // Set test environment variables
    envManager.SetEnv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
    envManager.SetEnv("JWT_SECRET", "test-secret-key")
    
    // Your test code here...
}
```

## Configuration Client Interface

The mock implements the actual configuration client interface used in the application:

```go
type ConfigClient interface {
    GetString(configName, key string) (string, error)
    GetIntD(configName, key string, defaultValue int64) int64
    GetMapD(configName, key string, defaultValue map[string]interface{}) map[string]interface{}
    GetBoolD(configName, key string, defaultValue bool) bool
    GetFloat64D(configName, key string, defaultValue float64) float64
    GetStringD(configName, key string, defaultValue string) string
    GetStringSliceD(configName, key string, defaultValue []string) []string
}
```

## Core Components

### 1. MockConfigClient

The basic mock implementation that can be used directly with testify expectations:

```go
mockClient := NewMockConfigClient()
mockClient.On("GetString", "database", "driverName").Return("postgres", nil)
mockClient.On("GetIntD", "database", "maxConnections", int64(10)).Return(int64(20))

// Use mockClient in your tests
driver, err := mockClient.GetString("database", "driverName")
maxConn := mockClient.GetIntD("database", "maxConnections", 10)
```

### 2. ConfigMockManager

High-level manager that simplifies setting up configuration values:

```go
manager := NewConfigMockManager()

// Set individual values
manager.SetValue("database", "driverName", "postgres")
manager.SetValue("database", "maxConnections", int64(20))

// Set database configuration in one call
manager.SetDatabaseConfig("postgres", "localhost:5432", 20, 10, time.Hour*2, time.Minute*45)

// Set API configuration
apiConfigs := map[string]interface{}{
    "getCatsFact": map[string]interface{}{
        "url":     "https://catfact.ninja/fact",
        "timeout": "30s",
        "retries": 3,
    },
}
manager.SetAPIConfig(apiConfigs)

// Set logger configuration
manager.SetLoggerConfig("info", map[string]interface{}{
    "console": true,
    "file":    true,
})

// Set errors for missing configurations
manager.SetError("database", "missingKey", errors.New("config not found"))
```

### 3. EnvironmentMockManager

Manages environment variable mocking with automatic restoration:

```go
envManager := NewEnvironmentMockManager()
defer envManager.RestoreEnvironment() // Always restore

// Set individual environment variables
envManager.SetEnv("DATABASE_URL", "postgres://localhost:5432/testdb")
envManager.SetEnv("JWT_SECRET", "test-secret")

// Unset environment variables
envManager.UnsetEnv("PRODUCTION_FLAG")

// Set up complete test environment
envManager.SetupTestEnvironment() // Sets APP_ENV=test, DB_DRIVER=sqlite3, etc.

// Get environment variables with defaults
dbURL := envManager.GetEnvWithDefault("DATABASE_URL", "sqlite://memory")
```

### 4. ConfigBuilder

Fluent interface for building configuration mocks:

```go
configManager := NewConfigBuilder().
    WithDatabaseConfig("sqlite3", ":memory:", 5, 2).
    WithLoggerConfig("debug").
    WithAPIConfig(map[string]interface{}{
        "timeout": "45s",
        "retries": 5,
    }).
    WithCustomConfig("application", "name", "test-app").
    WithError("missing", "key", errors.New("not found")).
    Build()

mockClient := configManager.GetMockClient()
```

### 5. ConfigTestHelper

Comprehensive helper that combines configuration and environment mocking:

```go
helper := NewConfigTestHelper()
defer helper.Cleanup() // Always cleanup

// Set up default test configuration
helper.SetupDefaultTestConfig()

// Access managers
configManager := helper.GetConfigManager()
envManager := helper.GetEnvManager()

// Add custom configuration
configManager.SetValue("custom", "feature", "enabled")

// Verify configuration access
accessed := helper.VerifyConfigurationAccess(map[string][]string{
    "database": {"driverName", "url"},
    "logger":   {"level"},
})
```

### 6. MockConfigProvider

Provider for setting up common test configurations:

```go
provider := NewMockConfigProvider()

// Set up all configurations for testing
provider.SetupAllConfigsForTesting()

// Or set up specific configurations
provider.SetupDatabaseConfigForTesting()
provider.SetupAPIConfigForTesting()
provider.SetupApplicationConfigForTesting()

mockClient := provider.GetMockClient()
```

## Common Usage Patterns

### Testing Database Initialization

```go
func TestDatabaseInit(t *testing.T) {
    configManager := NewConfigMockManager()
    configManager.SetDatabaseConfig("sqlite3", ":memory:", 10, 5, time.Hour, time.Minute*30)
    
    mockClient := configManager.GetMockClient()
    
    // Replace the global config client with the mock
    // (This would depend on your application's dependency injection setup)
    
    // Test your database initialization code
    err := initDatabase(context.Background())
    assert.NoError(t, err)
    
    // Verify that configuration was accessed
    assert.True(t, configManager.VerifyConfigAccessed("database", "driverName"))
    assert.True(t, configManager.VerifyConfigAccessed("database", "url"))
}
```

### Testing API Configuration

```go
func TestAPIClientInit(t *testing.T) {
    configManager := NewConfigMockManager()
    
    apiConfig := map[string]interface{}{
        "getCatsFact": map[string]interface{}{
            "url":     "https://test-api.example.com",
            "timeout": "30s",
            "retries": 3,
        },
    }
    configManager.SetAPIConfig(apiConfig)
    
    mockClient := configManager.GetMockClient()
    
    // Test your API client initialization
    initHttp(context.Background())
    
    // Verify configuration access
    assert.True(t, configManager.VerifyConfigAccessed("api", "getCatsFact"))
}
```

### Testing with Environment Variables

```go
func TestWithEnvironmentSubstitution(t *testing.T) {
    helper := NewConfigTestHelper()
    defer helper.Cleanup()
    
    // Set up environment variables
    envManager := helper.GetEnvManager()
    envManager.SetEnv("DB_HOST", "test-host")
    envManager.SetEnv("DB_PORT", "5432")
    
    // Set up configuration with environment variable placeholders
    configManager := helper.GetConfigManager()
    configManager.SetValue("database", "url", "postgres://${DB_HOST}:${DB_PORT}/testdb")
    
    // Test your code that uses environment variable substitution
    // (The actual substitution would happen in your application code)
}
```

### Testing Error Scenarios

```go
func TestConfigurationErrors(t *testing.T) {
    configManager := NewConfigMockManager()
    
    // Set up error for missing configuration
    configManager.SetError("database", "missingKey", errors.New("configuration not found"))
    
    mockClient := configManager.GetMockClient()
    
    // Test error handling
    _, err := mockClient.GetString("database", "missingKey")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "configuration not found")
}
```

## Best Practices

### 1. Always Clean Up

```go
func TestSomething(t *testing.T) {
    helper := NewConfigTestHelper()
    defer helper.Cleanup() // Always cleanup
    
    envManager := NewEnvironmentMockManager()
    defer envManager.RestoreEnvironment() // Always restore
    
    // Your test code...
}
```

### 2. Use Appropriate Abstraction Level

- Use `MockConfigClient` directly for simple, specific mocking
- Use `ConfigMockManager` for setting up multiple related configurations
- Use `ConfigBuilder` for complex configuration setups
- Use `ConfigTestHelper` for comprehensive test scenarios
- Use `MockConfigProvider` for common test configurations

### 3. Verify Configuration Access

```go
// Verify that your code accessed the expected configuration keys
assert.True(t, configManager.VerifyConfigAccessed("database", "driverName"))
assert.Equal(t, 1, configManager.GetAccessCount("database", "url"))
```

### 4. Use Realistic Test Data

```go
// Use realistic configuration values that match your application's needs
configManager.SetDatabaseConfig(
    "postgres",                    // Realistic driver
    "localhost:5432",             // Realistic connection string
    20,                           // Realistic max connections
    10,                           // Realistic idle connections
    time.Hour*2,                  // Realistic connection lifetime
    time.Minute*45,               // Realistic idle time
)
```

### 5. Test Both Success and Error Cases

```go
func TestConfigurationHandling(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        configManager := NewConfigMockManager()
        configManager.SetValue("database", "url", "postgres://localhost:5432/db")
        // Test success scenario
    })
    
    t.Run("error case", func(t *testing.T) {
        configManager := NewConfigMockManager()
        configManager.SetError("database", "url", errors.New("config error"))
        // Test error handling
    })
}
```

## Integration with Application Code

To use these mocks effectively, your application code should use dependency injection or a similar pattern to allow replacing the configuration client during testing:

```go
// In your application code
type MyService struct {
    configClient ConfigClient
}

func NewMyService(configClient ConfigClient) *MyService {
    return &MyService{configClient: configClient}
}

// In your tests
func TestMyService(t *testing.T) {
    configManager := NewConfigMockManager()
    configManager.SetValue("service", "timeout", "30s")
    
    service := NewMyService(configManager.GetMockClient())
    
    // Test your service...
}
```

## Error Types

The mocking framework provides common error types for testing:

```go
var (
    ErrMockConfigNotFound    = errors.New("mock: configuration key not found")
    ErrMockConfigInvalid     = errors.New("mock: invalid configuration value")
    ErrMockConfigTypeMismatch = errors.New("mock: configuration type mismatch")
    ErrMockConfigConnection  = errors.New("mock: configuration service connection failed")
)
```

## Type Conversion Utilities

The framework includes utilities for converting between different types:

```go
converter := NewTypeConverter()

// Convert to string
str := converter.ToString(123) // "123"

// Convert to int
num, err := converter.ToInt("456") // 456, nil

// Convert to bool
flag, err := converter.ToBool("true") // true, nil

// Convert to duration
duration, err := converter.ToDuration("30s") // 30*time.Second, nil
```

This comprehensive mocking system provides everything needed to test configuration-dependent code effectively while maintaining clean, readable, and maintainable tests.