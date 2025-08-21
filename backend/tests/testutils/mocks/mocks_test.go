package mocks

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestMockDatabase tests the database mock functionality
func TestMockDatabase(t *testing.T) {
	t.Run("QueryRaw success", func(t *testing.T) {
		mockDB := NewMockDatabase()
		ctx := context.Background()
		
		// Setup expectation
		mockDB.On("QueryRaw", ctx, mock.Anything, 1).Return(nil)
		
		// Execute
		err := mockDB.QueryRaw(ctx, nil, 1)
		
		// Assert
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	
	t.Run("QueryRaw error", func(t *testing.T) {
		mockDB := NewMockDatabase()
		ctx := context.Background()
		expectedErr := errors.New("database error")
		
		// Setup expectation
		mockDB.On("QueryRaw", ctx, mock.Anything, 1).Return(expectedErr)
		
		// Execute
		err := mockDB.QueryRaw(ctx, nil, 1)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockDB.AssertExpectations(t)
	})
}

// TestMockDatabaseManager tests the database manager functionality
func TestMockDatabaseManager(t *testing.T) {
	t.Run("SetupSuccessfulQuery", func(t *testing.T) {
		manager := NewMockDatabaseManager()
		ctx := context.Background()
		
		// Setup
		manager.SetupSuccessfulQuery(nil, 1)
		
		// Execute
		err := manager.GetMockDB().QueryRaw(ctx, nil, 1)
		
		// Assert
		assert.NoError(t, err)
	})
	
	t.Run("SetupFailedQuery", func(t *testing.T) {
		manager := NewMockDatabaseManager()
		ctx := context.Background()
		expectedErr := database.ErrNoRecords
		
		// Setup
		manager.SetupFailedQuery(nil, 1, expectedErr)
		
		// Execute
		err := manager.GetMockDB().QueryRaw(ctx, nil, 1)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestMockHTTPClient tests the HTTP client mock functionality
func TestMockHTTPClient(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		mockClient := NewMockHTTPClient()
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		expectedResp := &http.Response{StatusCode: 200}
		
		// Setup expectation
		mockClient.On("Do", req).Return(expectedResp, nil)
		
		// Execute
		resp, err := mockClient.Do(req)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
		mockClient.AssertExpectations(t)
	})
	
	t.Run("request error", func(t *testing.T) {
		mockClient := NewMockHTTPClient()
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		expectedErr := errors.New("connection failed")
		
		// Setup expectation
		mockClient.On("Do", req).Return((*http.Response)(nil), expectedErr)
		
		// Execute
		resp, err := mockClient.Do(req)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
		mockClient.AssertExpectations(t)
	})
}

// TestHTTPClientMockManager tests the HTTP client manager functionality
func TestHTTPClientMockManager(t *testing.T) {
	t.Run("SetupSuccessResponse", func(t *testing.T) {
		manager := NewHTTPClientMockManager()
		
		// Setup
		manager.SetupSuccessResponse("GET", "http://example.com", "success")
		
		// Create request
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		
		// Execute
		resp, err := manager.GetMockClient().Do(req)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	
	t.Run("SetupErrorResponse", func(t *testing.T) {
		manager := NewHTTPClientMockManager()
		expectedErr := errors.New("network error")
		
		// Setup
		manager.SetupErrorResponse("GET", "http://example.com", expectedErr)
		
		// Create request
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		
		// Execute
		resp, err := manager.GetMockClient().Do(req)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
	})
}

// TestMockConfigClient tests the configuration client mock functionality
func TestMockConfigClient(t *testing.T) {
	t.Run("GetString success", func(t *testing.T) {
		mockConfig := NewMockConfigClient()
		expectedValue := "test-value"
		
		// Setup expectation
		mockConfig.On("GetString", "database", "driverName").Return(expectedValue, nil)
		
		// Execute
		value, err := mockConfig.GetString("database", "driverName")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
		mockConfig.AssertExpectations(t)
	})
	
	t.Run("GetIntD success", func(t *testing.T) {
		mockConfig := NewMockConfigClient()
		expectedValue := int64(42)
		
		// Setup expectation
		mockConfig.On("GetIntD", "database", "maxConnections", int64(10)).Return(expectedValue)
		
		// Execute
		value := mockConfig.GetIntD("database", "maxConnections", 10)
		
		// Assert
		assert.Equal(t, expectedValue, value)
		mockConfig.AssertExpectations(t)
	})
	
	t.Run("GetBoolD success", func(t *testing.T) {
		mockConfig := NewMockConfigClient()
		expectedValue := true
		
		// Setup expectation
		mockConfig.On("GetBoolD", "application", "debug", false).Return(expectedValue)
		
		// Execute
		value := mockConfig.GetBoolD("application", "debug", false)
		
		// Assert
		assert.Equal(t, expectedValue, value)
		mockConfig.AssertExpectations(t)
	})
	
	t.Run("GetMapD success", func(t *testing.T) {
		mockConfig := NewMockConfigClient()
		expectedValue := map[string]interface{}{
			"url":     "http://example.com",
			"timeout": "30s",
		}
		
		// Setup expectation
		mockConfig.On("GetMapD", "api", "getCatsFact", mock.Anything).Return(expectedValue)
		
		// Execute
		value := mockConfig.GetMapD("api", "getCatsFact", nil)
		
		// Assert
		assert.Equal(t, expectedValue, value)
		mockConfig.AssertExpectations(t)
	})
}

// TestConfigMockManager tests the configuration manager functionality
func TestConfigMockManager(t *testing.T) {
	t.Run("SetValue string", func(t *testing.T) {
		manager := NewConfigMockManager()
		
		// Setup
		manager.SetValue("database", "driverName", "hello")
		
		// Execute
		value, err := manager.GetMockClient().GetString("database", "driverName")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "hello", value)
	})
	
	t.Run("SetValue int", func(t *testing.T) {
		manager := NewConfigMockManager()
		
		// Setup
		manager.SetValue("database", "maxConnections", 123)
		
		// Execute
		value := manager.GetMockClient().GetIntD("database", "maxConnections", 10)
		
		// Assert
		assert.Equal(t, int64(123), value)
	})
	
	t.Run("SetError", func(t *testing.T) {
		manager := NewConfigMockManager()
		expectedErr := errors.New("config error")
		
		// Setup
		manager.SetError("database", "missingKey", expectedErr)
		
		// Execute
		value, err := manager.GetMockClient().GetString("database", "missingKey")
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, value)
	})
	
	t.Run("SetDatabaseConfig", func(t *testing.T) {
		manager := NewConfigMockManager()
		
		// Setup
		manager.SetDatabaseConfig("postgres", "localhost:5432", 10, 5, time.Hour, time.Minute*30)
		
		// Execute and Assert
		driver, err := manager.GetMockClient().GetString("database", "driverName")
		assert.NoError(t, err)
		assert.Equal(t, "postgres", driver)
		
		url, err := manager.GetMockClient().GetString("database", "url")
		assert.NoError(t, err)
		assert.Equal(t, "localhost:5432", url)
		
		maxOpen := manager.GetMockClient().GetIntD("database", "maxOpenConnections", 1)
		assert.Equal(t, int64(10), maxOpen)
	})
}

// TestEnvironmentMockManager tests the environment manager functionality
func TestEnvironmentMockManager(t *testing.T) {
	t.Run("SetEnv and GetEnv", func(t *testing.T) {
		manager := NewEnvironmentMockManager()
		defer manager.RestoreEnvironment()
		
		// Setup
		manager.SetEnv("TEST_VAR", "test-value")
		
		// Execute
		value := manager.GetEnv("TEST_VAR")
		
		// Assert
		assert.Equal(t, "test-value", value)
	})
	
	t.Run("GetEnvWithDefault", func(t *testing.T) {
		manager := NewEnvironmentMockManager()
		defer manager.RestoreEnvironment()
		
		// Test with existing value
		manager.SetEnv("EXISTING_VAR", "existing-value")
		value := manager.GetEnvWithDefault("EXISTING_VAR", "default")
		assert.Equal(t, "existing-value", value)
		
		// Test with non-existing value
		value = manager.GetEnvWithDefault("NON_EXISTING_VAR", "default")
		assert.Equal(t, "default", value)
	})
	
	t.Run("SetupTestEnvironment", func(t *testing.T) {
		manager := NewEnvironmentMockManager()
		defer manager.RestoreEnvironment()
		
		// Setup
		manager.SetupTestEnvironment()
		
		// Assert
		assert.Equal(t, "test", manager.GetEnv("APP_ENV"))
		assert.Equal(t, "sqlite3", manager.GetEnv("DB_DRIVER"))
		assert.Equal(t, ":memory:", manager.GetEnv("DB_URL"))
		assert.Equal(t, "test-secret-key", manager.GetEnv("JWT_SECRET"))
	})
}

// TestConfigBuilder tests the configuration builder functionality
func TestConfigBuilder(t *testing.T) {
	t.Run("fluent interface", func(t *testing.T) {
		builder := NewConfigBuilder()
		
		// Build configuration
		manager := builder.
			WithDatabaseConfig("postgres", "localhost:5432", 10, 5).
			WithLoggerConfig("info").
			WithCustomConfig("custom", "key", "custom-value").
			Build()
		
		// Test database config
		driver, err := manager.GetMockClient().GetString("database", "driverName")
		assert.NoError(t, err)
		assert.Equal(t, "postgres", driver)
		
		// Test logger config
		level, err := manager.GetMockClient().GetString("logger", "level")
		assert.NoError(t, err)
		assert.Equal(t, "info", level)
		
		// Test custom config
		custom, err := manager.GetMockClient().GetString("custom", "key")
		assert.NoError(t, err)
		assert.Equal(t, "custom-value", custom)
	})
}

// TestConfigTestHelper tests the configuration test helper functionality
func TestConfigTestHelper(t *testing.T) {
	t.Run("SetupDefaultTestConfig", func(t *testing.T) {
		helper := NewConfigTestHelper()
		defer helper.Cleanup()
		
		// Setup
		helper.SetupDefaultTestConfig()
		
		// Test config values
		driver, err := helper.GetConfigManager().GetMockClient().GetString("database", "driverName")
		assert.NoError(t, err)
		assert.Equal(t, "sqlite3", driver)
		
		level, err := helper.GetConfigManager().GetMockClient().GetString("logger", "level")
		assert.NoError(t, err)
		assert.Equal(t, "debug", level)
		
		// Test environment values
		assert.Equal(t, "test", helper.GetEnvManager().GetEnv("APP_ENV"))
		assert.Equal(t, "sqlite3", helper.GetEnvManager().GetEnv("DB_DRIVER"))
	})
}

// TestMockConfigProvider tests the mock config provider functionality
func TestMockConfigProvider(t *testing.T) {
	t.Run("SetupDatabaseConfigForTesting", func(t *testing.T) {
		provider := NewMockConfigProvider()
		
		// Setup
		provider.SetupDatabaseConfigForTesting()
		
		// Test database configuration
		driver, err := provider.GetMockClient().GetString("database", "driverName")
		assert.NoError(t, err)
		assert.Equal(t, "sqlite3", driver)
		
		url, err := provider.GetMockClient().GetString("database", "url")
		assert.NoError(t, err)
		assert.Equal(t, ":memory:", url)
		
		maxOpen := provider.GetMockClient().GetIntD("database", "maxOpenConnections", 1)
		assert.Equal(t, int64(10), maxOpen)
	})
	
	t.Run("SetupAPIConfigForTesting", func(t *testing.T) {
		provider := NewMockConfigProvider()
		
		// Setup
		provider.SetupAPIConfigForTesting()
		
		// Test API configuration
		apiConfig := provider.GetMockClient().GetMapD("api", "getCatsFact", nil)
		assert.NotNil(t, apiConfig)
		assert.Equal(t, "http://test-api.example.com", apiConfig["url"])
		assert.Equal(t, "30s", apiConfig["timeout"])
		assert.Equal(t, 3, apiConfig["retries"])
	})
	
	t.Run("SetupAllConfigsForTesting", func(t *testing.T) {
		provider := NewMockConfigProvider()
		
		// Setup
		provider.SetupAllConfigsForTesting()
		
		// Test that all configurations are set up
		driver, err := provider.GetMockClient().GetString("database", "driverName")
		assert.NoError(t, err)
		assert.Equal(t, "sqlite3", driver)
		
		apiConfig := provider.GetMockClient().GetMapD("api", "getCatsFact", nil)
		assert.NotNil(t, apiConfig)
		
		appName := provider.GetMockClient().GetStringD("application", "name", "")
		assert.Equal(t, "post-eaze-test", appName)
	})
}