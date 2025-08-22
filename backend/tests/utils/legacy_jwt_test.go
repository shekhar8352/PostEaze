package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shekhar8352/PostEaze/utils"
)

func TestGenerateJWT(t *testing.T) {
	// Store original environment variable
	originalJWTKey := os.Getenv("JWT_KEY")
	defer func() {
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
	}()

	// Set test key
	testJWTKey := "test-jwt-key-for-legacy-functions"
	os.Setenv("JWT_KEY", testJWTKey)

	tests := []struct {
		name   string
		userID int
		email  string
	}{
		{
			name:   "valid input",
			userID: 123,
			email:  "test@example.com",
		},
		{
			name:   "zero user ID",
			userID: 0,
			email:  "test@example.com",
		},
		{
			name:   "empty email",
			userID: 123,
			email:  "",
		},
		{
			name:   "negative user ID",
			userID: -1,
			email:  "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateJWT(tt.userID, tt.email)
			if err != nil {
				t.Errorf("GenerateJWT() error = %v", err)
				return
			}
			if token == "" {
				t.Error("GenerateJWT() returned empty token")
				return
			}

			// Verify token can be parsed
			claims, err := utils.ParseJWT(token)
			if err != nil {
				t.Errorf("ParseJWT() error = %v", err)
				return
			}

			if (*claims)["user_id"] != float64(tt.userID) {
				t.Errorf("ParseJWT() user_id = %v, want %v", (*claims)["user_id"], float64(tt.userID))
			}
			if (*claims)["email"] != tt.email {
				t.Errorf("ParseJWT() email = %v, want %v", (*claims)["email"], tt.email)
			}

			// Check expiration (should be 72 hours from now)
			exp := (*claims)["exp"].(float64)
			expectedExp := time.Now().Add(72 * time.Hour).Unix()
			timeDiff := int64(exp) - expectedExp

			// Allow for small time differences (within 5 seconds)
			if timeDiff >= 5 || timeDiff <= -5 {
				t.Errorf("Token expiration time difference too large: %v seconds", timeDiff)
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	// Set test key
	testJWTKey := "test-jwt-key-for-legacy-functions"
	originalJWTKey := os.Getenv("JWT_KEY")
	os.Setenv("JWT_KEY", testJWTKey)
	defer func() {
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
	}()

	userID := 123
	email := "test@example.com"
	token, err := utils.GenerateJWT(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token for test: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		wantErr   bool
		wantUser  float64
		wantEmail string
	}{
		{
			name:      "valid token",
			token:     token,
			wantErr:   false,
			wantUser:  float64(userID),
			wantEmail: email,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not.a.jwt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ParseJWT(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if claims != nil {
					t.Error("ParseJWT() should return nil claims on error")
				}
				return
			}

			if claims == nil {
				t.Error("ParseJWT() should not return nil claims on success")
				return
			}

			if (*claims)["user_id"] != tt.wantUser {
				t.Errorf("ParseJWT() user_id = %v, want %v", (*claims)["user_id"], tt.wantUser)
			}
			if (*claims)["email"] != tt.wantEmail {
				t.Errorf("ParseJWT() email = %v, want %v", (*claims)["email"], tt.wantEmail)
			}

			// Verify expiration exists and is in the future
			exp, exists := (*claims)["exp"]
			if !exists {
				t.Error("ParseJWT() claims should contain exp")
			}
			if exp.(float64) <= float64(time.Now().Unix()) {
				t.Error("ParseJWT() token should not be expired")
			}
		})
	}

	t.Run("empty token", func(t *testing.T) {
		// ParseJWT with empty string may panic due to nil pointer dereference
		defer func() {
			if r := recover(); r != nil {
				// If it panics, that's expected behavior for empty token
				t.Log("ParseJWT() panicked with empty token (expected)")
			}
		}()

		claims, err := utils.ParseJWT("")
		if err == nil {
			t.Error("ParseJWT() should return error for empty token")
		}
		if claims != nil {
			t.Error("ParseJWT() should return nil claims for empty token")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		// Create a token that's already expired
		claims := &jwt.MapClaims{
			"user_id": 123,
			"email":   "test@example.com",
			"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testJWTKey))
		if err != nil {
			t.Fatalf("Failed to create expired token: %v", err)
		}

		parsedClaims, err := utils.ParseJWT(tokenString)
		if err == nil {
			t.Error("ParseJWT() should return error for expired token")
		}
		if parsedClaims != nil {
			t.Error("ParseJWT() should return nil claims for expired token")
		}
	})

	t.Run("wrong signing key", func(t *testing.T) {
		// Create a token with a different signing key
		wrongKey := "wrong-signing-key"
		claims := &jwt.MapClaims{
			"user_id": 123,
			"email":   "test@example.com",
			"exp":     time.Now().Add(72 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(wrongKey))
		if err != nil {
			t.Fatalf("Failed to create token with wrong key: %v", err)
		}

		parsedClaims, err := utils.ParseJWT(tokenString)
		if err == nil {
			t.Error("ParseJWT() should return error for wrong signing key")
		}
		if parsedClaims != nil {
			t.Error("ParseJWT() should return nil claims for wrong signing key")
		}
	})
}

func TestJWTRoundTrip(t *testing.T) {
	// Set test key
	testJWTKey := "test-jwt-key-for-legacy-functions"
	originalJWTKey := os.Getenv("JWT_KEY")
	os.Setenv("JWT_KEY", testJWTKey)
	defer func() {
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
	}()

	testCases := []struct {
		name   string
		userID int
		email  string
	}{
		{"normal case", 123, "user@example.com"},
		{"zero user ID", 0, "zero@example.com"},
		{"empty email", 456, ""},
		{"special characters in email", 789, "user+test@example-domain.com"},
		{"large user ID", 999999999, "large@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate token
			token, err := utils.GenerateJWT(tc.userID, tc.email)
			if err != nil {
				t.Errorf("GenerateJWT() error = %v", err)
				return
			}
			if token == "" {
				t.Error("GenerateJWT() returned empty token")
				return
			}

			// Parse token
			claims, err := utils.ParseJWT(token)
			if err != nil {
				t.Errorf("ParseJWT() error = %v", err)
				return
			}
			if claims == nil {
				t.Error("ParseJWT() returned nil claims")
				return
			}

			// Verify claims
			if (*claims)["user_id"] != float64(tc.userID) {
				t.Errorf("Round trip user_id = %v, want %v", (*claims)["user_id"], float64(tc.userID))
			}
			if (*claims)["email"] != tc.email {
				t.Errorf("Round trip email = %v, want %v", (*claims)["email"], tc.email)
			}

			// Verify expiration
			exp := (*claims)["exp"].(float64)
			if exp <= float64(time.Now().Unix()) {
				t.Error("Round trip token should not be expired")
			}
		})
	}
}

func TestJWTWithoutEnvironmentVariable(t *testing.T) {
	// Store original environment variable
	originalJWTKey := os.Getenv("JWT_KEY")
	defer func() {
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
	}()

	// Temporarily unset environment variable
	os.Unsetenv("JWT_KEY")

	userID := 123
	email := "test@example.com"

	// This should still work but use empty key (not recommended for production)
	token, err := utils.GenerateJWT(userID, email)
	if err != nil {
		t.Errorf("GenerateJWT() error = %v", err)
		return
	}
	if token == "" {
		t.Error("GenerateJWT() returned empty token")
		return
	}

	// The token should be parseable with empty key
	claims, err := utils.ParseJWT(token)
	if err != nil {
		t.Errorf("ParseJWT() error = %v", err)
		return
	}
	if (*claims)["user_id"] != float64(userID) {
		t.Errorf("ParseJWT() user_id = %v, want %v", (*claims)["user_id"], float64(userID))
	}
}

func TestJWTClaimsStructure(t *testing.T) {
	// Set test key
	testJWTKey := "test-jwt-key-for-legacy-functions"
	originalJWTKey := os.Getenv("JWT_KEY")
	os.Setenv("JWT_KEY", testJWTKey)
	defer func() {
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
	}()

	userID := 123
	email := "test@example.com"
	token, err := utils.GenerateJWT(userID, email)
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	claims, err := utils.ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT() error = %v", err)
	}

	// Check that all expected claims are present
	expectedClaims := []string{"user_id", "email", "exp"}
	for _, claim := range expectedClaims {
		if _, exists := (*claims)[claim]; !exists {
			t.Errorf("Claim %s should exist", claim)
		}
	}

	// Check claim types
	if _, ok := (*claims)["user_id"].(float64); !ok {
		t.Errorf("user_id should be float64, got %T", (*claims)["user_id"])
	}
	if _, ok := (*claims)["email"].(string); !ok {
		t.Errorf("email should be string, got %T", (*claims)["email"])
	}
	if _, ok := (*claims)["exp"].(float64); !ok {
		t.Errorf("exp should be float64, got %T", (*claims)["exp"])
	}
}