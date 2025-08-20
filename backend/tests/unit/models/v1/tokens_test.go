package modelsv1_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// TokenModelTestSuite tests token model validation and serialization
type TokenModelTestSuite struct {
	testutils.ModelTestSuite
}

// TestTokenModelTestSuite runs the token model test suite
func TestTokenModelTestSuite(t *testing.T) {
	suite.Run(t, new(TokenModelTestSuite))
}

// TestRefreshToken_JSONSerialization tests refresh token struct JSON serialization
func (s *TokenModelTestSuite) TestRefreshToken_JSONSerialization() {
	userID := uuid.New()
	tokenID := uuid.New()
	
	refreshToken := modelsv1.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     "test_refresh_token_123",
		ExpiresAt: time.Date(2024, 1, 22, 10, 30, 0, 0, time.UTC),
		Revoked:   false,
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(refreshToken)
	s.NoError(err, "Should marshal refresh token to JSON without error")

	// Test JSON unmarshaling
	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	s.NoError(err, "Should unmarshal JSON to refresh token struct")
	s.Equal(refreshToken.ID, unmarshaledToken.ID)
	s.Equal(refreshToken.UserID, unmarshaledToken.UserID)
	s.Equal(refreshToken.Token, unmarshaledToken.Token)
	s.Equal(refreshToken.ExpiresAt.Unix(), unmarshaledToken.ExpiresAt.Unix())
	s.Equal(refreshToken.Revoked, unmarshaledToken.Revoked)
}

// TestRefreshToken_WithUser tests refresh token with user relationship
func (s *TokenModelTestSuite) TestRefreshToken_WithUser() {
	userID := uuid.New()
	tokenID := uuid.New()
	
	user := modelsv1.User{
		ID:       userID.String(),
		Name:     "John Doe",
		Email:    "john@example.com",
		UserType: modelsv1.UserTypeIndividual,
	}

	refreshToken := modelsv1.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		User:      user,
		Token:     "test_refresh_token_123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	jsonData, err := json.Marshal(refreshToken)
	s.NoError(err, "Should marshal refresh token with user")

	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	s.NoError(err, "Should unmarshal JSON to refresh token with user")
	s.Equal(refreshToken.User.ID, unmarshaledToken.User.ID)
	s.Equal(refreshToken.User.Name, unmarshaledToken.User.Name)
	s.Equal(refreshToken.User.Email, unmarshaledToken.User.Email)
}

// TestRefreshToken_ExpirationLogic tests token expiration logic
func (s *TokenModelTestSuite) TestRefreshToken_ExpirationLogic() {
	now := time.Now()
	
	// Test valid (non-expired) token
	validToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "valid_token",
		ExpiresAt: now.Add(7 * 24 * time.Hour), // Expires in 7 days
		Revoked:   false,
	}
	
	s.True(validToken.ExpiresAt.After(now), "Valid token should not be expired")
	s.False(validToken.Revoked, "Valid token should not be revoked")

	// Test expired token
	expiredToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "expired_token",
		ExpiresAt: now.Add(-1 * time.Hour), // Expired 1 hour ago
		Revoked:   false,
	}
	
	s.True(expiredToken.ExpiresAt.Before(now), "Expired token should be expired")
	s.False(expiredToken.Revoked, "Expired token should not be revoked (just expired)")

	// Test revoked token
	revokedToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "revoked_token",
		ExpiresAt: now.Add(7 * 24 * time.Hour), // Valid expiration but revoked
		Revoked:   true,
	}
	
	s.True(revokedToken.ExpiresAt.After(now), "Revoked token expiration should still be valid")
	s.True(revokedToken.Revoked, "Revoked token should be marked as revoked")
}

// TestRefreshToken_RevocationLogic tests token revocation logic
func (s *TokenModelTestSuite) TestRefreshToken_RevocationLogic() {
	userID := uuid.New()
	
	// Test token before revocation
	token := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "test_token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	s.False(token.Revoked, "Token should not be revoked initially")

	// Simulate revocation with a small delay to ensure UpdatedAt is after CreatedAt
	time.Sleep(1 * time.Millisecond)
	token.Revoked = true
	token.UpdatedAt = time.Now()
	
	s.True(token.Revoked, "Token should be revoked after revocation")
	s.True(token.UpdatedAt.After(token.CreatedAt) || token.UpdatedAt.Equal(token.CreatedAt), "UpdatedAt should be after or equal to CreatedAt when revoked")
}

