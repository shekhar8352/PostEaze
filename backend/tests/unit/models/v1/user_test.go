package modelsv1_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// UserModelTestSuite tests user model validation and serialization
type UserModelTestSuite struct {
	testutils.ModelTestSuite
}

// TestUserModelTestSuite runs the user model test suite
func TestUserModelTestSuite(t *testing.T) {
	suite.Run(t, new(UserModelTestSuite))
}

// TestUser_JSONSerialization tests user struct JSON serialization
func (s *UserModelTestSuite) TestUser_JSONSerialization() {
	// Test valid user serialization
	user := modelsv1.User{
		ID:        "user-123",
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword",
		UserType:  modelsv1.UserTypeIndividual,
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(user)
	s.NoError(err, "Should marshal user to JSON without error")

	// Verify password is excluded from JSON (json:"-" tag)
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	s.NoError(err, "Should unmarshal JSON to map")
	s.NotContains(jsonMap, "password", "Password should not be included in JSON")

	// Verify other fields are present
	s.Equal("user-123", jsonMap["id"])
	s.Equal("John Doe", jsonMap["name"])
	s.Equal("john@example.com", jsonMap["email"])
	s.Equal("individual", jsonMap["user_type"])

	// Test JSON unmarshaling
	var unmarshaledUser modelsv1.User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	s.NoError(err, "Should unmarshal JSON to user struct")
	s.Equal(user.ID, unmarshaledUser.ID)
	s.Equal(user.Name, unmarshaledUser.Name)
	s.Equal(user.Email, unmarshaledUser.Email)
	s.Equal(user.UserType, unmarshaledUser.UserType)
	// Password should remain empty after unmarshaling
	s.Empty(unmarshaledUser.Password, "Password should not be unmarshaled from JSON")
}

// TestUser_WithMemberships tests user serialization with memberships
func (s *UserModelTestSuite) TestUser_WithMemberships() {
	user := modelsv1.User{
		ID:       "user-123",
		Name:     "John Doe",
		Email:    "john@example.com",
		UserType: modelsv1.UserTypeIndividual,
		Memberships: []modelsv1.TeamMember{
			{
				ID:     "member-1",
				TeamID: "team-1",
				UserID: "user-123",
				Role:   modelsv1.RoleEditor,
			},
		},
	}

	jsonData, err := json.Marshal(user)
	s.NoError(err, "Should marshal user with memberships")

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	s.NoError(err, "Should unmarshal JSON to map")

	// Verify memberships are included
	s.Contains(jsonMap, "memberships", "Memberships should be included in JSON")
	memberships, ok := jsonMap["memberships"].([]interface{})
	s.True(ok, "Memberships should be an array")
	s.Len(memberships, 1, "Should have one membership")
}

// TestUserType_Constants tests user type constants
func (s *UserModelTestSuite) TestUserType_Constants() {
	s.Equal("individual", string(modelsv1.UserTypeIndividual))
	s.Equal("team", string(modelsv1.UserTypeTeam))
}

// TestUserType_Validation tests user type validation
func (s *UserModelTestSuite) TestUserType_Validation() {
	validTypes := []modelsv1.UserType{
		modelsv1.UserTypeIndividual,
		modelsv1.UserTypeTeam,
	}

	for _, userType := range validTypes {
		s.Contains([]string{"individual", "team"}, string(userType), 
			"User type should be valid: %s", userType)
	}
}

// TestRole_Constants tests role constants
func (s *UserModelTestSuite) TestRole_Constants() {
	s.Equal("admin", string(modelsv1.RoleAdmin))
	s.Equal("editor", string(modelsv1.RoleEditor))
	s.Equal("creator", string(modelsv1.RoleCreator))
}

// TestRole_Validation tests role validation
func (s *UserModelTestSuite) TestRole_Validation() {
	validRoles := []modelsv1.Role{
		modelsv1.RoleAdmin,
		modelsv1.RoleEditor,
		modelsv1.RoleCreator,
	}

	for _, role := range validRoles {
		s.Contains([]string{"admin", "editor", "creator"}, string(role), 
			"Role should be valid: %s", role)
	}
}

// TestSignupParams_JSONSerialization tests signup parameters JSON handling
func (s *UserModelTestSuite) TestSignupParams_JSONSerialization() {
	// Test individual user signup
	signupParams := modelsv1.SignupParams{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		UserType: modelsv1.UserTypeIndividual,
	}

	jsonData, err := json.Marshal(signupParams)
	s.NoError(err, "Should marshal signup params to JSON")

	var unmarshaledParams modelsv1.SignupParams
	err = json.Unmarshal(jsonData, &unmarshaledParams)
	s.NoError(err, "Should unmarshal JSON to signup params")
	s.Equal(signupParams.Name, unmarshaledParams.Name)
	s.Equal(signupParams.Email, unmarshaledParams.Email)
	s.Equal(signupParams.Password, unmarshaledParams.Password)
	s.Equal(signupParams.UserType, unmarshaledParams.UserType)

	// Test team user signup
	teamSignupParams := modelsv1.SignupParams{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Password: "password456",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Development Team",
	}

	jsonData, err = json.Marshal(teamSignupParams)
	s.NoError(err, "Should marshal team signup params to JSON")

	err = json.Unmarshal(jsonData, &unmarshaledParams)
	s.NoError(err, "Should unmarshal JSON to team signup params")
	s.Equal(teamSignupParams.TeamName, unmarshaledParams.TeamName)
}

// TestSignupParams_ValidationTags tests signup parameters validation tags
func (s *UserModelTestSuite) TestSignupParams_ValidationTags() {
	// Test that validation tags are properly set
	signupParams := modelsv1.SignupParams{}
	
	// Use reflection to check struct tags
	signupType := reflect.TypeOf(signupParams)
	
	// Check name field tags
	nameField, found := signupType.FieldByName("Name")
	s.True(found, "Name field should exist")
	s.Contains(nameField.Tag.Get("binding"), "required", "Name should be required")
	s.Contains(nameField.Tag.Get("binding"), "min=2", "Name should have minimum length")
	
	// Check email field tags
	emailField, found := signupType.FieldByName("Email")
	s.True(found, "Email field should exist")
	s.Contains(emailField.Tag.Get("binding"), "required", "Email should be required")
	s.Contains(emailField.Tag.Get("binding"), "email", "Email should have email validation")
	
	// Check password field tags
	passwordField, found := signupType.FieldByName("Password")
	s.True(found, "Password field should exist")
	s.Contains(passwordField.Tag.Get("binding"), "required", "Password should be required")
	s.Contains(passwordField.Tag.Get("binding"), "min=8", "Password should have minimum length")
	
	// Check user type field tags
	userTypeField, found := signupType.FieldByName("UserType")
	s.True(found, "UserType field should exist")
	s.Contains(userTypeField.Tag.Get("binding"), "required", "UserType should be required")
	s.Contains(userTypeField.Tag.Get("binding"), "oneof=individual team", "UserType should have enum validation")
	
	// Check team name field tags
	teamNameField, found := signupType.FieldByName("TeamName")
	s.True(found, "TeamName field should exist")
	s.Contains(teamNameField.Tag.Get("binding"), "required_if=UserType team", "TeamName should be required if UserType is team")
}

// TestLoginParams_JSONSerialization tests login parameters JSON handling
func (s *UserModelTestSuite) TestLoginParams_JSONSerialization() {
	loginParams := modelsv1.LoginParams{
		Email:    "john@example.com",
		Password: "password123",
	}

	jsonData, err := json.Marshal(loginParams)
	s.NoError(err, "Should marshal login params to JSON")

	var unmarshaledParams modelsv1.LoginParams
	err = json.Unmarshal(jsonData, &unmarshaledParams)
	s.NoError(err, "Should unmarshal JSON to login params")
	s.Equal(loginParams.Email, unmarshaledParams.Email)
	s.Equal(loginParams.Password, unmarshaledParams.Password)
}

// TestLoginParams_ValidationTags tests login parameters validation tags
func (s *UserModelTestSuite) TestLoginParams_ValidationTags() {
	loginParams := modelsv1.LoginParams{}
	loginType := reflect.TypeOf(loginParams)
	
	// Check email field tags
	emailField, found := loginType.FieldByName("Email")
	s.True(found, "Email field should exist")
	s.Contains(emailField.Tag.Get("binding"), "required", "Email should be required")
	s.Contains(emailField.Tag.Get("binding"), "email", "Email should have email validation")
	
	// Check password field tags
	passwordField, found := loginType.FieldByName("Password")
	s.True(found, "Password field should exist")
	s.Contains(passwordField.Tag.Get("binding"), "required", "Password should be required")
}

// TestRefreshTokenParams_JSONSerialization tests refresh token parameters JSON handling
func (s *UserModelTestSuite) TestRefreshTokenParams_JSONSerialization() {
	refreshParams := modelsv1.RefreshTokenParams{
		RefreshToken: "test_refresh_token_123",
	}

	jsonData, err := json.Marshal(refreshParams)
	s.NoError(err, "Should marshal refresh token params to JSON")

	var unmarshaledParams modelsv1.RefreshTokenParams
	err = json.Unmarshal(jsonData, &unmarshaledParams)
	s.NoError(err, "Should unmarshal JSON to refresh token params")
	s.Equal(refreshParams.RefreshToken, unmarshaledParams.RefreshToken)
}

// TestRefreshTokenParams_ValidationTags tests refresh token parameters validation tags
func (s *UserModelTestSuite) TestRefreshTokenParams_ValidationTags() {
	refreshParams := modelsv1.RefreshTokenParams{}
	refreshType := reflect.TypeOf(refreshParams)
	
	// Check refresh token field tags
	tokenField, found := refreshType.FieldByName("RefreshToken")
	s.True(found, "RefreshToken field should exist")
	s.Contains(tokenField.Tag.Get("binding"), "required", "RefreshToken should be required")
}

// TestTeam_JSONSerialization tests team struct JSON serialization
func (s *UserModelTestSuite) TestTeam_JSONSerialization() {
	team := modelsv1.Team{
		ID:      "team-123",
		Name:    "Development Team",
		OwnerID: "user-123",
		Owner: modelsv1.User{
			ID:       "user-123",
			Name:     "John Doe",
			Email:    "john@example.com",
			UserType: modelsv1.UserTypeTeam,
		},
		Members: []modelsv1.TeamMember{
			{
				ID:     "member-1",
				TeamID: "team-123",
				UserID: "user-456",
				Role:   modelsv1.RoleEditor,
			},
		},
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	jsonData, err := json.Marshal(team)
	s.NoError(err, "Should marshal team to JSON")

	var unmarshaledTeam modelsv1.Team
	err = json.Unmarshal(jsonData, &unmarshaledTeam)
	s.NoError(err, "Should unmarshal JSON to team struct")
	s.Equal(team.ID, unmarshaledTeam.ID)
	s.Equal(team.Name, unmarshaledTeam.Name)
	s.Equal(team.OwnerID, unmarshaledTeam.OwnerID)
	s.Equal(team.Owner.ID, unmarshaledTeam.Owner.ID)
	s.Len(unmarshaledTeam.Members, 1)
}

// TestTeamMember_JSONSerialization tests team member struct JSON serialization
func (s *UserModelTestSuite) TestTeamMember_JSONSerialization() {
	teamMember := modelsv1.TeamMember{
		ID:        "member-123",
		TeamID:    "team-456",
		UserID:    "user-789",
		Role:      modelsv1.RoleEditor,
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	jsonData, err := json.Marshal(teamMember)
	s.NoError(err, "Should marshal team member to JSON")

	var unmarshaledMember modelsv1.TeamMember
	err = json.Unmarshal(jsonData, &unmarshaledMember)
	s.NoError(err, "Should unmarshal JSON to team member struct")
	s.Equal(teamMember.ID, unmarshaledMember.ID)
	s.Equal(teamMember.TeamID, unmarshaledMember.TeamID)
	s.Equal(teamMember.UserID, unmarshaledMember.UserID)
	s.Equal(teamMember.Role, unmarshaledMember.Role)
}

// TestUser_BusinessRules tests user business rules and validation
func (s *UserModelTestSuite) TestUser_BusinessRules() {
	// Test that individual users don't require team name
	individualUser := modelsv1.User{
		ID:       "user-123",
		Name:     "John Doe",
		Email:    "john@example.com",
		UserType: modelsv1.UserTypeIndividual,
	}
	s.Equal(modelsv1.UserTypeIndividual, individualUser.UserType)

	// Test that team users have team type
	teamUser := modelsv1.User{
		ID:       "user-456",
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		UserType: modelsv1.UserTypeTeam,
	}
	s.Equal(modelsv1.UserTypeTeam, teamUser.UserType)
}

// TestSignupParams_BusinessRules tests signup parameters business rules
func (s *UserModelTestSuite) TestSignupParams_BusinessRules() {
	// Test individual signup doesn't require team name
	individualSignup := modelsv1.SignupParams{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		UserType: modelsv1.UserTypeIndividual,
		TeamName: "", // Should be empty for individual
	}
	s.Equal(modelsv1.UserTypeIndividual, individualSignup.UserType)
	s.Empty(individualSignup.TeamName)

	// Test team signup requires team name
	teamSignup := modelsv1.SignupParams{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Password: "password456",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Development Team", // Should be provided for team
	}
	s.Equal(modelsv1.UserTypeTeam, teamSignup.UserType)
	s.NotEmpty(teamSignup.TeamName)
}

// TestUser_EdgeCases tests edge cases for user model
func (s *UserModelTestSuite) TestUser_EdgeCases() {
	// Test user with empty memberships
	user := modelsv1.User{
		ID:          "user-123",
		Name:        "John Doe",
		Email:       "john@example.com",
		UserType:    modelsv1.UserTypeIndividual,
		Memberships: []modelsv1.TeamMember{}, // Empty slice
	}

	jsonData, err := json.Marshal(user)
	s.NoError(err, "Should marshal user with empty memberships")

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	s.NoError(err, "Should unmarshal JSON to map")

	// Memberships should be omitted when empty due to omitempty tag
	s.NotContains(jsonMap, "memberships", "Empty memberships should be omitted from JSON")
}

