package mocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ExampleConfigMockUsage demonstrates how to use configuration mocks in tests
func ExampleConfigMockUsage(t *testing.T) {
	// Create a configuration mock manager
	configManager := NewConfigMockManager()
	
	// Set up database configuration
	configManager.SetDatabaseConfig("postgres", "localhost:5432", 20, 10, time.Hour*2, time.Minute*45)
	
	// Set up API configuration
	apiConfigs := map[string]interface{}{
		"getCatsFact": map[string]interface{}{
			"url":     "https://catfact.ninja/fact",
			"timeout": "30s",
			"retries": 3,
		},
	}
	configManager.SetAPIConfig(apiConfigs)
	
	// Set up logger configuration
	loggerParams := map[string]interface{}{
		"console": true,
		"file":    true,
		"level":   "info",
	}
	configManager.SetLoggerConfig("info", loggerParams)
	
	// Use the mock client
	mockClient := configManager.GetMockClient()
	
	// Test database configuration
	driver, err := mockClient.GetString("database", "driverName")
	assert.NoError(t, err)
	assert.Equal(t, "postgres", driver)
	
	maxConnections := mockClient.GetIntD("database", "maxOpenConnections", 1)
	assert.Equal(t, int64(20), maxConnections)
	
	// Test API configuration
	apiConfig := mockClient.GetMapD("api", "getCatsFact", nil)
	assert.NotNil(t, apiConfig)
	
	// Verify configuration access
	assert.True(t, configManager.VerifyConfigAccessed("database", "driverName"))
	assert.True(t, configManager.VerifyConfigAccessed("database", "maxOpenConnections"))
	assert.True(t, configManager.VerifyConfigAccessed("api", "getCatsFact"))
}

// ExampleEnvironmentMockUsage demonstrates how to use environment variable mocks
func ExampleEnvironmentMockUsage(t *testing.T) {
	// Create an environment mock manager
	envManager := NewEnvironmentMockManager()
	defer envManager.RestoreEnvironment() // Always restore after test
	
	// Set up test environment variables
	envManager.SetEnv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	envManager.SetEnv("JWT_SECRET", "test-secret-key-for-testing")
	envManager.SetEnv("LOG_LEVEL", "debug")
	envManager.SetEnv("API_TIMEOUT", "60s")
	
	// Test environment variables
	dbURL := envManager.GetEnv("DATABASE_URL")
	assert.Equal(t, "postgres://test:test@localhost:5432/testdb", dbURL)
	
	jwtSecret := envManager.GetEnv("JWT_SECRET")
	assert.Equal(t, "test-secret-key-for-testing", jwtSecret)
	
	// Test with default values
	port := envManager.GetEnvWithDefault("PORT", "8080")
	assert.Equal(t, "8080", port) // Should return default since PORT is not set
	
	logLevel := envManager.GetEnvWithDefault("LOG_LEVEL", "info")
	assert.Equal(t, "debug", logLevel) // Should return set value
}

// ExampleConfigBuilderUsage demonstrates the fluent configuration builder
func ExampleConfigBuilderUsage(t *testing.T) {
	// Build configuration using fluent interface
	configManager := NewConfigBuilder().
		WithDatabaseConfig("sqlite3", ":memory:", 5, 2).
		WithLoggerConfig("debug").
		WithAPIConfig(map[string]interface{}{
			"timeout": "45s",
			"retries": 5,
		}).
		WithCustomConfig("application", "name", "test-app").
		WithCustomConfig("application", "version", "1.0.0-test").
		Build()
	
	mockClient := configManager.GetMockClient()
	
	// Test the built configuration
	driver, err := mockClient.GetString("database", "driverName")
	assert.NoError(t, err)
	assert.Equal(t, "sqlite3", driver)
	
	logLevel, err := mockClient.GetString("logger", "level")
	assert.NoError(t, err)
	assert.Equal(t, "debug", logLevel)
	
	appName, err := mockClient.GetString("application", "name")
	assert.NoError(t, err)
	assert.Equal(t, "test-app", appName)
}

// ExampleConfigTestHelperUsage demonstrates the comprehensive test helper
func ExampleConfigTestHelperUsage(t *testing.T) {
	// Create a comprehensive test helper
	helper := NewConfigTestHelper()
	defer helper.Cleanup() // Always cleanup after test
	
	// Set up default test configuration
	helper.SetupDefaultTestConfig()
	
	// Access configuration manager
	configManager := helper.GetConfigManager()
	mockClient := configManager.GetMockClient()
	
	// Test database configuration
	driver, err := mockClient.GetString("database", "driverName")
	assert.NoError(t, err)
	assert.Equal(t, "sqlite3", driver)
	
	// Access environment manager
	envManager := helper.GetEnvManager()
	
	// Test environment variables
	appEnv := envManager.GetEnv("APP_ENV")
	assert.Equal(t, "test", appEnv)
	
	// Add custom configuration
	configManager.SetValue("custom", "feature", "enabled")
	
	customFeature, err := mockClient.GetString("custom", "feature")
	assert.NoError(t, err)
	assert.Equal(t, "enabled", customFeature)
}