// TestRefreshToken_UUIDHandling tests UUID field handling
func (s *TokenModelTestSuite) TestRefreshToken_UUIDHandling() {
	// Test with valid UUIDs
	userID := uuid.New()
	tokenID := uuid.New()
	
	token := modelsv1.RefreshToken{
		ID:     tokenID,
		UserID: userID,
		Token:  "test_token",
	}
	
	s.NotEqual(uuid.Nil, token.ID, "Token ID should not be nil UUID")
	s.NotEqual(uuid.Nil, token.UserID, "User ID should not be nil UUID")
	s.NotEqual(token.ID, token.UserID, "Token ID and User ID should be different")

	// Test JSON serialization with UUIDs
	jsonData, err := json.Marshal(token)
	s.NoError(err, "Should marshal token with UUIDs")

	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	s.NoError(err, "Should unmarshal token with UUIDs")
	s.Equal(token.ID, unmarshaledToken.ID, "Token ID should be preserved")
	s.Equal(token.UserID, unmarshaledToken.UserID, "User ID should be preserved")
}

// TestRefreshToken_GormTags tests GORM struct tags
func (s *TokenModelTestSuite) TestRefreshToken_GormTags() {
	tokenType := reflect.TypeOf(modelsv1.RefreshToken{})
	
	// Check ID field GORM tags
	idField, found := tokenType.FieldByName("ID")
	s.True(found, "ID field should exist")
	gormTag := idField.Tag.Get("gorm")
	s.Contains(gormTag, "type:uuid", "ID should have UUID type")
	s.Contains(gormTag, "default:gen_random_uuid()", "ID should have default UUID generation")
	s.Contains(gormTag, "primaryKey", "ID should be primary key")
	
	// Check UserID field GORM tags
	userIDField, found := tokenType.FieldByName("UserID")
	s.True(found, "UserID field should exist")
	userIDGormTag := userIDField.Tag.Get("gorm")
	s.Contains(userIDGormTag, "type:uuid", "UserID should have UUID type")
	s.Contains(userIDGormTag, "index", "UserID should have index")
	
	// Check User field GORM tags
	userField, found := tokenType.FieldByName("User")
	s.True(found, "User field should exist")
	userGormTag := userField.Tag.Get("gorm")
	s.Contains(userGormTag, "foreignKey:UserID", "User should have foreign key relationship")
	
	// Check Token field GORM tags
	tokenField, found := tokenType.FieldByName("Token")
	s.True(found, "Token field should exist")
	tokenGormTag := tokenField.Tag.Get("gorm")
	s.Contains(tokenGormTag, "uniqueIndex", "Token should have unique index")
}

// TestRefreshToken_Relationships tests token relationships
func (s *TokenModelTestSuite) TestRefreshToken_Relationships() {
	userID := uuid.New()
	
	// Create user
	user := modelsv1.User{
		ID:       userID.String(),
		Name:     "John Doe",
		Email:    "john@example.com",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	// Create token with user relationship
	token := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		User:      user,
		Token:     "test_token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	
	// Verify relationship
	s.Equal(userID.String(), token.User.ID, "Token user ID should match user ID")
	s.Equal(user.Name, token.User.Name, "Token user name should match")
	s.Equal(user.Email, token.User.Email, "Token user email should match")
}

// TestRefreshToken_EdgeCases tests edge cases for refresh token
func (s *TokenModelTestSuite) TestRefreshToken_EdgeCases() {
	// Test token with zero time values
	token := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "test_token",
		ExpiresAt: time.Time{}, // Zero time
		Revoked:   false,
		CreatedAt: time.Time{}, // Zero time
		UpdatedAt: time.Time{}, // Zero time
	}
	
	s.True(token.ExpiresAt.IsZero(), "Zero expiration time should be zero")
	s.True(token.CreatedAt.IsZero(), "Zero created time should be zero")
	s.True(token.UpdatedAt.IsZero(), "Zero updated time should be zero")

	// Test JSON serialization with zero times
	jsonData, err := json.Marshal(token)
	s.NoError(err, "Should marshal token with zero times")

	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	s.NoError(err, "Should unmarshal token with zero times")
}

