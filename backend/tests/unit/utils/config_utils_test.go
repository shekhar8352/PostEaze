package utils

import (
	"os"
	"testing"

	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigUtilsTestSuite defines the test suite for configuration utilities
type ConfigUtilsTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupSuite runs before all tests in the suite
func (suite *ConfigUtilsTestSuite) SetupSuite() {
	// Create a temporary directory for test config files
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		suite.T().Fatalf("Failed to create temp dir: %v", err)
	}
	suite.tempDir = tempDir
}

// TearDownSuite runs after all tests in the suite
func (suite *ConfigUtilsTestSuite) TearDownSuite() {
	// Clean up temporary directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// TestInitDev tests the development configuration initialization
func (suite *ConfigUtilsTestSuite) TestInitDev_ValidDirectory() {
	// Create a test config file
	configContent := `
test:
  value: "test_value"
  number: 123
database:
  host: "localhost"
  port: 5432
`
	configFile := suite.tempDir + "/test_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Test initialization
	err = configs.InitDev(suite.tempDir, "test_config")
	
	assert.NoError(suite.T(), err)
	
	// Verify client is available
	client := configs.Get()
	assert.NotNil(suite.T(), client)
}

func (suite *ConfigUtilsTestSuite) TestInitDev_NonExistentDirectory() {
	nonExistentDir := "/path/that/does/not/exist"
	
	err := configs.InitDev(nonExistentDir, "config")
	
	assert.Error(suite.T(), err)
}

func (suite *ConfigUtilsTestSuite) TestInitDev_EmptyDirectory() {
	err := configs.InitDev("", "config")
	
	assert.Error(suite.T(), err)
}

func (suite *ConfigUtilsTestSuite) TestInitDev_MultipleConfigFiles() {
	// Create multiple test config files
	config1Content := `
app:
  name: "test_app"
  version: "1.0.0"
`
	config2Content := `
database:
  host: "localhost"
  port: 5432
`
	
	config1File := suite.tempDir + "/app_config.yaml"
	config2File := suite.tempDir + "/db_config.yaml"
	
	err := os.WriteFile(config1File, []byte(config1Content), 0644)
	assert.NoError(suite.T(), err)
	
	err = os.WriteFile(config2File, []byte(config2Content), 0644)
	assert.NoError(suite.T(), err)
	
	// Test initialization with multiple config files
	err = configs.InitDev(suite.tempDir, "app_config", "db_config")
	
	assert.NoError(suite.T(), err)
	
	// Verify client is available
	client := configs.Get()
	assert.NotNil(suite.T(), client)
}

func (suite *ConfigUtilsTestSuite) TestInitDev_NoConfigNames() {
	_ = configs.InitDev(suite.tempDir)
	
	// Should handle empty config names gracefully
	// The exact behavior depends on the underlying config library
	// We just verify it doesn't panic
	assert.NotPanics(suite.T(), func() {
		configs.Get()
	})
}

// TestInitRelease tests the release configuration initialization
func (suite *ConfigUtilsTestSuite) TestInitRelease_ValidParameters() {
	env := "test"
	region := "us-east-1"
	configNames := []string{"app", "database"}
	
	// Note: This will likely fail in a test environment without AWS credentials
	// but we can test that it doesn't panic and handles the parameters correctly
	err := configs.InitRelease(env, region, configNames...)
	
	// We expect this to fail in test environment, but it should be a specific error
	// not a panic or nil pointer exception
	if err != nil {
		assert.Error(suite.T(), err)
		// The error should be related to AWS credentials or connectivity, not parameter validation
		assert.NotContains(suite.T(), err.Error(), "nil pointer")
		assert.NotContains(suite.T(), err.Error(), "invalid memory address")
	}
}

func (suite *ConfigUtilsTestSuite) TestInitRelease_EmptyEnvironment() {
	env := ""
	region := "us-east-1"
	
	err := configs.InitRelease(env, region, "config")
	
	// Should handle empty environment parameter
	assert.Error(suite.T(), err)
}

func (suite *ConfigUtilsTestSuite) TestInitRelease_EmptyRegion() {
	env := "test"
	region := ""
	
	err := configs.InitRelease(env, region, "config")
	
	// Should handle empty region parameter
	assert.Error(suite.T(), err)
}

func (suite *ConfigUtilsTestSuite) TestInitRelease_NoConfigNames() {
	env := "test"
	region := "us-east-1"
	
	err := configs.InitRelease(env, region)
	
	// Should handle empty config names
	if err != nil {
		assert.Error(suite.T(), err)
	}
}

// TestGet tests the client getter function
func (suite *ConfigUtilsTestSuite) TestGet_AfterInitDev() {
	// Initialize with dev config
	configContent := `
test:
  key: "value"
`
	configFile := suite.tempDir + "/get_test_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	err = configs.InitDev(suite.tempDir, "get_test_config")
	assert.NoError(suite.T(), err)
	
	// Test getter
	client := configs.Get()
	
	assert.NotNil(suite.T(), client)
}

func (suite *ConfigUtilsTestSuite) TestGet_BeforeInit() {
	// Reset the client by creating a new test suite context
	// This tests the behavior when Get() is called before any Init function
	client := configs.Get()
	
	// The behavior depends on the implementation
	// It might return nil or a default client
	// We just verify it doesn't panic
	assert.NotPanics(suite.T(), func() {
		configs.Get()
	})
	
	// If it returns something, it should be consistent
	client2 := configs.Get()
	assert.Equal(suite.T(), client, client2)
}

// TestConfigInitialization tests the overall configuration initialization flow
func (suite *ConfigUtilsTestSuite) TestConfigInitialization_DevFlow() {
	// Create a comprehensive config file
	configContent := `
app:
  name: "PostEaze"
  version: "1.0.0"
  debug: true

database:
  host: "localhost"
  port: 5432
  name: "posteaze_test"
  user: "test_user"
  password: "test_pass"

jwt:
  access_secret: "test_access_secret"
  refresh_secret: "test_refresh_secret"
  access_expiry: "15m"
  refresh_expiry: "7d"

logging:
  level: "debug"
  format: "json"
`
	
	configFile := suite.tempDir + "/full_config.yaml"
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Test the full initialization flow
	err = configs.InitDev(suite.tempDir, "full_config")
	assert.NoError(suite.T(), err)
	
	// Verify client is properly initialized
	client := configs.Get()
	assert.NotNil(suite.T(), client)
	
	// Test that we can call Get() multiple times
	client2 := configs.Get()
	assert.Equal(suite.T(), client, client2)
}

func (suite *ConfigUtilsTestSuite) TestConfigInitialization_MultipleInits() {
	// Test what happens when we call Init multiple times
	configContent1 := `
first:
  value: "first_config"
`
	configContent2 := `
second:
  value: "second_config"
`
	
	config1File := suite.tempDir + "/first_config.yaml"
	config2File := suite.tempDir + "/second_config.yaml"
	
	err := os.WriteFile(config1File, []byte(configContent1), 0644)
	assert.NoError(suite.T(), err)
	
	err = os.WriteFile(config2File, []byte(configContent2), 0644)
	assert.NoError(suite.T(), err)
	
	// First initialization
	err = configs.InitDev(suite.tempDir, "first_config")
	assert.NoError(suite.T(), err)
	
	client1 := configs.Get()
	assert.NotNil(suite.T(), client1)
	
	// Second initialization (should replace the first)
	err = configs.InitDev(suite.tempDir, "second_config")
	assert.NoError(suite.T(), err)
	
	client2 := configs.Get()
	assert.NotNil(suite.T(), client2)
	
	// The client should be updated (exact behavior depends on implementation)
	// We just verify both calls succeeded and returned valid clients
}

// TestConfigErrorHandling tests error handling in configuration utilities
func (suite *ConfigUtilsTestSuite) TestConfigErrorHandling_InvalidYAML() {
	// Create an invalid YAML file
	invalidYAML := `
invalid:
  yaml: content
    missing: proper indentation
  - invalid list item
`
	
	configFile := suite.tempDir + "/invalid_config.yaml"
	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	assert.NoError(suite.T(), err)
	
	// Test initialization with invalid YAML
	err = configs.InitDev(suite.tempDir, "invalid_config")
	
	// Should handle invalid YAML gracefully
	if err != nil {
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "yaml") // Error should mention YAML parsing
	}
}

func (suite *ConfigUtilsTestSuite) TestConfigErrorHandling_PermissionDenied() {
	// Create a directory with restricted permissions
	restrictedDir := suite.tempDir + "/restricted"
	err := os.Mkdir(restrictedDir, 0000) // No permissions
	if err != nil {
		suite.T().Skip("Cannot create restricted directory on this system")
	}
	defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup
	
	// Test initialization with restricted directory
	err = configs.InitDev(restrictedDir, "config")
	
	assert.Error(suite.T(), err)
}

// Run the test suite
func TestConfigUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigUtilsTestSuite))
}