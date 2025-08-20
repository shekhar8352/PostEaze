package testutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// TestJWTSecrets holds test JWT secrets
var (
	TestAccessSecret  = []byte("test-access-secret-key-for-testing-only")
	TestRefreshSecret = []byte("test-refresh-secret-key-for-testing-only")
)

// SetupTestJWTSecrets sets up test JWT secrets in environment variables
func SetupTestJWTSecrets() {
	os.Setenv("JWT_ACCESS_SECRET", string(TestAccessSecret))
	os.Setenv("JWT_REFRESH_SECRET", string(TestRefreshSecret))
}

// CleanupTestJWTSecrets removes test JWT secrets from environment variables
func CleanupTestJWTSecrets() {
	os.Unsetenv("JWT_ACCESS_SECRET")
	os.Unsetenv("JWT_REFRESH_SECRET")
}

// GenerateTestJWT generates a test JWT token with specified user ID and role
func GenerateTestJWT(userID, role string) (string, error) {
	claims := utils.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(TestAccessSecret)
}

// GenerateTestRefreshToken generates a test refresh token with specified user ID
func GenerateTestRefreshToken(userID string) (string, error) {
	claims := utils.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(TestRefreshSecret)
}

// GenerateExpiredTestJWT generates an expired test JWT token for testing expiration scenarios
func GenerateExpiredTestJWT(userID, role string) (string, error) {
	claims := utils.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(TestAccessSecret)
}

// GenerateInvalidTestJWT generates an invalid test JWT token signed with wrong secret
func GenerateInvalidTestJWT(userID, role string) (string, error) {
	claims := utils.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign with wrong secret to make it invalid
	wrongSecret := []byte("wrong-secret-key")
	return token.SignedString(wrongSecret)
}

// CreateAuthenticatedContext creates a Gin context with authentication headers set
func CreateAuthenticatedContext(method, url, userID, role string, body interface{}) (*gin.Context, *httptest.ResponseRecorder, error) {
	ctx, recorder := NewTestGinContext(method, url, body)
	
	token, err := GenerateTestJWT(userID, role)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate test JWT: %w", err)
	}
	
	SetAuthorizationHeader(ctx, token)
	
	// Set user context values that would normally be set by auth middleware
	ctx.Set("user_id", userID)
	ctx.Set("role", role)
	
	return ctx, recorder, nil
}

// CreateUnauthenticatedContext creates a Gin context without authentication headers
func CreateUnauthenticatedContext(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	return NewTestGinContext(method, url, body)
}

// MockAuthMiddleware returns a mock authentication middleware for testing
// This middleware simulates the behavior of the real auth middleware without actual JWT validation
func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "msg": "Unauthorized"})
			return
		}
		
		// For testing, we'll accept any Bearer token and extract mock user info
		if authHeader == "Bearer valid-token" {
			c.Set("user_id", "test-user-id")
			c.Set("role", "admin")
		} else if authHeader == "Bearer user-token" {
			c.Set("user_id", "test-user-id")
			c.Set("role", "creator")
		} else if authHeader == "Bearer expired-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "msg": "Invalid or expired token"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "msg": "Invalid token"})
			return
		}
		
		c.Next()
	}
}

// MockRequireRole returns a mock role-based authorization middleware for testing
func MockRequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No role in token"})
			return
		}
		
		userRole := roleIface.(string)
		for _, role := range allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}

// CreateTestUser creates a test user with specified parameters
func CreateTestUser(userType modelsv1.UserType) *modelsv1.User {
	return &modelsv1.User{
		ID:       "test-user-id",
		Name:     "Test User",
		Email:    "test@example.com",
		UserType: userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUserWithID creates a test user with a specific ID
func CreateTestUserWithID(userID string, userType modelsv1.UserType) *modelsv1.User {
	return &modelsv1.User{
		ID:       userID,
		Name:     fmt.Sprintf("Test User %s", userID),
		Email:    fmt.Sprintf("test-%s@example.com", userID),
		UserType: userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestTeam creates a test team with specified owner
func CreateTestTeam(ownerID string) *modelsv1.Team {
	return &modelsv1.Team{
		ID:      "test-team-id",
		Name:    "Test Team",
		OwnerID: ownerID,
		Owner:   *CreateTestUserWithID(ownerID, modelsv1.UserTypeIndividual),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestTeamMember creates a test team member
func CreateTestTeamMember(teamID, userID string, role modelsv1.Role) *modelsv1.TeamMember {
	return &modelsv1.TeamMember{
		ID:     "test-member-id",
		TeamID: teamID,
		UserID: userID,
		Role:   role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// TestAuthScenarios contains common authentication test scenarios
type TestAuthScenarios struct {
	ValidToken    string
	ExpiredToken  string
	InvalidToken  string
	UserID        string
	Role          string
}

// SetupTestAuthScenarios creates common authentication test scenarios
func SetupTestAuthScenarios(userID, role string) (*TestAuthScenarios, error) {
	SetupTestJWTSecrets()
	
	validToken, err := GenerateTestJWT(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate valid token: %w", err)
	}
	
	expiredToken, err := GenerateExpiredTestJWT(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate expired token: %w", err)
	}
	
	invalidToken, err := GenerateInvalidTestJWT(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invalid token: %w", err)
	}
	
	return &TestAuthScenarios{
		ValidToken:   validToken,
		ExpiredToken: expiredToken,
		InvalidToken: invalidToken,
		UserID:       userID,
		Role:         role,
	}, nil
}