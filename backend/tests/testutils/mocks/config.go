package mocks

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stretchr/testify/mock"
)

// ConfigClient interface for mocking configuration client
// This interface matches the actual go-config-client interface used in the application
type ConfigClient interface {
	GetString(configName, key string) (string, error)
	GetIntD(configName, key string, defaultValue int64) int64
	GetMapD(configName, key string, defaultValue map[string]interface{}) map[string]interface{}
	GetBoolD(configName, key string, defaultValue bool) bool
	GetFloat64D(configName, key string, defaultValue float64) float64
	GetStringD(configName, key string, defaultValue string) string
	GetStringSliceD(configName, key string, defaultValue []string) []string
}

// MockConfigClient is a mock implementation of ConfigClient
type MockConfigClient struct {
	mock.Mock
}

// GetString mocks the GetString method
func (m *MockConfigClient) GetString(configName, key string) (string, error) {
	args := m.Called(configName, key)
	return args.String(0), args.Error(1)
}

// GetIntD mocks the GetIntD method (returns int64 with default)
func (m *MockConfigClient) GetIntD(configName, key string, defaultValue int64) int64 {
	args := m.Called(configName, key, defaultValue)
	if len(args) > 0 && args.Get(0) != nil {
		return args.Get(0).(int64)
	}
	return defaultValue
}

// GetMapD mocks the GetMapD method (returns map with default)
func (m *MockConfigClient) GetMapD(configName, key string, defaultValue map[string]interface{}) map[string]interface{} {
	args := m.Called(configName, key, defaultValue)
	if len(args) > 0 && args.Get(0) != nil {
		return args.Get(0).(map[string]interface{})
	}
	return defaultValue
}

// GetBoolD mocks the GetBoolD method (returns bool with default)
func (m *MockConfigClient) GetBoolD(configName, key string, defaultValue bool) bool {
	args := m.Called(configName, key, defaultValue)
	if len(args) > 0 && args.Get(0) != nil {
		return args.Bool(0)
	}
	return defaultValue
}

// GetFloat64D mocks the GetFloat64D method (returns float64 with default)
func (m *MockConfigClient) GetFloat64D(configName, key string, defaultValue float64) float64 {
	args := m.Called(configName, key, defaultValue)
	if len(args) > 0 && args.Get(0) != nil {
		return args.Get(0).(float64)
	}
	return defaultValue
}

// GetStringD mocks the GetStringD method (returns string with default)
func (m *MockConfigClient) GetStringD(configName, key string, defaultValue string) string {
	args := m.Called(configName, key, defaultValue)
	if len(args) > 0 && args.Get(0) != nil {
		return args.String(0)
	}
	return defaultValue
}

// GetStringSliceD mocks the GetStringSliceD method (returns []string with default)
func (m *MockConfigClient) GetStringSliceD(configName, key string, defaultValue []string) []string {
	args := m.Called(configName, key, defaultValue)
	if len(args) > 0 && args.Get(0) != nil {
		return args.Get(0).([]string)
	}
	return defaultValue
}

// NewMockConfigClient creates a new mock config client
func NewMockConfigClient() *MockConfigClient {
	return &MockConfigClient{}
}

// ConfigMockManager provides utilities for managing configuration mocks
type ConfigMockManager struct {
	mockClient *MockConfigClient
	values     map[string]map[string]interface{} // configName -> key -> value
}

// NewConfigMockManager creates a new configuration mock manager
func NewConfigMockManager() *ConfigMockManager {
	return &ConfigMockManager{
		mockClient: NewMockConfigClient(),
		values:     make(map[string]map[string]interface{}),
	}
}

// GetMockClient returns the mock config client
func (m *ConfigMockManager) GetMockClient() *MockConfigClient {
	return m.mockClient
}

