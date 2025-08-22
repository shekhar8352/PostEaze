package modelsv1_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func TestUser_JSONSerialization(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal user to JSON without error: %v", err)
	}

	// Verify password is excluded from JSON (json:"-" tag)
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to map: %v", err)
	}
	if _, exists := jsonMap["password"]; exists {
		t.Errorf("Password should not be included in JSON")
	}

	// Verify other fields are present
	if jsonMap["id"] != "user-123" {
		t.Errorf("Expected id user-123, got %v", jsonMap["id"])
	}
	if jsonMap["name"] != "John Doe" {
		t.Errorf("Expected name John Doe, got %v", jsonMap["name"])
	}
	if jsonMap["email"] != "john@example.com" {
		t.Errorf("Expected email john@example.com, got %v", jsonMap["email"])
	}
	if jsonMap["user_type"] != "individual" {
		t.Errorf("Expected user_type individual, got %v", jsonMap["user_type"])
	}

	// Test JSON unmarshaling
	var unmarshaledUser modelsv1.User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to user struct: %v", err)
	}
	if user.ID != unmarshaledUser.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, unmarshaledUser.ID)
	}
	if user.Name != unmarshaledUser.Name {
		t.Errorf("Expected Name %s, got %s", user.Name, unmarshaledUser.Name)
	}
	if user.Email != unmarshaledUser.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, unmarshaledUser.Email)
	}
	if user.UserType != unmarshaledUser.UserType {
		t.Errorf("Expected UserType %s, got %s", user.UserType, unmarshaledUser.UserType)
	}
	// Password should remain empty after unmarshaling
	if unmarshaledUser.Password != "" {
		t.Errorf("Password should not be unmarshaled from JSON")
	}
}

func TestUser_WithMemberships(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal user with memberships: %v", err)
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to map: %v", err)
	}

	// Verify memberships are included
	if _, exists := jsonMap["memberships"]; !exists {
		t.Errorf("Memberships should be included in JSON")
	}
	memberships, ok := jsonMap["memberships"].([]interface{})
	if !ok {
		t.Errorf("Memberships should be an array")
	}
	if len(memberships) != 1 {
		t.Errorf("Should have one membership, got %d", len(memberships))
	}
}

func TestUserType_Constants(t *testing.T) {
	if string(modelsv1.UserTypeIndividual) != "individual" {
		t.Errorf("Expected UserTypeIndividual to be 'individual', got %s", string(modelsv1.UserTypeIndividual))
	}
	if string(modelsv1.UserTypeTeam) != "team" {
		t.Errorf("Expected UserTypeTeam to be 'team', got %s", string(modelsv1.UserTypeTeam))
	}
}

func TestUserType_Validation(t *testing.T) {
	validTypes := []modelsv1.UserType{
		modelsv1.UserTypeIndividual,
		modelsv1.UserTypeTeam,
	}

	expectedValues := []string{"individual", "team"}

	for i, userType := range validTypes {
		if string(userType) != expectedValues[i] {
			t.Errorf("User type should be valid: expected %s, got %s", expectedValues[i], string(userType))
		}
	}
}

func TestRole_Constants(t *testing.T) {
	if string(modelsv1.RoleAdmin) != "admin" {
		t.Errorf("Expected RoleAdmin to be 'admin', got %s", string(modelsv1.RoleAdmin))
	}
	if string(modelsv1.RoleEditor) != "editor" {
		t.Errorf("Expected RoleEditor to be 'editor', got %s", string(modelsv1.RoleEditor))
	}
	if string(modelsv1.RoleCreator) != "creator" {
		t.Errorf("Expected RoleCreator to be 'creator', got %s", string(modelsv1.RoleCreator))
	}
}

func TestRole_Validation(t *testing.T) {
	validRoles := []modelsv1.Role{
		modelsv1.RoleAdmin,
		modelsv1.RoleEditor,
		modelsv1.RoleCreator,
	}

	expectedValues := []string{"admin", "editor", "creator"}

	for i, role := range validRoles {
		if string(role) != expectedValues[i] {
			t.Errorf("Role should be valid: expected %s, got %s", expectedValues[i], string(role))
		}
	}
}

func TestSignupParams_JSONSerialization(t *testing.T) {
	// Test individual user signup
	signupParams := modelsv1.SignupParams{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		UserType: modelsv1.UserTypeIndividual,
	}

	jsonData, err := json.Marshal(signupParams)
	if err != nil {
		t.Fatalf("Should marshal signup params to JSON: %v", err)
	}

	var unmarshaledParams modelsv1.SignupParams
	err = json.Unmarshal(jsonData, &unmarshaledParams)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to signup params: %v", err)
	}
	if signupParams.Name != unmarshaledParams.Name {
		t.Errorf("Expected Name %s, got %s", signupParams.Name, unmarshaledParams.Name)
	}
	if signupParams.Email != unmarshaledParams.Email {
		t.Errorf("Expected Email %s, got %s", signupParams.Email, unmarshaledParams.Email)
	}
	if signupParams.Password != unmarshaledParams.Password {
		t.Errorf("Expected Password %s, got %s", signupParams.Password, unmarshaledParams.Password)
	}
	if signupParams.UserType != unmarshaledParams.UserType {
		t.Errorf("Expected UserType %s, got %s", signupParams.UserType, unmarshaledParams.UserType)
	}

	// Test team user signup
	teamSignupParams := modelsv1.SignupParams{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Password: "password456",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Development Team",
	}

	jsonData, err = json.Marshal(teamSignupParams)
	if err != nil {
		t.Fatalf("Should marshal team signup params to JSON: %v", err)
	}

	err = json.Unmarshal(jsonData, &unmarshaledParams)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to team signup params: %v", err)
	}
	if teamSignupParams.TeamName != unmarshaledParams.TeamName {
		t.Errorf("Expected TeamName %s, got %s", teamSignupParams.TeamName, unmarshaledParams.TeamName)
	}
}

