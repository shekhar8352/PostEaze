package helpers

import (
	"context"
	"testing"
	"time"
)

// Example test showing how to use the database helper
func TestDatabaseHelper_Example(t *testing.T) {
	// Create a new test database
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup() // Always cleanup after test

	// Example 1: Load predefined fixtures
	if err := db.LoadFixture(DefaultUsers); err != nil {
		t.Fatalf("Failed to load default users: %v", err)
	}

	// Example 2: Create custom test data
	customUser := NewTestUser(func(u *TestUser) {
		u.Name = "Custom User"
		u.Email = "custom@example.com"
		u.UserType = "team"
	})

	if err := db.LoadFixture(customUser); err != nil {
		t.Fatalf("Failed to load custom user: %v", err)
	}

	// Example 3: Query the database to verify data
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}

	expectedCount := len(DefaultUsers) + 1 // Default users + custom user
	if count != expectedCount {
		t.Errorf("Expected %d users, got %d", expectedCount, count)
	}

	// Example 4: Test with transactions
	ctx := context.Background()
	tx, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Insert data in transaction
	tempUser := NewTestUser(func(u *TestUser) {
		u.Name = "Temp User"
	})

	_, err = tx.Exec("INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		tempUser.ID, tempUser.Name, tempUser.Email, tempUser.Password, tempUser.UserType, tempUser.CreatedAt, tempUser.UpdatedAt)
	if err != nil {
		t.Fatalf("Failed to insert temp user: %v", err)
	}

	// Rollback to test isolation
	if err := tx.Rollback(); err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// Verify temp user was not committed
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", tempUser.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query temp user: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected temp user to be rolled back, but found %d", count)
	}

	// Example 5: Clean data between test cases
	if err := db.CleanData(); err != nil {
		t.Fatalf("Failed to clean data: %v", err)
	}

	// Verify data is cleaned
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users after cleanup: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users after cleanup, got %d", count)
	}
}

// Example test showing how to test with related data (foreign keys)
func TestDatabaseHelper_RelatedData(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// Load user first (required for team foreign key)
	user := NewTestUser(func(u *TestUser) {
		u.ID = "owner-123"
		u.Name = "Team Owner"
		u.UserType = "team"
	})

	if err := db.LoadFixture(user); err != nil {
		t.Fatalf("Failed to load user: %v", err)
	}

	// Load team with the user as owner
	team := NewTestTeam(user.ID, func(t *TestTeam) {
		t.Name = "Development Team"
	})

	if err := db.LoadFixture(team); err != nil {
		t.Fatalf("Failed to load team: %v", err)
	}

	// Load refresh token for the user
	token := NewTestToken(user.ID, func(t *TestToken) {
		t.Token = "valid_refresh_token"
		t.ExpiresAt = time.Now().Add(24 * time.Hour)
	})

	if err := db.LoadFixture(token); err != nil {
		t.Fatalf("Failed to load token: %v", err)
	}

	// Verify all related data exists
	var userCount, teamCount, tokenCount int

	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", user.ID).Scan(&userCount)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	err = db.DB.QueryRow("SELECT COUNT(*) FROM teams WHERE owner_id = ?", user.ID).Scan(&teamCount)
	if err != nil {
		t.Fatalf("Failed to query team: %v", err)
	}

	err = db.DB.QueryRow("SELECT COUNT(*) FROM refresh_tokens WHERE user_id = ?", user.ID).Scan(&tokenCount)
	if err != nil {
		t.Fatalf("Failed to query token: %v", err)
	}

	if userCount != 1 {
		t.Errorf("Expected 1 user, got %d", userCount)
	}
	if teamCount != 1 {
		t.Errorf("Expected 1 team, got %d", teamCount)
	}
	if tokenCount != 1 {
		t.Errorf("Expected 1 token, got %d", tokenCount)
	}
}

// Example test showing how to use the helper in a typical unit test scenario
func TestUserService_CreateUser_Example(t *testing.T) {
	// This is how you would use the database helper in a real unit test
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	// Test case: Create a new user
	testCases := []struct {
		name     string
		userData TestUser
		wantErr  bool
	}{
		{
			name: "valid user",
			userData: TestUser{
				ID:       "new-user-1",
				Name:     "New User",
				Email:    "newuser@example.com",
				Password: "hashedpassword",
				UserType: "individual",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			userData: TestUser{
				ID:       "new-user-2",
				Name:     "Another User",
				Email:    "newuser@example.com", // Same email as above
				Password: "hashedpassword",
				UserType: "individual",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean data before each test case
			if err := db.CleanData(); err != nil {
				t.Fatalf("Failed to clean data: %v", err)
			}

			// If this is the duplicate email test, load the first user
			if tc.name == "duplicate email" {
				firstUser := TestUser{
					ID:        "existing-user",
					Name:      "Existing User",
					Email:     "newuser@example.com",
					Password:  "hashedpassword",
					UserType:  "individual",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := db.LoadFixture(firstUser); err != nil {
					t.Fatalf("Failed to load existing user: %v", err)
				}
			}

			// Simulate creating a user (in real test, this would call your service)
			tc.userData.CreatedAt = time.Now()
			tc.userData.UpdatedAt = time.Now()
			
			err := db.LoadFixture(tc.userData)
			
			if tc.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}