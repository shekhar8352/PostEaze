package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// BaseTestSuiteTest tests the BaseTestSuite functionality
type BaseTestSuiteTest struct {
	BaseTestSuite
}

func (s *BaseTestSuiteTest) TestBaseSetup() {
	// Test that basic setup is working
	s.NotNil(s.Router, "Router should be initialized")
	s.NotNil(s.Fixtures, "Fixtures map should be initialized")
	s.NotNil(s.ctx, "Context should be initialized")
}

func (s *BaseTestSuiteTest) TestFixtureLoading() {
	// Test fixture loading
	s.LoadTestFixtures()
	
	// Check that fixtures are loaded
	s.Contains(s.Fixtures, "users", "Users fixtures should be loaded")
	s.Contains(s.Fixtures, "teams", "Teams fixtures should be loaded")
	s.Contains(s.Fixtures, "tokens", "Tokens fixtures should be loaded")
}

func (s *BaseTestSuiteTest) TestGetTestData() {
	// Test getting test data
	user := s.GetTestUser(0)
	s.NotEmpty(user.ID, "Test user should have ID")
	s.NotEmpty(user.Email, "Test user should have email")
	
	team := s.GetTestTeam(0)
	s.NotEmpty(team.ID, "Test team should have ID")
	s.NotEmpty(team.Name, "Test team should have name")
	
	token := s.GetTestToken(0)
	s.NotEmpty(token.Token, "Test token should have token value")
	s.NotEmpty(token.UserID, "Test token should have user ID")
}

func (s *BaseTestSuiteTest) TestCreateTestData() {
	// Test creating custom test data
	user := s.CreateTestUser(func(u *UserFixture) {
		u.Email = "custom@test.com"
		u.UserType = "admin"
	})
	
	s.Equal("custom@test.com", user.Email, "Custom user email should be set")
	s.Equal("admin", user.UserType, "Custom user type should be set")
	
	team := s.CreateTestTeam("user123", func(t *TeamFixture) {
		t.Name = "Custom Team"
	})
	
	s.Equal("Custom Team", team.Name, "Custom team name should be set")
	s.Equal("user123", team.OwnerID, "Team owner ID should be set")
	
	token := s.CreateTestToken("user456", func(t *RefreshTokenFixture) {
		t.Token = "custom-token"
	})
	
	s.Equal("custom-token", token.Token, "Custom token should be set")
	s.Equal("user456", token.UserID, "Token user ID should be set")
}

// APITestSuiteTest tests the APITestSuite functionality
type APITestSuiteTest struct {
	APITestSuite
}

func (s *APITestSuiteTest) TestAPISetup() {
	// Test that API-specific setup is working
	s.NotNil(s.Router, "Router should be initialized")
	
	// Test that common routes are set up (we can check middleware)
	handlers := s.Router.Handlers
	s.NotEmpty(handlers, "Router should have handlers/middleware")
}

func (s *APITestSuiteTest) TestCreateAuthenticatedRequest() {
	// Test creating authenticated requests
	ctx, recorder, err := s.CreateAuthenticatedRequest("GET", "/test", nil, "user123", "admin")
	
	s.NoError(err, "Should create authenticated request without error")
	s.NotNil(ctx, "Context should be created")
	s.NotNil(recorder, "Recorder should be created")
	
	// Check that authorization header is set
	authHeader := ctx.Request.Header.Get("Authorization")
	s.NotEmpty(authHeader, "Authorization header should be set")
	s.Contains(authHeader, "Bearer ", "Authorization header should contain Bearer token")
}

func (s *APITestSuiteTest) TestAPIResponseAssertions() {
	// Test API response assertion helpers
	ctx, recorder := NewTestGinContext("GET", "/test", nil)
	
	// Simulate a successful response
	ctx.JSON(200, map[string]interface{}{
		"message": "success",
		"data":    "test data",
	})
	
	// Test success assertion
	s.AssertSuccessAPIResponse(recorder, nil)
	
	// Test general API response assertion
	s.AssertAPIResponse(recorder, 200, "success")
}

func (s *APITestSuiteTest) TestErrorResponseAssertions() {
	// Test error response assertions
	ctx, recorder := NewTestGinContext("GET", "/test", nil)
	
	// Simulate an error response
	ctx.JSON(400, map[string]interface{}{
		"error": "validation failed",
	})
	
	// Test error assertion
	s.AssertErrorAPIResponse(recorder, 400, "validation failed")
}

// BusinessLogicTestSuiteTest tests the BusinessLogicTestSuite functionality
type BusinessLogicTestSuiteTest struct {
	BusinessLogicTestSuite
}

func (s *BusinessLogicTestSuiteTest) TestBusinessLogicSetup() {
	// Test that business logic setup works - fixtures are loaded on demand
	s.NotNil(s.Fixtures, "Fixtures map should be initialized")
	
	// Load fixtures manually to test
	s.LoadTestFixtures()
	s.Contains(s.Fixtures, "users", "Users fixtures should be loaded")
	s.Contains(s.Fixtures, "teams", "Teams fixtures should be loaded")
	s.Contains(s.Fixtures, "tokens", "Tokens fixtures should be loaded")
}

func (s *BusinessLogicTestSuiteTest) TestCreateMockContext() {
	// Test mock context creation
	ctx := s.CreateMockContext()
	s.NotNil(ctx, "Mock context should be created")
}

// ModelTestSuiteTest tests the ModelTestSuite functionality
type ModelTestSuiteTest struct {
	ModelTestSuite
}

func (s *ModelTestSuiteTest) TestModelSetup() {
	// Test that model suite setup works
	s.NotNil(s.Router, "Router should be initialized")
	s.NotNil(s.Fixtures, "Fixtures map should be initialized")
}

func (s *ModelTestSuiteTest) TestValidationAssertions() {
	// Create a custom error that contains the field name
	customError := &ValidationError{Field: "test_field", Message: "validation failed"}
	s.AssertValidationError(customError, "test_field")
	
	// Test no validation error
	s.AssertNoValidationError(nil)
}

// ValidationError is a test helper for validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// TestSuiteHelperTest tests the TestSuiteHelper functionality
type TestSuiteHelperTest struct {
	suite.Suite
	helper *TestSuiteHelper
}

func (s *TestSuiteHelperTest) SetupTest() {
	s.helper = NewTestSuiteHelper()
}

func (s *TestSuiteHelperTest) TestRunWithTimeout() {
	// Test function that completes within timeout
	s.helper.RunWithTimeout(s.T(), 100*time.Millisecond, func() {
		time.Sleep(10 * time.Millisecond)
	})
	
	// This should pass without timing out
}

func (s *TestSuiteHelperTest) TestAssertEventuallyTrue() {
	// Test condition that becomes true
	counter := 0
	condition := func() bool {
		counter++
		return counter >= 3
	}
	
	s.helper.AssertEventuallyTrue(s.T(), condition, 100*time.Millisecond, "Counter should reach 3")
}

// Test suite runner functions
func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuiteTest))
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuiteTest))
}

func TestBusinessLogicTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessLogicTestSuiteTest))
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuiteTest))
}

func TestSuiteHelperFunctionality(t *testing.T) {
	suite.Run(t, new(TestSuiteHelperTest))
}

// TestRunTestSuite tests the RunTestSuite helper function
func TestRunTestSuite(t *testing.T) {
	// Test that the helper function works
	testSuite := new(BaseTestSuiteTest)
	
	// This should not panic
	RunTestSuite(t, testSuite)
}