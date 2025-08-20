package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/env"
)

// ComprehensiveIntegrationTestSuite tests all integration requirements from task 12
type ComprehensiveIntegrationTestSuite struct {
	suite.Suite
	router      *gin.Engine
	testLogDir  string
	cleanup     func()
	ctx         context.Context
	originalDir string
}

// SetupSuite initializes the complete test environment
func (s *ComprehensiveIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	gin.SetMode(gin.TestMode)
	
	// Setup test environment
	s.setupTestEnvironment()
	s.initializeConfigs()
	
	// Setup test database with actual database connections
	_, cleanup, err := testutils.SetupTestDBWithDatabase(s.ctx)
	require.NoError(s.T(), err, "Failed to setup test database")
	s.cleanup = cleanup
	
	// Setup test log directory for real file system operations
	s.testLogDir, err = os.MkdirTemp("", "comprehensive_test_logs_*")
	require.NoError(s.T(), err, "Failed to create temp log directory")
	s.originalDir = utils.GetLogDirectory()
	s.setTestLogDirectory()
	
	// Setup complete router with all routes
	s.router = s.setupCompleteRouter()
	
	// Create comprehensive test log files
	s.createComprehensiveTestLogFiles()
	
	s.T().Log("Comprehensive integration test suite setup completed")
}

// TearDownSuite cleans up test environment
func (s *ComprehensiveIntegrationTestSuite) TearDownSuite() {
	// Restore original log directory
	if s.originalDir != "" {
		os.Setenv("LOG_DIR", s.originalDir)
	} else {
		os.Unsetenv("LOG_DIR")
	}
	
	// Clean up temporary directory
	if s.testLogDir != "" {
		os.RemoveAll(s.testLogDir)
	}
	
	if s.cleanup != nil {
		s.cleanup()
	}
	
	s.T().Log("Comprehensive integration test suite teardown completed")
}

// SetupTest prepares each test with clean state
func (s *ComprehensiveIntegrationTestSuite) SetupTest() {
	// Clean up database - we'll use a direct SQL connection for cleanup
	// This is a simplified approach for the integration test
	
	// Recreate test log files
	s.createComprehensiveTestLogFiles()
}

// setupTestEnvironment configures environment variables for testing
func (s *ComprehensiveIntegrationTestSuite) setupTestEnvironment() {
	os.Setenv("MODE", "dev")
	os.Setenv("BASE_CONFIG_PATH", "../config")
	os.Setenv("JWT_ACCESS_SECRET", "comprehensive-test-access-secret")
	os.Setenv("JWT_REFRESH_SECRET", "comprehensive-test-refresh-secret")
	env.InitEnv()
}

// initializeConfigs initializes application configurations for testing
func (s *ComprehensiveIntegrationTestSuite) initializeConfigs() {
	configNames := []string{"api", "application", "database"}
	err := configs.InitDev("../config", configNames...)
	require.NoError(s.T(), err, "Failed to initialize configs")
}

// setTestLogDirectory configures the log directory for testing
func (s *ComprehensiveIntegrationTestSuite) setTestLogDirectory() {
	os.Setenv("LOG_DIR", s.testLogDir)
}

// setupCompleteRouter creates a complete router with all authentication and log routes
func (s *ComprehensiveIntegrationTestSuite) setupCompleteRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	api := router.Group("/api")
	v1 := api.Group("/v1")
	
	// Authentication routes
	authv1 := v1.Group("/auth")
	authv1.POST("/signup", s.signupHandler)
	authv1.POST("/login", s.loginHandler)
	authv1.POST("/refresh", s.refreshTokenHandler)
	authv1.POST("/logout", s.logoutHandler)
	
	// Log routes
	logv1 := v1.Group("/logs")
	logv1.GET("/date/:date", s.getLogsByDateHandler)
	logv1.GET("/id/:log_id", s.getLogByIDHandler)
	
	return router
}

// Authentication handlers using real business logic
func (s *ComprehensiveIntegrationTestSuite) signupHandler(c *gin.Context) {
	var body modelsv1.SignupParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid signup data"})
		return
	}
	
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

func (s *ComprehensiveIntegrationTestSuite) loginHandler(c *gin.Context) {
	var body modelsv1.LoginParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Email and password are required"})
		return
	}
	
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

