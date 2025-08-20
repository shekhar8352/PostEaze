package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/shekhar8352/PostEaze/utils/database"
)

// TestDBConfig holds configuration for test database
type TestDBConfig struct {
	Driver   string
	URL      string
	InMemory bool
}

// GetTestDBConfig returns configuration for test database
func GetTestDBConfig() TestDBConfig {
	// Use PostgreSQL for testing to avoid CGO compilation issues with SQLite
	// This assumes a test database is available
	return TestDBConfig{
		Driver:   "postgres",
		URL:      "postgres://postgres:password@localhost:5432/postease_test?sslmode=disable",
		InMemory: false,
	}
}

// SetupTestDB creates and configures a test database connection
// Returns the database connection and a cleanup function
func SetupTestDB(ctx context.Context) (*sql.DB, func(), error) {
	config := GetTestDBConfig()
	
	db, err := sql.Open(config.Driver, config.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	// Configure connection pool for testing
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Create tables
	if err := CreateTestTables(ctx, db); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create test tables: %w", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing test database: %v", err)
		}
	}

	return db, cleanup, nil
}

// CreateTestTables creates all necessary tables for testing
func CreateTestTables(ctx context.Context, db *sql.DB) error {
	// Read the migration file
	migrationPath := filepath.Join("..", "..", "migrations", "001_create_initial_tables.up.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		// If migration file is not found, create tables manually
		return createTablesManually(ctx, db)
	}

	// Convert PostgreSQL migration to SQLite compatible
	sqliteSQL := convertPostgresToSQLite(string(migrationSQL))
	
	// Execute the migration
	if _, err := db.ExecContext(ctx, sqliteSQL); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// createTablesManually creates tables when migration file is not available
func createTablesManually(ctx context.Context, db *sql.DB) error {
	// Check if we're using PostgreSQL or SQLite and adjust accordingly
	config := GetTestDBConfig()
	
	var tables []string
	if config.Driver == "postgres" {
		tables = []string{
			`CREATE TABLE IF NOT EXISTS users (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name TEXT NOT NULL,
				email TEXT UNIQUE NOT NULL,
				password TEXT NOT NULL,
				user_type TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)`,
			`CREATE TABLE IF NOT EXISTS teams (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name TEXT NOT NULL,
				owner_id UUID NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
			)`,
			`CREATE TABLE IF NOT EXISTS team_members (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				team_id UUID NOT NULL,
				user_id UUID NOT NULL,
				role TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT fk_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
				CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				UNIQUE (team_id, user_id)
			)`,
			`CREATE TABLE IF NOT EXISTS refresh_tokens (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				token TEXT NOT NULL UNIQUE,
				expires_at TIMESTAMP NOT NULL,
				revoked BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT fk_refresh_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			)`,
		}
	} else {
		// SQLite tables
		tables = []string{
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
	}

	for _, table := range tables {
		if _, err := db.ExecContext(ctx, table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// convertPostgresToSQLite converts PostgreSQL specific syntax to SQLite
func convertPostgresToSQLite(postgresSQL string) string {
	// Basic conversions for SQLite compatibility
	// Remove PostgreSQL extensions
	sqliteSQL := postgresSQL
	
	// Remove CREATE EXTENSION statements
	sqliteSQL = removeLines(sqliteSQL, "CREATE EXTENSION")
	
	// Convert UUID to TEXT
	sqliteSQL = replaceAll(sqliteSQL, "UUID", "TEXT")
	
	// Convert gen_random_uuid() to a simple text placeholder
	sqliteSQL = replaceAll(sqliteSQL, "DEFAULT gen_random_uuid()", "")
	
	// Convert TIMESTAMP to DATETIME
	sqliteSQL = replaceAll(sqliteSQL, "TIMESTAMP", "DATETIME")
	
	// Convert CURRENT_TIMESTAMP to CURRENT_TIMESTAMP (SQLite compatible)
	// No change needed for CURRENT_TIMESTAMP
	
	return sqliteSQL
}

// Helper functions for string manipulation
func removeLines(text, pattern string) string {
	lines := splitLines(text)
	var result []string
	for _, line := range lines {
		if !contains(line, pattern) {
			result = append(result, line)
		}
	}
	return joinLines(result)
}

func replaceAll(text, old, new string) string {
	// Simple string replacement
	result := text
	for contains(result, old) {
		result = replace(result, old, new)
	}
	return result
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) && findSubstring(text, substr) >= 0
}

func replace(text, old, new string) string {
	pos := findSubstring(text, old)
	if pos < 0 {
		return text
	}
	return text[:pos] + new + text[pos+len(old):]
}

func findSubstring(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func splitLines(text string) []string {
	var lines []string
	var current string
	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

// CleanupTestData removes all data from test tables
func CleanupTestData(ctx context.Context, db *sql.DB) error {
	tables := []string{
		"refresh_tokens",
		"team_members", 
		"teams",
		"users",
	}

	for _, table := range tables {
		if _, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to cleanup table %s: %w", table, err)
		}
	}

	return nil
}

// LoadFixtures loads test fixtures into the database
func LoadFixtures(ctx context.Context, db *sql.DB, fixtures ...interface{}) error {
	for _, fixture := range fixtures {
		if err := loadFixture(ctx, db, fixture); err != nil {
			return fmt.Errorf("failed to load fixture: %w", err)
		}
	}
	return nil
}

// loadFixture loads a single fixture based on its type
func loadFixture(ctx context.Context, db *sql.DB, fixture interface{}) error {
	switch f := fixture.(type) {
	case UserFixture:
		return loadUserFixture(ctx, db, f)
	case TeamFixture:
		return loadTeamFixture(ctx, db, f)
	case RefreshTokenFixture:
		return loadRefreshTokenFixture(ctx, db, f)
	case []UserFixture:
		for _, user := range f {
			if err := loadUserFixture(ctx, db, user); err != nil {
				return err
			}
		}
	case []TeamFixture:
		for _, team := range f {
			if err := loadTeamFixture(ctx, db, team); err != nil {
				return err
			}
		}
	case []RefreshTokenFixture:
		for _, token := range f {
			if err := loadRefreshTokenFixture(ctx, db, token); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported fixture type: %T", fixture)
	}
	return nil
}

// loadUserFixture loads a user fixture
func loadUserFixture(ctx context.Context, db *sql.DB, user UserFixture) error {
	query := `INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.ExecContext(ctx, query, 
		user.ID, user.Name, user.Email, user.Password, 
		user.UserType, user.CreatedAt, user.UpdatedAt)
	
	return err
}

// loadTeamFixture loads a team fixture
func loadTeamFixture(ctx context.Context, db *sql.DB, team TeamFixture) error {
	query := `INSERT INTO teams (id, name, owner_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`
	
	_, err := db.ExecContext(ctx, query, 
		team.ID, team.Name, team.OwnerID, team.CreatedAt, team.UpdatedAt)
	
	return err
}

// loadRefreshTokenFixture loads a refresh token fixture
func loadRefreshTokenFixture(ctx context.Context, db *sql.DB, token RefreshTokenFixture) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, revoked, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.ExecContext(ctx, query, 
		token.ID, token.UserID, token.Token, token.ExpiresAt, 
		token.Revoked, token.CreatedAt, token.UpdatedAt)
	
	return err
}

// SetupTestDBWithDatabase initializes the database package with test database
func SetupTestDBWithDatabase(ctx context.Context) (*sql.DB, func(), error) {
	db, cleanup, err := SetupTestDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Initialize the database package with test configuration
	testConfig := GetTestDBConfig()
	config := database.Config{
		DriverName:            testConfig.Driver,
		URL:                   testConfig.URL,
		MaxOpenConnections:    1,
		MaxIdleConnections:    1,
		ConnectionMaxLifetime: time.Hour,
		ConnectionMaxIdleTime: time.Minute * 10,
	}

	if err := database.Init(ctx, config); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failed to initialize database package: %w", err)
	}

	enhancedCleanup := func() {
		database.Close()
		cleanup()
	}

	return db, enhancedCleanup, nil
}

// BeginTestTransaction starts a transaction for testing
func BeginTestTransaction(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	return db.BeginTx(ctx, nil)
}

// RollbackTestTransaction rolls back a test transaction
func RollbackTestTransaction(tx *sql.Tx) {
	if tx != nil {
		_ = tx.Rollback()
	}
}