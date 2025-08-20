package testutils

import (
	"testing"
	"time"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserFixture(t *testing.T) {
	t.Run("creates user fixture with defaults", func(t *testing.T) {
		user := CreateUserFixture()

		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "$2a$10$hashedpassword", user.Password)
		assert.Equal(t, string(modelsv1.UserTypeIndividual), user.UserType)
		assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
	})

	t.Run("creates user fixture with overrides", func(t *testing.T) {
		user := CreateUserFixture(func(u *UserFixture) {
			u.Name = "Custom User"
			u.Email = "custom@example.com"
			u.UserType = string(modelsv1.UserTypeTeam)
		})

		assert.Equal(t, "Custom User", user.Name)
		assert.Equal(t, "custom@example.com", user.Email)
		assert.Equal(t, string(modelsv1.UserTypeTeam), user.UserType)
	})
}

func TestCreateTeamFixture(t *testing.T) {
	t.Run("creates team fixture with defaults", func(t *testing.T) {
		ownerID := "test-owner-123"
		team := CreateTeamFixture(ownerID)

		assert.NotEmpty(t, team.ID)
		assert.Equal(t, "Test Team", team.Name)
		assert.Equal(t, ownerID, team.OwnerID)
		assert.WithinDuration(t, time.Now(), team.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), team.UpdatedAt, time.Second)
	})

	t.Run("creates team fixture with overrides", func(t *testing.T) {
		ownerID := "test-owner-123"
		team := CreateTeamFixture(ownerID, func(t *TeamFixture) {
			t.Name = "Custom Team"
		})

		assert.Equal(t, "Custom Team", team.Name)
		assert.Equal(t, ownerID, team.OwnerID)
	})
}

func TestCreateRefreshTokenFixture(t *testing.T) {
	t.Run("creates refresh token fixture with defaults", func(t *testing.T) {
		userID := "test-user-123"
		token := CreateRefreshTokenFixture(userID)

		assert.NotEmpty(t, token.ID)
		assert.Equal(t, userID, token.UserID)
		assert.NotEmpty(t, token.Token)
		assert.False(t, token.Revoked)
		assert.True(t, token.ExpiresAt.After(time.Now()))
		assert.WithinDuration(t, time.Now(), token.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), token.UpdatedAt, time.Second)
	})

	t.Run("creates refresh token fixture with overrides", func(t *testing.T) {
		userID := "test-user-123"
		token := CreateRefreshTokenFixture(userID, func(rt *RefreshTokenFixture) {
			rt.Token = "custom-token"
			rt.Revoked = true
		})

		assert.Equal(t, "custom-token", token.Token)
		assert.True(t, token.Revoked)
	})
}

func TestCreateTeamMemberFixture(t *testing.T) {
	t.Run("creates team member fixture with defaults", func(t *testing.T) {
		teamID := "test-team-123"
		userID := "test-user-123"
		member := CreateTeamMemberFixture(teamID, userID)

		assert.NotEmpty(t, member.ID)
		assert.Equal(t, teamID, member.TeamID)
		assert.Equal(t, userID, member.UserID)
		assert.Equal(t, string(modelsv1.RoleEditor), member.Role)
		assert.WithinDuration(t, time.Now(), member.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), member.UpdatedAt, time.Second)
	})

	t.Run("creates team member fixture with overrides", func(t *testing.T) {
		teamID := "test-team-123"
		userID := "test-user-123"
		member := CreateTeamMemberFixture(teamID, userID, func(tm *TeamMemberFixture) {
			tm.Role = string(modelsv1.RoleAdmin)
		})

		assert.Equal(t, string(modelsv1.RoleAdmin), member.Role)
	})
}

func TestCreateLogEntryFixture(t *testing.T) {
	t.Run("creates log entry fixture with defaults", func(t *testing.T) {
		logEntry := CreateLogEntryFixture()

		assert.NotEmpty(t, logEntry.Timestamp)
		assert.Equal(t, "INFO", logEntry.Level)
		assert.Equal(t, "Test log message", logEntry.Message)
		assert.NotEmpty(t, logEntry.LogID)
		assert.Equal(t, "GET", logEntry.Method)
		assert.Equal(t, "/api/v1/test", logEntry.Path)
		assert.Equal(t, 200, logEntry.Status)
		assert.Equal(t, "10ms", logEntry.Duration)
		assert.Equal(t, "127.0.0.1", logEntry.IP)
		assert.Equal(t, "test-agent/1.0", logEntry.UserAgent)
		assert.NotNil(t, logEntry.Extra)
		assert.Equal(t, "true", logEntry.Extra["test"])
	})

	t.Run("creates log entry fixture with overrides", func(t *testing.T) {
		logEntry := CreateLogEntryFixture(func(le *LogEntryFixture) {
			le.Level = "ERROR"
			le.Message = "Custom error message"
			le.Status = 500
			le.Extra = map[string]string{
				"error": "custom_error",
			}
		})

		assert.Equal(t, "ERROR", logEntry.Level)
		assert.Equal(t, "Custom error message", logEntry.Message)
		assert.Equal(t, 500, logEntry.Status)
		assert.Equal(t, "custom_error", logEntry.Extra["error"])
	})
}

