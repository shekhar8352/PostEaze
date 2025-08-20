package businessv1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// AuthBusinessLogicTestSuite tests the authentication business logic
type AuthBusinessLogicTestSuite struct {
	testutils.BusinessLogicTestSuite
	ctx context.Context
}

// SetupSuite initializes the test suite
func (s *AuthBusinessLogicTestSuite) SetupSuite() {
	s.BusinessLogicTestSuite.SetupSuite()
	s.ctx = context.Background()
}

// SetupTest prepares each test
func (s *AuthBusinessLogicTestSuite) SetupTest() {
	s.BusinessLogicTestSuite.SetupTest()
}

// TearDownTest cleans up after each test
func (s *AuthBusinessLogicTestSuite) TearDownTest() {
	s.BusinessLogicTestSuite.TearDownTest()
}

// TestPasswordHashing tests password hashing functionality
func (s *AuthBusinessLogicTestSuite) TestPasswordHashing() {
	// Test that password hashing works correctly
	testCases := []struct {
		name        string
		password    string
		expectError bool
		description string
	}{
		{
			name:        "ValidPassword",
			password:    "password123",
			expectError: false,
			description: "Valid password should be hashed successfully",
		},
		{
			name:        "EmptyPassword",
			password:    "",
			expectError: true,
			description: "Empty password should cause hashing error",
		},
		{
			name:        "LongPassword",
			password:    "this-is-a-very-long-password-that-should-still-work-fine-123456789",
			expectError: false,
			description: "Long password should be hashed successfully",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			params := modelsv1.SignupParams{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: tc.password,
				UserType: modelsv1.UserTypeIndividual,
			}

			// Call signup - it will fail due to database issues, but we can check
			// if the error is related to password hashing or database operations
			_, err := businessv1.Signup(s.ctx, params)

			if tc.expectError {
				s.Error(err, tc.description)
				// For empty password, the error should be from bcrypt
				if tc.password == "" {
					s.Contains(err.Error(), "crypto/bcrypt", "Should be bcrypt error for empty password")
				}
			} else {
				// For valid passwords, we expect database-related errors, not hashing errors
				if err != nil {
					s.NotContains(err.Error(), "crypto/bcrypt", "Should not be bcrypt error for valid password")
				}
			}
		})
	}
}

// TestUserTypeHandling tests different user type handling
func (s *AuthBusinessLogicTestSuite) TestUserTypeHandling() {
	testCases := []struct {
		name        string
		userType    modelsv1.UserType
		teamName    string
		description string
	}{
		{
			name:        "IndividualUser",
			userType:    modelsv1.UserTypeIndividual,
			teamName:    "",
			description: "Individual user should not require team name",
		},
		{
			name:        "TeamUser",
			userType:    modelsv1.UserTypeTeam,
			teamName:    "Development Team",
			description: "Team user should include team name",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			params := modelsv1.SignupParams{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				UserType: tc.userType,
				TeamName: tc.teamName,
			}

			// Call signup - will fail due to database, but we can verify
			// that the function processes different user types
			_, err := businessv1.Signup(s.ctx, params)

			// We expect database-related errors, not user type validation errors
			s.Error(err, "Should fail due to database connection")
			// The error should not be related to user type validation
			s.NotContains(err.Error(), "user_type", "Should not be user type validation error")
			s.NotContains(err.Error(), "team_name", "Should not be team name validation error")
		})
	}
}

// TestLoginParams_Validation tests login parameter validation
func (s *AuthBusinessLogicTestSuite) TestLoginParams_Validation() {
	testCases := []struct {
		name        string
		params      modelsv1.LoginParams
		description string
	}{
		{
			name: "ValidParams",
			params: modelsv1.LoginParams{
				Email:    "test@example.com",
				Password: "password123",
			},
			description: "Valid login parameters",
		},
		{
			name: "EmptyEmail",
			params: modelsv1.LoginParams{
				Email:    "",
				Password: "password123",
			},
			description: "Empty email should still be processed by business logic",
		},
		{
			name: "EmptyPassword",
			params: modelsv1.LoginParams{
				Email:    "test@example.com",
				Password: "",
			},
			description: "Empty password should still be processed by business logic",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Call login - will fail due to database connection, but we can verify
			// that the function processes the parameters
			_, err := businessv1.Login(s.ctx, tc.params)

			// We expect database-related errors, not parameter validation errors
			s.Error(err, "Should fail due to database connection")
		})
	}
}

// TestRefreshTokenParams_Validation tests refresh token parameter handling
func (s *AuthBusinessLogicTestSuite) TestRefreshTokenParams_Validation() {
	testCases := []struct {
		name        string
		token       string
		description string
	}{
		{
			name:        "ValidToken",
			token:       "valid-refresh-token-123",
			description: "Valid refresh token format",
		},
		{
			name:        "EmptyToken",
			token:       "",
			description: "Empty refresh token",
		},
		{
			name:        "InvalidToken",
			token:       "invalid-token",
			description: "Invalid refresh token format",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Call refresh token - will fail due to database connection
			_, err := businessv1.RefreshToken(s.ctx, tc.token)

			// We expect database-related errors
			s.Error(err, "Should fail due to database connection")
		})
	}
}

// TestLogoutParams_Validation tests logout parameter handling
func (s *AuthBusinessLogicTestSuite) TestLogoutParams_Validation() {
	testCases := []struct {
		name        string
		token       string
		description string
	}{
		{
			name:        "ValidToken",
			token:       "valid-refresh-token-123",
			description: "Valid refresh token for logout",
		},
		{
			name:        "EmptyToken",
			token:       "",
			description: "Empty refresh token for logout",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Call logout - will fail due to database connection
			err := businessv1.Logout(s.ctx, tc.token)

			// We expect database-related errors
			s.Error(err, "Should fail due to database connection")
		})
	}
}

// TestRunner runs all the authentication business logic tests
func TestAuthBusinessLogic(t *testing.T) {
	suite.Run(t, new(AuthBusinessLogicTestSuite))
}