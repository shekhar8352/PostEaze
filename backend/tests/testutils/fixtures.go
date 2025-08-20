package testutils

import (
	"database/sql"
	"fmt"
	"time"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// UserFixture represents test user data
type UserFixture struct {
	ID        string
	Name      string
	Email     string
	Password  string
	UserType  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TeamFixture represents test team data
type TeamFixture struct {
	ID        string
	Name      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RefreshTokenFixture represents test refresh token data
type RefreshTokenFixture struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TeamMemberFixture represents test team member data
type TeamMemberFixture struct {
	ID        string
	TeamID    string
	UserID    string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// LogEntryFixture represents test log entry data
type LogEntryFixture struct {
	Timestamp string
	Level     string
	Message   string
	LogID     string
	Method    string
	Path      string
	Status    int
	Duration  string
	IP        string
	UserAgent string
	Extra     map[string]string
}

// Predefined test fixtures
var (
	// TestUsers contains predefined test user data
	TestUsers = []UserFixture{
		{
			ID:        "user-1",
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			Password:  "$2a$10$hashedpassword1", // bcrypt hash of "password123"
			UserType:  string(modelsv1.UserTypeIndividual),
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "user-2",
			Name:      "Jane Smith",
			Email:     "jane.smith@example.com",
			Password:  "$2a$10$hashedpassword2", // bcrypt hash of "password456"
			UserType:  string(modelsv1.UserTypeTeam),
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        "user-3",
			Name:      "Bob Wilson",
			Email:     "bob.wilson@example.com",
			Password:  "$2a$10$hashedpassword3", // bcrypt hash of "password789"
			UserType:  string(modelsv1.UserTypeIndividual),
			CreatedAt: time.Now().Add(-6 * time.Hour),
			UpdatedAt: time.Now().Add(-6 * time.Hour),
		},
	}

	// TestTeams contains predefined test team data
	TestTeams = []TeamFixture{
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

	// TestTokens contains predefined test refresh token data
	TestTokens = []RefreshTokenFixture{
		{
			ID:        "token-1",
			UserID:    "user-1",
			Token:     "refresh_token_1_valid",
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Valid for 7 days
			Revoked:   false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "token-2",
			UserID:    "user-2",
			Token:     "refresh_token_2_expired",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			Revoked:   false,
			CreatedAt: time.Now().Add(-25 * time.Hour),
			UpdatedAt: time.Now().Add(-25 * time.Hour),
		},
		{
			ID:        "token-3",
			UserID:    "user-3",
			Token:     "refresh_token_3_revoked",
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Valid but revoked
			Revoked:   true,
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	// TestTeamMembers contains predefined test team member data
	TestTeamMembers = []TeamMemberFixture{
		{
			ID:        "member-1",
			TeamID:    "team-1",
			UserID:    "user-1",
			Role:      string(modelsv1.RoleEditor),
			CreatedAt: time.Now().Add(-10 * time.Hour),
			UpdatedAt: time.Now().Add(-10 * time.Hour),
		},
		{
			ID:        "member-2",
			TeamID:    "team-1",
			UserID:    "user-3",
			Role:      string(modelsv1.RoleCreator),
			CreatedAt: time.Now().Add(-8 * time.Hour),
			UpdatedAt: time.Now().Add(-8 * time.Hour),
		},
	}

	// TestLogs contains predefined test log entry data
	TestLogs = []LogEntryFixture{
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
		{
			Timestamp: "2024-01-15T11:15:00Z",
			Level:     "WARN",
			Message:   "Rate limit exceeded",
			LogID:     "log-4",
			Method:    "GET",
			Path:      "/api/v1/logs",
			Status:    429,
			Duration:  "5ms",
			IP:        "192.168.1.103",
			UserAgent: "Python/3.9 requests/2.28.1",
			Extra: map[string]string{
				"limit":     "100",
				"remaining": "0",
				"reset_at":  "2024-01-15T11:16:00Z",
			},
		},
		{
			Timestamp: "2024-01-15T12:00:00Z",
			Level:     "INFO",
			Message:   "Logs retrieved successfully",
			LogID:     "log-5",
			Method:    "GET",
			Path:      "/api/v1/logs/2024-01-15",
			Status:    200,
			Duration:  "25ms",
			IP:        "192.168.1.100",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			Extra: map[string]string{
				"user_id": "user-1",
				"count":   "150",
			},
		},
	}
)

// CreateUserFixture creates a user fixture with optional overrides
func CreateUserFixture(overrides ...func(*UserFixture)) UserFixture {
	now := time.Now()
	user := UserFixture{
		ID:        "test-user-" + generateID(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "$2a$10$hashedpassword", // bcrypt hash of "testpassword"
		UserType:  string(modelsv1.UserTypeIndividual),
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&user)
	}

	return user
}

// CreateTeamFixture creates a team fixture with optional overrides
func CreateTeamFixture(ownerID string, overrides ...func(*TeamFixture)) TeamFixture {
	now := time.Now()
	team := TeamFixture{
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

// CreateRefreshTokenFixture creates a refresh token fixture with optional overrides
func CreateRefreshTokenFixture(userID string, overrides ...func(*RefreshTokenFixture)) RefreshTokenFixture {
	now := time.Now()
	token := RefreshTokenFixture{
		ID:        "test-token-" + generateID(),
		UserID:    userID,
		Token:     "test_refresh_token_" + generateID(),
		ExpiresAt: now.Add(7 * 24 * time.Hour), // Valid for 7 days
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&token)
	}

	return token
}

// CreateTeamMemberFixture creates a team member fixture with optional overrides
func CreateTeamMemberFixture(teamID, userID string, overrides ...func(*TeamMemberFixture)) TeamMemberFixture {
	now := time.Now()
	member := TeamMemberFixture{
		ID:        "test-member-" + generateID(),
		TeamID:    teamID,
		UserID:    userID,
		Role:      string(modelsv1.RoleEditor),
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&member)
	}

	return member
}

// CreateLogEntryFixture creates a log entry fixture with optional overrides
func CreateLogEntryFixture(overrides ...func(*LogEntryFixture)) LogEntryFixture {
	now := time.Now()
	logEntry := LogEntryFixture{
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

// Helper function to generate simple IDs for testing
func generateID() string {
	return time.Now().Format("20060102150405.000000")
}

// ToModel converts fixtures to model structs

// ToUserModel converts UserFixture to modelsv1.User
func (u UserFixture) ToUserModel() modelsv1.User {
	return modelsv1.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		UserType:  modelsv1.UserType(u.UserType),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToTeamModel converts TeamFixture to modelsv1.Team
func (t TeamFixture) ToTeamModel() modelsv1.Team {
	return modelsv1.Team{
		ID:        t.ID,
		Name:      t.Name,
		OwnerID:   t.OwnerID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ToRefreshTokenModel converts RefreshTokenFixture to modelsv1.RefreshToken
func (rt RefreshTokenFixture) ToRefreshTokenModel() modelsv1.RefreshToken {
	// For testing, we'll create a simple UUID-like structure
	// In a real implementation, you'd parse the UUID properly
	return modelsv1.RefreshToken{
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
		Revoked:   rt.Revoked,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
	}
}

// ToTeamMemberModel converts TeamMemberFixture to modelsv1.TeamMember
func (tm TeamMemberFixture) ToTeamMemberModel() modelsv1.TeamMember {
	return modelsv1.TeamMember{
		ID:        tm.ID,
		TeamID:    tm.TeamID,
		UserID:    tm.UserID,
		Role:      modelsv1.Role(tm.Role),
		CreatedAt: tm.CreatedAt,
		UpdatedAt: tm.UpdatedAt,
	}
}

// ToLogEntryModel converts LogEntryFixture to modelsv1.LogEntry
func (le LogEntryFixture) ToLogEntryModel() modelsv1.LogEntry {
	return modelsv1.LogEntry{
		Timestamp: le.Timestamp,
		Level:     le.Level,
		Message:   le.Message,
		LogID:     le.LogID,
		Method:    le.Method,
		Path:      le.Path,
		Status:    le.Status,
		Duration:  le.Duration,
		IP:        le.IP,
		UserAgent: le.UserAgent,
		Extra:     le.Extra,
	}
}

// Helper function to parse UUID (simplified for testing)
func parseUUID(s string) (interface{}, error) {
	// For testing purposes, we'll return the string as-is
	// In a real implementation, you might want to parse it properly
	return s, nil
}

// Database fixture loading utilities

// LoadUserFixtures loads user fixtures into the test database
func LoadUserFixtures(db *sql.DB, users []UserFixture) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	stmt, err := db.Prepare(`
		INSERT INTO users (id, name, email, password, user_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare user insert statement: %w", err)
	}
	defer stmt.Close()

	for _, user := range users {
		_, err := stmt.Exec(
			user.ID,
			user.Name,
			user.Email,
			user.Password,
			user.UserType,
			user.CreatedAt,
			user.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user fixture %s: %w", user.ID, err)
		}
	}

	return nil
}

// LoadTeamFixtures loads team fixtures into the test database
func LoadTeamFixtures(db *sql.DB, teams []TeamFixture) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	stmt, err := db.Prepare(`
		INSERT INTO teams (id, name, owner_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare team insert statement: %w", err)
	}
	defer stmt.Close()

	for _, team := range teams {
		_, err := stmt.Exec(
			team.ID,
			team.Name,
			team.OwnerID,
			team.CreatedAt,
			team.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert team fixture %s: %w", team.ID, err)
		}
	}

	return nil
}

// LoadRefreshTokenFixtures loads refresh token fixtures into the test database
func LoadRefreshTokenFixtures(db *sql.DB, tokens []RefreshTokenFixture) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	stmt, err := db.Prepare(`
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, revoked, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare refresh token insert statement: %w", err)
	}
	defer stmt.Close()

	for _, token := range tokens {
		_, err := stmt.Exec(
			token.ID,
			token.UserID,
			token.Token,
			token.ExpiresAt,
			token.Revoked,
			token.CreatedAt,
			token.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert refresh token fixture %s: %w", token.ID, err)
		}
	}

	return nil
}

// LoadTeamMemberFixtures loads team member fixtures into the test database
func LoadTeamMemberFixtures(db *sql.DB, members []TeamMemberFixture) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	stmt, err := db.Prepare(`
		INSERT INTO team_members (id, team_id, user_id, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare team member insert statement: %w", err)
	}
	defer stmt.Close()

	for _, member := range members {
		_, err := stmt.Exec(
			member.ID,
			member.TeamID,
			member.UserID,
			member.Role,
			member.CreatedAt,
			member.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert team member fixture %s: %w", member.ID, err)
		}
	}

	return nil
}

// LoadAllFixtures loads all predefined fixtures into the test database
func LoadAllFixtures(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	// Load fixtures in dependency order
	if err := LoadUserFixtures(db, TestUsers); err != nil {
		return fmt.Errorf("failed to load user fixtures: %w", err)
	}

	if err := LoadTeamFixtures(db, TestTeams); err != nil {
		return fmt.Errorf("failed to load team fixtures: %w", err)
	}

	if err := LoadRefreshTokenFixtures(db, TestTokens); err != nil {
		return fmt.Errorf("failed to load refresh token fixtures: %w", err)
	}

	if err := LoadTeamMemberFixtures(db, TestTeamMembers); err != nil {
		return fmt.Errorf("failed to load team member fixtures: %w", err)
	}

	return nil
}

// LoadCustomFixtures loads custom fixtures into the test database
func LoadCustomFixtures(db *sql.DB, fixtures interface{}) error {
	switch f := fixtures.(type) {
	case []UserFixture:
		return LoadUserFixtures(db, f)
	case []TeamFixture:
		return LoadTeamFixtures(db, f)
	case []RefreshTokenFixture:
		return LoadRefreshTokenFixtures(db, f)
	case []TeamMemberFixture:
		return LoadTeamMemberFixtures(db, f)
	case UserFixture:
		return LoadUserFixtures(db, []UserFixture{f})
	case TeamFixture:
		return LoadTeamFixtures(db, []TeamFixture{f})
	case RefreshTokenFixture:
		return LoadRefreshTokenFixtures(db, []RefreshTokenFixture{f})
	case TeamMemberFixture:
		return LoadTeamMemberFixtures(db, []TeamMemberFixture{f})
	default:
		return fmt.Errorf("unsupported fixture type: %T", fixtures)
	}
}

// ClearAllFixtures removes all fixture data from the test database
func ClearAllFixtures(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	// Clear in reverse dependency order
	tables := []string{
		"team_members",
		"refresh_tokens",
		"teams",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}
	}

	return nil
}