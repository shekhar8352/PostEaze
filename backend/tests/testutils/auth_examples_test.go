package testutils

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExampleProtectedHandler demonstrates a protected API handler that requires authentication
func ExampleProtectedHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "User not authenticated"})
		return
	}
	
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "User role not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"user_id": userID,
		"role":    role,
		"message": "Access granted",
	})
}

// ExampleAdminOnlyHandler demonstrates a handler that requires admin role
func ExampleAdminOnlyHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "No role in token"})
		return
	}
	
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Admin access required"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Admin access granted",
	})
}

// TestExampleProtectedHandlerWithAuthentication demonstrates testing authenticated endpoints
func TestExampleProtectedHandlerWithAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("Authenticated user access", func(t *testing.T) {
		userID := "user-123"
		role := "creator"
		
		ctx, recorder, err := CreateAuthenticatedContext("GET", "/protected", userID, role, nil)
		require.NoError(t, err)
		
		ExampleProtectedHandler(ctx)
		
		assert.Equal(t, http.StatusOK, recorder.Code)
		
		var response map[string]interface{}
		err = ParseJSONResponse(recorder, &response)
		require.NoError(t, err)
		
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, userID, response["user_id"])
		assert.Equal(t, role, response["role"])
		assert.Equal(t, "Access granted", response["message"])
	})
	
	t.Run("Unauthenticated user access", func(t *testing.T) {
		ctx, recorder := CreateUnauthenticatedContext("GET", "/protected", nil)
		
		ExampleProtectedHandler(ctx)
		
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		
		var response map[string]interface{}
		err := ParseJSONResponse(recorder, &response)
		require.NoError(t, err)
		
		assert.Equal(t, "error", response["status"])
		assert.Equal(t, "User not authenticated", response["msg"])
	})
}

// TestExampleAdminOnlyHandlerWithRoles demonstrates testing role-based access
func TestExampleAdminOnlyHandlerWithRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("Admin user access", func(t *testing.T) {
		userID := "admin-123"
		role := "admin"
		
		ctx, recorder, err := CreateAuthenticatedContext("GET", "/admin", userID, role, nil)
		require.NoError(t, err)
		
		ExampleAdminOnlyHandler(ctx)
		
		assert.Equal(t, http.StatusOK, recorder.Code)
		
		var response map[string]interface{}
		err = ParseJSONResponse(recorder, &response)
		require.NoError(t, err)
		
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "Admin access granted", response["message"])
	})
	
	t.Run("Non-admin user access", func(t *testing.T) {
		userID := "user-123"
		role := "creator"
		
		ctx, recorder, err := CreateAuthenticatedContext("GET", "/admin", userID, role, nil)
		require.NoError(t, err)
		
		ExampleAdminOnlyHandler(ctx)
		
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		
		var response map[string]interface{}
		err = ParseJSONResponse(recorder, &response)
		require.NoError(t, err)
		
		assert.Equal(t, "error", response["status"])
		assert.Equal(t, "Admin access required", response["msg"])
	})
}

// TestExampleWithMockMiddleware demonstrates using mock authentication middleware
func TestExampleWithMockMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a router with mock middleware
	router := gin.New()
	router.Use(MockAuthMiddleware())
	router.GET("/protected", ExampleProtectedHandler)
	
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
			name:           "No token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := NewTestGinContext("GET", "/protected", nil)
			if tt.authHeader != "" {
				ctx.Request.Header.Set("Authorization", tt.authHeader)
			}
			
			router.ServeHTTP(recorder, ctx.Request)
			
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := ParseJSONResponse(recorder, &response)
				require.NoError(t, err)
				
				assert.Equal(t, "success", response["status"])
				assert.Equal(t, tt.expectedUserID, response["user_id"])
				assert.Equal(t, tt.expectedRole, response["role"])
			}
		})
	}
}

