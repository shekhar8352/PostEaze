package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/env"
)

// AuthIntegrationTestSuite tests complete authentication flows end-to-end
type AuthIntegrationTestSuite struct {
	suite.Suite
	router  *gin.Engine
	db      *sql.DB
	cleanup func()
	ctx     context.Context
}

// SetupSuite initializes the test environment with real database and router
func (s *AuthIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	
	// Set test mode
	gin.SetMode(gin.TestMode)
	
	// Initialize environment variables for testing
	s.setupTestEnvironment()
	
	// Initialize configurations
	s.initializeConfigs()
	
	// Setup test database with real database connection
	db, cleanup, err := testutils.SetupTestDBWithDatabase(s.ctx)
	require.NoError(s.T(), err, "Failed to setup test database")
	s.db = db
	s.cleanup = cleanup
	
	// Initialize router with all middleware and routes
	s.router = s.setupTestRouter()
	
	s.T().Log("Auth integration test suite setup completed")
}

// TearDownSuite cleans up test environment
func (s *AuthIntegrationTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
	s.T().Log("Auth integration test suite teardown completed")
}

// SetupTest prepares each test with clean database state
func (s *AuthIntegrationTestSuite) SetupTest() {
	// Clean up any existing test data
	err := testutils.CleanupTestData(s.ctx, s.db)
	require.NoError(s.T(), err, "Failed to cleanup test data")
}

// setupTestEnvironment configures environment variables for testing
func (s *AuthIntegrationTestSuite) setupTestEnvironment() {
	// Set test environment variables
	os.Setenv("MODE", "dev")
	os.Setenv("BASE_CONFIG_PATH", "../config")
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-key-for-integration-testing")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key-for-integration-testing")
	
	// Initialize environment
	env.InitEnv()
}

// initializeConfigs initializes application configurations for testing
func (s *AuthIntegrationTestSuite) initializeConfigs() {
	configNames := []string{"api", "application", "database"}
	err := configs.InitDev("../config", configNames...)
	require.NoError(s.T(), err, "Failed to initialize configs")
}

// setupTestRouter creates a test router with all routes and middleware
func (s *AuthIntegrationTestSuite) setupTestRouter() *gin.Engine {
	router := gin.New()
	
	// Add logging middleware for debugging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Setup API routes similar to main application
	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	v1 := api.Group("/v1")
	s.addAuthRoutes(v1)
	
	return router
}

// addAuthRoutes adds authentication routes to the router
func (s *AuthIntegrationTestSuite) addAuthRoutes(v1 *gin.RouterGroup) {
	// Import the actual API handlers
	authv1 := v1.Group("/auth")
	
	// Note: We would need to import the actual handlers here
	// For now, we'll create simplified test handlers that call the business logic
	authv1.POST("/signup", s.signupHandler)
	authv1.POST("/login", s.loginHandler)
	authv1.POST("/refresh", s.refreshTokenHandler)
	authv1.POST("/logout", s.logoutHandler)
}

// Test handlers that call actual business logic
func (s *AuthIntegrationTestSuite) signupHandler(c *gin.Context) {
	var body modelsv1.SignupParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid signup data"})
		return
	}
	
	// Call actual business logic
	result, err := businessv1.Signup(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Signed up successfully",
		"data":   result,
	})
}

func (s *AuthIntegrationTestSuite) loginHandler(c *gin.Context) {
	var body modelsv1.LoginParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Email and password are required"})
		return
	}
	
	// Call actual business logic
	result, err := businessv1.Login(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Logged in successfully",
		"data":   result,
	})
}

func (s *AuthIntegrationTestSuite) refreshTokenHandler(c *gin.Context) {
	var body modelsv1.RefreshTokenParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Refresh token is required"})
		return
	}
	
	// Call actual business logic
	result, err := businessv1.RefreshToken(c.Request.Context(), body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Refreshed token successfully",
		"data":   result,
	})
}

func (s *AuthIntegrationTestSuite) logoutHandler(c *gin.Context) {
	// Get refresh token from request body
	var body modelsv1.RefreshTokenParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Refresh token is required"})
		return
	}
	
	// Call actual business logic
	err := businessv1.Logout(c.Request.Context(), body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Logged out successfully",
		"data":   nil,
	})
}