// SetValue sets a configuration value for mocking
func (m *ConfigMockManager) SetValue(configName, key string, value interface{}) {
	if m.values[configName] == nil {
		m.values[configName] = make(map[string]interface{})
	}
	m.values[configName][key] = value
	
	// Set up type-specific methods based on value type
	switch v := value.(type) {
	case string:
		m.mockClient.On("GetString", configName, key).Return(v, nil)
		m.mockClient.On("GetStringD", configName, key, mock.Anything).Return(v)
	case int:
		m.mockClient.On("GetIntD", configName, key, mock.Anything).Return(int64(v))
	case int64:
		m.mockClient.On("GetIntD", configName, key, mock.Anything).Return(v)
	case bool:
		m.mockClient.On("GetBoolD", configName, key, mock.Anything).Return(v)
	case float64:
		m.mockClient.On("GetFloat64D", configName, key, mock.Anything).Return(v)
	case []string:
		m.mockClient.On("GetStringSliceD", configName, key, mock.Anything).Return(v)
	case map[string]interface{}:
		m.mockClient.On("GetMapD", configName, key, mock.Anything).Return(v)
	}
}

// SetError sets an error for a specific configuration key
func (m *ConfigMockManager) SetError(configName, key string, err error) {
	m.mockClient.On("GetString", configName, key).Return("", err)
}

// SetDatabaseConfig sets up database configuration values
func (m *ConfigMockManager) SetDatabaseConfig(driverName, url string, maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) {
	m.SetValue("database", "driverName", driverName)
	m.SetValue("database", "url", url)
	m.SetValue("database", "maxOpenConnections", int64(maxOpen))
	m.SetValue("database", "maxIdleConnections", int64(maxIdle))
	m.SetValue("database", "maxConnectionLifetimeInSeconds", int64(maxLifetime.Seconds()))
	m.SetValue("database", "maxConnectionIdleTimeInSeconds", int64(maxIdleTime.Seconds()))
}

// SetLoggerConfig sets up logger configuration values
func (m *ConfigMockManager) SetLoggerConfig(level string, params map[string]interface{}) {
	m.SetValue("logger", "level", level)
	m.SetValue("logger", "params", params)
}

// SetAPIConfig sets up API configuration values
func (m *ConfigMockManager) SetAPIConfig(configs map[string]interface{}) {
	for key, value := range configs {
		m.SetValue("api", key, value)
	}
}

// SetApplicationConfig sets up application configuration values
func (m *ConfigMockManager) SetApplicationConfig(configs map[string]interface{}) {
	for key, value := range configs {
		m.SetValue("application", key, value)
	}
}

// VerifyConfigAccessed verifies that a configuration key was accessed
func (m *ConfigMockManager) VerifyConfigAccessed(configName, key string) bool {
	for _, call := range m.mockClient.Calls {
		if len(call.Arguments) >= 2 && call.Arguments[0] == configName && call.Arguments[1] == key {
			return true
		}
	}
	return false
}

// GetAccessCount returns the number of times a configuration key was accessed
func (m *ConfigMockManager) GetAccessCount(configName, key string) int {
	count := 0
	for _, call := range m.mockClient.Calls {
		if len(call.Arguments) >= 2 && call.Arguments[0] == configName && call.Arguments[1] == key {
			count++
		}
	}
	return count
}

// AssertExpectations asserts that all expectations were met
func (m *ConfigMockManager) AssertExpectations(t mock.TestingT) {
	m.mockClient.AssertExpectations(t)
}

// Reset clears all expectations and call history
func (m *ConfigMockManager) Reset() {
	m.mockClient.ExpectedCalls = nil
	m.mockClient.Calls = nil
	m.values = make(map[string]map[string]interface{})
}

// EnvironmentMockManager provides utilities for mocking environment variables
type EnvironmentMockManager struct {
	originalEnv map[string]string
	mockEnv     map[string]string
}

// NewEnvironmentMockManager creates a new environment mock manager
func NewEnvironmentMockManager() *EnvironmentMockManager {
	return &EnvironmentMockManager{
		originalEnv: make(map[string]string),
		mockEnv:     make(map[string]string),
	}
}

// SetEnv sets an environment variable for testing
func (m *EnvironmentMockManager) SetEnv(key, value string) {
	// Store original value if it exists
	if originalValue, exists := os.LookupEnv(key); exists {
		m.originalEnv[key] = originalValue
	}
	
	// Set the mock value
	m.mockEnv[key] = value
	os.Setenv(key, value)
}

// UnsetEnv unsets an environment variable for testing
func (m *EnvironmentMockManager) UnsetEnv(key string) {
	// Store original value if it exists
	if originalValue, exists := os.LookupEnv(key); exists {
		m.originalEnv[key] = originalValue
	}
	
	os.Unsetenv(key)
}