func (s *ComprehensiveIntegrationTestSuite) refreshTokenHandler(c *gin.Context) {
	var body modelsv1.RefreshTokenParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Refresh token is required"})
		return
	}
	
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

func (s *ComprehensiveIntegrationTestSuite) logoutHandler(c *gin.Context) {
	var body modelsv1.RefreshTokenParams
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Refresh token is required"})
		return
	}
	
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

// Log handlers using real business logic with file system operations
func (s *ComprehensiveIntegrationTestSuite) getLogsByDateHandler(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Date is required"})
		return
	}
	
	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Invalid date format"})
		return
	}
	
	logs, total, err := businessv1.ReadLogsByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Read logs by date successfully",
		"data":   gin.H{"logs": logs, "total": total},
	})
}

func (s *ComprehensiveIntegrationTestSuite) getLogByIDHandler(c *gin.Context) {
	logID := c.Param("log_id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "Log ID is required"})
		return
	}
	
	logs, err := businessv1.ReadLogsByLogID(c.Request.Context(), logID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "Read log by ID successfully",
		"data":   logs,
	})
}

// createComprehensiveTestLogFiles creates comprehensive test log files
func (s *ComprehensiveIntegrationTestSuite) createComprehensiveTestLogFiles() {
	dates := []string{
		time.Now().Format("2006-01-02"),
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
	}
	
	for _, date := range dates {
		entries := s.generateComprehensiveLogEntries(date, 20)
		s.createLogFile(date, entries)
	}
}

// createLogFile creates a log file with specified content
func (s *ComprehensiveIntegrationTestSuite) createLogFile(date string, entries []modelsv1.LogEntry) {
	filename := fmt.Sprintf("app-%s.log", date)
	filepath := filepath.Join(s.testLogDir, filename)
	
	file, err := os.Create(filepath)
	require.NoError(s.T(), err, "Failed to create log file: %s", filename)
	defer file.Close()
	
	for _, entry := range entries {
		jsonData, err := json.Marshal(entry)
		require.NoError(s.T(), err, "Failed to marshal log entry")
		
		_, err = file.WriteString(string(jsonData) + "\n")
		require.NoError(s.T(), err, "Failed to write log entry")
	}
}

// generateComprehensiveLogEntries generates comprehensive log entries for testing
func (s *ComprehensiveIntegrationTestSuite) generateComprehensiveLogEntries(date string, count int) []modelsv1.LogEntry {
	entries := make([]modelsv1.LogEntry, count)
	
	for i := 0; i < count; i++ {
		timestamp := fmt.Sprintf("%sT%02d:%02d:%02d.000Z", date, i%24, i%60, i%60)
		logID := fmt.Sprintf("comprehensive-log-%s-%03d", date, i+1)
		
		entries[i] = modelsv1.LogEntry{
			Timestamp: timestamp,
			Level:     s.getLogLevel(i),
			Message:   fmt.Sprintf("Comprehensive test log message %d for date %s", i+1, date),
			LogID:     logID,
			Method:    s.getHTTPMethod(i),
			Path:      fmt.Sprintf("/api/v1/comprehensive/%d", i+1),
			Status:    s.getHTTPStatus(i),
			Duration:  fmt.Sprintf("%dms", (i+1)*3),
			IP:        fmt.Sprintf("10.0.0.%d", (i%254)+1),
			UserAgent: "Comprehensive-Test-Client/1.0",
			Extra: map[string]string{
				"request_id": fmt.Sprintf("req-%s-%03d", date, i+1),
				"user_id":    fmt.Sprintf("user-%03d", i+1),
				"session_id": fmt.Sprintf("session-%03d", i+1),
			},
		}
	}
	
	return entries
}

func (s *ComprehensiveIntegrationTestSuite) getLogLevel(index int) string {
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR", "TRACE"}
	return levels[index%len(levels)]
}

func (s *ComprehensiveIntegrationTestSuite) getHTTPMethod(index int) string {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	return methods[index%len(methods)]
}

func (s *ComprehensiveIntegrationTestSuite) getHTTPStatus(index int) int {
	statuses := []int{200, 201, 400, 401, 404, 500}
	return statuses[index%len(statuses)]
}