func TestFixtureToModelConversions(t *testing.T) {
	t.Run("converts user fixture to model", func(t *testing.T) {
		fixture := CreateUserFixture()
		model := fixture.ToUserModel()

		assert.Equal(t, fixture.ID, model.ID)
		assert.Equal(t, fixture.Name, model.Name)
		assert.Equal(t, fixture.Email, model.Email)
		assert.Equal(t, fixture.Password, model.Password)
		assert.Equal(t, modelsv1.UserType(fixture.UserType), model.UserType)
		assert.Equal(t, fixture.CreatedAt, model.CreatedAt)
		assert.Equal(t, fixture.UpdatedAt, model.UpdatedAt)
	})

	t.Run("converts team fixture to model", func(t *testing.T) {
		fixture := CreateTeamFixture("owner-123")
		model := fixture.ToTeamModel()

		assert.Equal(t, fixture.ID, model.ID)
		assert.Equal(t, fixture.Name, model.Name)
		assert.Equal(t, fixture.OwnerID, model.OwnerID)
		assert.Equal(t, fixture.CreatedAt, model.CreatedAt)
		assert.Equal(t, fixture.UpdatedAt, model.UpdatedAt)
	})

	t.Run("converts refresh token fixture to model", func(t *testing.T) {
		fixture := CreateRefreshTokenFixture("user-123")
		model := fixture.ToRefreshTokenModel()

		assert.Equal(t, fixture.Token, model.Token)
		assert.Equal(t, fixture.ExpiresAt, model.ExpiresAt)
		assert.Equal(t, fixture.Revoked, model.Revoked)
		assert.Equal(t, fixture.CreatedAt, model.CreatedAt)
		assert.Equal(t, fixture.UpdatedAt, model.UpdatedAt)
	})

	t.Run("converts team member fixture to model", func(t *testing.T) {
		fixture := CreateTeamMemberFixture("team-123", "user-123")
		model := fixture.ToTeamMemberModel()

		assert.Equal(t, fixture.ID, model.ID)
		assert.Equal(t, fixture.TeamID, model.TeamID)
		assert.Equal(t, fixture.UserID, model.UserID)
		assert.Equal(t, modelsv1.Role(fixture.Role), model.Role)
		assert.Equal(t, fixture.CreatedAt, model.CreatedAt)
		assert.Equal(t, fixture.UpdatedAt, model.UpdatedAt)
	})

	t.Run("converts log entry fixture to model", func(t *testing.T) {
		fixture := CreateLogEntryFixture()
		model := fixture.ToLogEntryModel()

		assert.Equal(t, fixture.Timestamp, model.Timestamp)
		assert.Equal(t, fixture.Level, model.Level)
		assert.Equal(t, fixture.Message, model.Message)
		assert.Equal(t, fixture.LogID, model.LogID)
		assert.Equal(t, fixture.Method, model.Method)
		assert.Equal(t, fixture.Path, model.Path)
		assert.Equal(t, fixture.Status, model.Status)
		assert.Equal(t, fixture.Duration, model.Duration)
		assert.Equal(t, fixture.IP, model.IP)
		assert.Equal(t, fixture.UserAgent, model.UserAgent)
		assert.Equal(t, fixture.Extra, model.Extra)
	})
}