// ExampleMockConfigProviderUsage demonstrates the mock config provider
func ExampleMockConfigProviderUsage(t *testing.T) {
	// Create a mock config provider
	provider := NewMockConfigProvider()
	
	// Set up all configurations for testing
	provider.SetupAllConfigsForTesting()
	
	mockClient := provider.GetMockClient()
	
	// Test database configuration
	driver, err := mockClient.GetString("database", "driverName")
	assert.NoError(t, err)
	assert.Equal(t, "sqlite3", driver)
	
	// Test API configuration
	apiConfig := mockClient.GetMapD("api", "getCatsFact", nil)
	assert.NotNil(t, apiConfig)
	assert.Equal(t, "http://test-api.example.com", apiConfig["url"])
	
	// Test application configuration
	appName := mockClient.GetStringD("application", "name", "")
	assert.Equal(t, "post-eaze-test", appName)
}

// ExampleErrorHandlingInConfigMocks demonstrates error handling in config mocks
func ExampleErrorHandlingInConfigMocks(t *testing.T) {
	configManager := NewConfigMockManager()
	
	// Set up an error for a missing configuration
	configManager.SetError("database", "missingKey", ErrMockConfigNotFound)
	
	mockClient := configManager.GetMockClient()
	
	// Test error handling
	value, err := mockClient.GetString("database", "missingKey")
	assert.Error(t, err)
	assert.Equal(t, ErrMockConfigNotFound, err)
	assert.Empty(t, value)
}

// ExampleComplexConfigurationScenario demonstrates a complex testing scenario
func ExampleComplexConfigurationScenario(t *testing.T) {
	// Set up a complex test scenario with multiple configurations
	helper := NewConfigTestHelper()
	defer helper.Cleanup()
	
	// Set up environment variables
	envManager := helper.GetEnvManager()
	envManager.SetEnv("DATABASE_HOST", "test-db-host")
	envManager.SetEnv("DATABASE_PORT", "5432")
	envManager.SetEnv("API_KEY", "test-api-key-12345")
	
	// Set up configuration values
	configManager := helper.GetConfigManager()
	
	// Database configuration with environment variable substitution
	configManager.SetValue("database", "host", "${DATABASE_HOST}")
	configManager.SetValue("database", "port", "${DATABASE_PORT}")
	configManager.SetValue("database", "driverName", "postgres")
	
	// API configuration
	configManager.SetValue("api", "key", "${API_KEY}")
	configManager.SetValue("api", "baseURL", "https://api.example.com")
	
	// Application configuration
	configManager.SetValue("application", "environment", "test")
	configManager.SetValue("application", "debug", true)
	
	mockClient := configManager.GetMockClient()
	
	// Test the complex configuration
	host, err := mockClient.GetString("database", "host")
	assert.NoError(t, err)
	assert.Equal(t, "${DATABASE_HOST}", host) // Mock returns as-is, env substitution would happen in real code
	
	apiKey, err := mockClient.GetString("api", "key")
	assert.NoError(t, err)
	assert.Equal(t, "${API_KEY}", apiKey)
	
	debug := mockClient.GetBoolD("application", "debug", false)
	assert.True(t, debug)
	
	// Verify that configurations were accessed
	assert.True(t, configManager.VerifyConfigAccessed("database", "host"))
	assert.True(t, configManager.VerifyConfigAccessed("api", "key"))
	assert.True(t, configManager.VerifyConfigAccessed("application", "debug"))
	
	// Verify access counts
	assert.Equal(t, 1, configManager.GetAccessCount("database", "host"))
	assert.Equal(t, 1, configManager.GetAccessCount("api", "key"))
}

// ExampleTypeConversionInMocks demonstrates type conversion utilities
func ExampleTypeConversionInMocks(t *testing.T) {
	converter := NewTypeConverter()
	
	// Test string conversion
	strValue := converter.ToString(123)
	assert.Equal(t, "123", strValue)
	
	strValue = converter.ToString(true)
	assert.Equal(t, "true", strValue)
	
	// Test int conversion
	intValue, err := converter.ToInt("456")
	assert.NoError(t, err)
	assert.Equal(t, 456, intValue)
	
	intValue, err = converter.ToInt(789.0)
	assert.NoError(t, err)
	assert.Equal(t, 789, intValue)
	
	// Test bool conversion
	boolValue, err := converter.ToBool("true")
	assert.NoError(t, err)
	assert.True(t, boolValue)
	
	boolValue, err = converter.ToBool(1)
	assert.NoError(t, err)
	assert.True(t, boolValue)
	
	// Test duration conversion
	duration, err := converter.ToDuration("30s")
	assert.NoError(t, err)
	assert.Equal(t, 30*time.Second, duration)
	
	duration, err = converter.ToDuration(60)
	assert.NoError(t, err)
	assert.Equal(t, 60*time.Second, duration)
}

// TestConfigMockExamples runs all the example functions as tests
func TestConfigMockExamples(t *testing.T) {
	t.Run("ConfigMockUsage", ExampleConfigMockUsage)
	t.Run("EnvironmentMockUsage", ExampleEnvironmentMockUsage)
	t.Run("ConfigBuilderUsage", ExampleConfigBuilderUsage)
	t.Run("ConfigTestHelperUsage", ExampleConfigTestHelperUsage)
	t.Run("MockConfigProviderUsage", ExampleMockConfigProviderUsage)
	t.Run("ErrorHandlingInConfigMocks", ExampleErrorHandlingInConfigMocks)
	t.Run("ComplexConfigurationScenario", ExampleComplexConfigurationScenario)
	t.Run("TypeConversionInMocks", ExampleTypeConversionInMocks)
}