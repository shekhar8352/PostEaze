package examples_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// IntegrationTestSuite demonstrates comprehensive integration testing
type IntegrationTestSuite struct {
	suite.Suite
	server   *httptest.Server
	client   *http.Client
	dbHelper *testutils.DatabaseTestHelper
}

// SetupSuite runs once before all tests in the suite
func (s *IntegrationTestSuite) SetupSuite() {
	// Skip integration tests in short mode
	if testing.Short() {
		s.T().Skip("Skipping integration tests in short mode")
	}

	// Setup test database
	s.dbHelper = testutils.NewDatabaseTestHelper()
	err := s.dbHelper.SetupIntegrationDB()
	s.Require().NoError(err)

	// Setup test server
	router := s.setupTestRouter()
	s.server = httptest.NewServer(router)

	// Setup HTTP client with timeout
	s.client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

// TearDownSuite runs once after all tests in the suite
func (s *IntegrationTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
	if s.dbHelper != nil {
		s.dbHelper.CleanupIntegrationDB()
	}
}

// SetupTest runs before each test method
func (s *IntegrationTestSuite) SetupTest() {
	// Clean database state between tests
	err := s.dbHelper.CleanupTestData()
	s.Require().NoError(err)

	// Load fresh test fixtures
	err = s.dbHelper.LoadTestFixtures()
	s.Require().NoError(err)
}

// Example 1: Complete authentication flow integration test
func (s *IntegrationTestSuite) TestAuthenticationFlow_CompleteWorkflow_Success() {
	// Test user signup
	signupData := map[string]interface{}{
		"name":     "Integration Test User",
		"email":    "integration@example.com",
		"password": "password123",
		"userType": "individual",
	}

	signupResp := s.postJSON("/api/v1/auth/signup", signupData)
	s.Equal(http.StatusCreated, signupResp.StatusCode)

	var signupResult map[string]interface{}
	err := json.NewDecoder(signupResp.Body).Decode(&signupResult)
	s.NoError(err)
	s.Equal("success", signupResult["status"])

	userID := signupResult["data"].(map[string]interface{})["user_id"].(string)
	s.NotEmpty(userID)

	// Test user login
	loginData := map[string]interface{}{
		"email":    "integration@example.com",
		"password": "password123",
	}

	loginResp := s.postJSON("/api/v1/auth/login", loginData)
	s.Equal(http.StatusOK, loginResp.StatusCode)

	var loginResult map[string]interface{}
	err = json.NewDecoder(loginResp.Body).Decode(&loginResult)
	s.NoError(err)
	s.Equal("success", loginResult["status"])

	// Extract tokens
	data := loginResult["data"].(map[string]interface{})
	accessToken := data["access_token"].(string)
	refreshToken := data["refresh_token"].(string)
	s.NotEmpty(accessToken)
	s.NotEmpty(refreshToken)

	// Test protected endpoint with access token
	profileResp := s.getWithAuth("/api/v1/users/profile", accessToken)
	s.Equal(http.StatusOK, profileResp.StatusCode)

	var profileResult map[string]interface{}
	err = json.NewDecoder(profileResp.Body).Decode(&profileResult)
	s.NoError(err)
	s.Equal("success", profileResult["status"])

	profileData := profileResult["data"].(map[string]interface{})
	s.Equal(userID, profileData["id"])
	s.Equal("Integration Test User", profileData["name"])

	// Test token refresh
	refreshData := map[string]interface{}{
		"refresh_token": refreshToken,
	}

	refreshResp := s.postJSON("/api/v1/auth/refresh", refreshData)
	s.Equal(http.StatusOK, refreshResp.StatusCode)

	var refreshResult map[string]interface{}
	err = json.NewDecoder(refreshResp.Body).Decode(&refreshResult)
	s.NoError(err)
	s.Equal("success", refreshResult["status"])

	// Test logout
	logoutResp := s.postJSONWithAuth("/api/v1/auth/logout", map[string]interface{}{}, accessToken)
	s.Equal(http.StatusOK, logoutResp.StatusCode)

	// Verify token is invalidated
	profileResp2 := s.getWithAuth("/api/v1/users/profile", accessToken)
	s.Equal(http.StatusUnauthorized, profileResp2.StatusCode)
}

// Example 2: Database transaction integration test
func (s *IntegrationTestSuite) TestUserCreation_DatabaseTransaction_ConsistentState() {
	// Create team user (requires both user and team creation)
	teamUserData := map[string]interface{}{
		"name":      "Team Owner",
		"email":     "owner@company.com",
		"password":  "password123",
		"userType":  "team",
		"teamName":  "Test Company",
	}

	resp := s.postJSON("/api/v1/auth/signup", teamUserData)
	s.Equal(http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.NoError(err)

	data := result["data"].(map[string]interface{})
	userID := data["user_id"].(string)
	teamID := data["team_id"].(string)

	// Verify user exists in database
	user, err := s.dbHelper.GetUserByID(userID)
	s.NoError(err)
	s.NotNil(user)
	s.Equal("Team Owner", user.Name)
	s.Equal("owner@company.com", user.Email)
	s.Equal(teamID, user.TeamID)

	// Verify team exists in database
	team, err := s.dbHelper.GetTeamByID(teamID)
	s.NoError(err)
	s.NotNil(team)
	s.Equal("Test Company", team.Name)
	s.Equal(userID, team.OwnerID)

	// Verify team membership
	membership, err := s.dbHelper.GetTeamMembership(teamID, userID)
	s.NoError(err)
	s.NotNil(membership)
	s.Equal("owner", membership.Role)
}

// Example 3: Error handling and rollback integration test
func (s *IntegrationTestSuite) TestUserCreation_DuplicateEmail_NoPartialState() {
	// Create first user
	userData1 := map[string]interface{}{
		"name":     "First User",
		"email":    "duplicate@example.com",
		"password": "password123",
		"userType": "individual",
	}

	resp1 := s.postJSON("/api/v1/auth/signup", userData1)
	s.Equal(http.StatusCreated, resp1.StatusCode)

	// Attempt to create second user with same email
	userData2 := map[string]interface{}{
		"name":     "Second User",
		"email":    "duplicate@example.com", // Same email
		"password": "password456",
		"userType": "individual",
	}

	resp2 := s.postJSON("/api/v1/auth/signup", userData2)
	s.Equal(http.StatusConflict, resp2.StatusCode)

	var errorResult map[string]interface{}
	err := json.NewDecoder(resp2.Body).Decode(&errorResult)
	s.NoError(err)
	s.Equal("error", errorResult["status"])
	s.Contains(errorResult["message"], "email already exists")

	// Verify only one user exists with that email
	users, err := s.dbHelper.GetUsersByEmail("duplicate@example.com")
	s.NoError(err)
	s.Len(users, 1)
	s.Equal("First User", users[0].Name)
}

// Example 4: Concurrent request handling integration test
func (s *IntegrationTestSuite) TestConcurrentUserCreation_NoRaceConditions() {
	numConcurrentRequests := 10
	results := make(chan *http.Response, numConcurrentRequests)
	errors := make(chan error, numConcurrentRequests)

	// Send concurrent signup requests
	for i := 0; i < numConcurrentRequests; i++ {
		go func(index int) {
			userData := map[string]interface{}{
				"name":     fmt.Sprintf("Concurrent User %d", index),
				"email":    fmt.Sprintf("concurrent%d@example.com", index),
				"password": "password123",
				"userType": "individual",
			}

			resp := s.postJSON("/api/v1/auth/signup", userData)
			results <- resp
		}(i)
	}

	// Collect all responses
	successCount := 0
	for i := 0; i < numConcurrentRequests; i++ {
		select {
		case resp := <-results:
			if resp.StatusCode == http.StatusCreated {
				successCount++
			}
		case err := <-errors:
			s.Fail("Unexpected error in concurrent request", err.Error())
		case <-time.After(30 * time.Second):
			s.Fail("Timeout waiting for concurrent requests")
		}
	}

	// All requests should succeed
	s.Equal(numConcurrentRequests, successCount)

	// Verify all users were created in database
	users, err := s.dbHelper.GetAllUsers()
	s.NoError(err)
	s.GreaterOrEqual(len(users), numConcurrentRequests)
}

// Example 5: Log management integration test
func (s *IntegrationTestSuite) TestLogManagement_CompleteWorkflow_Success() {
	// First, authenticate to get access token
	loginData := map[string]interface{}{
		"email":    "admin@example.com", // From test fixtures
		"password": "admin123",
	}

	loginResp := s.postJSON("/api/v1/auth/login", loginData)
	s.Equal(http.StatusOK, loginResp.StatusCode)

	var loginResult map[string]interface{}
	err := json.NewDecoder(loginResp.Body).Decode(&loginResult)
	s.NoError(err)

	accessToken := loginResult["data"].(map[string]interface{})["access_token"].(string)

	// Test getting logs by date
	today := time.Now().Format("2006-01-02")
	logsResp := s.getWithAuth(fmt.Sprintf("/api/v1/logs/date/%s", today), accessToken)
	s.Equal(http.StatusOK, logsResp.StatusCode)

	var logsResult map[string]interface{}
	err = json.NewDecoder(logsResp.Body).Decode(&logsResult)
	s.NoError(err)
	s.Equal("success", logsResult["status"])

	logs := logsResult["data"].([]interface{})
	s.GreaterOrEqual(len(logs), 0) // May be empty if no logs for today

	// Test getting specific log by ID (if logs exist)
	if len(logs) > 0 {
		firstLog := logs[0].(map[string]interface{})
		logID := firstLog["id"].(string)

		logResp := s.getWithAuth(fmt.Sprintf("/api/v1/logs/%s", logID), accessToken)
		s.Equal(http.StatusOK, logResp.StatusCode)

		var logResult map[string]interface{}
		err = json.NewDecoder(logResp.Body).Decode(&logResult)
		s.NoError(err)
		s.Equal("success", logResult["status"])

		logData := logResult["data"].(map[string]interface{})
		s.Equal(logID, logData["id"])
	}
}

// Example 6: Performance and load integration test
func (s *IntegrationTestSuite) TestAPIPerformance_ResponseTimes_WithinLimits() {
	// Authenticate first
	loginData := map[string]interface{}{
		"email":    "admin@example.com",
		"password": "admin123",
	}

	loginResp := s.postJSON("/api/v1/auth/login", loginData)
	s.Equal(http.StatusOK, loginResp.StatusCode)

	var loginResult map[string]interface{}
	err := json.NewDecoder(loginResp.Body).Decode(&loginResult)
	s.NoError(err)

	accessToken := loginResult["data"].(map[string]interface{})["access_token"].(string)

	// Test multiple endpoints for performance
	endpoints := []struct {
		name   string
		method string
		path   string
		maxDuration time.Duration
	}{
		{"user profile", "GET", "/api/v1/users/profile", 500 * time.Millisecond},
		{"logs by date", "GET", "/api/v1/logs/date/" + time.Now().Format("2006-01-02"), 1 * time.Second},
	}

	for _, endpoint := range endpoints {
		s.Run(endpoint.name, func() {
			start := time.Now()
			
			var resp *http.Response
			if endpoint.method == "GET" {
				resp = s.getWithAuth(endpoint.path, accessToken)
			}
			
			duration := time.Since(start)
			
			s.Equal(http.StatusOK, resp.StatusCode)
			s.Less(duration, endpoint.maxDuration, 
				"Endpoint %s took %v, expected less than %v", 
				endpoint.path, duration, endpoint.maxDuration)
		})
	}
}

// Example 7: Data consistency across multiple operations
func (s *IntegrationTestSuite) TestDataConsistency_MultipleOperations_ConsistentState() {
	// Create team user
	teamUserData := map[string]interface{}{
		"name":      "Team Owner",
		"email":     "consistency@example.com",
		"password":  "password123",
		"userType":  "team",
		"teamName":  "Consistency Test Team",
	}

	signupResp := s.postJSON("/api/v1/auth/signup", teamUserData)
	s.Equal(http.StatusCreated, signupResp.StatusCode)

	var signupResult map[string]interface{}
	err := json.NewDecoder(signupResp.Body).Decode(&signupResult)
	s.NoError(err)

	userID := signupResult["data"].(map[string]interface{})["user_id"].(string)
	teamID := signupResult["data"].(map[string]interface{})["team_id"].(string)

	// Login to get access token
	loginData := map[string]interface{}{
		"email":    "consistency@example.com",
		"password": "password123",
	}

	loginResp := s.postJSON("/api/v1/auth/login", loginData)
	s.Equal(http.StatusOK, loginResp.StatusCode)

	var loginResult map[string]interface{}
	err = json.NewDecoder(loginResp.Body).Decode(&loginResult)
	s.NoError(err)

	accessToken := loginResult["data"].(map[string]interface{})["access_token"].(string)

	// Update user profile
	updateData := map[string]interface{}{
		"name": "Updated Team Owner",
		"bio":  "Updated bio for consistency test",
	}

	updateResp := s.putJSONWithAuth("/api/v1/users/profile", updateData, accessToken)
	s.Equal(http.StatusOK, updateResp.StatusCode)

	// Verify consistency across different endpoints
	
	// 1. Check user profile endpoint
	profileResp := s.getWithAuth("/api/v1/users/profile", accessToken)
	s.Equal(http.StatusOK, profileResp.StatusCode)

	var profileResult map[string]interface{}
	err = json.NewDecoder(profileResp.Body).Decode(&profileResult)
	s.NoError(err)

	profileData := profileResult["data"].(map[string]interface{})
	s.Equal("Updated Team Owner", profileData["name"])
	s.Equal("Updated bio for consistency test", profileData["bio"])

	// 2. Check team members endpoint
	membersResp := s.getWithAuth(fmt.Sprintf("/api/v1/teams/%s/members", teamID), accessToken)
	s.Equal(http.StatusOK, membersResp.StatusCode)

	var membersResult map[string]interface{}
	err = json.NewDecoder(membersResp.Body).Decode(&membersResult)
	s.NoError(err)

	members := membersResult["data"].([]interface{})
	s.Len(members, 1)

	ownerMember := members[0].(map[string]interface{})
	s.Equal(userID, ownerMember["user_id"])
	s.Equal("Updated Team Owner", ownerMember["name"]) // Should reflect the update

	// 3. Verify in database directly
	user, err := s.dbHelper.GetUserByID(userID)
	s.NoError(err)
	s.Equal("Updated Team Owner", user.Name)
	s.Equal("Updated bio for consistency test", user.Bio)
}

// Helper methods for making HTTP requests

func (s *IntegrationTestSuite) postJSON(path string, data interface{}) *http.Response {
	jsonData, err := json.Marshal(data)
	s.Require().NoError(err)

	resp, err := s.client.Post(s.server.URL+path, "application/json", bytes.NewBuffer(jsonData))
	s.Require().NoError(err)

	return resp
}

func (s *IntegrationTestSuite) postJSONWithAuth(path string, data interface{}, token string) *http.Response {
	jsonData, err := json.Marshal(data)
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.server.URL+path, bytes.NewBuffer(jsonData))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *IntegrationTestSuite) putJSONWithAuth(path string, data interface{}, token string) *http.Response {
	jsonData, err := json.Marshal(data)
	s.Require().NoError(err)

	req, err := http.NewRequest("PUT", s.server.URL+path, bytes.NewBuffer(jsonData))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *IntegrationTestSuite) getWithAuth(path, token string) *http.Response {
	req, err := http.NewRequest("GET", s.server.URL+path, nil)
	s.Require().NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *IntegrationTestSuite) setupTestRouter() *gin.Engine {
	// This would set up your actual Gin router with all middleware and handlers
	// For this example, we'll create a minimal router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add your actual routes here
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", func(c *gin.Context) { /* handler implementation */ })
			auth.POST("/login", func(c *gin.Context) { /* handler implementation */ })
			auth.POST("/refresh", func(c *gin.Context) { /* handler implementation */ })
			auth.POST("/logout", func(c *gin.Context) { /* handler implementation */ })
		}

		users := api.Group("/users")
		{
			users.GET("/profile", func(c *gin.Context) { /* handler implementation */ })
			users.PUT("/profile", func(c *gin.Context) { /* handler implementation */ })
		}

		logs := api.Group("/logs")
		{
			logs.GET("/date/:date", func(c *gin.Context) { /* handler implementation */ })
			logs.GET("/:id", func(c *gin.Context) { /* handler implementation */ })
		}

		teams := api.Group("/teams")
		{
			teams.GET("/:id/members", func(c *gin.Context) { /* handler implementation */ })
		}
	}

	return router
}

// Run the integration test suite
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Example of a simple integration test without a test suite
func TestHealthCheck_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test server
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "healthy", result["status"])
}

// Helper types (these would normally be in your actual code)

type DatabaseUser struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Bio       string    `db:"bio"`
	UserType  string    `db:"user_type"`
	TeamID    string    `db:"team_id"`
	CreatedAt time.Time `db:"created_at"`
}

type DatabaseTeam struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	OwnerID string `db:"owner_id"`
}

type DatabaseTeamMembership struct {
	TeamID string `db:"team_id"`
	UserID string `db:"user_id"`
	Role   string `db:"role"`
}