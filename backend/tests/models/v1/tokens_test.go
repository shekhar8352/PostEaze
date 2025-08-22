package modelsv1_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func TestRefreshToken_JSONSerialization(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal refresh token to JSON without error: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to refresh token struct: %v", err)
	}

	if refreshToken.ID != unmarshaledToken.ID {
		t.Errorf("Expected ID %s, got %s", refreshToken.ID, unmarshaledToken.ID)
	}
	if refreshToken.UserID != unmarshaledToken.UserID {
		t.Errorf("Expected UserID %s, got %s", refreshToken.UserID, unmarshaledToken.UserID)
	}
	if refreshToken.Token != unmarshaledToken.Token {
		t.Errorf("Expected Token %s, got %s", refreshToken.Token, unmarshaledToken.Token)
	}
	if refreshToken.ExpiresAt.Unix() != unmarshaledToken.ExpiresAt.Unix() {
		t.Errorf("Expected ExpiresAt %v, got %v", refreshToken.ExpiresAt, unmarshaledToken.ExpiresAt)
	}
	if refreshToken.Revoked != unmarshaledToken.Revoked {
		t.Errorf("Expected Revoked %t, got %t", refreshToken.Revoked, unmarshaledToken.Revoked)
	}
}

func TestRefreshToken_WithUser(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Should marshal refresh token with user: %v", err)
	}

	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	if err != nil {
		t.Fatalf("Should unmarshal JSON to refresh token with user: %v", err)
	}

	if refreshToken.User.ID != unmarshaledToken.User.ID {
		t.Errorf("Expected User.ID %s, got %s", refreshToken.User.ID, unmarshaledToken.User.ID)
	}
	if refreshToken.User.Name != unmarshaledToken.User.Name {
		t.Errorf("Expected User.Name %s, got %s", refreshToken.User.Name, unmarshaledToken.User.Name)
	}
	if refreshToken.User.Email != unmarshaledToken.User.Email {
		t.Errorf("Expected User.Email %s, got %s", refreshToken.User.Email, unmarshaledToken.User.Email)
	}
}

func TestRefreshToken_ExpirationLogic(t *testing.T) {
	now := time.Now()

	// Test valid (non-expired) token
	validToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "valid_token",
		ExpiresAt: now.Add(7 * 24 * time.Hour), // Expires in 7 days
		Revoked:   false,
	}

	if !validToken.ExpiresAt.After(now) {
		t.Errorf("Valid token should not be expired")
	}
	if validToken.Revoked {
		t.Errorf("Valid token should not be revoked")
	}

	// Test expired token
	expiredToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "expired_token",
		ExpiresAt: now.Add(-1 * time.Hour), // Expired 1 hour ago
		Revoked:   false,
	}

	if !expiredToken.ExpiresAt.Before(now) {
		t.Errorf("Expired token should be expired")
	}
	if expiredToken.Revoked {
		t.Errorf("Expired token should not be revoked (just expired)")
	}

	// Test revoked token
	revokedToken := modelsv1.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "revoked_token",
		ExpiresAt: now.Add(7 * 24 * time.Hour), // Valid expiration but revoked
		Revoked:   true,
	}

	if !revokedToken.ExpiresAt.After(now) {
		t.Errorf("Revoked token expiration should still be valid")
	}
	if !revokedToken.Revoked {
		t.Errorf("Revoked token should be marked as revoked")
	}
}

func TestRefreshToken_RevocationLogic(t *testing.T) {
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

	if token.Revoked {
		t.Errorf("Token should not be revoked initially")
	}

	// Simulate revocation with a small delay to ensure UpdatedAt is after CreatedAt
	time.Sleep(1 * time.Millisecond)
	token.Revoked = true
	token.UpdatedAt = time.Now()

	if !token.Revoked {
		t.Errorf("Token should be revoked after revocation")
	}
	if !(token.UpdatedAt.After(token.CreatedAt) || token.UpdatedAt.Equal(token.CreatedAt)) {
		t.Errorf("UpdatedAt should be after or equal to CreatedAt when revoked")
	}
}