// TEST 1: Complete Authentication Flows (Requirement 1.3, 2.1)
func (s *ComprehensiveIntegrationTestSuite) TestCompleteAuthenticationFlows() {
	s.T().Log("=== Testing Complete Authentication Flows ===")
	
	// Test Individual User Authentication Flow
	s.testIndividualUserAuthFlow()
	
	// Test Team User Authentication Flow  
	s.testTeamUserAuthFlow()
	
	// Test Authentication Error Scenarios
	s.testAuthenticationErrorScenarios()
	
	s.T().Log("=== Complete Authentication Flows Test Completed ===")
}

func (s *ComprehensiveIntegrationTestSuite) testIndividualUserAuthFlow() {
	s.T().Log("--- Testing Individual User Authentication Flow ---")
	
	// Step 1: Individual User Signup
	signupData := modelsv1.SignupParams{
		Name:     "Individual Test User",
		Email:    "individual@comprehensive.test",
		Password: "individualpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	signupResponse := s.performRequest("POST", "/api/v1/auth/signup", signupData)
	assert.Equal(s.T(), http.StatusOK, signupResponse.Code, "Individual signup should succeed")
	
	var signupResult map[string]interface{}
	err := json.Unmarshal(signupResponse.Body.Bytes(), &signupResult)
	require.NoError(s.T(), err, "Should parse individual signup response")
	
	data := signupResult["data"].(map[string]interface{})
	user := data["user"].(*modelsv1.User)
	refreshToken := data["refresh_token"].(string)
	
	// Verify user in database
	dbUser, err := repositories.GetUserByEmail(s.ctx, signupData.Email)
	require.NoError(s.T(), err, "Should retrieve individual user from database")
	assert.Equal(s.T(), user.ID, dbUser.ID, "Database user should match")
	
	s.T().Logf("Individual user created: %s", user.ID)
	
	// Step 2: Login
	loginData := modelsv1.LoginParams{
		Email:    signupData.Email,
		Password: signupData.Password,
	}
	
	loginResponse := s.performRequest("POST", "/api/v1/auth/login", loginData)
	assert.Equal(s.T(), http.StatusOK, loginResponse.Code, "Individual login should succeed")
	
	// Step 3: Refresh Token
	refreshData := modelsv1.RefreshTokenParams{RefreshToken: refreshToken}
	refreshResponse := s.performRequest("POST", "/api/v1/auth/refresh", refreshData)
	assert.Equal(s.T(), http.StatusOK, refreshResponse.Code, "Individual refresh should succeed")
	
	// Step 4: Logout
	logoutData := modelsv1.RefreshTokenParams{RefreshToken: refreshToken}
	logoutResponse := s.performRequest("POST", "/api/v1/auth/logout", logoutData)
	assert.Equal(s.T(), http.StatusOK, logoutResponse.Code, "Individual logout should succeed")
	
	s.T().Log("Individual user authentication flow completed successfully")
}

func (s *ComprehensiveIntegrationTestSuite) testTeamUserAuthFlow() {
	s.T().Log("--- Testing Team User Authentication Flow ---")
	
	// Step 1: Team User Signup
	signupData := modelsv1.SignupParams{
		Name:     "Team Owner User",
		Email:    "teamowner@comprehensive.test",
		Password: "teampassword123",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Comprehensive Test Team",
	}
	
	signupResponse := s.performRequest("POST", "/api/v1/auth/signup", signupData)
	assert.Equal(s.T(), http.StatusOK, signupResponse.Code, "Team signup should succeed")
	
	var signupResult map[string]interface{}
	err := json.Unmarshal(signupResponse.Body.Bytes(), &signupResult)
	require.NoError(s.T(), err, "Should parse team signup response")
	
	data := signupResult["data"].(map[string]interface{})
	user := data["user"].(*modelsv1.User)
	
	// Verify team user in database
	dbUser, err := repositories.GetUserByEmail(s.ctx, signupData.Email)
	require.NoError(s.T(), err, "Should retrieve team user from database")
	assert.Equal(s.T(), string(modelsv1.UserTypeTeam), dbUser.UserType, "Should be team user")
	
	s.T().Logf("Team user created: %s", user.ID)
	
	// Complete authentication flow for team user
	loginData := modelsv1.LoginParams{
		Email:    signupData.Email,
		Password: signupData.Password,
	}
	
	loginResponse := s.performRequest("POST", "/api/v1/auth/login", loginData)
	assert.Equal(s.T(), http.StatusOK, loginResponse.Code, "Team login should succeed")
	
	s.T().Log("Team user authentication flow completed successfully")
}

func (s *ComprehensiveIntegrationTestSuite) testAuthenticationErrorScenarios() {
	s.T().Log("--- Testing Authentication Error Scenarios ---")
	
	// Test invalid login
	invalidLogin := modelsv1.LoginParams{
		Email:    "nonexistent@test.com",
		Password: "wrongpassword",
	}
	
	loginResponse := s.performRequest("POST", "/api/v1/auth/login", invalidLogin)
	assert.Equal(s.T(), http.StatusUnauthorized, loginResponse.Code, "Invalid login should fail")
	
	// Test invalid refresh token
	invalidRefresh := modelsv1.RefreshTokenParams{RefreshToken: "invalid-token"}
	refreshResponse := s.performRequest("POST", "/api/v1/auth/refresh", invalidRefresh)
	assert.Equal(s.T(), http.StatusUnauthorized, refreshResponse.Code, "Invalid refresh should fail")
	
	s.T().Log("Authentication error scenarios tested successfully")
}

// TEST 2: Log Retrieval with Real File System Operations (Requirement 1.3, 2.1)
func (s *ComprehensiveIntegrationTestSuite) TestLogRetrievalWithFileSystemOperations() {
	s.T().Log("=== Testing Log Retrieval with File System Operations ===")
	
	// Test log retrieval by date
	s.testLogRetrievalByDate()
	
	// Test log retrieval by ID
	s.testLogRetrievalByID()
	
	// Test log file system error handling
	s.testLogFileSystemErrorHandling()
	
	s.T().Log("=== Log Retrieval with File System Operations Test Completed ===")
}

func (s *ComprehensiveIntegrationTestSuite) testLogRetrievalByDate() {
	s.T().Log("--- Testing Log Retrieval by Date ---")
	
	today := time.Now().Format("2006-01-02")
	
	// Test retrieving logs for today
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", today), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Log retrieval by date should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse log response")
	
	assert.Equal(s.T(), "success", result["status"])
	data := result["data"].(map[string]interface{})
	logs := data["logs"].([]interface{})
	total := int(data["total"].(float64))
	
	assert.Greater(s.T(), len(logs), 0, "Should retrieve logs from file system")
	assert.Equal(s.T(), len(logs), total, "Total should match log count")
	
	// Verify log structure
	if len(logs) > 0 {
		firstLog := logs[0].(map[string]interface{})
		assert.Contains(s.T(), firstLog, "timestamp")
		assert.Contains(s.T(), firstLog, "log_id")
		assert.Contains(s.T(), firstLog, "message")
	}
	
	s.T().Logf("Retrieved %d logs for date %s from file system", len(logs), today)
}

