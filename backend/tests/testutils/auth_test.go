package testutils

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestJWTSecrets(t *testing.T) {
	// Clean up any existing secrets
	CleanupTestJWTSecrets()
	
	// Test setup
	SetupTestJWTSecrets()
	
	assert.Equal(t, string(TestAccessSecret), os.Getenv("JWT_ACCESS_SECRET"))
	assert.Equal(t, string(TestRefreshSecret), os.Getenv("JWT_REFRESH_SECRET"))
	
	// Clean up
	CleanupTestJWTSecrets()
}

func TestCleanupTestJWTSecrets(t *testing.T) {
	// Set up secrets first
	SetupTestJWTSecrets()
	
	// Test cleanup
	CleanupTestJWTSecrets()
	
	assert.Empty(t, os.Getenv("JWT_ACCESS_SECRET"))
	assert.Empty(t, os.Getenv("JWT_REFRESH_SECRET"))
}

func TestGenerateTestJWT(t *testing.T) {
	SetupTestJWTSecrets()
	defer CleanupTestJWTSecrets()
	
	userID := "test-user-123"
	role := "admin"
	
	token, err := GenerateTestJWT(userID, role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Verify the token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestAccessSecret, nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	
	claims, ok := parsedToken.Claims.(*utils.JWTClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestGenerateTestRefreshToken(t *testing.T) {
	SetupTestJWTSecrets()
	defer CleanupTestJWTSecrets()
	
	userID := "test-user-123"
	
	token, err := GenerateTestRefreshToken(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Verify the token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestRefreshSecret, nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	
	claims, ok := parsedToken.Claims.(*utils.JWTClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Empty(t, claims.Role) // Refresh tokens don't have roles
	assert.True(t, claims.ExpiresAt.After(time.Now().Add(6*24*time.Hour))) // Should expire in ~7 days
}

func TestGenerateExpiredTestJWT(t *testing.T) {
	SetupTestJWTSecrets()
	defer CleanupTestJWTSecrets()
	
	userID := "test-user-123"
	role := "admin"
	
	token, err := GenerateExpiredTestJWT(userID, role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Verify the token is expired
	parsedToken, err := jwt.ParseWithClaims(token, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestAccessSecret, nil
	})
	
	// Token should be parseable but invalid due to expiration
	assert.Error(t, err)
	assert.False(t, parsedToken.Valid)
	
	claims, ok := parsedToken.Claims.(*utils.JWTClaims)
	require.True(t, ok)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.True(t, claims.ExpiresAt.Before(time.Now()))
}

func TestGenerateInvalidTestJWT(t *testing.T) {
	SetupTestJWTSecrets()
	defer CleanupTestJWTSecrets()
	
	userID := "test-user-123"
	role := "admin"
	
	token, err := GenerateInvalidTestJWT(userID, role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Verify the token cannot be parsed with the correct secret
	parsedToken, err := jwt.ParseWithClaims(token, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestAccessSecret, nil
	})
	
	// Token should fail to parse due to wrong signature
	assert.Error(t, err)
	assert.False(t, parsedToken.Valid)
}

func TestCreateAuthenticatedContext(t *testing.T) {
	SetupTestJWTSecrets()
	defer CleanupTestJWTSecrets()
	
	userID := "test-user-123"
	role := "admin"
	
	ctx, recorder, err := CreateAuthenticatedContext("GET", "/test", userID, role, nil)
	require.NoError(t, err)
	assert.NotNil(t, ctx)
	assert.NotNil(t, recorder)
	
	// Verify authorization header is set
	authHeader := ctx.Request.Header.Get("Authorization")
	assert.True(t, len(authHeader) > 7) // "Bearer " + token
	assert.Contains(t, authHeader, "Bearer ")
	
	// Verify context values are set
	contextUserID, exists := ctx.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, userID, contextUserID)
	
	contextRole, exists := ctx.Get("role")
	assert.True(t, exists)
	assert.Equal(t, role, contextRole)
}

func TestCreateUnauthenticatedContext(t *testing.T) {
	ctx, recorder := CreateUnauthenticatedContext("GET", "/test", nil)
	assert.NotNil(t, ctx)
	assert.NotNil(t, recorder)
	
	// Verify no authorization header is set
	authHeader := ctx.Request.Header.Get("Authorization")
	assert.Empty(t, authHeader)
	
	// Verify no context values are set
	_, exists := ctx.Get("user_id")
	assert.False(t, exists)
	
	_, exists = ctx.Get("role")
	assert.False(t, exists)
}

func TestMockAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUserID string
		expectedRole   string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer valid-token",
			expectedStatus: http.StatusOK,
			expectedUserID: "test-user-id",
			expectedRole:   "admin",
		},
		{
			name:           "User token",
			authHeader:     "Bearer user-token",
			expectedStatus: http.StatusOK,
			expectedUserID: "test-user-id",
			expectedRole:   "creator",
		},
		{
			name:           "Expired token",
			authHeader:     "Bearer expired-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := NewTestGinContext("GET", "/test", nil)
			if tt.authHeader != "" {
				ctx.Request.Header.Set("Authorization", tt.authHeader)
			}
			
			// Create a test handler that uses the mock middleware
			handler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
			
			// Apply middleware and handler
			MockAuthMiddleware()(ctx)
			if !ctx.IsAborted() {
				handler(ctx)
			}
			
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			
			if tt.expectedStatus == http.StatusOK {
				userID, exists := ctx.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, tt.expectedUserID, userID)
				
				role, exists := ctx.Get("role")
				assert.True(t, exists)
				assert.Equal(t, tt.expectedRole, role)
			}
		})
	}
}

func TestMockRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		userRole       string
		allowedRoles   []string
		expectedStatus int
	}{
		{
			name:           "Admin role allowed",
			userRole:       "admin",
			allowedRoles:   []string{"admin", "editor"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Editor role allowed",
			userRole:       "editor",
			allowedRoles:   []string{"admin", "editor"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Creator role not allowed",
			userRole:       "creator",
			allowedRoles:   []string{"admin", "editor"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No role in context",
			userRole:       "",
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := NewTestGinContext("GET", "/test", nil)
			
			if tt.userRole != "" {
				ctx.Set("role", tt.userRole)
			}
			
			// Create a test handler that uses the mock role middleware
			handler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			}
			
			// Apply middleware and handler
			MockRequireRole(tt.allowedRoles...)(ctx)
			if !ctx.IsAborted() {
				handler(ctx)
			}
			
			assert.Equal(t, tt.expectedStatus, recorder.Code)
		})
	}
}

func TestCreateTestUser(t *testing.T) {
	user := CreateTestUser(modelsv1.UserTypeIndividual)
	
	assert.NotNil(t, user)
	assert.Equal(t, "test-user-id", user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, modelsv1.UserTypeIndividual, user.UserType)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestCreateTestUserWithID(t *testing.T) {
	userID := "custom-user-id"
	user := CreateTestUserWithID(userID, modelsv1.UserTypeTeam)
	
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "Test User custom-user-id", user.Name)
	assert.Equal(t, "test-custom-user-id@example.com", user.Email)
	assert.Equal(t, modelsv1.UserTypeTeam, user.UserType)
}

func TestCreateTestTeam(t *testing.T) {
	ownerID := "owner-123"
	team := CreateTestTeam(ownerID)
	
	assert.NotNil(t, team)
	assert.Equal(t, "test-team-id", team.ID)
	assert.Equal(t, "Test Team", team.Name)
	assert.Equal(t, ownerID, team.OwnerID)
	assert.Equal(t, ownerID, team.Owner.ID)
	assert.False(t, team.CreatedAt.IsZero())
	assert.False(t, team.UpdatedAt.IsZero())
}

func TestCreateTestTeamMember(t *testing.T) {
	teamID := "team-123"
	userID := "user-456"
	role := modelsv1.RoleEditor
	
	member := CreateTestTeamMember(teamID, userID, role)
	
	assert.NotNil(t, member)
	assert.Equal(t, "test-member-id", member.ID)
	assert.Equal(t, teamID, member.TeamID)
	assert.Equal(t, userID, member.UserID)
	assert.Equal(t, role, member.Role)
	assert.False(t, member.CreatedAt.IsZero())
	assert.False(t, member.UpdatedAt.IsZero())
}

func TestSetupTestAuthScenarios(t *testing.T) {
	userID := "test-user-123"
	role := "admin"
	
	scenarios, err := SetupTestAuthScenarios(userID, role)
	require.NoError(t, err)
	assert.NotNil(t, scenarios)
	
	assert.NotEmpty(t, scenarios.ValidToken)
	assert.NotEmpty(t, scenarios.ExpiredToken)
	assert.NotEmpty(t, scenarios.InvalidToken)
	assert.Equal(t, userID, scenarios.UserID)
	assert.Equal(t, role, scenarios.Role)
	
	// Verify valid token works
	parsedToken, err := jwt.ParseWithClaims(scenarios.ValidToken, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestAccessSecret, nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	
	// Verify expired token is expired
	expiredToken, err := jwt.ParseWithClaims(scenarios.ExpiredToken, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestAccessSecret, nil
	})
	assert.Error(t, err)
	assert.False(t, expiredToken.Valid)
	
	// Verify invalid token fails with correct secret
	invalidToken, err := jwt.ParseWithClaims(scenarios.InvalidToken, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TestAccessSecret, nil
	})
	assert.Error(t, err)
	assert.False(t, invalidToken.Valid)
	
	CleanupTestJWTSecrets()
}