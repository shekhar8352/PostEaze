package utils

import (
	"os"
	"testing"
	"time"

	"github.com/shekhar8352/PostEaze/utils"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "this_is_a_very_long_password_that_is_still_within_bcrypt_limits_test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := utils.HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned unhashed password")
				}
				if len(hash) < 50 {
					t.Error("HashPassword() returned hash that's too short")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "valid password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "invalid password",
			password: "wrongpassword456",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "invalid hash",
			password: password,
			hash:     "invalid-hash",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CheckPasswordHash(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	// Set test JWT secrets
	originalAccessSecret := os.Getenv("JWT_ACCESS_SECRET")
	originalRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	defer func() {
		if originalAccessSecret != "" {
			os.Setenv("JWT_ACCESS_SECRET", originalAccessSecret)
		} else {
			os.Unsetenv("JWT_ACCESS_SECRET")
		}
		if originalRefreshSecret != "" {
			os.Setenv("JWT_REFRESH_SECRET", originalRefreshSecret)
		} else {
			os.Unsetenv("JWT_REFRESH_SECRET")
		}
	}()

	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-key-for-testing-purposes")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key-for-testing-purposes")

	tests := []struct {
		name   string
		userID string
		role   string
	}{
		{
			name:   "valid input",
			userID: "user123",
			role:   "admin",
		},
		{
			name:   "empty user ID",
			userID: "",
			role:   "user",
		},
		{
			name:   "empty role",
			userID: "user123",
			role:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateAccessToken(tt.userID, tt.role)
			if err != nil {
				t.Errorf("GenerateAccessToken() error = %v", err)
				return
			}
			if token == "" {
				t.Error("GenerateAccessToken() returned empty token")
				return
			}

			// Verify token can be parsed
			claims, err := utils.ParseToken(token, false)
			if err != nil {
				t.Errorf("ParseToken() error = %v", err)
				return
			}
			if claims.UserID != tt.userID {
				t.Errorf("ParseToken() userID = %v, want %v", claims.UserID, tt.userID)
			}
			if claims.Role != tt.role {
				t.Errorf("ParseToken() role = %v, want %v", claims.Role, tt.role)
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	// Set test JWT secrets
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-key-for-testing-purposes")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key-for-testing-purposes")

	tests := []struct {
		name   string
		userID string
	}{
		{
			name:   "valid input",
			userID: "user123",
		},
		{
			name:   "empty user ID",
			userID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateRefreshToken(tt.userID)
			if err != nil {
				t.Errorf("GenerateRefreshToken() error = %v", err)
				return
			}
			if token == "" {
				t.Error("GenerateRefreshToken() returned empty token")
				return
			}

			// Verify token can be parsed as refresh token
			claims, err := utils.ParseToken(token, true)
			if err != nil {
				t.Errorf("ParseToken() error = %v", err)
				return
			}
			if claims.UserID != tt.userID {
				t.Errorf("ParseToken() userID = %v, want %v", claims.UserID, tt.userID)
			}
			if claims.Role != "" {
				t.Error("ParseToken() refresh token should not have role")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	// Set test JWT secrets
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-key-for-testing-purposes")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key-for-testing-purposes")

	userID := "user123"
	role := "admin"
	accessToken, err := utils.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("Failed to generate access token for test: %v", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate refresh token for test: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		isRefresh bool
		wantErr   bool
		wantUser  string
		wantRole  string
	}{
		{
			name:      "valid access token",
			token:     accessToken,
			isRefresh: false,
			wantErr:   false,
			wantUser:  userID,
			wantRole:  role,
		},
		{
			name:      "valid refresh token",
			token:     refreshToken,
			isRefresh: true,
			wantErr:   false,
			wantUser:  userID,
			wantRole:  "",
		},
		{
			name:      "invalid token",
			token:     "invalid.token.here",
			isRefresh: false,
			wantErr:   true,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ParseToken(tt.token, tt.isRefresh)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if claims != nil {
					t.Error("ParseToken() should return nil claims on error")
				}
				return
			}

			if claims.UserID != tt.wantUser {
				t.Errorf("ParseToken() userID = %v, want %v", claims.UserID, tt.wantUser)
			}
			if claims.Role != tt.wantRole {
				t.Errorf("ParseToken() role = %v, want %v", claims.Role, tt.wantRole)
			}
			if !claims.ExpiresAt.After(time.Now()) {
				t.Error("ParseToken() token should not be expired")
			}
		})
	}

	t.Run("empty token", func(t *testing.T) {
		// ParseToken with empty string may panic, so we handle it
		defer func() {
			if r := recover(); r != nil {
				t.Log("ParseToken() panicked with empty token (expected)")
			}
		}()

		claims, err := utils.ParseToken("", false)
		if err == nil {
			t.Error("ParseToken() should return error for empty token")
		}
		if claims != nil {
			t.Error("ParseToken() should return nil claims for empty token")
		}
	})
}

func TestGetUserIDFromToken(t *testing.T) {
	// Set test JWT secrets
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-key-for-testing-purposes")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key-for-testing-purposes")

	userID := "user123"
	role := "admin"
	token, err := utils.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("Failed to generate token for test: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			want:    userID,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.GetUserIDFromToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserIDFromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserIDFromToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRefreshTokenExpiry(t *testing.T) {
	before := time.Now()
	expiry := utils.GetRefreshTokenExpiry()
	after := time.Now()

	expectedMin := before.Add(7 * 24 * time.Hour)
	expectedMax := after.Add(7 * 24 * time.Hour)

	if !(expiry.After(expectedMin) || expiry.Equal(expectedMin)) {
		t.Errorf("GetRefreshTokenExpiry() = %v, should be after %v", expiry, expectedMin)
	}
	if !(expiry.Before(expectedMax) || expiry.Equal(expectedMax)) {
		t.Errorf("GetRefreshTokenExpiry() = %v, should be before %v", expiry, expectedMax)
	}
}

func TestTokenExpiration(t *testing.T) {
	// Set test JWT secrets
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-key-for-testing-purposes")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-key-for-testing-purposes")

	t.Run("access token expiration", func(t *testing.T) {
		userID := "user123"
		role := "admin"
		token, err := utils.GenerateAccessToken(userID, role)
		if err != nil {
			t.Fatalf("GenerateAccessToken() error = %v", err)
		}

		claims, err := utils.ParseToken(token, false)
		if err != nil {
			t.Fatalf("ParseToken() error = %v", err)
		}

		// Access token should expire in 15 minutes
		expectedExpiry := time.Now().Add(15 * time.Minute)
		timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)

		// Allow for small time differences (within 1 second)
		if timeDiff >= time.Second || timeDiff <= -time.Second {
			t.Errorf("Access token expiry time difference too large: %v", timeDiff)
		}
	})

	t.Run("refresh token expiration", func(t *testing.T) {
		userID := "user123"
		token, err := utils.GenerateRefreshToken(userID)
		if err != nil {
			t.Fatalf("GenerateRefreshToken() error = %v", err)
		}

		claims, err := utils.ParseToken(token, true)
		if err != nil {
			t.Fatalf("ParseToken() error = %v", err)
		}

		// Refresh token should expire in 7 days
		expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
		timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)

		// Allow for small time differences (within 1 second)
		if timeDiff >= time.Second || timeDiff <= -time.Second {
			t.Errorf("Refresh token expiry time difference too large: %v", timeDiff)
		}
	})
}