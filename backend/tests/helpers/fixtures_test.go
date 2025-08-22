package helpers

import (
	"testing"
	"time"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

func TestGetUser(t *testing.T) {
	// Test getting existing user
	user := GetUser("user-1")
	if user.ID != "user-1" {
		t.Errorf("expected user ID 'user-1', got '%s'", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected user name 'John Doe', got '%s'", user.Name)
	}

	// Test getting non-existing user (should return first user)
	user = GetUser("non-existing")
	if user.ID != "user-1" {
		t.Errorf("expected fallback to first user with ID 'user-1', got '%s'", user.ID)
	}
}

func TestGetTeam(t *testing.T) {
	// Test getting existing team
	team := GetTeam("team-1")
	if team.ID != "team-1" {
		t.Errorf("expected team ID 'team-1', got '%s'", team.ID)
	}
	if team.Name != "Development Team" {
		t.Errorf("expected team name 'Development Team', got '%s'", team.Name)
	}

	// Test getting non-existing team (should return first team)
	team = GetTeam("non-existing")
	if team.ID != "team-1" {
		t.Errorf("expected fallback to first team with ID 'team-1', got '%s'", team.ID)
	}
}

func TestGetTeamMember(t *testing.T) {
	// Test getting existing team member
	member := GetTeamMember("member-1")
	if member.ID != "member-1" {
		t.Errorf("expected member ID 'member-1', got '%s'", member.ID)
	}
	if member.Role != modelsv1.RoleEditor {
		t.Errorf("expected member role 'editor', got '%s'", member.Role)
	}

	// Test getting non-existing member (should return first member)
	member = GetTeamMember("non-existing")
	if member.ID != "member-1" {
		t.Errorf("expected fallback to first member with ID 'member-1', got '%s'", member.ID)
	}
}

func TestCreateUser(t *testing.T) {
	// Test creating user without overrides
	user := CreateUser()
	if user.Name != "Test User" {
		t.Errorf("expected default name 'Test User', got '%s'", user.Name)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected default email 'test@example.com', got '%s'", user.Email)
	}
	if user.UserType != modelsv1.UserTypeIndividual {
		t.Errorf("expected default user type 'individual', got '%s'", user.UserType)
	}

	// Test creating user with overrides
	user = CreateUser(func(u *modelsv1.User) {
		u.Name = "Custom User"
		u.Email = "custom@example.com"
		u.UserType = modelsv1.UserTypeTeam
	})
	if user.Name != "Custom User" {
		t.Errorf("expected overridden name 'Custom User', got '%s'", user.Name)
	}
	if user.Email != "custom@example.com" {
		t.Errorf("expected overridden email 'custom@example.com', got '%s'", user.Email)
	}
	if user.UserType != modelsv1.UserTypeTeam {
		t.Errorf("expected overridden user type 'team', got '%s'", user.UserType)
	}
}

func TestCreateTeam(t *testing.T) {
	ownerID := "test-owner-123"

	// Test creating team without overrides
	team := CreateTeam(ownerID)
	if team.Name != "Test Team" {
		t.Errorf("expected default name 'Test Team', got '%s'", team.Name)
	}
	if team.OwnerID != ownerID {
		t.Errorf("expected owner ID '%s', got '%s'", ownerID, team.OwnerID)
	}

	// Test creating team with overrides
	team = CreateTeam(ownerID, func(t *modelsv1.Team) {
		t.Name = "Custom Team"
	})
	if team.Name != "Custom Team" {
		t.Errorf("expected overridden name 'Custom Team', got '%s'", team.Name)
	}
	if team.OwnerID != ownerID {
		t.Errorf("expected owner ID '%s', got '%s'", ownerID, team.OwnerID)
	}
}

func TestCreateTeamMember(t *testing.T) {
	teamID := "test-team-123"
	userID := "test-user-456"

	// Test creating team member without overrides
	member := CreateTeamMember(teamID, userID)
	if member.TeamID != teamID {
		t.Errorf("expected team ID '%s', got '%s'", teamID, member.TeamID)
	}
	if member.UserID != userID {
		t.Errorf("expected user ID '%s', got '%s'", userID, member.UserID)
	}
	if member.Role != modelsv1.RoleEditor {
		t.Errorf("expected default role 'editor', got '%s'", member.Role)
	}

	// Test creating team member with overrides
	member = CreateTeamMember(teamID, userID, func(m *modelsv1.TeamMember) {
		m.Role = modelsv1.RoleAdmin
	})
	if member.Role != modelsv1.RoleAdmin {
		t.Errorf("expected overridden role 'admin', got '%s'", member.Role)
	}
}

func TestCreateLogEntry(t *testing.T) {
	// Test creating log entry without overrides
	logEntry := CreateLogEntry()
	if logEntry.Level != "INFO" {
		t.Errorf("expected default level 'INFO', got '%s'", logEntry.Level)
	}
	if logEntry.Message != "Test log message" {
		t.Errorf("expected default message 'Test log message', got '%s'", logEntry.Message)
	}
	if logEntry.Status != 200 {
		t.Errorf("expected default status 200, got %d", logEntry.Status)
	}

	// Test creating log entry with overrides
	logEntry = CreateLogEntry(func(l *modelsv1.LogEntry) {
		l.Level = "ERROR"
		l.Message = "Custom error message"
		l.Status = 500
	})
	if logEntry.Level != "ERROR" {
		t.Errorf("expected overridden level 'ERROR', got '%s'", logEntry.Level)
	}
	if logEntry.Message != "Custom error message" {
		t.Errorf("expected overridden message 'Custom error message', got '%s'", logEntry.Message)
	}
	if logEntry.Status != 500 {
		t.Errorf("expected overridden status 500, got %d", logEntry.Status)
	}
}

func TestGenerateID(t *testing.T) {
	// Test that generateID creates unique IDs
	id1 := generateID()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	id2 := generateID()

	if id1 == id2 {
		t.Errorf("expected unique IDs, but got same ID: %s", id1)
	}

	// Test ID format (should be timestamp-based)
	if len(id1) != len("20060102150405.000000") {
		t.Errorf("expected ID length %d, got %d", len("20060102150405.000000"), len(id1))
	}
}

func TestTestDataConsistency(t *testing.T) {
	// Test that predefined test data is consistent
	if len(TestUsers) == 0 {
		t.Error("TestUsers should not be empty")
	}

	if len(TestTeams) == 0 {
		t.Error("TestTeams should not be empty")
	}

	if len(TestTokens) == 0 {
		t.Error("TestTokens should not be empty")
	}

	if len(TestTeamMembers) == 0 {
		t.Error("TestTeamMembers should not be empty")
	}

	if len(TestLogs) == 0 {
		t.Error("TestLogs should not be empty")
	}

	// Test that team members reference valid teams and users
	for _, member := range TestTeamMembers {
		teamExists := false
		for _, team := range TestTeams {
			if team.ID == member.TeamID {
				teamExists = true
				break
			}
		}
		if !teamExists {
			t.Errorf("team member %s references non-existing team %s", member.ID, member.TeamID)
		}

		userExists := false
		for _, user := range TestUsers {
			if user.ID == member.UserID {
				userExists = true
				break
			}
		}
		if !userExists {
			t.Errorf("team member %s references non-existing user %s", member.ID, member.UserID)
		}
	}

	// Test that teams reference valid owners
	for _, team := range TestTeams {
		ownerExists := false
		for _, user := range TestUsers {
			if user.ID == team.OwnerID {
				ownerExists = true
				break
			}
		}
		if !ownerExists {
			t.Errorf("team %s references non-existing owner %s", team.ID, team.OwnerID)
		}
	}
}