// TestCompleteAuthenticationFlow tests the complete authentication flow from signup to logout
func (s *AuthIntegrationTestSuite) TestCompleteAuthenticationFlow() {
	// Step 1: Test user signup with real business logic
	signupData := modelsv1.SignupParams{
		Name:     "Integration Test User",
		Email:    "integration@test.com",
		Password: "testpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	signupResponse := s.performRequest("POST", "/api/v1/auth/signup", signupData)
	assert.Equal(s.T(), http.StatusOK, signupResponse.Code, "Signup should succeed")
	
	var signupResult map[string]interface{}
	err := json.Unmarshal(signupResponse.Body.Bytes(), &signupResult)
	require.NoError(s.T(), err, "Should parse signup response")
	
	assert.Equal(s.T(), "success", signupResult["status"])
	assert.Contains(s.T(), signupResult, "data")
	
	data := signupResult["data"].(map[string]interface{})
	assert.Contains(s.T(), data, "user")
	assert.Contains(s.T(), data, "access_token")
	assert.Contains(s.T(), data, "refresh_token")
	
	user := data["user"].(map[string]interface{})
	userID := user["id"].(string)
	accessToken := data["access_token"].(string)
	refreshToken := data["refresh_token"].(string)
	
	// Verify user was actually created in database
	createdUser, err := repositories.GetUserByEmail(s.ctx, signupData.Email)
	require.NoError(s.T(), err, "Should retrieve user from database")
	assert.Equal(s.T(), userID, createdUser.ID, "Database user ID should match")
	assert.Equal(s.T(), signupData.Name, createdUser.Name, "Database user name should match")
	assert.Equal(s.T(), signupData.Email, createdUser.Email, "Database user email should match")
	
	s.T().Logf("Signup successful - User ID: %s, Access Token: %s, Refresh Token: %s", userID, accessToken, refreshToken)
	
	// Step 2: Test user login with same credentials
	loginData := modelsv1.LoginParams{
		Email:    signupData.Email,
		Password: signupData.Password,
	}
	
	loginResponse := s.performRequest("POST", "/api/v1/auth/login", loginData)
	assert.Equal(s.T(), http.StatusOK, loginResponse.Code, "Login should succeed")
	
	var loginResult map[string]interface{}
	err = json.Unmarshal(loginResponse.Body.Bytes(), &loginResult)
	require.NoError(s.T(), err, "Should parse login response")
	
	assert.Equal(s.T(), "success", loginResult["status"])
	loginData2 := loginResult["data"].(map[string]interface{})
	newAccessToken := loginData2["access_token"].(string)
	newRefreshToken := loginData2["refresh_token"].(string)
	
	// Verify refresh token was stored in database
	tokenUser, err := repositories.GetUserbyToken(s.ctx, newRefreshToken)
	require.NoError(s.T(), err, "Should retrieve user by refresh token")
	assert.Equal(s.T(), userID, tokenUser.ID, "Token should belong to correct user")
	
	s.T().Logf("Login successful - New Access Token: %s, New Refresh Token: %s", newAccessToken, newRefreshToken)
	
	// Step 3: Test refresh token
	refreshData := modelsv1.RefreshTokenParams{
		RefreshToken: newRefreshToken,
	}
	
	refreshResponse := s.performRequest("POST", "/api/v1/auth/refresh", refreshData)
	assert.Equal(s.T(), http.StatusOK, refreshResponse.Code, "Refresh token should succeed")
	
	var refreshResult map[string]interface{}
	err = json.Unmarshal(refreshResponse.Body.Bytes(), &refreshResult)
	require.NoError(s.T(), err, "Should parse refresh response")
	
	assert.Equal(s.T(), "success", refreshResult["status"])
	refreshedData := refreshResult["data"].(map[string]interface{})
	refreshedAccessToken := refreshedData["access_token"].(string)
	
	// Verify new access token is different
	assert.NotEqual(s.T(), newAccessToken, refreshedAccessToken, "Refreshed access token should be different")
	
	s.T().Logf("Token refresh successful - Refreshed Access Token: %s", refreshedAccessToken)
	
	// Step 4: Test logout with refresh token revocation
	logoutResponse := s.performRequestWithRefreshToken("POST", "/api/v1/auth/logout", nil, newRefreshToken)
	assert.Equal(s.T(), http.StatusOK, logoutResponse.Code, "Logout should succeed")
	
	var logoutResult map[string]interface{}
	err = json.Unmarshal(logoutResponse.Body.Bytes(), &logoutResult)
	require.NoError(s.T(), err, "Should parse logout response")
	
	assert.Equal(s.T(), "success", logoutResult["status"])
	
	// Verify refresh token was revoked in database
	_, err = repositories.GetUserbyToken(s.ctx, newRefreshToken)
	assert.Error(s.T(), err, "Should not be able to use revoked refresh token")
	
	s.T().Log("Logout successful - refresh token revoked")
}

// TestTeamUserSignupFlow tests the complete team user signup flow
func (s *AuthIntegrationTestSuite) TestTeamUserSignupFlow() {
	signupData := modelsv1.SignupParams{
		Name:     "Team Owner",
		Email:    "teamowner@test.com",
		Password: "teampassword123",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Test Team",
	}
	
	response := s.performRequest("POST", "/api/v1/auth/signup", signupData)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Team signup should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse team signup response")
	
	assert.Equal(s.T(), "success", result["status"])
	data := result["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	userID := user["id"].(string)
	
	assert.Equal(s.T(), "team", user["user_type"])
	assert.Equal(s.T(), signupData.Name, user["name"])
	assert.Equal(s.T(), signupData.Email, user["email"])
	
	// Verify user was created in database
	createdUser, err := repositories.GetUserByEmail(s.ctx, signupData.Email)
	require.NoError(s.T(), err, "Should retrieve team user from database")
	assert.Equal(s.T(), userID, createdUser.ID, "Database user ID should match")
	assert.Equal(s.T(), string(modelsv1.UserTypeTeam), createdUser.UserType, "User should be team type")
	
	// Note: Team verification would require additional repository methods
	// For now, we verify the user exists and has correct type
	
	s.T().Logf("Team user signup flow completed successfully - User ID: %s", userID)
}

// TestInvalidAuthenticationScenarios tests various invalid authentication scenarios
func (s *AuthIntegrationTestSuite) TestInvalidAuthenticationScenarios() {
	// Test invalid signup data
	invalidSignupData := modelsv1.SignupParams{
		Name:     "", // Invalid: empty name
		Email:    "invalid-email", // Invalid: not a valid email
		Password: "123", // Invalid: too short
		UserType: modelsv1.UserTypeIndividual,
	}
	
	response := s.performRequest("POST", "/api/v1/auth/signup", invalidSignupData)
	assert.Equal(s.T(), http.StatusBadRequest, response.Code, "Invalid signup should fail")
	
	// Test invalid login data
	invalidLoginData := modelsv1.LoginParams{
		Email:    "", // Invalid: empty email
		Password: "", // Invalid: empty password
	}
	
	response = s.performRequest("POST", "/api/v1/auth/login", invalidLoginData)
	assert.Equal(s.T(), http.StatusBadRequest, response.Code, "Invalid login should fail")
	
	// Test invalid refresh token
	invalidRefreshData := modelsv1.RefreshTokenParams{
		RefreshToken: "", // Invalid: empty refresh token
	}
	
	response = s.performRequest("POST", "/api/v1/auth/refresh", invalidRefreshData)
	assert.Equal(s.T(), http.StatusBadRequest, response.Code, "Invalid refresh token should fail")
	
	s.T().Log("Invalid authentication scenarios tested successfully")
}

// TestConcurrentAuthenticationRequests tests concurrent authentication requests
func (s *AuthIntegrationTestSuite) TestConcurrentAuthenticationRequests() {
	const numConcurrentRequests = 5
	
	signupData := modelsv1.SignupParams{
		Name:     "Concurrent Test User",
		Email:    "concurrent@test.com",
		Password: "concurrentpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	// Create channels to collect results
	responses := make(chan *httptest.ResponseRecorder, numConcurrentRequests)
	
	// Launch concurrent signup requests
	for i := 0; i < numConcurrentRequests; i++ {
		go func(index int) {
			// Modify email to make each request unique
			data := signupData
			data.Email = fmt.Sprintf("concurrent%d@test.com", index)
			response := s.performRequest("POST", "/api/v1/auth/signup", data)
			responses <- response
		}(i)
	}
	
	// Collect and verify responses
	successCount := 0
	for i := 0; i < numConcurrentRequests; i++ {
		response := <-responses
		if response.Code == http.StatusOK {
			successCount++
		}
	}
	
	assert.Equal(s.T(), numConcurrentRequests, successCount, "All concurrent requests should succeed")
	s.T().Logf("Concurrent authentication test completed - %d/%d requests succeeded", successCount, numConcurrentRequests)
}

// TestAuthenticationWithDatabaseTransactions tests authentication with database transaction handling
func (s *AuthIntegrationTestSuite) TestAuthenticationWithDatabaseTransactions() {
	// This test verifies that database transactions work correctly during authentication
	signupData := modelsv1.SignupParams{
		Name:     "Transaction Test User",
		Email:    "transaction@test.com",
		Password: "transactionpassword123",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Transaction Test Team",
	}
	
	// Perform signup which involves multiple database operations in a transaction
	response := s.performRequest("POST", "/api/v1/auth/signup", signupData)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Transaction-based signup should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse transaction signup response")
	
	assert.Equal(s.T(), "success", result["status"])
	
	// Verify that the user was created successfully
	data := result["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.NotEmpty(s.T(), user["id"], "User ID should be generated")
	
	s.T().Log("Database transaction test completed successfully")
}

// performRequest is a helper method to perform HTTP requests on the test router
func (s *AuthIntegrationTestSuite) performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	var err error
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(s.T(), err, "Failed to marshal request body")
		req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		require.NoError(s.T(), err, "Failed to create request")
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, path, nil)
		require.NoError(s.T(), err, "Failed to create request")
	}
	
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, req)
	
	return recorder
}

// performRequestWithRefreshToken is a helper method to perform HTTP requests with refresh token in body
func (s *AuthIntegrationTestSuite) performRequestWithRefreshToken(method, path string, body interface{}, refreshToken string) *httptest.ResponseRecorder {
	// If no body provided, create one with refresh token
	if body == nil {
		body = modelsv1.RefreshTokenParams{RefreshToken: refreshToken}
	}
	
	return s.performRequest(method, path, body)
}

// TestAuthIntegrationSuite runs the authentication integration test suite
func TestAuthIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}