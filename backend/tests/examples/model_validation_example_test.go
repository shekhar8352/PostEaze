package examples_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// ModelValidationExampleTestSuite demonstrates comprehensive model validation testing
type ModelValidationExampleTestSuite struct {
	testutils.ModelTestSuite
}

// Example 1: Testing basic model validation
func (s *ModelValidationExampleTestSuite) TestUser_Validate_ValidUser_NoError() {
	// Arrange
	user := &User{
		ID:       "user-123",
		Name:     "John Doe",
		Email:    "john@example.com",
		UserType: UserTypeIndividual,
		CreatedAt: time.Now(),
	}

	// Act
	err := user.Validate()

	// Assert
	s.NoError(err)
}

// Example 2: Testing validation errors for required fields
func (s *ModelValidationExampleTestSuite) TestUser_Validate_MissingRequiredFields_ReturnsErrors() {
	testCases := []struct {
		name          string
		user          *User
		expectedError string
	}{
		{
			name: "missing name",
			user: &User{
				ID:       "user-123",
				Email:    "john@example.com",
				UserType: UserTypeIndividual,
			},
			expectedError: "name is required",
		},
		{
			name: "missing email",
			user: &User{
				ID:       "user-123",
				Name:     "John Doe",
				UserType: UserTypeIndividual,
			},
			expectedError: "email is required",
		},
		{
			name: "missing user type",
			user: &User{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expectedError: "user type is required",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			err := tc.user.Validate()

			// Assert
			s.Error(err)
			s.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Example 3: Testing format validation
func (s *ModelValidationExampleTestSuite) TestUser_Validate_InvalidFormats_ReturnsErrors() {
	testCases := []struct {
		name          string
		user          *User
		expectedError string
	}{
		{
			name: "invalid email format",
			user: &User{
				ID:       "user-123",
				Name:     "John Doe",
				Email:    "invalid-email",
				UserType: UserTypeIndividual,
			},
			expectedError: "invalid email format",
		},
		{
			name: "empty name",
			user: &User{
				ID:       "user-123",
				Name:     "",
				Email:    "john@example.com",
				UserType: UserTypeIndividual,
			},
			expectedError: "name cannot be empty",
		},
		{
			name: "invalid user type",
			user: &User{
				ID:       "user-123",
				Name:     "John Doe",
				Email:    "john@example.com",
				UserType: "invalid-type",
			},
			expectedError: "invalid user type",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			err := tc.user.Validate()

			// Assert
			s.Error(err)
			s.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Example 4: Testing business rule validation
func (s *ModelValidationExampleTestSuite) TestUser_Validate_BusinessRules_ReturnsErrors() {
	testCases := []struct {
		name          string
		user          *User
		expectedError string
	}{
		{
			name: "team user without team ID",
			user: &User{
				ID:       "user-123",
				Name:     "Team User",
				Email:    "team@example.com",
				UserType: UserTypeTeam,
				// TeamID is missing
			},
			expectedError: "team users must have a team ID",
		},
		{
			name: "individual user with team ID",
			user: &User{
				ID:       "user-123",
				Name:     "Individual User",
				Email:    "individual@example.com",
				UserType: UserTypeIndividual,
				TeamID:   "team-456", // Should not have team ID
			},
			expectedError: "individual users cannot have a team ID",
		},
		{
			name: "name too long",
			user: &User{
				ID:       "user-123",
				Name:     "This is a very long name that exceeds the maximum allowed length for user names in the system",
				Email:    "john@example.com",
				UserType: UserTypeIndividual,
			},
			expectedError: "name exceeds maximum length",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			err := tc.user.Validate()

			// Assert
			s.Error(err)
			s.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Example 5: Testing JSON serialization and deserialization
func (s *ModelValidationExampleTestSuite) TestUser_JSONSerialization_ValidUser_Success() {
	// Arrange
	originalUser := &User{
		ID:        "user-123",
		Name:      "John Doe",
		Email:     "john@example.com",
		UserType:  UserTypeIndividual,
		Bio:       "Software developer",
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	// Act - Serialize to JSON
	jsonData, err := json.Marshal(originalUser)
	s.NoError(err)

	// Act - Deserialize from JSON
	var deserializedUser User
	err = json.Unmarshal(jsonData, &deserializedUser)
	s.NoError(err)

	// Assert
	s.Equal(originalUser.ID, deserializedUser.ID)
	s.Equal(originalUser.Name, deserializedUser.Name)
	s.Equal(originalUser.Email, deserializedUser.Email)
	s.Equal(originalUser.UserType, deserializedUser.UserType)
	s.Equal(originalUser.Bio, deserializedUser.Bio)
	s.True(originalUser.CreatedAt.Equal(deserializedUser.CreatedAt))
	s.True(originalUser.UpdatedAt.Equal(deserializedUser.UpdatedAt))
}

// Example 6: Testing JSON field validation
func (s *ModelValidationExampleTestSuite) TestUser_JSONDeserialization_InvalidJSON_ReturnsError() {
	testCases := []struct {
		name        string
		jsonData    string
		expectError bool
	}{
		{
			name:        "valid JSON",
			jsonData:    `{"id":"user-123","name":"John","email":"john@example.com","user_type":"individual"}`,
			expectError: false,
		},
		{
			name:        "invalid JSON syntax",
			jsonData:    `{"id":"user-123","name":"John","email":"john@example.com"`,
			expectError: true,
		},
		{
			name:        "wrong field types",
			jsonData:    `{"id":123,"name":"John","email":"john@example.com","user_type":"individual"}`,
			expectError: true,
		},
		{
			name:        "invalid enum value",
			jsonData:    `{"id":"user-123","name":"John","email":"john@example.com","user_type":"invalid"}`,
			expectError: false, // JSON unmarshaling succeeds, but validation should fail
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			var user User
			err := json.Unmarshal([]byte(tc.jsonData), &user)

			// Assert JSON unmarshaling
			if tc.expectError {
				s.Error(err)
				return
			}
			s.NoError(err)

			// Validate the unmarshaled user
			validationErr := user.Validate()
			if tc.name == "invalid enum value" {
				s.Error(validationErr)
			}
		})
	}
}

// Example 7: Testing model relationships
func (s *ModelValidationExampleTestSuite) TestTeam_Validate_WithMembers_Success() {
	// Arrange
	owner := &User{
		ID:       "user-123",
		Name:     "Team Owner",
		Email:    "owner@example.com",
		UserType: UserTypeTeam,
		TeamID:   "team-456",
	}

	member := &User{
		ID:       "user-789",
		Name:     "Team Member",
		Email:    "member@example.com",
		UserType: UserTypeTeam,
		TeamID:   "team-456",
	}

	team := &Team{
		ID:      "team-456",
		Name:    "Test Team",
		OwnerID: owner.ID,
		Members: []*User{owner, member},
	}

	// Act
	err := team.Validate()

	// Assert
	s.NoError(err)
	s.Len(team.Members, 2)
	s.Equal(owner.ID, team.OwnerID)
}

// Example 8: Testing model constraints
func (s *ModelValidationExampleTestSuite) TestTeam_Validate_InvalidConstraints_ReturnsErrors() {
	testCases := []struct {
		name          string
		team          *Team
		expectedError string
	}{
		{
			name: "owner not in members",
			team: &Team{
				ID:      "team-123",
				Name:    "Test Team",
				OwnerID: "user-123",
				Members: []*User{
					{ID: "user-456", TeamID: "team-123"},
				},
			},
			expectedError: "team owner must be a member of the team",
		},
		{
			name: "member with different team ID",
			team: &Team{
				ID:      "team-123",
				Name:    "Test Team",
				OwnerID: "user-123",
				Members: []*User{
					{ID: "user-123", TeamID: "team-123"},
					{ID: "user-456", TeamID: "team-789"}, // Different team ID
				},
			},
			expectedError: "all team members must belong to this team",
		},
		{
			name: "duplicate members",
			team: &Team{
				ID:      "team-123",
				Name:    "Test Team",
				OwnerID: "user-123",
				Members: []*User{
					{ID: "user-123", TeamID: "team-123"},
					{ID: "user-123", TeamID: "team-123"}, // Duplicate
				},
			},
			expectedError: "duplicate team members are not allowed",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			err := tc.team.Validate()

			// Assert
			s.Error(err)
			s.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Example 9: Testing custom validation methods
func (s *ModelValidationExampleTestSuite) TestUser_ValidateForUpdate_PartialValidation_Success() {
	// Arrange
	updateData := &UserUpdateRequest{
		Name:  StringPtr("Updated Name"),
		Email: StringPtr("updated@example.com"),
		// Bio is nil, should not be validated
	}

	// Act
	err := updateData.ValidateForUpdate()

	// Assert
	s.NoError(err)
}

func (s *ModelValidationExampleTestSuite) TestUser_ValidateForUpdate_InvalidFields_ReturnsErrors() {
	testCases := []struct {
		name          string
		updateData    *UserUpdateRequest
		expectedError string
	}{
		{
			name: "invalid email format",
			updateData: &UserUpdateRequest{
				Email: StringPtr("invalid-email"),
			},
			expectedError: "invalid email format",
		},
		{
			name: "empty name",
			updateData: &UserUpdateRequest{
				Name: StringPtr(""),
			},
			expectedError: "name cannot be empty",
		},
		{
			name: "bio too long",
			updateData: &UserUpdateRequest{
				Bio: StringPtr("This is a very long bio that exceeds the maximum allowed length for user bios in the system and should trigger a validation error"),
			},
			expectedError: "bio exceeds maximum length",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			err := tc.updateData.ValidateForUpdate()

			// Assert
			s.Error(err)
			s.Contains(err.Error(), tc.expectedError)
		})
	}
}

// Example 10: Testing model transformation methods
func (s *ModelValidationExampleTestSuite) TestUser_ToPublicUser_RemovesSensitiveData() {
	// Arrange
	user := &User{
		ID:           "user-123",
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashed-password",
		UserType:     UserTypeIndividual,
		Bio:          "Software developer",
		CreatedAt:    time.Now(),
	}

	// Act
	publicUser := user.ToPublicUser()

	// Assert
	s.Equal(user.ID, publicUser.ID)
	s.Equal(user.Name, publicUser.Name)
	s.Equal(user.Email, publicUser.Email)
	s.Equal(user.UserType, publicUser.UserType)
	s.Equal(user.Bio, publicUser.Bio)
	s.True(user.CreatedAt.Equal(publicUser.CreatedAt))
	
	// Sensitive data should be removed
	s.Empty(publicUser.PasswordHash)
}

// Run the test suite
func TestModelValidationExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ModelValidationExampleTestSuite))
}

// Example of testing a simple validation function without a test suite
func TestValidateEmail_Examples(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectValid bool
	}{
		{
			name:        "valid email",
			email:       "user@example.com",
			expectValid: true,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			expectValid: true,
		},
		{
			name:        "valid email with plus",
			email:       "user+tag@example.com",
			expectValid: true,
		},
		{
			name:        "missing @ symbol",
			email:       "userexample.com",
			expectValid: false,
		},
		{
			name:        "missing domain",
			email:       "user@",
			expectValid: false,
		},
		{
			name:        "missing local part",
			email:       "@example.com",
			expectValid: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectValid: false,
		},
		{
			name:        "spaces in email",
			email:       "user @example.com",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			valid := ValidateEmail(tt.email)

			// Assert
			assert.Equal(t, tt.expectValid, valid)
		})
	}
}

// Example of testing model factory functions
func TestCreateUser_Examples(t *testing.T) {
	t.Run("create individual user", func(t *testing.T) {
		// Arrange
		params := &CreateUserParams{
			Name:     "John Doe",
			Email:    "john@example.com",
			UserType: UserTypeIndividual,
		}

		// Act
		user, err := CreateUser(params)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, params.Name, user.Name)
		assert.Equal(t, params.Email, user.Email)
		assert.Equal(t, params.UserType, user.UserType)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("create team user", func(t *testing.T) {
		// Arrange
		params := &CreateUserParams{
			Name:     "Team Owner",
			Email:    "owner@company.com",
			UserType: UserTypeTeam,
		}

		// Act
		user, err := CreateUser(params)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, UserTypeTeam, user.UserType)
		// Team ID should be set for team users
		assert.NotEmpty(t, user.TeamID)
	})
}

// Helper types and functions (these would normally be in your actual code)

type UserType string

const (
	UserTypeIndividual UserType = "individual"
	UserTypeTeam       UserType = "team"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON
	UserType     UserType  `json:"user_type" db:"user_type"`
	TeamID       string    `json:"team_id,omitempty" db:"team_id"`
	Bio          string    `json:"bio" db:"bio"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type PublicUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserType  UserType  `json:"user_type"`
	TeamID    string    `json:"team_id,omitempty"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	ID      string  `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	OwnerID string  `json:"owner_id" db:"owner_id"`
	Members []*User `json:"members,omitempty"`
}

type UserUpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Bio  *string `json:"bio,omitempty"`
}

type CreateUserParams struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`
}

// Validation methods
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	
	if len(u.Name) > 100 {
		return errors.New("name exceeds maximum length")
	}
	
	if u.Email == "" {
		return errors.New("email is required")
	}
	
	if !ValidateEmail(u.Email) {
		return errors.New("invalid email format")
	}
	
	if u.UserType == "" {
		return errors.New("user type is required")
	}
	
	if u.UserType != UserTypeIndividual && u.UserType != UserTypeTeam {
		return errors.New("invalid user type")
	}
	
	// Business rules
	if u.UserType == UserTypeTeam && u.TeamID == "" {
		return errors.New("team users must have a team ID")
	}
	
	if u.UserType == UserTypeIndividual && u.TeamID != "" {
		return errors.New("individual users cannot have a team ID")
	}
	
	return nil
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return errors.New("team name is required")
	}
	
	if t.OwnerID == "" {
		return errors.New("team owner ID is required")
	}
	
	// Check if owner is in members
	ownerFound := false
	memberIDs := make(map[string]bool)
	
	for _, member := range t.Members {
		if member.ID == t.OwnerID {
			ownerFound = true
		}
		
		if member.TeamID != t.ID {
			return errors.New("all team members must belong to this team")
		}
		
		if memberIDs[member.ID] {
			return errors.New("duplicate team members are not allowed")
		}
		memberIDs[member.ID] = true
	}
	
	if !ownerFound && len(t.Members) > 0 {
		return errors.New("team owner must be a member of the team")
	}
	
	return nil
}

func (u *UserUpdateRequest) ValidateForUpdate() error {
	if u.Name != nil {
		if *u.Name == "" {
			return errors.New("name cannot be empty")
		}
		if len(*u.Name) > 100 {
			return errors.New("name exceeds maximum length")
		}
	}
	
	if u.Email != nil {
		if !ValidateEmail(*u.Email) {
			return errors.New("invalid email format")
		}
	}
	
	if u.Bio != nil {
		if len(*u.Bio) > 500 {
			return errors.New("bio exceeds maximum length")
		}
	}
	
	return nil
}

func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		UserType:  u.UserType,
		TeamID:    u.TeamID,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
	}
}

func ValidateEmail(email string) bool {
	// Simple email validation for example
	if email == "" {
		return false
	}
	
	// Check for @ symbol
	atIndex := strings.Index(email, "@")
	if atIndex <= 0 || atIndex == len(email)-1 {
		return false
	}
	
	// Check for spaces
	if strings.Contains(email, " ") {
		return false
	}
	
	// More comprehensive validation would use regex or email parsing library
	return true
}

func CreateUser(params *CreateUserParams) (*User, error) {
	user := &User{
		ID:        generateID(),
		Name:      params.Name,
		Email:     params.Email,
		UserType:  params.UserType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set team ID for team users
	if params.UserType == UserTypeTeam {
		user.TeamID = generateID()
	}
	
	// Validate the created user
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	return user, nil
}

func generateID() string {
	// Simple ID generation for example
	return fmt.Sprintf("id-%d", time.Now().UnixNano())
}

func StringPtr(s string) *string {
	return &s
}