func (s *ComprehensiveIntegrationTestSuite) testLogRetrievalByID() {
	s.T().Log("--- Testing Log Retrieval by ID ---")
	
	today := time.Now().Format("2006-01-02")
	testLogID := fmt.Sprintf("comprehensive-log-%s-001", today)
	
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/id/%s", testLogID), nil)
	assert.Equal(s.T(), http.StatusOK, response.Code, "Log retrieval by ID should succeed")
	
	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(s.T(), err, "Should parse log by ID response")
	
	logs := result["data"].([]interface{})
	assert.GreaterOrEqual(s.T(), len(logs), 1, "Should retrieve log by ID from file system")
	
	// Verify correct log ID
	if len(logs) > 0 {
		retrievedLog := logs[0].(map[string]interface{})
		assert.Equal(s.T(), testLogID, retrievedLog["log_id"], "Should retrieve correct log ID")
	}
	
	s.T().Logf("Retrieved log by ID %s from file system", testLogID)
}

func (s *ComprehensiveIntegrationTestSuite) testLogFileSystemErrorHandling() {
	s.T().Log("--- Testing Log File System Error Handling ---")
	
	// Test non-existent date
	nonExistentDate := "2020-01-01"
	response := s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", nonExistentDate), nil)
	assert.Equal(s.T(), http.StatusInternalServerError, response.Code, "Should handle missing log file")
	
	// Test invalid date format
	invalidDate := "invalid-date"
	response = s.performRequest("GET", fmt.Sprintf("/api/v1/logs/date/%s", invalidDate), nil)
	assert.Equal(s.T(), http.StatusBadRequest, response.Code, "Should handle invalid date format")
	
	s.T().Log("Log file system error handling tested successfully")
}

