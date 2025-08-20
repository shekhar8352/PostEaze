package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// ExampleBaseTestSuite demonstrates basic usage of BaseTestSuite
type ExampleBaseTestSuite struct {
	BaseTestSuite
}

func (s *ExampleBaseTestSuite) TestBasicFunctionality() {
	// Test that basic setup works
	s.T().Log("Testing basic test suite functionality")
	
	// Router should be available
	s.NotNil(s.Router, "Router should be initialized")
	
	// Fixtures map should be available
	s.NotNil(s.Fixtures, "Fixtures should be initialized")
	
	// Load test fixtures
	s.LoadTestFixtures()
	
	// Get predefined test data
	user := s.GetTestUser(0)
	s.NotEmpty(user.ID, "Test user should have ID")
	s.NotEmpty(user.Email, "Test user should have email")
	
	s.T().Log("Basic functionality test completed successfully")
}

func (s *ExampleBaseTestSuite) TestCustomFixtures() {
	// Create custom test data
	customUser := s.CreateTestUser(func(u *UserFixture) {
		u.Email = "example@test.com"
		u.UserType = "admin"
	})
	
	s.Equal("example@test.com", customUser.Email)
	s.Equal("admin", customUser.UserType)
	
	// Create custom team
	customTeam := s.CreateTestTeam("user123", func(t *TeamFixture) {
		t.Name = "Example Team"
	})
	
	s.Equal("Example Team", customTeam.Name)
	s.Equal("user123", customTeam.OwnerID)
	
	s.T().Log("Custom fixtures test completed successfully")
}

// ExampleAPITestSuite demonstrates API testing functionality
type ExampleAPITestSuite struct {
	APITestSuite
}

func (s *ExampleAPITestSuite) TestAPIUtilities() {
	s.T().Log("Testing API utilities")
	
	// Test creating authenticated request
	ctx, recorder, err := s.CreateAuthenticatedRequest(
		"GET", "/api/test", nil, "user123", "admin",
	)
	
	s.NoError(err, "Should create authenticated request without error")
	s.NotNil(ctx, "Context should be created")
	s.NotNil(recorder, "Recorder should be created")
	
	// Check authorization header
	authHeader := ctx.Request.Header.Get("Authorization")
	s.NotEmpty(authHeader, "Authorization header should be set")
	s.Contains(authHeader, "Bearer ", "Should contain Bearer token")
	
	s.T().Log("API utilities test completed successfully")
}

func (s *ExampleAPITestSuite) TestResponseAssertions() {
	s.T().Log("Testing response assertions")
	
	// Create a test context and simulate a response
	ctx, recorder := NewTestGinContext("GET", "/test", nil)
	
	// Simulate successful response
	ctx.JSON(200, map[string]interface{}{
		"message": "success",
		"data":    "test data",
	})
	
	// Test response assertions
	s.AssertSuccessAPIResponse(recorder, nil)
	s.AssertAPIResponse(recorder, 200, "success")
	
	s.T().Log("Response assertions test completed successfully")
}

// ExampleBusinessLogicTestSuite demonstrates business logic testing
type ExampleBusinessLogicTestSuite struct {
	BusinessLogicTestSuite
}

func (s *ExampleBusinessLogicTestSuite) TestBusinessLogicUtilities() {
	s.T().Log("Testing business logic utilities")
	
	// Test mock context creation
	ctx := s.CreateMockContext()
	s.NotNil(ctx, "Mock context should be created")
	
	// Note: BusinessLogicTestSuite loads fixtures in SetupSuite, but they may not be
	// loaded into the Fixtures map unless LoadTestFixtures is called explicitly
	// Let's verify the fixtures are available by loading them
	s.LoadTestFixtures()
	s.Contains(s.Fixtures, "users", "Users fixtures should be loaded")
	s.Contains(s.Fixtures, "teams", "Teams fixtures should be loaded")
	s.Contains(s.Fixtures, "tokens", "Tokens fixtures should be loaded")
	
	s.T().Log("Business logic utilities test completed successfully")
}

// ExampleModelTestSuite demonstrates model testing functionality
type ExampleModelTestSuite struct {
	ModelTestSuite
}

func (s *ExampleModelTestSuite) TestModelUtilities() {
	s.T().Log("Testing model utilities")
	
	// Test validation assertions with custom error
	validationError := &ValidationError{
		Field:   "email",
		Message: "invalid email format",
	}
	
	s.AssertValidationError(validationError, "email")
	s.AssertNoValidationError(nil)
	
	s.T().Log("Model utilities test completed successfully")
}

// ExampleHelperTest demonstrates TestSuiteHelper functionality
type ExampleHelperTest struct {
	suite.Suite
	helper *TestSuiteHelper
}

func (s *ExampleHelperTest) SetupTest() {
	s.helper = NewTestSuiteHelper()
}

func (s *ExampleHelperTest) TestHelperUtilities() {
	s.T().Log("Testing helper utilities")
	
	// Test timeout functionality
	s.helper.RunWithTimeout(s.T(), 100*time.Millisecond, func() {
		// Quick operation that should complete within timeout
		s.T().Log("Operation completed within timeout")
	})
	
	// Test eventually true assertion
	counter := 0
	s.helper.AssertEventuallyTrue(s.T(), func() bool {
		counter++
		return counter >= 2
	}, 100*time.Millisecond, "Counter should reach 2")
	
	s.T().Log("Helper utilities test completed successfully")
}

// Test runner functions
func TestExampleBaseTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleBaseTestSuite))
}

func TestExampleAPITestSuite(t *testing.T) {
	suite.Run(t, new(ExampleAPITestSuite))
}

func TestExampleBusinessLogicTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleBusinessLogicTestSuite))
}

func TestExampleModelTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleModelTestSuite))
}

func TestExampleHelper(t *testing.T) {
	suite.Run(t, new(ExampleHelperTest))
}

// TestRunTestSuiteHelper demonstrates the RunTestSuite helper function
func TestRunTestSuiteHelper(t *testing.T) {
	t.Log("Testing RunTestSuite helper function")
	
	// This demonstrates how to use the helper function
	testSuite := new(ExampleBaseTestSuite)
	RunTestSuite(t, testSuite)
	
	t.Log("RunTestSuite helper test completed successfully")
}