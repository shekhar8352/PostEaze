package helpers

import (
	"time"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// Simple test data structures - using actual model types directly
var (
	// TestUsers contains predefined test users
	TestUsers = []modelsv1.User{
		{
			ID:        "user-1",
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			Password:  "$2a$10$hashedpassword1", // bcrypt hash of "password123"
			UserType:  modelsv1.UserTypeIndividual,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "user-2",
			Name:      "Jane Smith",
			Email:     "jane.smith@example.com",
			Password:  "$2a$10$hashedpassword2", // bcrypt hash of "password456"
			UserType:  modelsv1.UserTypeTeam,
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        "user-3",
			Name:      "Bob Wilson",
			Email:     "bob.wilson@example.com",
			Password:  "$2a$10$hashedpassword3", // bcrypt hash of "password789"
			UserType:  modelsv1.UserTypeIndividual,
			CreatedAt: time.Now().Add(-6 * time.Hour),
			UpdatedAt: time.Now().Add(-6 * time.Hour),
		},
	}

	// TestTeams contains predefined test teams
	TestTeams = []modelsv1.Team{
		{
			ID:        "team-1",
			Name:      "Development Team",
			OwnerID:   "user-2",
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        "team-2",
			Name:      "Marketing Team",
			OwnerID:   "user-1",
			CreatedAt: time.Now().Add(-6 * time.Hour),
			UpdatedAt: time.Now().Add(-6 * time.Hour),
		},
	}

	// TestTokens contains predefined test refresh tokens
	TestTokens = []string{
		"refresh_token_1_valid",
		"refresh_token_2_expired",
		"refresh_token_3_revoked",
	}

	// TestTeamMembers contains predefined test team members
	TestTeamMembers = []modelsv1.TeamMember{
		{
			ID:        "member-1",
			TeamID:    "team-1",
			UserID:    "user-1",
			Role:      modelsv1.RoleEditor,
			CreatedAt: time.Now().Add(-10 * time.Hour),
			UpdatedAt: time.Now().Add(-10 * time.Hour),
		},
		{
			ID:        "member-2",
			TeamID:    "team-1",
			UserID:    "user-3",
			Role:      modelsv1.RoleCreator,
			CreatedAt: time.Now().Add(-8 * time.Hour),
			UpdatedAt: time.Now().Add(-8 * time.Hour),
		},
	}

	// TestLogs contains predefined test log entries
	TestLogs = []modelsv1.LogEntry{
		{
			Timestamp: "2024-01-15T10:30:00Z",
			Level:     "INFO",
			Message:   "User login successful",
			LogID:     "log-1",
			Method:    "POST",
			Path:      "/api/v1/auth/login",
			Status:    200,
			Duration:  "45ms",
			IP:        "192.168.1.100",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Extra: map[string]string{
				"user_id": "user-1",
				"email":   "john.doe@example.com",
			},
		},
		{
			Timestamp: "2024-01-15T10:35:00Z",
			Level:     "ERROR",
			Message:   "Authentication failed - invalid credentials",
			LogID:     "log-2",
			Method:    "POST",
			Path:      "/api/v1/auth/login",
			Status:    401,
			Duration:  "12ms",
			IP:        "192.168.1.101",
			UserAgent: "curl/7.68.0",
			Extra: map[string]string{
				"email":  "invalid@example.com",
				"reason": "invalid_password",
			},
		},
		{
			Timestamp: "2024-01-15T11:00:00Z",
			Level:     "INFO",
			Message:   "User signup completed",
			LogID:     "log-3",
			Method:    "POST",
			Path:      "/api/v1/auth/signup",
			Status:    201,
			Duration:  "120ms",
			IP:        "192.168.1.102",
			UserAgent: "PostmanRuntime/7.32.3",
			Extra: map[string]string{
				"user_id":   "user-3",
				"user_type": "individual",
			},
		},
	}
)

// Simple helper functions for test data creation with optional overrides

// GetUser returns a test user by ID, or the first user if ID not found
func GetUser(id string) modelsv1.User {
	for _, user := range TestUsers {
		if user.ID == id {
			return user
		}
	}
	return TestUsers[0] // Return first user as default
}

// GetTeam returns a test team by ID, or the first team if ID not found
func GetTeam(id string) modelsv1.Team {
	for _, team := range TestTeams {
		if team.ID == id {
			return team
		}
	}
	return TestTeams[0] // Return first team as default
}

// GetTeamMember returns a test team member by ID, or the first member if ID not found
func GetTeamMember(id string) modelsv1.TeamMember {
	for _, member := range TestTeamMembers {
		if member.ID == id {
			return member
		}
	}
	return TestTeamMembers[0] // Return first member as default
}

// CreateUser creates a test user with optional overrides
func CreateUser(overrides ...func(*modelsv1.User)) modelsv1.User {
	now := time.Now()
	user := modelsv1.User{
		ID:        "test-user-" + generateID(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "$2a$10$hashedpassword", // bcrypt hash of "testpassword"
		UserType:  modelsv1.UserTypeIndividual,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&user)
	}

	return user
}

// CreateTeam creates a test team with optional overrides
func CreateTeam(ownerID string, overrides ...func(*modelsv1.Team)) modelsv1.Team {
	now := time.Now()
	team := modelsv1.Team{
		ID:        "test-team-" + generateID(),
		Name:      "Test Team",
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&team)
	}

	return team
}

// CreateTeamMember creates a test team member with optional overrides
func CreateTeamMember(teamID, userID string, overrides ...func(*modelsv1.TeamMember)) modelsv1.TeamMember {
	now := time.Now()
	member := modelsv1.TeamMember{
		ID:        "test-member-" + generateID(),
		TeamID:    teamID,
		UserID:    userID,
		Role:      modelsv1.RoleEditor,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&member)
	}

	return member
}

// CreateLogEntry creates a test log entry with optional overrides
func CreateLogEntry(overrides ...func(*modelsv1.LogEntry)) modelsv1.LogEntry {
	now := time.Now()
	logEntry := modelsv1.LogEntry{
		Timestamp: now.Format(time.RFC3339),
		Level:     "INFO",
		Message:   "Test log message",
		LogID:     "test-log-" + generateID(),
		Method:    "GET",
		Path:      "/api/v1/test",
		Status:    200,
		Duration:  "10ms",
		IP:        "127.0.0.1",
		UserAgent: "test-agent/1.0",
		Extra: map[string]string{
			"test": "true",
		},
	}

	for _, override := range overrides {
		override(&logEntry)
	}

	return logEntry
}

// generateID creates a simple ID for testing
func generateID() string {
	return time.Now().Format("20060102150405.000000")
}