// GetEnv gets an environment variable value
func (m *EnvironmentMockManager) GetEnv(key string) string {
	return os.Getenv(key)
}

// GetEnvWithDefault gets an environment variable value with a default
func (m *EnvironmentMockManager) GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestEnvironment sets up a complete test environment
func (m *EnvironmentMockManager) SetupTestEnvironment() {
	m.SetEnv("APP_ENV", "test")
	m.SetEnv("DB_DRIVER", "sqlite3")
	m.SetEnv("DB_URL", ":memory:")
	m.SetEnv("JWT_SECRET", "test-secret-key")
	m.SetEnv("LOG_LEVEL", "debug")
	m.SetEnv("LOG_DIR", "test-logs")
}

// RestoreEnvironment restores the original environment variables
func (m *EnvironmentMockManager) RestoreEnvironment() {
	// Restore original values
	for key, value := range m.originalEnv {
		os.Setenv(key, value)
	}
	
	// Unset mock-only values
	for key := range m.mockEnv {
		if _, exists := m.originalEnv[key]; !exists {
			os.Unsetenv(key)
		}
	}
	
	// Clear tracking maps
	m.originalEnv = make(map[string]string)
	m.mockEnv = make(map[string]string)
}

// ConfigBuilder provides a fluent interface for building configuration mocks
type ConfigBuilder struct {
	manager *ConfigMockManager
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		manager: NewConfigMockManager(),
	}
}

// WithDatabaseConfig adds database configuration
func (b *ConfigBuilder) WithDatabaseConfig(driverName, url string, maxOpen, maxIdle int) *ConfigBuilder {
	b.manager.SetDatabaseConfig(driverName, url, maxOpen, maxIdle, time.Hour, time.Minute*30)
	return b
}

// WithLoggerConfig adds logger configuration
func (b *ConfigBuilder) WithLoggerConfig(level string) *ConfigBuilder {
	b.manager.SetLoggerConfig(level, map[string]interface{}{
		"console": true,
		"file":    true,
	})
	return b
}

// WithAPIConfig adds API configuration
func (b *ConfigBuilder) WithAPIConfig(configs map[string]interface{}) *ConfigBuilder {
	b.manager.SetAPIConfig(configs)
	return b
}

// WithCustomConfig adds custom configuration
func (b *ConfigBuilder) WithCustomConfig(configName, key string, value interface{}) *ConfigBuilder {
	b.manager.SetValue(configName, key, value)
	return b
}

// WithError adds error for a specific key
func (b *ConfigBuilder) WithError(configName, key string, err error) *ConfigBuilder {
	b.manager.SetError(configName, key, err)
	return b
}

// Build returns the configured mock manager
func (b *ConfigBuilder) Build() *ConfigMockManager {
	return b.manager
}

// Common configuration errors for testing
var (
	ErrMockConfigNotFound    = errors.New("mock: configuration key not found")
	ErrMockConfigInvalid     = errors.New("mock: invalid configuration value")
	ErrMockConfigTypeMismatch = errors.New("mock: configuration type mismatch")
	ErrMockConfigConnection  = errors.New("mock: configuration service connection failed")
)

// ConfigTestHelper provides high-level utilities for configuration testing
type ConfigTestHelper struct {
	configManager *ConfigMockManager
	envManager    *EnvironmentMockManager
}

// NewConfigTestHelper creates a new configuration test helper
func NewConfigTestHelper() *ConfigTestHelper {
	return &ConfigTestHelper{
		configManager: NewConfigMockManager(),
		envManager:    NewEnvironmentMockManager(),
	}
}

// GetConfigManager returns the configuration manager
func (h *ConfigTestHelper) GetConfigManager() *ConfigMockManager {
	return h.configManager
}

// GetEnvManager returns the environment manager
func (h *ConfigTestHelper) GetEnvManager() *EnvironmentMockManager {
	return h.envManager
}

