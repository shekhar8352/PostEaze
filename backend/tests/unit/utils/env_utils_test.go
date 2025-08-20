package utils

import (
	"os"
	"testing"

	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// EnvUtilsTestSuite defines the test suite for environment utilities
type EnvUtilsTestSuite struct {
	suite.Suite
	originalEnvVars map[string]string
	testEnvFile     string
}

// SetupSuite runs before all tests in the suite
func (suite *EnvUtilsTestSuite) SetupSuite() {
	// Store original environment variables that we'll modify
	suite.originalEnvVars = make(map[string]string)
	testVars := []string{"TEST_ENV_VAR", "DB_HOST", "API_KEY", "PORT", "DEBUG", "EMPTY_VAR"}
	
	for _, varName := range testVars {
		if val := os.Getenv(varName); val != "" {
			suite.originalEnvVars[varName] = val
		}
	}
	
	// Create a test .env file
	suite.testEnvFile = ".env.test"
	envContent := `# Test environment file
TEST_ENV_VAR=test_value_from_file
DB_HOST=localhost
API_KEY=secret123
PORT=8080
DEBUG=true
EMPTY_VAR=
MULTILINE_VAR=line1
line2
line3
SPECIAL_CHARS=!@#$%^&*()
SPACES_VAR=value with spaces
`
	
	err := os.WriteFile(suite.testEnvFile, []byte(envContent), 0644)
	if err != nil {
		suite.T().Fatalf("Failed to create test env file: %v", err)
	}
}

// TearDownSuite runs after all tests in the suite
func (suite *EnvUtilsTestSuite) TearDownSuite() {
	// Restore original environment variables
	testVars := []string{"TEST_ENV_VAR", "DB_HOST", "API_KEY", "PORT", "DEBUG", "EMPTY_VAR", "MULTILINE_VAR", "SPECIAL_CHARS", "SPACES_VAR"}
	
	for _, varName := range testVars {
		if originalVal, exists := suite.originalEnvVars[varName]; exists {
			os.Setenv(varName, originalVal)
		} else {
			os.Unsetenv(varName)
		}
	}
	
	// Clean up test env file
	if suite.testEnvFile != "" {
		os.Remove(suite.testEnvFile)
	}
}

// TestInitEnv tests the environment initialization function
func (suite *EnvUtilsTestSuite) TestInitEnv_LoadsEnvironmentVariables() {
	// Set some environment variables before initialization
	os.Setenv("TEST_INIT_VAR", "init_value")
	os.Setenv("ANOTHER_VAR", "another_value")
	
	// Initialize environment
	env.InitEnv()
	
	// The function should load all environment variables
	// We can't directly test the internal envObj map, but we can test the ApplyEnvironmentToString function
	result := env.ApplyEnvironmentToString("Value is ${TEST_INIT_VAR}")
	assert.Equal(suite.T(), "Value is init_value", result)
	
	result = env.ApplyEnvironmentToString("Another is ${ANOTHER_VAR}")
	assert.Equal(suite.T(), "Another is another_value", result)
	
	// Clean up
	os.Unsetenv("TEST_INIT_VAR")
	os.Unsetenv("ANOTHER_VAR")
}

func (suite *EnvUtilsTestSuite) TestInitEnv_HandlesEmptyEnvironment() {
	// Clear some environment variables
	os.Unsetenv("TEST_ENV_VAR")
	os.Unsetenv("DB_HOST")
	
	// Initialize environment
	env.InitEnv()
	
	// Should handle missing variables gracefully
	result := env.ApplyEnvironmentToString("Value is ${TEST_ENV_VAR}")
	assert.Equal(suite.T(), "Value is ", result) // Empty replacement for missing var
}

func (suite *EnvUtilsTestSuite) TestInitEnv_LoadsDotEnvFile() {
	// Change to a temporary directory and create a .env file
	originalDir, _ := os.Getwd()
	tempDir, err := os.MkdirTemp("", "env_test")
	assert.NoError(suite.T(), err)
	defer os.RemoveAll(tempDir)
	
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)
	
	// Create .env file in temp directory
	envContent := `TEST_DOTENV_VAR=dotenv_value
ANOTHER_DOTENV_VAR=another_dotenv_value`
	
	err = os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Initialize environment (should load .env file)
	env.InitEnv()
	
	// Test that .env variables are loaded
	result := env.ApplyEnvironmentToString("Value is ${TEST_DOTENV_VAR}")
	assert.Equal(suite.T(), "Value is dotenv_value", result)
}

func (suite *EnvUtilsTestSuite) TestInitEnv_HandlesInvalidDotEnvFile() {
	// This test verifies that InitEnv doesn't crash with invalid .env files
	originalDir, _ := os.Getwd()
	tempDir, err := os.MkdirTemp("", "env_test")
	assert.NoError(suite.T(), err)
	defer os.RemoveAll(tempDir)
	
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)
	
	// Create invalid .env file
	invalidEnvContent := `INVALID LINE WITHOUT EQUALS
=INVALID_EMPTY_KEY
VALID_VAR=valid_value`
	
	err = os.WriteFile(".env", []byte(invalidEnvContent), 0644)
	assert.NoError(suite.T(), err)
	
	// Should not panic
	assert.NotPanics(suite.T(), func() {
		env.InitEnv()
	})
}

