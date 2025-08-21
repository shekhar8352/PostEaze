package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// BaseTestSuite provides common setup and teardown for all test suites
type BaseTestSuite struct {
	suite.Suite
	DB       *sql.DB
	Router   *gin.Engine
	Cleanup  func()
	Fixtures map[string]interface{}
	ctx      context.Context
}

// SetupSuite runs once before all tests in the suite
func (s *BaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Initialize fixtures map
	s.Fixtures = make(map[string]interface{})
	
	// Try to setup test database, but don't fail if it's not available
	db, cleanup, err := SetupTestDB(s.ctx)
	if err != nil {
		s.T().Logf("Warning: Failed to setup test database: %v", err)
		s.T().Log("Tests will run without database connection")
		s.DB = nil
		s.Cleanup = func() {}
	} else {
		s.DB = db
		s.Cleanup = cleanup
	}
	
	// Setup Gin router for HTTP testing
	s.Router = gin.New()
	
	s.T().Log("BaseTestSuite setup completed")
}

// TearDownSuite runs once after all tests in the suite
func (s *BaseTestSuite) TearDownSuite() {
	if s.Cleanup != nil {
		s.Cleanup()
	}
	s.T().Log("BaseTestSuite teardown completed")
}

// SetupTest runs before each individual test
func (s *BaseTestSuite) SetupTest() {
	// Clean up test data before each test (only if database is available)
	if s.DB != nil {
		err := CleanupTestData(s.ctx, s.DB)
		if err != nil {
			s.T().Logf("Warning: Failed to cleanup test data: %v", err)
		}
	}
	
	// Clear fixtures
	s.Fixtures = make(map[string]interface{})
}

// TearDownTest runs after each individual test
func (s *BaseTestSuite) TearDownTest() {
	// Additional cleanup if needed
	// This is called after each test method
}

// LoadTestFixtures loads predefined test fixtures
func (s *BaseTestSuite) LoadTestFixtures() {
	// Store fixtures in the suite for easy access (always available)
	s.Fixtures["users"] = TestUsers
	s.Fixtures["teams"] = TestTeams
	s.Fixtures["tokens"] = TestTokens
	
	// Load basic test data into database (only if database is available)
	if s.DB != nil {
		err := LoadFixtures(s.ctx, s.DB, TestUsers, TestTeams, TestTokens)
		if err != nil {
			s.T().Logf("Warning: Failed to load test fixtures into database: %v", err)
		}
	}
}

// LoadCustomFixtures loads custom fixtures
func (s *BaseTestSuite) LoadCustomFixtures(fixtures ...interface{}) {
	// Only load into database if available
	if s.DB != nil {
		err := LoadFixtures(s.ctx, s.DB, fixtures...)
		if err != nil {
			s.T().Logf("Warning: Failed to load custom fixtures into database: %v", err)
		}
	}
}

// GetTestUser returns a test user by index
func (s *BaseTestSuite) GetTestUser(index int) UserFixture {
	if index < 0 || index >= len(TestUsers) {
		s.T().Fatalf("Invalid user index: %d", index)
	}
	return TestUsers[index]
}

// GetTestTeam returns a test team by index
func (s *BaseTestSuite) GetTestTeam(index int) TeamFixture {
	if index < 0 || index >= len(TestTeams) {
		s.T().Fatalf("Invalid team index: %d", index)
	}
	return TestTeams[index]
}

// GetTestToken returns a test token by index
func (s *BaseTestSuite) GetTestToken(index int) RefreshTokenFixture {
	if index < 0 || index >= len(TestTokens) {
		s.T().Fatalf("Invalid token index: %d", index)
	}
	return TestTokens[index]
}

// CreateTestUser creates and loads a new test user
func (s *BaseTestSuite) CreateTestUser(overrides ...func(*UserFixture)) UserFixture {
	user := CreateUserFixture(overrides...)
	s.LoadCustomFixtures(user)
	return user
}

// CreateTestTeam creates and loads a new test team
func (s *BaseTestSuite) CreateTestTeam(ownerID string, overrides ...func(*TeamFixture)) TeamFixture {
	team := CreateTeamFixture(ownerID, overrides...)
	s.LoadCustomFixtures(team)
	return team
}

// CreateTestToken creates and loads a new test refresh token
func (s *BaseTestSuite) CreateTestToken(userID string, overrides ...func(*RefreshTokenFixture)) RefreshTokenFixture {
	token := CreateRefreshTokenFixture(userID, overrides...)
	s.LoadCustomFixtures(token)
	return token
}

// DatabaseTestSuite extends BaseTestSuite with database-specific utilities
type DatabaseTestSuite struct {
	BaseTestSuite
}

// SetupSuite extends the base setup with database package initialization
func (s *DatabaseTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Initialize the database package for integration testing
	db, cleanup, err := SetupTestDBWithDatabase(s.ctx)
	s.Require().NoError(err, "Failed to setup database package")
	
	// Replace the base cleanup with enhanced cleanup
	oldCleanup := s.Cleanup
	s.Cleanup = func() {
		cleanup()
		if oldCleanup != nil {
			oldCleanup()
		}
	}
	
	s.DB = db
}

// BeginTransaction starts a test transaction
func (s *DatabaseTestSuite) BeginTransaction() *sql.Tx {
	tx, err := BeginTestTransaction(s.ctx, s.DB)
	s.Require().NoError(err, "Failed to begin test transaction")
	return tx
}