func TestPredefinedFixtures(t *testing.T) {
	t.Run("has predefined test users", func(t *testing.T) {
		require.Len(t, TestUsers, 3)
		
		assert.Equal(t, "user-1", TestUsers[0].ID)
		assert.Equal(t, "John Doe", TestUsers[0].Name)
		assert.Equal(t, "john.doe@example.com", TestUsers[0].Email)
		assert.Equal(t, string(modelsv1.UserTypeIndividual), TestUsers[0].UserType)
	})

	t.Run("has predefined test teams", func(t *testing.T) {
		require.Len(t, TestTeams, 2)
		
		assert.Equal(t, "team-1", TestTeams[0].ID)
		assert.Equal(t, "Development Team", TestTeams[0].Name)
		assert.Equal(t, "user-2", TestTeams[0].OwnerID)
	})

	t.Run("has predefined test tokens", func(t *testing.T) {
		require.Len(t, TestTokens, 3)
		
		assert.Equal(t, "token-1", TestTokens[0].ID)
		assert.Equal(t, "user-1", TestTokens[0].UserID)
		assert.Equal(t, "refresh_token_1_valid", TestTokens[0].Token)
		assert.False(t, TestTokens[0].Revoked)
	})

	t.Run("has predefined test team members", func(t *testing.T) {
		require.Len(t, TestTeamMembers, 2)
		
		assert.Equal(t, "member-1", TestTeamMembers[0].ID)
		assert.Equal(t, "team-1", TestTeamMembers[0].TeamID)
		assert.Equal(t, "user-1", TestTeamMembers[0].UserID)
		assert.Equal(t, string(modelsv1.RoleEditor), TestTeamMembers[0].Role)
	})

	t.Run("has predefined test logs", func(t *testing.T) {
		require.Len(t, TestLogs, 5)
		
		assert.Equal(t, "log-1", TestLogs[0].LogID)
		assert.Equal(t, "INFO", TestLogs[0].Level)
		assert.Equal(t, "User login successful", TestLogs[0].Message)
		assert.Equal(t, "POST", TestLogs[0].Method)
		assert.Equal(t, "/api/v1/auth/login", TestLogs[0].Path)
		assert.Equal(t, 200, TestLogs[0].Status)
	})
}

func TestGenerateID(t *testing.T) {
	t.Run("generates unique IDs", func(t *testing.T) {
		id1 := generateID()
		time.Sleep(time.Millisecond) // Ensure different timestamps
		id2 := generateID()

		assert.NotEmpty(t, id1)
		assert.NotEmpty(t, id2)
		assert.NotEqual(t, id1, id2)
	})
}

func TestLoadCustomFixtures(t *testing.T) {
	t.Run("supports different fixture types", func(t *testing.T) {
		// Test that LoadCustomFixtures can handle different types
		// Since we don't have a real database connection in tests, we'll test the type checking
		
		userFixture := CreateUserFixture()
		teamFixture := CreateTeamFixture("owner-123")
		tokenFixture := CreateRefreshTokenFixture("user-123")
		memberFixture := CreateTeamMemberFixture("team-123", "user-123")
		
		userFixtures := []UserFixture{userFixture}
		teamFixtures := []TeamFixture{teamFixture}
		tokenFixtures := []RefreshTokenFixture{tokenFixture}
		memberFixtures := []TeamMemberFixture{memberFixture}
		
		// These would normally interact with a database, but we're testing the interface
		// In a real test environment with a database, these would actually load data
		
		// Test single fixtures
		err := LoadCustomFixtures(nil, userFixture)
		assert.Error(t, err) // Should error because db is nil
		
		err = LoadCustomFixtures(nil, teamFixture)
		assert.Error(t, err) // Should error because db is nil
		
		err = LoadCustomFixtures(nil, tokenFixture)
		assert.Error(t, err) // Should error because db is nil
		
		err = LoadCustomFixtures(nil, memberFixture)
		assert.Error(t, err) // Should error because db is nil
		
		// Test slice fixtures
		err = LoadCustomFixtures(nil, userFixtures)
		assert.Error(t, err) // Should error because db is nil
		
		err = LoadCustomFixtures(nil, teamFixtures)
		assert.Error(t, err) // Should error because db is nil
		
		err = LoadCustomFixtures(nil, tokenFixtures)
		assert.Error(t, err) // Should error because db is nil
		
		err = LoadCustomFixtures(nil, memberFixtures)
		assert.Error(t, err) // Should error because db is nil
		
		// Test unsupported type
		err = LoadCustomFixtures(nil, "unsupported")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported fixture type")
	})
}

func TestFixtureLoadingUtilities(t *testing.T) {
	t.Run("fixture loading functions exist and handle nil database", func(t *testing.T) {
		// Test that all loading functions exist and handle nil database gracefully
		
		err := LoadUserFixtures(nil, TestUsers)
		assert.Error(t, err)
		
		err = LoadTeamFixtures(nil, TestTeams)
		assert.Error(t, err)
		
		err = LoadRefreshTokenFixtures(nil, TestTokens)
		assert.Error(t, err)
		
		err = LoadTeamMemberFixtures(nil, TestTeamMembers)
		assert.Error(t, err)
		
		err = LoadAllFixtures(nil)
		assert.Error(t, err)
		
		err = ClearAllFixtures(nil)
		assert.Error(t, err)
	})
}