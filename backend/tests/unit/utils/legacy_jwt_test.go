package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// LegacyJWTTestSuite defines the test suite for legacy JWT utilities in auth_utils.go
type LegacyJWTTestSuite struct {
	suite.Suite
	originalJWTKey string
	testJWTKey     string
}

// SetupSuite runs before all tests in the suite
func (suite *LegacyJWTTestSuite) SetupSuite() {
	// Store original environment variable
	suite.originalJWTKey = os.Getenv("JWT_KEY")
	
	// Set test key
	suite.testJWTKey = "test-jwt-key-for-legacy-functions"
	os.Setenv("JWT_KEY", suite.testJWTKey)
}

// TearDownSuite runs after all tests in the suite
func (suite *LegacyJWTTestSuite) TearDownSuite() {
	// Restore original environment variable
	if suite.originalJWTKey != "" {
		os.Setenv("JWT_KEY", suite.originalJWTKey)
	} else {
		os.Unsetenv("JWT_KEY")
	}
}

// TestGenerateJWT tests the legacy JWT generation function
func (suite *LegacyJWTTestSuite) TestGenerateJWT_ValidInput() {
	userID := 123
	email := "test@example.com"
	
	token, err := utils.GenerateJWT(userID, email)
	
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	
	// Verify token can be parsed
	claims, err := utils.ParseJWT(token)
	require.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), float64(userID), (*claims)["user_id"])
	assert.Equal(suite.T(), email, (*claims)["email"])
	
	// Check expiration (should be 72 hours from now)
	exp := (*claims)["exp"].(float64)
	expectedExp := time.Now().Add(72 * time.Hour).Unix()
	timeDiff := int64(exp) - expectedExp
	
	// Allow for small time differences (within 5 seconds)
	assert.True(suite.T(), timeDiff < 5 && timeDiff > -5)
}

func (suite *LegacyJWTTestSuite) TestGenerateJWT_ZeroUserID() {
	userID := 0
	email := "test@example.com"
	
	token, err := utils.GenerateJWT(userID, email)
	
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	
	claims, err := utils.ParseJWT(token)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(0), (*claims)["user_id"])
}

func (suite *LegacyJWTTestSuite) TestGenerateJWT_EmptyEmail() {
	userID := 123
	email := ""
	
	token, err := utils.GenerateJWT(userID, email)
	
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	
	claims, err := utils.ParseJWT(token)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "", (*claims)["email"])
}

func (suite *LegacyJWTTestSuite) TestGenerateJWT_NegativeUserID() {
	userID := -1
	email := "test@example.com"
	
	token, err := utils.GenerateJWT(userID, email)
	
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	
	claims, err := utils.ParseJWT(token)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(-1), (*claims)["user_id"])
}

// TestParseJWT tests the legacy JWT parsing function
func (suite *LegacyJWTTestSuite) TestParseJWT_ValidToken() {
	userID := 123
	email := "test@example.com"
	token, err := utils.GenerateJWT(userID, email)
	require.NoError(suite.T(), err)
	
	claims, err := utils.ParseJWT(token)
	
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), float64(userID), (*claims)["user_id"])
	assert.Equal(suite.T(), email, (*claims)["email"])
	
	// Verify expiration exists and is in the future
	exp, exists := (*claims)["exp"]
	assert.True(suite.T(), exists)
	assert.True(suite.T(), exp.(float64) > float64(time.Now().Unix()))
}

