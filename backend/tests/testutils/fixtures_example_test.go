package testutils

import (
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/stretchr/testify/assert"
)

// Example test demonstrating how to use predefined fixtures
func TestExampleUsingPredefinedFixtures(t *testing.T) {
	t.Run("using predefined user fixtures", func(t *testing.T) {
		// Get a predefined user
		user := TestUsers[0]
		
		// Convert to model
		userModel := user.ToUserModel()
		
		// Use in your tests
		assert.Equal(t, "user-1", userModel.ID)
		assert.Equal(t, "John Doe", userModel.Name)
		assert.Equal(t, modelsv1.UserTypeIndividual, userModel.UserType)
	})

	t.Run("using predefined team fixtures", func(t *testing.T) {
		// Get a predefined team
		team := TestTeams[0]
		
		// Convert to model
		teamModel := team.ToTeamModel()
		
		// Use in your tests
		assert.Equal(t, "team-1", teamModel.ID)
		assert.Equal(t, "Development Team", teamModel.Name)
		assert.Equal(t, "user-2", teamModel.OwnerID)
	})

	t.Run("using predefined log fixtures", func(t *testing.T) {
		// Get a predefined log entry
		logEntry := TestLogs[0]
		
		// Convert to model
		logModel := logEntry.ToLogEntryModel()
		
		// Use in your tests
		assert.Equal(t, "log-1", logModel.LogID)
		assert.Equal(t, "INFO", logModel.Level)
		assert.Equal(t, "User login successful", logModel.Message)
		assert.Equal(t, 200, logModel.Status)
	})
}

// Example test demonstrating how to create custom fixtures
func TestExampleUsingCustomFixtures(t *testing.T) {
	t.Run("creating custom user fixture", func(t *testing.T) {
		// Create a custom user fixture
		user := CreateUserFixture(func(u *UserFixture) {
			u.Name = "Test Admin"
			u.Email = "admin@test.com"
			u.UserType = string(modelsv1.UserTypeTeam)
		})
		
		// Convert to model and use
		userModel := user.ToUserModel()
		assert.Equal(t, "Test Admin", userModel.Name)
		assert.Equal(t, "admin@test.com", userModel.Email)
		assert.Equal(t, modelsv1.UserTypeTeam, userModel.UserType)
	})

	t.Run("creating custom log entry fixture", func(t *testing.T) {
		// Create a custom log entry fixture
		logEntry := CreateLogEntryFixture(func(le *LogEntryFixture) {
			le.Level = "ERROR"
			le.Message = "Database connection failed"
			le.Status = 500
			le.Method = "POST"
			le.Path = "/api/v1/users"
			le.Extra = map[string]string{
				"error_code": "DB_CONNECTION_FAILED",
				"retry_count": "3",
			}
		})
		
		// Convert to model and use
		logModel := logEntry.ToLogEntryModel()
		assert.Equal(t, "ERROR", logModel.Level)
		assert.Equal(t, "Database connection failed", logModel.Message)
		assert.Equal(t, 500, logModel.Status)
		assert.Equal(t, "POST", logModel.Method)
		assert.Equal(t, "/api/v1/users", logModel.Path)
		assert.Equal(t, "DB_CONNECTION_FAILED", logModel.Extra["error_code"])
		assert.Equal(t, "3", logModel.Extra["retry_count"])
	})

	t.Run("creating related fixtures", func(t *testing.T) {
		// Create a user first
		user := CreateUserFixture(func(u *UserFixture) {
			u.Name = "Team Owner"
			u.UserType = string(modelsv1.UserTypeTeam)
		})
		
		// Create a team owned by that user
		team := CreateTeamFixture(user.ID, func(t *TeamFixture) {
			t.Name = "Engineering Team"
		})
		
		// Create another user to be a team member
		member := CreateUserFixture(func(u *UserFixture) {
			u.Name = "Team Member"
			u.UserType = string(modelsv1.UserTypeIndividual)
		})
		
		// Create team membership
		teamMember := CreateTeamMemberFixture(team.ID, member.ID, func(tm *TeamMemberFixture) {
			tm.Role = string(modelsv1.RoleAdmin)
		})
		
		// Create a refresh token for the member
		token := CreateRefreshTokenFixture(member.ID, func(rt *RefreshTokenFixture) {
			rt.Token = "custom_refresh_token_123"
		})
		
		// Verify relationships
		assert.Equal(t, user.ID, team.OwnerID)
		assert.Equal(t, team.ID, teamMember.TeamID)
		assert.Equal(t, member.ID, teamMember.UserID)
		assert.Equal(t, member.ID, token.UserID)
		assert.Equal(t, string(modelsv1.RoleAdmin), teamMember.Role)
		assert.Equal(t, "custom_refresh_token_123", token.Token)
	})
}

// Example test demonstrating fixture usage patterns for different scenarios
func TestExampleFixtureUsagePatterns(t *testing.T) {
	t.Run("authentication test scenario", func(t *testing.T) {
		// Create a user for authentication tests
		user := CreateUserFixture(func(u *UserFixture) {
			u.Email = "auth@test.com"
			u.Password = "$2a$10$hashedpassword" // Pre-hashed password
		})
		
		// Create valid and expired tokens
		validToken := CreateRefreshTokenFixture(user.ID, func(rt *RefreshTokenFixture) {
			rt.Token = "valid_token_123"
			rt.Revoked = false
		})
		
		expiredToken := CreateRefreshTokenFixture(user.ID, func(rt *RefreshTokenFixture) {
			rt.Token = "expired_token_456"
			rt.Revoked = false
			rt.ExpiresAt = rt.CreatedAt.Add(-1 * 24 * 60 * 60 * 1000000000) // 1 day ago
		})
		
		revokedToken := CreateRefreshTokenFixture(user.ID, func(rt *RefreshTokenFixture) {
			rt.Token = "revoked_token_789"
			rt.Revoked = true
		})
		
		// Use in authentication tests
		assert.Equal(t, user.ID, validToken.UserID)
		assert.False(t, validToken.Revoked)
		assert.True(t, validToken.ExpiresAt.After(validToken.CreatedAt))
		
		assert.True(t, expiredToken.ExpiresAt.Before(expiredToken.CreatedAt))
		assert.True(t, revokedToken.Revoked)
	})

	t.Run("logging test scenario", func(t *testing.T) {
		// Create different types of log entries for testing
		infoLog := CreateLogEntryFixture(func(le *LogEntryFixture) {
			le.Level = "INFO"
			le.Message = "User successfully logged in"
			le.Status = 200
			le.Method = "POST"
			le.Path = "/api/v1/auth/login"
		})
		
		errorLog := CreateLogEntryFixture(func(le *LogEntryFixture) {
			le.Level = "ERROR"
			le.Message = "Invalid credentials provided"
			le.Status = 401
			le.Method = "POST"
			le.Path = "/api/v1/auth/login"
		})
		
		warningLog := CreateLogEntryFixture(func(le *LogEntryFixture) {
			le.Level = "WARN"
			le.Message = "Rate limit approaching"
			le.Status = 200
			le.Method = "GET"
			le.Path = "/api/v1/logs"
		})
		
		// Use in logging tests
		logs := []LogEntryFixture{infoLog, errorLog, warningLog}
		
		assert.Len(t, logs, 3)
		assert.Equal(t, "INFO", logs[0].Level)
		assert.Equal(t, "ERROR", logs[1].Level)
		assert.Equal(t, "WARN", logs[2].Level)
	})
}