// TestRefreshToken_TokenUniqueness tests token uniqueness constraints
func (s *TokenModelTestSuite) TestRefreshToken_TokenUniqueness() {
	userID1 := uuid.New()
	userID2 := uuid.New()
	
	// Create two tokens with same token string (should be unique)
	token1 := modelsv1.RefreshToken{
		ID:     uuid.New(),
		UserID: userID1,
		Token:  "duplicate_token",
	}
	
	token2 := modelsv1.RefreshToken{
		ID:     uuid.New(),
		UserID: userID2,
		Token:  "duplicate_token", // Same token string
	}
	
	// In a real database, this would violate the unique constraint
	// Here we just verify the tokens have the same token string
	s.Equal(token1.Token, token2.Token, "Tokens should have same token string")
	s.NotEqual(token1.ID, token2.ID, "Tokens should have different IDs")
	s.NotEqual(token1.UserID, token2.UserID, "Tokens should have different user IDs")
}

// TestRefreshToken_BusinessRules tests refresh token business rules
func (s *TokenModelTestSuite) TestRefreshToken_BusinessRules() {
	now := time.Now()
	userID := uuid.New()
	
	// Test valid token business rules
	validToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "valid_business_token",
		ExpiresAt: now.Add(7 * 24 * time.Hour), // 7 days from now
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Business rule: Valid token should not be expired and not revoked
	isTokenValid := !validToken.Revoked && validToken.ExpiresAt.After(now)
	s.True(isTokenValid, "Valid token should pass business rules")
	
	// Business rule: Token should have reasonable expiration (not too far in future)
	maxExpiration := now.Add(30 * 24 * time.Hour) // 30 days max
	s.True(validToken.ExpiresAt.Before(maxExpiration), "Token expiration should be reasonable")
	
	// Business rule: Token should have been created before or at the same time as updated
	s.True(validToken.CreatedAt.Equal(validToken.UpdatedAt) || validToken.CreatedAt.Before(validToken.UpdatedAt),
		"CreatedAt should be before or equal to UpdatedAt")
}

// TestRefreshToken_ValidationScenarios tests various validation scenarios
func (s *TokenModelTestSuite) TestRefreshToken_ValidationScenarios() {
	now := time.Now()
	userID := uuid.New()
	
	testCases := []struct {
		name        string
		token       modelsv1.RefreshToken
		shouldBeValid bool
		description string
	}{
		{
			name: "ValidToken",
			token: modelsv1.RefreshToken{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "valid_token_123",
				ExpiresAt: now.Add(7 * 24 * time.Hour),
				Revoked:   false,
			},
			shouldBeValid: true,
			description:   "Valid token should pass validation",
		},
		{
			name: "ExpiredToken",
			token: modelsv1.RefreshToken{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "expired_token_123",
				ExpiresAt: now.Add(-1 * time.Hour),
				Revoked:   false,
			},
			shouldBeValid: false,
			description:   "Expired token should fail validation",
		},
		{
			name: "RevokedToken",
			token: modelsv1.RefreshToken{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "revoked_token_123",
				ExpiresAt: now.Add(7 * 24 * time.Hour),
				Revoked:   true,
			},
			shouldBeValid: false,
			description:   "Revoked token should fail validation",
		},
		{
			name: "EmptyToken",
			token: modelsv1.RefreshToken{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "",
				ExpiresAt: now.Add(7 * 24 * time.Hour),
				Revoked:   false,
			},
			shouldBeValid: false,
			description:   "Empty token string should fail validation",
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Simulate validation logic
			isValid := !tc.token.Revoked && 
					  tc.token.ExpiresAt.After(now) && 
					  tc.token.Token != "" &&
					  tc.token.ID != uuid.Nil &&
					  tc.token.UserID != uuid.Nil
			
			if tc.shouldBeValid {
				s.True(isValid, tc.description)
			} else {
				s.False(isValid, tc.description)
			}
		})
	}
}