// SetupDefaultTestConfig sets up default configuration for testing
func (h *ConfigTestHelper) SetupDefaultTestConfig() {
	// Database config
	h.configManager.SetDatabaseConfig("sqlite3", ":memory:", 10, 5, time.Hour, time.Minute*30)
	
	// Logger config
	h.configManager.SetLoggerConfig("debug", map[string]interface{}{
		"console": true,
		"file":    false,
	})
	
	// API config
	h.configManager.SetAPIConfig(map[string]interface{}{
		"port":    8080,
		"timeout": "30s",
	})
	
	// Environment variables
	h.envManager.SetupTestEnvironment()
}

// Cleanup cleans up all mocks and restores original state
func (h *ConfigTestHelper) Cleanup() {
	h.configManager.Reset()
	h.envManager.RestoreEnvironment()
}

// VerifyConfigurationAccess verifies that specific configuration keys were accessed
func (h *ConfigTestHelper) VerifyConfigurationAccess(configKeys map[string][]string) map[string]map[string]bool {
	results := make(map[string]map[string]bool)
	for configName, keys := range configKeys {
		results[configName] = make(map[string]bool)
		for _, key := range keys {
			results[configName][key] = h.configManager.VerifyConfigAccessed(configName, key)
		}
	}
	return results
}

// TypeConverter provides utilities for converting configuration values
type TypeConverter struct{}

// NewTypeConverter creates a new type converter
func NewTypeConverter() *TypeConverter {
	return &TypeConverter{}
}

// ToString converts a value to string
func (c *TypeConverter) ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// ToInt converts a value to int
func (c *TypeConverter) ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// ToBool converts a value to bool
func (c *TypeConverter) ToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// ToDuration converts a value to time.Duration
func (c *TypeConverter) ToDuration(value interface{}) (time.Duration, error) {
	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(v)
	case int:
		return time.Duration(v) * time.Second, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to duration", value)
	}
}

// MockConfigProvider provides a way to replace the global config client for testing
type MockConfigProvider struct {
	mockClient *MockConfigClient
	original   interface{} // Store original config client if needed
}

// NewMockConfigProvider creates a new mock config provider
func NewMockConfigProvider() *MockConfigProvider {
	return &MockConfigProvider{
		mockClient: NewMockConfigClient(),
	}
}

// GetMockClient returns the mock config client
func (p *MockConfigProvider) GetMockClient() *MockConfigClient {
	return p.mockClient
}

// SetupDatabaseConfigForTesting sets up database configuration values for testing
func (p *MockConfigProvider) SetupDatabaseConfigForTesting() {
	// Setup typical database configuration values for testing
	p.mockClient.On("GetString", "database", "driverName").Return("sqlite3", nil)
	p.mockClient.On("GetString", "database", "url").Return(":memory:", nil)
	p.mockClient.On("GetIntD", "database", "maxOpenConnections", int64(1)).Return(int64(10))
	p.mockClient.On("GetIntD", "database", "maxIdleConnections", int64(0)).Return(int64(5))
	p.mockClient.On("GetIntD", "database", "maxConnectionLifetimeInSeconds", int64(30)).Return(int64(3600))
	p.mockClient.On("GetIntD", "database", "maxConnectionIdleTimeInSeconds", int64(10)).Return(int64(1800))
}

// SetupAPIConfigForTesting sets up API configuration values for testing
func (p *MockConfigProvider) SetupAPIConfigForTesting() {
	// Setup typical API configuration values for testing
	testAPIConfig := map[string]interface{}{
		"url":     "http://test-api.example.com",
		"timeout": "30s",
		"retries": 3,
	}
	p.mockClient.On("GetMapD", "api", "getCatsFact", mock.Anything).Return(testAPIConfig)
}

// SetupApplicationConfigForTesting sets up application configuration values for testing
func (p *MockConfigProvider) SetupApplicationConfigForTesting() {
	// Setup typical application configuration values for testing
	p.mockClient.On("GetStringD", "application", "name", mock.Anything).Return("post-eaze-test")
	p.mockClient.On("GetStringD", "application", "version", mock.Anything).Return("1.0.0-test")
	p.mockClient.On("GetBoolD", "application", "debug", mock.Anything).Return(true)
}

// SetupAllConfigsForTesting sets up all configuration values for testing
func (p *MockConfigProvider) SetupAllConfigsForTesting() {
	p.SetupDatabaseConfigForTesting()
	p.SetupAPIConfigForTesting()
	p.SetupApplicationConfigForTesting()
}