// TEST 3: Database Operations with Actual Database Connections (Requirement 2.1)
func (s *ComprehensiveIntegrationTestSuite) TestDatabaseOperationsWithActualConnections() {
	s.T().Log("=== Testing Database Operations with Actual Connections ===")
	
	// Test database CRUD operations
	s.testDatabaseCRUDOperations()
	
	// Test database transaction handling
	s.testDatabaseTransactionHandling()
	
	// Test database connection pooling
	s.testDatabaseConnectionPooling()
	
	s.T().Log("=== Database Operations with Actual Connections Test Completed ===")
}

func (s *ComprehensiveIntegrationTestSuite) testDatabaseCRUDOperations() {
	s.T().Log("--- Testing Database CRUD Operations ---")
	
	// Create user
	user := modelsv1.User{
		Name:     "CRUD Test User",
		Email:    "crud@comprehensive.test",
		Password: "crudpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Should start database transaction")
	
	createdUser, err := repositories.CreateUser(s.ctx, tx, user)
	require.NoError(s.T(), err, "Should create user in database")
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Should commit database transaction")
	
	// Read user
	retrievedUser, err := repositories.GetUserByEmail(s.ctx, user.Email)
	require.NoError(s.T(), err, "Should retrieve user from database")
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID, "Retrieved user should match")
	
	s.T().Logf("Database CRUD operations completed for user: %s", createdUser.ID)
}

func (s *ComprehensiveIntegrationTestSuite) testDatabaseTransactionHandling() {
	s.T().Log("--- Testing Database Transaction Handling ---")
	
	// Test successful transaction
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Should start transaction")
	
	user := modelsv1.User{
		Name:     "Transaction Test User",
		Email:    "transaction@comprehensive.test",
		Password: "transactionpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	_, err = repositories.CreateUser(s.ctx, tx, user)
	require.NoError(s.T(), err, "Should create user in transaction")
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Should commit transaction")
	
	// Verify user exists
	_, err = repositories.GetUserByEmail(s.ctx, user.Email)
	require.NoError(s.T(), err, "User should exist after commit")
	
	// Test rollback transaction
	tx2, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Should start second transaction")
	
	rollbackUser := modelsv1.User{
		Name:     "Rollback Test User",
		Email:    "rollback@comprehensive.test",
		Password: "rollbackpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	_, err = repositories.CreateUser(s.ctx, tx2, rollbackUser)
	require.NoError(s.T(), err, "Should create user in rollback transaction")
	
	database.RollbackTx(tx2)
	
	// Verify user doesn't exist
	_, err = repositories.GetUserByEmail(s.ctx, rollbackUser.Email)
	assert.Error(s.T(), err, "User should not exist after rollback")
	
	s.T().Log("Database transaction handling tested successfully")
}

func (s *ComprehensiveIntegrationTestSuite) testDatabaseConnectionPooling() {
	s.T().Log("--- Testing Database Connection Pooling ---")
	
	const numConcurrentOps = 5
	results := make(chan error, numConcurrentOps)
	
	// Test concurrent database operations
	for i := 0; i < numConcurrentOps; i++ {
		go func(index int) {
			user := modelsv1.User{
				Name:     fmt.Sprintf("Pool Test User %d", index),
				Email:    fmt.Sprintf("pool%d@comprehensive.test", index),
				Password: "poolpassword123",
				UserType: modelsv1.UserTypeIndividual,
			}
			
			tx, err := database.GetTx(s.ctx, nil)
			if err != nil {
				results <- err
				return
			}
			
			_, err = repositories.CreateUser(s.ctx, tx, user)
			if err != nil {
				database.RollbackTx(tx)
				results <- err
				return
			}
			
			err = database.CommitTx(tx)
			results <- err
		}(i)
	}
	
	// Collect results
	successCount := 0
	for i := 0; i < numConcurrentOps; i++ {
		if <-results == nil {
			successCount++
		}
	}
	
	assert.Equal(s.T(), numConcurrentOps, successCount, "All concurrent operations should succeed")
	s.T().Logf("Database connection pooling tested with %d concurrent operations", numConcurrentOps)
}

// Helper method to perform HTTP requests
func (s *ComprehensiveIntegrationTestSuite) performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	req := testutils.CreateTestRequest(method, path, body)
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, req)
	return recorder
}

// TestComprehensiveIntegrationSuite runs the comprehensive integration test suite
func TestComprehensiveIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ComprehensiveIntegrationTestSuite))
}