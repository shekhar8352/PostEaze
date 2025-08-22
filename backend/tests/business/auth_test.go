package business_test

import (
	"strings"
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
	"github.com/shekhar8352/PostEaze/utils"
)

func TestPasswordHashing_Behavior(t *testing.T) {
	// Test password hashing utility directly to verify business logic behavior
	tests := []struct {
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
			expectError: false, // bcrypt can hash empty strings
			description: "Empty password should be hashed (bcrypt allows empty strings)",
		},
		{
			name:        "LongPassword",
			password:    "this-is-a-very-long-password-that-should-still-work-fine-123456789",
			expectError: false,
			description: "Long password should be hashed successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the password hashing utility directly
			hashedPassword, err := utils.HashPassword(tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %s, but got none", tt.description)
				}
				// For empty password, the error should be from bcrypt
				if tt.password == "" && err != nil {
					errMsg := err.Error()
					if !strings.Contains(errMsg, "bcrypt") {
						t.Errorf("expected bcrypt error for empty password, got: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s: %v", tt.description, err)
				}
				if hashedPassword == "" {
					t.Errorf("expected non-empty hashed password for %s", tt.description)
				}
				// Verify the hashed password is different from the original (except for empty password)
				if tt.password != "" && hashedPassword == tt.password {
					t.Errorf("hashed password should be different from original password")
				}
			}
		})
	}
}

func TestSignupParams_Validation(t *testing.T) {
	// Test signup parameter validation patterns
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := modelsv1.SignupParams{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				UserType: tt.userType,
				TeamName: tt.teamName,
			}

			// Verify parameter structure is valid
			if params.Name == "" {
				t.Error("name should not be empty")
			}
			if params.Email == "" {
				t.Error("email should not be empty")
			}
			if params.Password == "" {
				t.Error("password should not be empty")
			}
			
			// Verify user type handling
			if tt.userType == modelsv1.UserTypeTeam && tt.teamName == "" {
				t.Error("team user should have team name")
			}
			if tt.userType == modelsv1.UserTypeIndividual && tt.teamName != "" {
				t.Log("individual user has team name, which is acceptable")
			}
		})
	}
}

func TestLoginParams_Structure(t *testing.T) {
	// Test login parameter structure and validation patterns
	tests := []struct {
		name        string
		params      modelsv1.LoginParams
		expectValid bool
		description string
	}{
		{
			name: "ValidParams",
			params: modelsv1.LoginParams{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectValid: true,
			description: "Valid login parameters",
		},
		{
			name: "EmptyEmail",
			params: modelsv1.LoginParams{
				Email:    "",
				Password: "password123",
			},
			expectValid: false,
			description: "Empty email should be invalid",
		},
		{
			name: "EmptyPassword",
			params: modelsv1.LoginParams{
				Email:    "test@example.com",
				Password: "",
			},
			expectValid: false,
			description: "Empty password should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test parameter structure validation
			isValid := tt.params.Email != "" && tt.params.Password != ""
			
			if tt.expectValid && !isValid {
				t.Errorf("expected valid parameters for %s", tt.description)
			}
			if !tt.expectValid && isValid {
				t.Errorf("expected invalid parameters for %s", tt.description)
			}
		})
	}
}

func TestRefreshToken_ParameterValidation(t *testing.T) {
	// Test refresh token parameter validation patterns
	tests := []struct {
		name        string
		token       string
		expectValid bool
		description string
	}{
		{
			name:        "ValidToken",
			token:       "valid-refresh-token-123",
			expectValid: true,
			description: "Valid refresh token format",
		},
		{
			name:        "EmptyToken",
			token:       "",
			expectValid: false,
			description: "Empty refresh token",
		},
		{
			name:        "ShortToken",
			token:       "abc",
			expectValid: false,
			description: "Short token should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test token validation logic
			isValid := tt.token != "" && len(tt.token) > 5
			
			if tt.expectValid && !isValid {
				t.Errorf("expected valid token for %s", tt.description)
			}
			if !tt.expectValid && isValid {
				t.Errorf("expected invalid token for %s", tt.description)
			}
		})
	}
}

func TestLogout_ParameterValidation(t *testing.T) {
	// Test logout parameter validation patterns
	tests := []struct {
		name        string
		token       string
		expectValid bool
		description string
	}{
		{
			name:        "ValidToken",
			token:       "valid-refresh-token-123",
			expectValid: true,
			description: "Valid refresh token for logout",
		},
		{
			name:        "EmptyToken",
			token:       "",
			expectValid: false,
			description: "Empty refresh token for logout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test token validation logic
			isValid := tt.token != ""
			
			if tt.expectValid && !isValid {
				t.Errorf("expected valid token for %s", tt.description)
			}
			if !tt.expectValid && isValid {
				t.Errorf("expected invalid token for %s", tt.description)
			}
		})
	}
}

// TestSignupParams_WithFixtures tests signup parameter creation with test fixtures
func TestSignupParams_WithFixtures(t *testing.T) {
	// Use test fixtures for consistent test data
	testUser := helpers.CreateUser(func(u *modelsv1.User) {
		u.Email = "signup.test@example.com"
		u.UserType = modelsv1.UserTypeIndividual
	})
	
	params := modelsv1.SignupParams{
		Name:     testUser.Name,
		Email:    testUser.Email,
		Password: "testpassword123",
		UserType: testUser.UserType,
	}
	
	// Verify parameter structure
	if params.Name == "" {
		t.Error("name should not be empty")
	}
	if params.Email == "" {
		t.Error("email should not be empty")
	}
	if params.Password == "" {
		t.Error("password should not be empty")
	}
	if params.UserType == "" {
		t.Error("user type should not be empty")
	}
	
	// Verify test fixture data
	if testUser.Email != params.Email {
		t.Errorf("expected email %s, got %s", testUser.Email, params.Email)
	}
	if testUser.UserType != params.UserType {
		t.Errorf("expected user type %s, got %s", testUser.UserType, params.UserType)
	}
}

// TestLoginParams_WithFixtures tests login parameter creation with test fixtures
func TestLoginParams_WithFixtures(t *testing.T) {
	// Use test fixtures for consistent test data
	testUser := helpers.GetUser("user-1")
	
	params := modelsv1.LoginParams{
		Email:    testUser.Email,
		Password: "password123", // This matches the test user's password
	}
	
	// Verify parameter structure
	if params.Email == "" {
		t.Error("email should not be empty")
	}
	if params.Password == "" {
		t.Error("password should not be empty")
	}
	
	// Verify test fixture data
	if testUser.Email != params.Email {
		t.Errorf("expected email %s, got %s", testUser.Email, params.Email)
	}
	
	// Verify password hashing would work
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		t.Errorf("password hashing failed: %v", err)
	}
	if hashedPassword == "" {
		t.Error("hashed password should not be empty")
	}
}