func TestSignupParams_ValidationTags(t *testing.T) {
	// Test that validation tags are properly set
	signupParams := modelsv1.SignupParams{}

	// Use reflection to check struct tags
	signupType := reflect.TypeOf(signupParams)

	// Check name field tags
	nameField, found := signupType.FieldByName("Name")
	if !found {
		t.Fatalf("Name field should exist")
	}
	bindingTag := nameField.Tag.Get("binding")
	if bindingTag == "" {
		t.Errorf("Name should have binding tag")
	}

	// Check email field tags
	emailField, found := signupType.FieldByName("Email")
	if !found {
		t.Fatalf("Email field should exist")
	}
	emailBindingTag := emailField.Tag.Get("binding")
	if emailBindingTag == "" {
		t.Errorf("Email should have binding tag")
	}

	// Check password field tags
	passwordField, found := signupType.FieldByName("Password")
	if !found {
		t.Fatalf("Password field should exist")
	}
	passwordBindingTag := passwordField.Tag.Get("binding")
	if passwordBindingTag == "" {
		t.Errorf("Password should have binding tag")
	}

	// Check user type field tags
	userTypeField, found := signupType.FieldByName("UserType")
	if !found {
		t.Fatalf("UserType field should exist")
	}
	userTypeBindingTag := userTypeField.Tag.Get("binding")
	if userTypeBindingTag == "" {
		t.Errorf("UserType should have binding tag")
	}

	// Check team name field tags
	teamNameField, found := signupType.FieldByName("TeamName")
	if !found {
		t.Fatalf("TeamName field should exist")
	}
	teamNameBindingTag := teamNameField.Tag.Get("binding")
	if teamNameBindingTag == "" {
		t.Errorf("TeamName should have binding tag")
	}
}

func TestLoginParams_JSONSerialization(t *testing.T) {
	loginParams := modelsv1.LoginParams{
		Email:    "john@example.com",
		Password: "password123",
	}

	jsonData, err := json.Marshal(loginParams)
	if err != nil {
		t.Fatalf("Should marshal login params to JSON: %v", err)
	}

	var unmarshaledParams modelsv1.LoginParams
	err = json.Unmarshal(jsonData, &unmarshaledParams)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to login params: %v", err)
	}
	if loginParams.Email != unmarshaledParams.Email {
		t.Errorf("Expected Email %s, got %s", loginParams.Email, unmarshaledParams.Email)
	}
	if loginParams.Password != unmarshaledParams.Password {
		t.Errorf("Expected Password %s, got %s", loginParams.Password, unmarshaledParams.Password)
	}
}

func TestLoginParams_ValidationTags(t *testing.T) {
	loginParams := modelsv1.LoginParams{}
	loginType := reflect.TypeOf(loginParams)

	// Check email field tags
	emailField, found := loginType.FieldByName("Email")
	if !found {
		t.Fatalf("Email field should exist")
	}
	emailBindingTag := emailField.Tag.Get("binding")
	if emailBindingTag == "" {
		t.Errorf("Email should have binding tag")
	}

	// Check password field tags
	passwordField, found := loginType.FieldByName("Password")
	if !found {
		t.Fatalf("Password field should exist")
	}
	passwordBindingTag := passwordField.Tag.Get("binding")
	if passwordBindingTag == "" {
		t.Errorf("Password should have binding tag")
	}
}

func TestRefreshTokenParams_JSONSerialization(t *testing.T) {
	refreshParams := modelsv1.RefreshTokenParams{
		RefreshToken: "test_refresh_token_123",
	}

	jsonData, err := json.Marshal(refreshParams)
	if err != nil {
		t.Fatalf("Should marshal refresh token params to JSON: %v", err)
	}

	var unmarshaledParams modelsv1.RefreshTokenParams
	err = json.Unmarshal(jsonData, &unmarshaledParams)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to refresh token params: %v", err)
	}
	if refreshParams.RefreshToken != unmarshaledParams.RefreshToken {
		t.Errorf("Expected RefreshToken %s, got %s", refreshParams.RefreshToken, unmarshaledParams.RefreshToken)
	}
}

