package helpers

import (
	"context"
	"testing"
	"time"
)

func TestNewTestDB(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// Test that database is accessible
	if err := db.DB.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func TestLoadUserFixture(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// Create test user
	user := NewTestUser(func(u *TestUser) {
		u.ID = "test-user-1"
		u.Name = "Test User"
		u.Email = "test@example.com"
	})

	// Load fixture
	if err := db.LoadFixture(user); err != nil {
		t.Fatalf("Failed to load user fixture: %v", err)
	}

	// Verify user was loaded
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", user.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func TestLoadMultipleFixtures(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// Load default users
	if err := db.LoadFixture(DefaultUsers); err != nil {
		t.Fatalf("Failed to load user fixtures: %v", err)
	}

	// Verify users were loaded
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}

	if count != len(DefaultUsers) {
		t.Errorf("Expected %d users, got %d", len(DefaultUsers), count)
	}
}

func TestLoadTeamFixture(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// First load a user (required for foreign key)
	user := DefaultUsers[1] // user-2 who owns the team
	if err := db.LoadFixture(user); err != nil {
		t.Fatalf("Failed to load user fixture: %v", err)
	}

	// Load team
	team := DefaultTeams[0]
	if err := db.LoadFixture(team); err != nil {
		t.Fatalf("Failed to load team fixture: %v", err)
	}

	// Verify team was loaded
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM teams WHERE id = ?", team.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query team: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 team, got %d", count)
	}
}

func TestCleanData(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// Load some data
	if err := db.LoadFixture(DefaultUsers); err != nil {
		t.Fatalf("Failed to load fixtures: %v", err)
	}

	// Verify data exists
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}
	if count == 0 {
		t.Fatal("No users found after loading fixtures")
	}

	// Clean data
	if err := db.CleanData(); err != nil {
		t.Fatalf("Failed to clean data: %v", err)
	}

	// Verify data is gone
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users after cleanup: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users after cleanup, got %d", count)
	}
}

func TestBeginTx(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Insert data in transaction
	user := NewTestUser()
	_, err = tx.Exec("INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Name, user.Email, user.Password, user.UserType, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		t.Fatalf("Failed to insert user in transaction: %v", err)
	}

	// Rollback transaction
	if err := tx.Rollback(); err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// Verify data was not committed
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", user.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users after rollback, got %d", count)
	}
}

func TestNewTestUserWithOverrides(t *testing.T) {
	user := NewTestUser(func(u *TestUser) {
		u.Name = "Custom Name"
		u.Email = "custom@example.com"
		u.UserType = "team"
	})

	if user.Name != "Custom Name" {
		t.Errorf("Expected name 'Custom Name', got '%s'", user.Name)
	}
	if user.Email != "custom@example.com" {
		t.Errorf("Expected email 'custom@example.com', got '%s'", user.Email)
	}
	if user.UserType != "team" {
		t.Errorf("Expected user type 'team', got '%s'", user.UserType)
	}
}

func TestNewTestTeamWithOverrides(t *testing.T) {
	team := NewTestTeam("owner-123", func(t *TestTeam) {
		t.Name = "Custom Team"
	})

	if team.Name != "Custom Team" {
		t.Errorf("Expected name 'Custom Team', got '%s'", team.Name)
	}
	if team.OwnerID != "owner-123" {
		t.Errorf("Expected owner ID 'owner-123', got '%s'", team.OwnerID)
	}
}

func TestNewTestTokenWithOverrides(t *testing.T) {
	futureTime := time.Now().Add(30 * 24 * time.Hour)
	token := NewTestToken("user-123", func(t *TestToken) {
		t.Token = "custom_token"
		t.ExpiresAt = futureTime
		t.Revoked = true
	})

	if token.Token != "custom_token" {
		t.Errorf("Expected token 'custom_token', got '%s'", token.Token)
	}
	if token.UserID != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", token.UserID)
	}
	if !token.Revoked {
		t.Error("Expected token to be revoked")
	}
	if !token.ExpiresAt.Equal(futureTime) {
		t.Errorf("Expected expires at %v, got %v", futureTime, token.ExpiresAt)
	}
}