func TestRefreshToken_UUIDHandling(t *testing.T) {
	// Test with valid UUIDs
	userID := uuid.New()
	tokenID := uuid.New()

	token := modelsv1.RefreshToken{
		ID:     tokenID,
		UserID: userID,
		Token:  "test_token",
	}

	if token.ID == uuid.Nil {
		t.Errorf("Token ID should not be nil UUID")
	}
	if token.UserID == uuid.Nil {
		t.Errorf("User ID should not be nil UUID")
	}
	if token.ID == token.UserID {
		t.Errorf("Token ID and User ID should be different")
	}

	// Test JSON serialization with UUIDs
	jsonData, err := json.Marshal(token)
	if err != nil {
		t.Fatalf("Should marshal token with UUIDs: %v", err)
	}

	var unmarshaledToken modelsv1.RefreshToken
	err = json.Unmarshal(jsonData, &unmarshaledToken)
	if err != nil {
		t.Fatalf("Should unmarshal token with UUIDs: %v", err)
	}
	if token.ID != unmarshaledToken.ID {
		t.Errorf("Token ID should be preserved: expected %s, got %s", token.ID, unmarshaledToken.ID)
	}
	if token.UserID != unmarshaledToken.UserID {
		t.Errorf("User ID should be preserved: expected %s, got %s", token.UserID, unmarshaledToken.UserID)
	}
}

func TestRefreshToken_GormTags(t *testing.T) {
	tokenType := reflect.TypeOf(modelsv1.RefreshToken{})

	// Check ID field GORM tags
	idField, found := tokenType.FieldByName("ID")
	if !found {
		t.Fatalf("ID field should exist")
	}
	gormTag := idField.Tag.Get("gorm")
	if gormTag == "" {
		t.Errorf("ID field should have GORM tags")
	}
	// Note: We can't easily test the exact content of GORM tags in a unit test
	// without more complex reflection, but we can verify they exist

	// Check UserID field GORM tags
	userIDField, found := tokenType.FieldByName("UserID")
	if !found {
		t.Fatalf("UserID field should exist")
	}
	userIDGormTag := userIDField.Tag.Get("gorm")
	if userIDGormTag == "" {
		t.Errorf("UserID field should have GORM tags")
	}

	// Check User field GORM tags
	userField, found := tokenType.FieldByName("User")
	if !found {
		t.Fatalf("User field should exist")
	}
	userGormTag := userField.Tag.Get("gorm")
	if userGormTag == "" {
		t.Errorf("User field should have GORM tags")
	}

	// Check Token field GORM tags
	tokenField, found := tokenType.FieldByName("Token")
	if !found {
		t.Fatalf("Token field should exist")
	}
	tokenGormTag := tokenField.Tag.Get("gorm")
	if tokenGormTag == "" {
		t.Errorf("Token field should have GORM tags")
	}
}

func TestRefreshToken_Relationships(t *testing.T) {
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
	if token.User.ID != userID.String() {
		t.Errorf("Token user ID should match user ID: expected %s, got %s", userID.String(), token.User.ID)
	}
	if token.User.Name != user.Name {
		t.Errorf("Token user name should match: expected %s, got %s", user.Name, token.User.Name)
	}
	if token.User.Email != user.Email {
		t.Errorf("Token user email should match: expected %s, got %s", user.Email, token.User.Email)
	}
}

func TestRefreshToken_ValidationScenarios(t *testing.T) {
	now := time.Now()
	userID := uuid.New()

	testCases := []struct {
		name          string
		token         modelsv1.RefreshToken
		shouldBeValid bool
		description   string
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
		t.Run(tc.name, func(t *testing.T) {
			// Simulate validation logic
			isValid := !tc.token.Revoked &&
				tc.token.ExpiresAt.After(now) &&
				tc.token.Token != "" &&
				tc.token.ID != uuid.Nil &&
				tc.token.UserID != uuid.Nil

			if tc.shouldBeValid && !isValid {
				t.Errorf("%s: expected valid but got invalid", tc.description)
			}
			if !tc.shouldBeValid && isValid {
				t.Errorf("%s: expected invalid but got valid", tc.description)
			}
		})
	}
}

func TestRefreshToken_BusinessRules(t *testing.T) {
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
	if !isTokenValid {
		t.Errorf("Valid token should pass business rules")
	}

	// Business rule: Token should have reasonable expiration (not too far in future)
	maxExpiration := now.Add(30 * 24 * time.Hour) // 30 days max
	if !validToken.ExpiresAt.Before(maxExpiration) {
		t.Errorf("Token expiration should be reasonable")
	}

	// Business rule: Token should have been created before or at the same time as updated
	if !(validToken.CreatedAt.Equal(validToken.UpdatedAt) || validToken.CreatedAt.Before(validToken.UpdatedAt)) {
		t.Errorf("CreatedAt should be before or equal to UpdatedAt")
	}
}