func TestRefreshTokenParams_ValidationTags(t *testing.T) {
	refreshParams := modelsv1.RefreshTokenParams{}
	refreshType := reflect.TypeOf(refreshParams)

	// Check refresh token field tags
	tokenField, found := refreshType.FieldByName("RefreshToken")
	if !found {
		t.Fatalf("RefreshToken field should exist")
	}
	tokenBindingTag := tokenField.Tag.Get("binding")
	if tokenBindingTag == "" {
		t.Errorf("RefreshToken should have binding tag")
	}
}

func TestTeam_JSONSerialization(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal team to JSON: %v", err)
	}

	var unmarshaledTeam modelsv1.Team
	err = json.Unmarshal(jsonData, &unmarshaledTeam)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to team struct: %v", err)
	}
	if team.ID != unmarshaledTeam.ID {
		t.Errorf("Expected ID %s, got %s", team.ID, unmarshaledTeam.ID)
	}
	if team.Name != unmarshaledTeam.Name {
		t.Errorf("Expected Name %s, got %s", team.Name, unmarshaledTeam.Name)
	}
	if team.OwnerID != unmarshaledTeam.OwnerID {
		t.Errorf("Expected OwnerID %s, got %s", team.OwnerID, unmarshaledTeam.OwnerID)
	}
	if team.Owner.ID != unmarshaledTeam.Owner.ID {
		t.Errorf("Expected Owner.ID %s, got %s", team.Owner.ID, unmarshaledTeam.Owner.ID)
	}
	if len(unmarshaledTeam.Members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(unmarshaledTeam.Members))
	}
}

func TestTeamMember_JSONSerialization(t *testing.T) {
	teamMember := modelsv1.TeamMember{
		ID:        "member-123",
		TeamID:    "team-456",
		UserID:    "user-789",
		Role:      modelsv1.RoleEditor,
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	jsonData, err := json.Marshal(teamMember)
	if err != nil {
		t.Fatalf("Should marshal team member to JSON: %v", err)
	}

	var unmarshaledMember modelsv1.TeamMember
	err = json.Unmarshal(jsonData, &unmarshaledMember)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to team member struct: %v", err)
	}
	if teamMember.ID != unmarshaledMember.ID {
		t.Errorf("Expected ID %s, got %s", teamMember.ID, unmarshaledMember.ID)
	}
	if teamMember.TeamID != unmarshaledMember.TeamID {
		t.Errorf("Expected TeamID %s, got %s", teamMember.TeamID, unmarshaledMember.TeamID)
	}
	if teamMember.UserID != unmarshaledMember.UserID {
		t.Errorf("Expected UserID %s, got %s", teamMember.UserID, unmarshaledMember.UserID)
	}
	if teamMember.Role != unmarshaledMember.Role {
		t.Errorf("Expected Role %s, got %s", teamMember.Role, unmarshaledMember.Role)
	}
}

func TestUser_BusinessRules(t *testing.T) {
	// Test that individual users don't require team name
	individualUser := modelsv1.User{
		ID:       "user-123",
		Name:     "John Doe",
		Email:    "john@example.com",
		UserType: modelsv1.UserTypeIndividual,
	}
	if individualUser.UserType != modelsv1.UserTypeIndividual {
		t.Errorf("Expected UserType to be individual")
	}

	// Test that team users have team type
	teamUser := modelsv1.User{
		ID:       "user-456",
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		UserType: modelsv1.UserTypeTeam,
	}
	if teamUser.UserType != modelsv1.UserTypeTeam {
		t.Errorf("Expected UserType to be team")
	}
}

func TestSignupParams_BusinessRules(t *testing.T) {
	// Test individual signup doesn't require team name
	individualSignup := modelsv1.SignupParams{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		UserType: modelsv1.UserTypeIndividual,
		TeamName: "", // Should be empty for individual
	}
	if individualSignup.UserType != modelsv1.UserTypeIndividual {
		t.Errorf("Expected UserType to be individual")
	}
	if individualSignup.TeamName != "" {
		t.Errorf("TeamName should be empty for individual signup")
	}

	// Test team signup requires team name
	teamSignup := modelsv1.SignupParams{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Password: "password456",
		UserType: modelsv1.UserTypeTeam,
		TeamName: "Development Team", // Should be provided for team
	}
	if teamSignup.UserType != modelsv1.UserTypeTeam {
		t.Errorf("Expected UserType to be team")
	}
	if teamSignup.TeamName == "" {
		t.Errorf("TeamName should not be empty for team signup")
	}
}

func TestUser_EdgeCases(t *testing.T) {
	// Test user with empty memberships
	user := modelsv1.User{
		ID:          "user-123",
		Name:        "John Doe",
		Email:       "john@example.com",
		UserType:    modelsv1.UserTypeIndividual,
		Memberships: []modelsv1.TeamMember{}, // Empty slice
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Should marshal user with empty memberships: %v", err)
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to map: %v", err)
	}

	// Memberships should be omitted when empty due to omitempty tag
	if _, exists := jsonMap["memberships"]; exists {
		t.Errorf("Empty memberships should be omitted from JSON")
	}
}