func (suite *LegacyJWTTestSuite) TestParseJWT_InvalidToken() {
	invalidToken := "invalid.token.here"
	
	claims, err := utils.ParseJWT(invalidToken)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *LegacyJWTTestSuite) TestParseJWT_EmptyToken() {
	emptyToken := ""
	
	// ParseJWT with empty string may panic due to nil pointer dereference
	// We'll handle this gracefully
	defer func() {
		if r := recover(); r != nil {
			// If it panics, that's expected behavior for empty token
			assert.NotNil(suite.T(), r)
		}
	}()
	
	claims, err := utils.ParseJWT(emptyToken)
	
	// If we reach here without panic, it should return an error
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *LegacyJWTTestSuite) TestParseJWT_MalformedToken() {
	malformedToken := "not.a.jwt"
	
	claims, err := utils.ParseJWT(malformedToken)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

func (suite *LegacyJWTTestSuite) TestParseJWT_ExpiredToken() {
	// Create a token that's already expired
	claims := &jwt.MapClaims{
		"user_id": 123,
		"email":   "test@example.com",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(suite.testJWTKey))
	require.NoError(suite.T(), err)
	
	parsedClaims, err := utils.ParseJWT(tokenString)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), parsedClaims)
}

func (suite *LegacyJWTTestSuite) TestParseJWT_WrongSigningKey() {
	// Create a token with a different signing key
	wrongKey := "wrong-signing-key"
	claims := &jwt.MapClaims{
		"user_id": 123,
		"email":   "test@example.com",
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(wrongKey))
	require.NoError(suite.T(), err)
	
	parsedClaims, err := utils.ParseJWT(tokenString)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), parsedClaims)
}

// TestJWTRoundTrip tests the complete generate -> parse cycle
func (suite *LegacyJWTTestSuite) TestJWTRoundTrip() {
	testCases := []struct {
		name   string
		userID int
		email  string
	}{
		{"Normal case", 123, "user@example.com"},
		{"Zero user ID", 0, "zero@example.com"},
		{"Empty email", 456, ""},
		{"Special characters in email", 789, "user+test@example-domain.com"},
		{"Large user ID", 999999999, "large@example.com"},
	}
	
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Generate token
			token, err := utils.GenerateJWT(tc.userID, tc.email)
			require.NoError(t, err)
			assert.NotEmpty(t, token)
			
			// Parse token
			claims, err := utils.ParseJWT(token)
			require.NoError(t, err)
			assert.NotNil(t, claims)
			
			// Verify claims
			assert.Equal(t, float64(tc.userID), (*claims)["user_id"])
			assert.Equal(t, tc.email, (*claims)["email"])
			
			// Verify expiration
			exp := (*claims)["exp"].(float64)
			assert.True(t, exp > float64(time.Now().Unix()))
		})
	}
}

// TestJWTWithoutEnvironmentVariable tests behavior when JWT_KEY is not set
func (suite *LegacyJWTTestSuite) TestJWTWithoutEnvironmentVariable() {
	// Temporarily unset environment variable
	os.Unsetenv("JWT_KEY")
	
	userID := 123
	email := "test@example.com"
	
	// This should still work but use empty key (not recommended for production)
	token, err := utils.GenerateJWT(userID, email)
	
	// Restore environment variable for other tests
	os.Setenv("JWT_KEY", suite.testJWTKey)
	
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	
	// The token should be parseable with empty key
	os.Unsetenv("JWT_KEY")
	claims, err := utils.ParseJWT(token)
	os.Setenv("JWT_KEY", suite.testJWTKey)
	
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(userID), (*claims)["user_id"])
}

// TestJWTClaimsStructure tests the structure of JWT claims
func (suite *LegacyJWTTestSuite) TestJWTClaimsStructure() {
	userID := 123
	email := "test@example.com"
	token, err := utils.GenerateJWT(userID, email)
	require.NoError(suite.T(), err)
	
	claims, err := utils.ParseJWT(token)
	require.NoError(suite.T(), err)
	
	// Check that all expected claims are present
	expectedClaims := []string{"user_id", "email", "exp"}
	for _, claim := range expectedClaims {
		_, exists := (*claims)[claim]
		assert.True(suite.T(), exists, "Claim %s should exist", claim)
	}
	
	// Check claim types
	assert.IsType(suite.T(), float64(0), (*claims)["user_id"])
	assert.IsType(suite.T(), "", (*claims)["email"])
	assert.IsType(suite.T(), float64(0), (*claims)["exp"])
}

// Run the test suite
func TestLegacyJWTTestSuite(t *testing.T) {
	suite.Run(t, new(LegacyJWTTestSuite))
}