package helpers

import (
	"testing"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// Example tests demonstrating how to use the simplified fixture system

func TestFixtureUsageExamples(t *testing.T) {
	t.Run("using predefined test users", func(t *testing.T) {
		// Get a predefined test user
		user := GetUser("user-1")
		
		// Use the user in your test
		if user.Email != "john.doe@example.com" {
			t.Errorf("expected john.doe@example.com, got %s", user.Email)
		}
		
		// Test user type
		if user.UserType != modelsv1.UserTypeIndividual {
			t.Errorf("expected individual user type, got %s", user.UserType)
		}
	})

	t.Run("creating custom test user", func(t *testing.T) {
		// Create a custom user with overrides
		user := CreateUser(func(u *modelsv1.User) {
			u.Name = "Custom Test User"
			u.Email = "custom@test.com"
			u.UserType = modelsv1.UserTypeTeam
		})
		
		// Verify the overrides were applied
		if user.Name != "Custom Test User" {
			t.Errorf("expected Custom Test User, got %s", user.Name)
		}
		if user.Email != "custom@test.com" {
			t.Errorf("expected custom@test.com, got %s", user.Email)
		}
		if user.UserType != modelsv1.UserTypeTeam {
			t.Errorf("expected team user type, got %s", user.UserType)
		}
	})

	t.Run("using predefined test teams", func(t *testing.T) {
		// Get a predefined test team
		team := GetTeam("team-1")
		
		// Use the team in your test
		if team.Name != "Development Team" {
			t.Errorf("expected Development Team, got %s", team.Name)
		}
		
		// Verify team owner
		if team.OwnerID != "user-2" {
			t.Errorf("expected user-2 as owner, got %s", team.OwnerID)
		}
	})

	t.Run("creating custom test team", func(t *testing.T) {
		ownerID := "test-owner-123"
		
		// Create a custom team with overrides
		team := CreateTeam(ownerID, func(t *modelsv1.Team) {
			t.Name = "Custom Test Team"
		})
		
		// Verify the overrides were applied
		if team.Name != "Custom Test Team" {
			t.Errorf("expected Custom Test Team, got %s", team.Name)
		}
		if team.OwnerID != ownerID {
			t.Errorf("expected %s as owner, got %s", ownerID, team.OwnerID)
		}
	})

	t.Run("using predefined team members", func(t *testing.T) {
		// Get a predefined team member
		member := GetTeamMember("member-1")
		
		// Use the member in your test
		if member.TeamID != "team-1" {
			t.Errorf("expected team-1, got %s", member.TeamID)
		}
		if member.UserID != "user-1" {
			t.Errorf("expected user-1, got %s", member.UserID)
		}
		if member.Role != modelsv1.RoleEditor {
			t.Errorf("expected editor role, got %s", member.Role)
		}
	})

	t.Run("creating custom team member", func(t *testing.T) {
		teamID := "test-team-456"
		userID := "test-user-789"
		
		// Create a custom team member with overrides
		member := CreateTeamMember(teamID, userID, func(m *modelsv1.TeamMember) {
			m.Role = modelsv1.RoleAdmin
		})
		
		// Verify the overrides were applied
		if member.TeamID != teamID {
			t.Errorf("expected %s as team ID, got %s", teamID, member.TeamID)
		}
		if member.UserID != userID {
			t.Errorf("expected %s as user ID, got %s", userID, member.UserID)
		}
		if member.Role != modelsv1.RoleAdmin {
			t.Errorf("expected admin role, got %s", member.Role)
		}
	})

	t.Run("using predefined log entries", func(t *testing.T) {
		// Use predefined log entries
		if len(TestLogs) == 0 {
			t.Error("expected test logs to be available")
		}
		
		// Get first log entry
		logEntry := TestLogs[0]
		
		// Verify log entry structure
		if logEntry.Level != "INFO" {
			t.Errorf("expected INFO level, got %s", logEntry.Level)
		}
		if logEntry.Status != 200 {
			t.Errorf("expected status 200, got %d", logEntry.Status)
		}
	})

	t.Run("creating custom log entry", func(t *testing.T) {
		// Create a custom log entry with overrides
		logEntry := CreateLogEntry(func(l *modelsv1.LogEntry) {
			l.Level = "ERROR"
			l.Message = "Custom error occurred"
			l.Status = 500
			l.Method = "POST"
			l.Path = "/api/v1/custom"
		})
		
		// Verify the overrides were applied
		if logEntry.Level != "ERROR" {
			t.Errorf("expected ERROR level, got %s", logEntry.Level)
		}
		if logEntry.Message != "Custom error occurred" {
			t.Errorf("expected custom message, got %s", logEntry.Message)
		}
		if logEntry.Status != 500 {
			t.Errorf("expected status 500, got %d", logEntry.Status)
		}
		if logEntry.Method != "POST" {
			t.Errorf("expected POST method, got %s", logEntry.Method)
		}
		if logEntry.Path != "/api/v1/custom" {
			t.Errorf("expected custom path, got %s", logEntry.Path)
		}
	})

	t.Run("using test tokens", func(t *testing.T) {
		// Use predefined test tokens
		if len(TestTokens) == 0 {
			t.Error("expected test tokens to be available")
		}
		
		// Get first token
		token := TestTokens[0]
		
		// Verify token format
		if token != "refresh_token_1_valid" {
			t.Errorf("expected refresh_token_1_valid, got %s", token)
		}
	})
}

// Example of how to use fixtures in a more complex test scenario
func TestComplexFixtureUsage(t *testing.T) {
	t.Run("team workflow with fixtures", func(t *testing.T) {
		// Get a team owner
		owner := GetUser("user-2")
		
		// Get the team owned by this user
		team := GetTeam("team-1")
		
		// Verify the relationship
		if team.OwnerID != owner.ID {
			t.Errorf("expected team owner to be %s, got %s", owner.ID, team.OwnerID)
		}
		
		// Get team members
		member := GetTeamMember("member-1")
		
		// Verify member belongs to the team
		if member.TeamID != team.ID {
			t.Errorf("expected member to belong to team %s, got %s", team.ID, member.TeamID)
		}
		
		// Create additional test data for this scenario
		newMember := CreateTeamMember(team.ID, "new-user-id", func(m *modelsv1.TeamMember) {
			m.Role = modelsv1.RoleCreator
		})
		
		// Verify the new member
		if newMember.TeamID != team.ID {
			t.Errorf("expected new member to belong to team %s, got %s", team.ID, newMember.TeamID)
		}
		if newMember.Role != modelsv1.RoleCreator {
			t.Errorf("expected creator role, got %s", newMember.Role)
		}
	})

	t.Run("authentication workflow with fixtures", func(t *testing.T) {
		// Get a test user for authentication
		user := GetUser("user-1")
		
		// Get a test token
		token := TestTokens[0]
		
		// Create a log entry for the authentication attempt
		authLog := CreateLogEntry(func(l *modelsv1.LogEntry) {
			l.Level = "INFO"
			l.Message = "User authentication successful"
			l.Method = "POST"
			l.Path = "/api/v1/auth/login"
			l.Status = 200
			l.Extra = map[string]string{
				"user_id": user.ID,
				"email":   user.Email,
				"token":   token,
			}
		})
		
		// Verify the log entry
		if authLog.Extra["user_id"] != user.ID {
			t.Errorf("expected user ID %s in log, got %s", user.ID, authLog.Extra["user_id"])
		}
		if authLog.Extra["email"] != user.Email {
			t.Errorf("expected email %s in log, got %s", user.Email, authLog.Extra["email"])
		}
		if authLog.Extra["token"] != token {
			t.Errorf("expected token %s in log, got %s", token, authLog.Extra["token"])
		}
	})
}