// TestApplyEnvironmentToString tests the environment variable replacement function
func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_SingleVariable() {
	os.Setenv("SINGLE_VAR", "single_value")
	env.InitEnv()
	
	input := "The value is ${SINGLE_VAR}"
	expected := "The value is single_value"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("SINGLE_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_MultipleVariables() {
	os.Setenv("VAR1", "value1")
	os.Setenv("VAR2", "value2")
	os.Setenv("VAR3", "value3")
	env.InitEnv()
	
	input := "${VAR1} and ${VAR2} and ${VAR3}"
	expected := "value1 and value2 and value3"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("VAR1")
	os.Unsetenv("VAR2")
	os.Unsetenv("VAR3")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_NoVariables() {
	env.InitEnv()
	
	input := "This string has no environment variables"
	expected := "This string has no environment variables"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_EmptyString() {
	env.InitEnv()
	
	input := ""
	expected := ""
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_NonExistentVariable() {
	env.InitEnv()
	
	input := "Value is ${NON_EXISTENT_VAR}"
	expected := "Value is " // Empty replacement for non-existent var
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_EmptyVariable() {
	os.Setenv("EMPTY_TEST_VAR", "")
	env.InitEnv()
	
	input := "Value is ${EMPTY_TEST_VAR}"
	expected := "Value is " // Empty replacement for empty var
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("EMPTY_TEST_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_SameVariableMultipleTimes() {
	os.Setenv("REPEAT_VAR", "repeated")
	env.InitEnv()
	
	input := "${REPEAT_VAR} and ${REPEAT_VAR} again"
	expected := "repeated and repeated again"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("REPEAT_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_VariableAtStart() {
	os.Setenv("START_VAR", "start")
	env.InitEnv()
	
	input := "${START_VAR} is at the beginning"
	expected := "start is at the beginning"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("START_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_VariableAtEnd() {
	os.Setenv("END_VAR", "end")
	env.InitEnv()
	
	input := "This is at the ${END_VAR}"
	expected := "This is at the end"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("END_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_OnlyVariable() {
	os.Setenv("ONLY_VAR", "only")
	env.InitEnv()
	
	input := "${ONLY_VAR}"
	expected := "only"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("ONLY_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_SpecialCharactersInValue() {
	os.Setenv("SPECIAL_TEST_VAR", "value with spaces and !@#$%^&*()")
	env.InitEnv()
	
	input := "Special: ${SPECIAL_TEST_VAR}"
	expected := "Special: value with spaces and !@#$%^&*()"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("SPECIAL_TEST_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_URLTemplate() {
	os.Setenv("PROTOCOL", "https")
	os.Setenv("DOMAIN", "api.example.com")
	os.Setenv("VERSION", "v1")
	env.InitEnv()
	
	input := "${PROTOCOL}://${DOMAIN}/${VERSION}/users"
	expected := "https://api.example.com/v1/users"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("PROTOCOL")
	os.Unsetenv("DOMAIN")
	os.Unsetenv("VERSION")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_DatabaseConnectionString() {
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASS", "secret")
	os.Setenv("DB_HOST_TEST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "myapp")
	env.InitEnv()
	
	input := "postgres://${DB_USER}:${DB_PASS}@${DB_HOST_TEST}:${DB_PORT}/${DB_NAME}"
	expected := "postgres://admin:secret@localhost:5432/myapp"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_HOST_TEST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_NestedReplacements() {
	// Test that the function doesn't do nested replacements
	os.Setenv("OUTER_VAR", "${INNER_VAR}")
	os.Setenv("INNER_VAR", "inner_value")
	env.InitEnv()
	
	input := "Value is ${OUTER_VAR}"
	expected := "Value is ${INNER_VAR}" // Should not resolve nested variables
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("OUTER_VAR")
	os.Unsetenv("INNER_VAR")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_CaseSensitive() {
	os.Setenv("CASE_VAR", "lowercase")
	os.Setenv("case_var", "uppercase")
	env.InitEnv()
	
	input := "${CASE_VAR} vs ${case_var}"
	expected := "lowercase vs uppercase"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("CASE_VAR")
	os.Unsetenv("case_var")
}

func (suite *EnvUtilsTestSuite) TestApplyEnvironmentToString_EnvironmentVariableWithEquals() {
	// Test environment variables that contain equals signs in their values
	os.Setenv("EQUALS_VAR", "key=value")
	env.InitEnv()
	
	input := "Config: ${EQUALS_VAR}"
	expected := "Config: key=value"
	
	result := env.ApplyEnvironmentToString(input)
	
	assert.Equal(suite.T(), expected, result)
	
	os.Unsetenv("EQUALS_VAR")
}

// TestEnvironmentVariableHandling tests edge cases in environment variable handling
func (suite *EnvUtilsTestSuite) TestEnvironmentVariableHandling_MultipleEquals() {
	// Test that InitEnv handles environment variables with multiple equals signs correctly
	os.Setenv("MULTI_EQUALS", "key=value=another=value")
	env.InitEnv()
	
	result := env.ApplyEnvironmentToString("${MULTI_EQUALS}")
	assert.Equal(suite.T(), "key=value=another=value", result)
	
	os.Unsetenv("MULTI_EQUALS")
}

func (suite *EnvUtilsTestSuite) TestEnvironmentVariableHandling_WhitespaceInValues() {
	os.Setenv("WHITESPACE_VAR", "  value with leading and trailing spaces  ")
	env.InitEnv()
	
	result := env.ApplyEnvironmentToString("${WHITESPACE_VAR}")
	// The function should preserve whitespace in values
	assert.Equal(suite.T(), "  value with leading and trailing spaces  ", result)
	
	os.Unsetenv("WHITESPACE_VAR")
}

func (suite *EnvUtilsTestSuite) TestEnvironmentVariableHandling_NewlinesInValues() {
	os.Setenv("NEWLINE_VAR", "line1\nline2\nline3")
	env.InitEnv()
	
	result := env.ApplyEnvironmentToString("${NEWLINE_VAR}")
	assert.Equal(suite.T(), "line1\nline2\nline3", result)
	
	os.Unsetenv("NEWLINE_VAR")
}

// Run the test suite
func TestEnvUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(EnvUtilsTestSuite))
}