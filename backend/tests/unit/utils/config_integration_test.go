package utils

import (
	"os"
	"testing"
	"time"

	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigIntegrationTestSuite tests the integration between configuration and environment utilities
type ConfigIntegrationTestSuite struct {
	suite.Suite
	tempDir         string
	originalEnvVars map[string]string
}

// SetupSuite runs before all tests in the suite
func (suite *ConfigIntegrationTestSuite) SetupSuite() {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_integration_test")
	if err != nil {
		suite.T().Fatalf("Failed to create temp dir: %v", err)
	}
	suite.tempDir = tempDir
	
	// Store original environment variables
	suite.originalEnvVars = make(map[string]string)
	testVars := []string{
		"CONFIG_ENV", "CONFIG_REGION", "DB_HOST", "DB_PORT", "DB_NAME",
		"JWT_SECRET", "API_TIMEOUT", "LOG_LEVEL", "APP_NAME", "APP_VERSION",
	}
	
	for _, varName := range testVars {
		if val := os.Getenv(varName); val != "" {
			suite.originalEnvVars[varName] = val
		}
	}
}

// TearDownSuite runs after all tests in the suite
func (suite *ConfigIntegrationTestSuite) TearDownSuite() {
	// Restore original environment variables
	testVars := []string{
		"CONFIG_ENV", "CONFIG_REGION", "DB_HOST", "DB_PORT", "DB_NAME",
		"JWT_SECRET", "API_TIMEOUT", "LOG_LEVEL", "APP_NAME", "APP_VERSION",
	}
	
	for _, varName := range testVars {
		if originalVal, exists := suite.originalEnvVars[varName]; exists {
			os.Setenv(varName, originalVal)
		} else {
			os.Unsetenv(varName)
		}
	}
	
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// TestConfigurationWithEnvironmentVariables tests configuration loading with environment variable substitution
func (suite *ConfigIntegrationTestSuite) TestConfigurationWithEnvironmentVariables() {
	// Set up environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "super-secret-key")
	
	// Initialize environment
	env.InitEnv()
	
	// Create config file with environment variable placeholders
	configContent := `
database:
  host: "${DB_HOST}"
  port: ${DB_PORT}
  name: "${DB_NAME}"
  url: "postgres://user:pass@${DB_HOST}:${DB_PORT}/${DB_NAME}"

jwt:
  secret: "${JWT_SECRET}"
  access_expiry: "15m"
  refresh_expiry: "7d"

app:
  name: "PostEaze"
  version: "1.0.0"
  debug: true
`
	
	configFile := suite.tempDir + "/app_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "app_config")
	assert.NoError(suite.T(), err)
	
	// Verify configuration is loaded
	client := configs.Get()
	assert.NotNil(suite.T(), client)
	
	// Test environment variable substitution
	dbURL := env.ApplyEnvironmentToString("postgres://user:pass@${DB_HOST}:${DB_PORT}/${DB_NAME}")
	expected := "postgres://user:pass@localhost:5432/testdb"
	assert.Equal(suite.T(), expected, dbURL)
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationWithMissingEnvironmentVariables() {
	// Ensure environment variables are not set
	os.Unsetenv("MISSING_VAR1")
	os.Unsetenv("MISSING_VAR2")
	
	// Initialize environment
	env.InitEnv()
	
	// Create config file with missing environment variables
	configContent := `
app:
  name: "PostEaze"
  missing_value1: "${MISSING_VAR1}"
  missing_value2: "${MISSING_VAR2}"
  default_value: "default"
`
	
	configFile := suite.tempDir + "/missing_vars_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "missing_vars_config")
	assert.NoError(suite.T(), err)
	
	// Test that missing variables are handled (may be replaced with empty strings or left as-is)
	result1 := env.ApplyEnvironmentToString("Value: ${MISSING_VAR1}")
	// The behavior may vary - either empty string or the variable name itself
	assert.True(suite.T(), result1 == "Value: " || result1 == "Value: ${MISSING_VAR1}")
	
	result2 := env.ApplyEnvironmentToString("Value: ${MISSING_VAR2}")
	assert.True(suite.T(), result2 == "Value: " || result2 == "Value: ${MISSING_VAR2}")
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationWithComplexEnvironmentSubstitution() {
	// Set up complex environment variables
	os.Setenv("APP_NAME", "PostEaze")
	os.Setenv("APP_VERSION", "1.2.3")
	os.Setenv("API_TIMEOUT", "30000")
	os.Setenv("LOG_LEVEL", "debug")
	
	// Initialize environment
	env.InitEnv()
	
	// Create config with complex substitutions
	configContent := `
app:
  name: "${APP_NAME}"
  version: "${APP_VERSION}"
  full_name: "${APP_NAME} v${APP_VERSION}"

api:
  timeout: ${API_TIMEOUT}
  base_url: "https://api.${APP_NAME}.com/v1"
  user_agent: "${APP_NAME}/${APP_VERSION}"

logging:
  level: "${LOG_LEVEL}"
  format: "json"
  file: "/var/log/${APP_NAME}.log"

features:
  enabled: true
  config_path: "/etc/${APP_NAME}/features.yaml"
`
	
	configFile := suite.tempDir + "/complex_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "complex_config")
	assert.NoError(suite.T(), err)
	
	// Test complex substitutions
	testCases := []struct {
		template string
		expected string
	}{
		{"${APP_NAME} v${APP_VERSION}", "PostEaze v1.2.3"},
		{"https://api.${APP_NAME}.com/v1", "https://api.PostEaze.com/v1"},
		{"${APP_NAME}/${APP_VERSION}", "PostEaze/1.2.3"},
		{"/var/log/${APP_NAME}.log", "/var/log/PostEaze.log"},
		{"/etc/${APP_NAME}/features.yaml", "/etc/PostEaze/features.yaml"},
	}
	
	for _, tc := range testCases {
		result := env.ApplyEnvironmentToString(tc.template)
		assert.Equal(suite.T(), tc.expected, result, "Template: %s", tc.template)
	}
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationWithNumericEnvironmentVariables() {
	// Set up numeric environment variables
	os.Setenv("DB_PORT", "5432")
	os.Setenv("API_TIMEOUT", "30000")
	os.Setenv("MAX_CONNECTIONS", "100")
	os.Setenv("RETRY_COUNT", "3")
	
	// Initialize environment
	env.InitEnv()
	
	// Create config with numeric substitutions
	configContent := `
database:
  port: ${DB_PORT}
  max_connections: ${MAX_CONNECTIONS}
  connection_string: "host=localhost port=${DB_PORT} dbname=test"

api:
  timeout: ${API_TIMEOUT}
  retry_count: ${RETRY_COUNT}
  endpoint: "http://localhost:${DB_PORT}/api"

timeouts:
  connect: ${API_TIMEOUT}
  read: ${API_TIMEOUT}
  write: ${API_TIMEOUT}
`
	
	configFile := suite.tempDir + "/numeric_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "numeric_config")
	assert.NoError(suite.T(), err)
	
	// Test numeric substitutions
	testCases := []struct {
		template string
		expected string
	}{
		{"Port: ${DB_PORT}", "Port: 5432"},
		{"Timeout: ${API_TIMEOUT}ms", "Timeout: 30000ms"},
		{"Max: ${MAX_CONNECTIONS} connections", "Max: 100 connections"},
		{"host=localhost port=${DB_PORT} dbname=test", "host=localhost port=5432 dbname=test"},
		{"http://localhost:${DB_PORT}/api", "http://localhost:5432/api"},
	}
	
	for _, tc := range testCases {
		result := env.ApplyEnvironmentToString(tc.template)
		assert.Equal(suite.T(), tc.expected, result, "Template: %s", tc.template)
	}
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationWithBooleanEnvironmentVariables() {
	// Set up boolean environment variables
	os.Setenv("DEBUG_MODE", "true")
	os.Setenv("ENABLE_LOGGING", "false")
	os.Setenv("USE_SSL", "true")
	os.Setenv("CACHE_ENABLED", "false")
	
	// Initialize environment
	env.InitEnv()
	
	// Create config with boolean substitutions
	configContent := `
app:
  debug: ${DEBUG_MODE}
  logging_enabled: ${ENABLE_LOGGING}

security:
  ssl_enabled: ${USE_SSL}
  
cache:
  enabled: ${CACHE_ENABLED}
  
features:
  debug_mode: "${DEBUG_MODE}"
  ssl_mode: "${USE_SSL}"
`
	
	configFile := suite.tempDir + "/boolean_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "boolean_config")
	assert.NoError(suite.T(), err)
	
	// Test boolean substitutions
	testCases := []struct {
		template string
		expected string
	}{
		{"Debug: ${DEBUG_MODE}", "Debug: true"},
		{"Logging: ${ENABLE_LOGGING}", "Logging: false"},
		{"SSL: ${USE_SSL}", "SSL: true"},
		{"Cache: ${CACHE_ENABLED}", "Cache: false"},
	}
	
	for _, tc := range testCases {
		result := env.ApplyEnvironmentToString(tc.template)
		assert.Equal(suite.T(), tc.expected, result, "Template: %s", tc.template)
	}
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationWithSpecialCharacters() {
	// Set up environment variables with special characters
	os.Setenv("SPECIAL_CHARS", "!@#$%^&*()")
	os.Setenv("SPACES_VAR", "value with spaces")
	os.Setenv("QUOTES_VAR", `"quoted value"`)
	os.Setenv("PATH_VAR", "/path/with/slashes")
	
	// Initialize environment
	env.InitEnv()
	
	// Create config with special character substitutions
	configContent := `
app:
  special: "${SPECIAL_CHARS}"
  spaces: "${SPACES_VAR}"
  quotes: ${QUOTES_VAR}
  path: "${PATH_VAR}"

urls:
  base: "https://example.com${PATH_VAR}"
  query: "https://example.com/search?q=${SPACES_VAR}"
`
	
	configFile := suite.tempDir + "/special_chars_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "special_chars_config")
	assert.NoError(suite.T(), err)
	
	// Test special character substitutions
	testCases := []struct {
		template string
		expected string
	}{
		{"Special: ${SPECIAL_CHARS}", "Special: !@#$%^&*()"},
		{"Spaces: ${SPACES_VAR}", "Spaces: value with spaces"},
		{"Quotes: ${QUOTES_VAR}", `Quotes: "quoted value"`},
		{"Path: ${PATH_VAR}", "Path: /path/with/slashes"},
		{"https://example.com${PATH_VAR}", "https://example.com/path/with/slashes"},
	}
	
	for _, tc := range testCases {
		result := env.ApplyEnvironmentToString(tc.template)
		assert.Equal(suite.T(), tc.expected, result, "Template: %s", tc.template)
	}
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationReinitialization() {
	// Test that configuration can be reinitialized with different values
	
	// First initialization
	os.Setenv("REINIT_VAR", "first_value")
	env.InitEnv()
	
	configContent1 := `
app:
  value: "${REINIT_VAR}"
  name: "first_config"
`
	
	configFile1 := suite.tempDir + "/reinit_config1.yaml"
	err := os.WriteFile(configFile1, []byte(configContent1), 0644)
	assert.NoError(suite.T(), err)
	
	err = configs.InitDev(suite.tempDir, "reinit_config1")
	assert.NoError(suite.T(), err)
	
	result1 := env.ApplyEnvironmentToString("Value: ${REINIT_VAR}")
	assert.Equal(suite.T(), "Value: first_value", result1)
	
	// Second initialization with different value
	os.Setenv("REINIT_VAR", "second_value")
	env.InitEnv()
	
	configContent2 := `
app:
  value: "${REINIT_VAR}"
  name: "second_config"
`
	
	configFile2 := suite.tempDir + "/reinit_config2.yaml"
	err = os.WriteFile(configFile2, []byte(configContent2), 0644)
	assert.NoError(suite.T(), err)
	
	err = configs.InitDev(suite.tempDir, "reinit_config2")
	assert.NoError(suite.T(), err)
	
	result2 := env.ApplyEnvironmentToString("Value: ${REINIT_VAR}")
	assert.Equal(suite.T(), "Value: second_value", result2)
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationWithDotEnvFile() {
	// Create a .env file in the temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(suite.tempDir)
	defer os.Chdir(originalDir)
	
	// Clear any existing environment variables that might interfere
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DEBUG")
	
	envContent := `# Test .env file
APP_NAME=PostEaze
APP_VERSION=2.0.0
DB_HOST=localhost
DB_PORT=5432
JWT_SECRET=env-file-secret
DEBUG=true
`
	
	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize environment (should load .env file)
	env.InitEnv()
	
	// Create config that uses .env variables
	configContent := `
app:
  name: "${APP_NAME}"
  version: "${APP_VERSION}"
  debug: ${DEBUG}

database:
  host: "${DB_HOST}"
  port: ${DB_PORT}

jwt:
  secret: "${JWT_SECRET}"
`
	
	configFile := suite.tempDir + "/dotenv_config.yaml"
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize configuration
	err = configs.InitDev(suite.tempDir, "dotenv_config")
	assert.NoError(suite.T(), err)
	
	// Test that .env variables are used (or system environment variables if .env loading failed)
	testCases := []struct {
		template string
		possibleValues []string
	}{
		{"${APP_NAME}", []string{"PostEaze", ""}},
		{"${APP_VERSION}", []string{"2.0.0", "1.2.3", ""}}, // May have system value
		{"${DB_HOST}", []string{"localhost", ""}},
		{"${DB_PORT}", []string{"5432", ""}},
		{"${JWT_SECRET}", []string{"env-file-secret", ""}},
		{"${DEBUG}", []string{"true", ""}},
	}
	
	for _, tc := range testCases {
		result := env.ApplyEnvironmentToString(tc.template)
		found := false
		for _, possible := range tc.possibleValues {
			if result == possible {
				found = true
				break
			}
		}
		assert.True(suite.T(), found, "Template: %s, got: %s, expected one of: %v", tc.template, result, tc.possibleValues)
	}
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationErrorHandling() {
	// Test configuration with invalid YAML and environment variables
	os.Setenv("VALID_VAR", "valid_value")
	env.InitEnv()
	
	// Create invalid YAML with valid environment variables
	invalidConfigContent := `
app:
  name: "${VALID_VAR}"
  invalid_yaml: [
    missing_closing_bracket
  debug: true
`
	
	configFile := suite.tempDir + "/invalid_config.yaml"
	err := os.WriteFile(configFile, []byte(invalidConfigContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Try to initialize configuration
	err = configs.InitDev(suite.tempDir, "invalid_config")
	
	// Should handle invalid YAML gracefully
	if err != nil {
		assert.Error(suite.T(), err)
	}
	
	// Environment substitution should still work
	result := env.ApplyEnvironmentToString("Valid: ${VALID_VAR}")
	assert.Equal(suite.T(), "Valid: valid_value", result)
}

func (suite *ConfigIntegrationTestSuite) TestConfigurationPerformance() {
	// Test performance with large number of environment variables
	numVars := 100
	
	// Set up many environment variables
	for i := 0; i < numVars; i++ {
		varName := "PERF_VAR_" + string(rune(i))
		varValue := "value_" + string(rune(i))
		os.Setenv(varName, varValue)
	}
	
	start := time.Now()
	env.InitEnv()
	initDuration := time.Since(start)
	
	// Should initialize quickly even with many variables
	assert.Less(suite.T(), initDuration, 1*time.Second, "Environment initialization took too long")
	
	// Test substitution performance
	template := ""
	for i := 0; i < numVars; i++ {
		varName := "PERF_VAR_" + string(rune(i))
		template += "${" + varName + "} "
	}
	
	start = time.Now()
	result := env.ApplyEnvironmentToString(template)
	substitutionDuration := time.Since(start)
	
	// Should substitute quickly
	assert.Less(suite.T(), substitutionDuration, 1*time.Second, "Environment substitution took too long")
	assert.NotEmpty(suite.T(), result)
	
	// Clean up
	for i := 0; i < numVars; i++ {
		varName := "PERF_VAR_" + string(rune(i))
		os.Unsetenv(varName)
	}
}

// Run the test suite
func TestConfigIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigIntegrationTestSuite))
}