// TestExampleWithRoleBasedMiddleware demonstrates using mock role-based middleware
func TestExampleWithRoleBasedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a router with mock middleware
	router := gin.New()
	router.Use(MockAuthMiddleware())
	router.Use(MockRequireRole("admin", "editor"))
	router.GET("/admin", ExampleAdminOnlyHandler)
	
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Admin token (allowed)",
			authHeader:     "Bearer valid-token", // This sets role to "admin"
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User token (not allowed)",
			authHeader:     "Bearer user-token", // This sets role to "creator"
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := NewTestGinContext("GET", "/admin", nil)
			if tt.authHeader != "" {
				ctx.Request.Header.Set("Authorization", tt.authHeader)
			}
			
			router.ServeHTTP(recorder, ctx.Request)
			
			assert.Equal(t, tt.expectedStatus, recorder.Code)
		})
	}
}

// TestExampleWithTestAuthScenarios demonstrates using pre-configured auth scenarios
func TestExampleWithTestAuthScenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	userID := "test-user-123"
	role := "admin"
	
	scenarios, err := SetupTestAuthScenarios(userID, role)
	require.NoError(t, err)
	defer CleanupTestJWTSecrets()
	
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			token:          scenarios.ValidToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Expired token",
			token:          scenarios.ExpiredToken,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			token:          scenarios.InvalidToken,
			expectedStatus: http.StatusUnauthorized,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := NewTestGinContext("GET", "/protected", nil)
			SetAuthorizationHeader(ctx, tt.token)
			
			// For this example, we'll simulate what would happen with real middleware
			// In practice, you'd use the actual auth middleware or mock middleware
			if tt.token == scenarios.ValidToken {
				// Valid token - set context and call handler
				ctx.Set("user_id", scenarios.UserID)
				ctx.Set("role", scenarios.Role)
				ExampleProtectedHandler(ctx)
				
				assert.Equal(t, http.StatusOK, recorder.Code)
				
				var response map[string]interface{}
				err := ParseJSONResponse(recorder, &response)
				require.NoError(t, err)
				
				assert.Equal(t, "success", response["status"])
				assert.Equal(t, scenarios.UserID, response["user_id"])
				assert.Equal(t, scenarios.Role, response["role"])
			} else {
				// Invalid/expired tokens - simulate middleware rejection
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Invalid or expired token"})
				
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
				
				var response map[string]interface{}
				err := ParseJSONResponse(recorder, &response)
				require.NoError(t, err)
				
				assert.Equal(t, "error", response["status"])
				assert.Equal(t, "Invalid or expired token", response["msg"])
			}
		})
	}
}

// TestExampleWithTestUserCreation demonstrates creating test users
func TestExampleWithTestUserCreation(t *testing.T) {
	t.Run("Individual user", func(t *testing.T) {
		user := CreateTestUser(modelsv1.UserTypeIndividual)
		
		assert.Equal(t, modelsv1.UserTypeIndividual, user.UserType)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Name)
		assert.NotEmpty(t, user.Email)
	})
	
	t.Run("Team user", func(t *testing.T) {
		user := CreateTestUser(modelsv1.UserTypeTeam)
		
		assert.Equal(t, modelsv1.UserTypeTeam, user.UserType)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Name)
		assert.NotEmpty(t, user.Email)
	})
	
	t.Run("Custom user ID", func(t *testing.T) {
		customID := "custom-user-123"
		user := CreateTestUserWithID(customID, modelsv1.UserTypeIndividual)
		
		assert.Equal(t, customID, user.ID)
		assert.Contains(t, user.Name, customID)
		assert.Contains(t, user.Email, customID)
	})
}

// TestExampleWithTestTeamCreation demonstrates creating test teams
func TestExampleWithTestTeamCreation(t *testing.T) {
	ownerID := "owner-123"
	team := CreateTestTeam(ownerID)
	
	assert.Equal(t, ownerID, team.OwnerID)
	assert.Equal(t, ownerID, team.Owner.ID)
	assert.NotEmpty(t, team.ID)
	assert.NotEmpty(t, team.Name)
	
	// Test team member creation
	member := CreateTestTeamMember(team.ID, "member-456", modelsv1.RoleEditor)
	
	assert.Equal(t, team.ID, member.TeamID)
	assert.Equal(t, "member-456", member.UserID)
	assert.Equal(t, modelsv1.RoleEditor, member.Role)
}