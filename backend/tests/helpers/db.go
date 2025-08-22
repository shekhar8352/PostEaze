package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shekhar8352/PostEaze/utils/database"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DBHelper provides simple database testing utilities
type DBHelper struct {
	DB *sql.DB
}

// Implement Database interface methods
func (h *DBHelper) QueryRaw(ctx context.Context, entity database.RawEntity, code int) error {
	row := h.DB.QueryRowContext(ctx, entity.GetQuery(code), entity.GetQueryValues(code)...)
	err := entity.BindRawRow(code, row)
	if err == sql.ErrNoRows {
		return database.ErrNoRecords
	}
	return err
}

func (h *DBHelper) QueryMultiRaw(ctx context.Context, entity database.RawEntity, code int) ([]database.RawEntity, error) {
	rows, err := h.DB.QueryContext(ctx, entity.GetMultiQuery(code), entity.GetMultiQueryValues(code)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]database.RawEntity, 0)
	for rows.Next() {
		err = entity.BindRawRow(code, rows)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
		entity = entity.GetNextRaw()
	}
	if len(result) == 0 {
		return nil, database.ErrNoRecords
	}
	return result, nil
}

func (h *DBHelper) ExecRaws(ctx context.Context, source string, execs ...database.RawExec) error {
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, exec := range execs {
		_, err = tx.ExecContext(ctx, exec.Entity.GetExec(exec.Code), exec.Entity.GetExecValues(exec.Code, source)...)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (h *DBHelper) ExecRawsConsistent(ctx context.Context, source string, execs ...database.RawExec) error {
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, exec := range execs {
		result, err := tx.ExecContext(ctx, exec.Entity.GetExec(exec.Code), exec.Entity.GetExecValues(exec.Code, source)...)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return database.ErrNoRowsAffected
		}
	}
	return tx.Commit()
}

// NewTestDB creates a new in-memory SQLite database for testing
func NewTestDB() (*DBHelper, error) {
	// Use in-memory SQLite for fast unit tests
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	helper := &DBHelper{DB: db}

	// Create tables
	if err := helper.createTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return helper, nil
}

// createTables creates all necessary tables for testing
func (h *DBHelper) createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			user_type TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teams (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			owner_id TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS team_members (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE (team_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			revoked BOOLEAN DEFAULT FALSE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, table := range tables {
		if _, err := h.DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// LoadFixture loads test data into the database
func (h *DBHelper) LoadFixture(data interface{}) error {
	switch fixture := data.(type) {
	case TestUser:
		return h.loadUser(fixture)
	case []TestUser:
		for _, user := range fixture {
			if err := h.loadUser(user); err != nil {
				return err
			}
		}
	case TestTeam:
		return h.loadTeam(fixture)
	case []TestTeam:
		for _, team := range fixture {
			if err := h.loadTeam(team); err != nil {
				return err
			}
		}
	case TestToken:
		return h.loadToken(fixture)
	case []TestToken:
		for _, token := range fixture {
			if err := h.loadToken(token); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported fixture type: %T", fixture)
	}
	return nil
}

// loadUser loads a user fixture
func (h *DBHelper) loadUser(user TestUser) error {
	query := `INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := h.DB.Exec(query, 
		user.ID, user.Name, user.Email, user.Password, 
		user.UserType, user.CreatedAt, user.UpdatedAt)
	
	return err
}

// loadTeam loads a team fixture
func (h *DBHelper) loadTeam(team TestTeam) error {
	query := `INSERT INTO teams (id, name, owner_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`
	
	_, err := h.DB.Exec(query, 
		team.ID, team.Name, team.OwnerID, team.CreatedAt, team.UpdatedAt)
	
	return err
}

// loadToken loads a refresh token fixture
func (h *DBHelper) loadToken(token TestToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, revoked, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := h.DB.Exec(query, 
		token.ID, token.UserID, token.Token, token.ExpiresAt, 
		token.Revoked, token.CreatedAt, token.UpdatedAt)
	
	return err
}

// Cleanup removes all data from test tables and closes the database
func (h *DBHelper) Cleanup() error {
	if h.DB == nil {
		return nil
	}

	// Clear all tables in reverse dependency order
	tables := []string{
		"team_members",
		"refresh_tokens", 
		"teams",
		"users",
	}

	for _, table := range tables {
		if _, err := h.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			// Log error but continue cleanup
			fmt.Printf("Warning: failed to cleanup table %s: %v\n", table, err)
		}
	}

	return h.DB.Close()
}

// CleanData removes all data from test tables without closing the database
func (h *DBHelper) CleanData() error {
	if h.DB == nil {
		return nil
	}

	tables := []string{
		"team_members",
		"refresh_tokens", 
		"teams",
		"users",
	}

	for _, table := range tables {
		if _, err := h.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to cleanup table %s: %w", table, err)
		}
	}

	return nil
}

// BeginTx starts a transaction for testing
func (h *DBHelper) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return h.DB.BeginTx(ctx, nil)
}

// Simple test data structures
type TestUser struct {
	ID        string
	Name      string
	Email     string
	Password  string
	UserType  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TestTeam struct {
	ID        string
	Name      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TestToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Predefined test data
var (
	DefaultUsers = []TestUser{
		{
			ID:        "user-1",
			Name:      "John Doe",
			Email:     "john.doe@test.com",
			Password:  "$2a$10$hashedpassword1",
			UserType:  "individual",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "user-2", 
			Name:      "Jane Smith",
			Email:     "jane.smith@test.com",
			Password:  "$2a$10$hashedpassword2",
			UserType:  "team",
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
	}

	DefaultTeams = []TestTeam{
		{
			ID:        "team-1",
			Name:      "Test Team",
			OwnerID:   "user-2",
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
	}

	DefaultTokens = []TestToken{
		{
			ID:        "token-1",
			UserID:    "user-1",
			Token:     "valid_refresh_token",
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}
)

// Helper functions for creating test data with overrides
func NewTestUser(overrides ...func(*TestUser)) TestUser {
	now := time.Now()
	user := TestUser{
		ID:        fmt.Sprintf("test-user-%d", now.UnixNano()),
		Name:      "Test User",
		Email:     fmt.Sprintf("test-%d@example.com", now.UnixNano()),
		Password:  "$2a$10$hashedpassword",
		UserType:  "individual",
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&user)
	}

	return user
}

func NewTestTeam(ownerID string, overrides ...func(*TestTeam)) TestTeam {
	now := time.Now()
	team := TestTeam{
		ID:        fmt.Sprintf("test-team-%d", now.UnixNano()),
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

func NewTestToken(userID string, overrides ...func(*TestToken)) TestToken {
	now := time.Now()
	token := TestToken{
		ID:        fmt.Sprintf("test-token-%d", now.UnixNano()),
		UserID:    userID,
		Token:     fmt.Sprintf("test_token_%d", now.UnixNano()),
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, override := range overrides {
		override(&token)
	}

	return token
}