// RollbackTransaction rolls back a test transaction
func (s *DatabaseTestSuite) RollbackTransaction(tx *sql.Tx) {
	RollbackTestTransaction(tx)
}

// RunTestSuite is a helper function to run a test suite
func RunTestSuite(t *testing.T, testSuite suite.TestingSuite) {
	suite.Run(t, testSuite)
}

// APITestSuite extends BaseTestSuite with API-specific utilities
type APITestSuite struct {
	BaseTestSuite
}

// SetupSuite extends the base setup with API-specific initialization
func (s *APITestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup common API routes for testing
	s.setupCommonRoutes()
}

// setupCommonRoutes sets up common API routes for testing
func (s *APITestSuite) setupCommonRoutes() {
	// Add common middleware and routes that are used across tests
	// This can be extended as needed
	s.Router.Use(gin.Recovery())
	s.Router.Use(gin.Logger())
}

// CreateAuthenticatedRequest creates an HTTP request with authentication
func (s *APITestSuite) CreateAuthenticatedRequest(method, url string, body interface{}, userID, userType string) (*gin.Context, *httptest.ResponseRecorder, error) {
	// Generate test JWT token
	token, err := GenerateTestJWT(userID, userType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate test JWT: %w", err)
	}
	
	// Create context with authentication
	ctx, recorder := NewTestGinContext(method, url, body)
	ctx.Request.Header.Set("Authorization", "Bearer "+token)
	
	return ctx, recorder, nil
}

// AssertAPIResponse provides common API response assertions
func (s *APITestSuite) AssertAPIResponse(recorder *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	s.Equal(expectedStatus, recorder.Code, "Status code should match")
	
	if expectedMessage != "" {
		var response map[string]interface{}
		err := ParseJSONResponse(recorder, &response)
		s.NoError(err, "Should be able to parse JSON response")
		
		if message, exists := response["message"]; exists {
			s.Equal(expectedMessage, message, "Response message should match")
		}
	}
}

// AssertSuccessAPIResponse asserts a successful API response
func (s *APITestSuite) AssertSuccessAPIResponse(recorder *httptest.ResponseRecorder, expectedData interface{}) {
	s.Equal(200, recorder.Code, "Should return 200 OK")
	
	if expectedData != nil {
		err := ParseJSONResponse(recorder, expectedData)
		s.NoError(err, "Should be able to parse JSON response")
	}
}

// AssertErrorAPIResponse asserts an error API response
func (s *APITestSuite) AssertErrorAPIResponse(recorder *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	s.Equal(expectedStatus, recorder.Code, "Status code should match")
	
	var response map[string]interface{}
	err := ParseJSONResponse(recorder, &response)
	s.NoError(err, "Should be able to parse JSON response")
	
	if expectedMessage != "" {
		if message, exists := response["message"]; exists {
			s.Contains(message.(string), expectedMessage, "Error message should contain expected text")
		} else if errorMsg, exists := response["error"]; exists {
			s.Contains(errorMsg.(string), expectedMessage, "Error message should contain expected text")
		}
	}
}

// BusinessLogicTestSuite extends BaseTestSuite for business logic testing
type BusinessLogicTestSuite struct {
	BaseTestSuite
}

// SetupSuite extends the base setup with business logic specific initialization
func (s *BusinessLogicTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Load common fixtures for business logic testing
	s.LoadTestFixtures()
}

// CreateMockContext creates a mock context for business logic testing
func (s *BusinessLogicTestSuite) CreateMockContext() context.Context {
	return context.Background()
}

// ModelTestSuite extends BaseTestSuite for model testing
type ModelTestSuite struct {
	BaseTestSuite
}

// SetupSuite extends the base setup with model-specific initialization
func (s *ModelTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	// Model tests typically don't need database or router setup
	// but we keep the base structure for consistency
}

// AssertValidationError asserts that a validation error occurred
func (s *ModelTestSuite) AssertValidationError(err error, fieldName string) {
	s.Error(err, "Should have validation error")
	s.Contains(err.Error(), fieldName, "Error should mention the field name")
}

// AssertNoValidationError asserts that no validation error occurred
func (s *ModelTestSuite) AssertNoValidationError(err error) {
	s.NoError(err, "Should not have validation error")
}

// TestSuiteHelper provides additional helper methods for all test suites
type TestSuiteHelper struct{}

// NewTestSuiteHelper creates a new test suite helper
func NewTestSuiteHelper() *TestSuiteHelper {
	return &TestSuiteHelper{}
}

// RunWithTimeout runs a function with a timeout
func (h *TestSuiteHelper) RunWithTimeout(t *testing.T, timeout time.Duration, fn func()) {
	done := make(chan bool, 1)
	
	go func() {
		fn()
		done <- true
	}()
	
	select {
	case <-done:
		// Function completed successfully
	case <-time.After(timeout):
		t.Fatalf("Test timed out after %v", timeout)
	}
}

// AssertEventuallyTrue asserts that a condition becomes true within a timeout
func (h *TestSuiteHelper) AssertEventuallyTrue(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	
	timeoutChan := time.After(timeout)
	
	for {
		select {
		case <-ticker.C:
			if condition() {
				return
			}
		case <-timeoutChan:
			t.Fatalf("Condition was not met within timeout: %s", message)
		}
	}
}