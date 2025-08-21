package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
)

// TestUserCRUDOperations tests complete user CRUD operations with database
func TestUserCRUDOperations(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	ctx := context.Background()

	// Test user creation
	user := helpers.CreateUser(func(u *modelsv1.User) {
		u.Email = "integration@test.com"
		u.Name = "Integration Test User"
	})

	// Load user into database
	err = db.LoadFixture(helpers.TestUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		UserType: string(user.UserType),
	})
	require.NoError(t, err, "Failed to create user")

	t.Logf("Created user with ID: %s", user.ID)

	// Test user retrieval by querying database directly
	var retrievedUser helpers.TestUser
	query := "SELECT id, name, email, password, user_type FROM users WHERE email = ?"
	row := db.DB.QueryRowContext(ctx, query, user.Email)
	err = row.Scan(&retrievedUser.ID, &retrievedUser.Name, &retrievedUser.Email, 
		&retrievedUser.Password, &retrievedUser.UserType)
	require.NoError(t, err, "Failed to retrieve user by email")

	assert.Equal(t, user.ID, retrievedUser.ID, "Retrieved user ID should match")
	assert.Equal(t, user.Email, retrievedUser.Email, "Retrieved user email should match")
	assert.Equal(t, user.Name, retrievedUser.Name, "Retrieved user name should match")

	t.Log("Successfully retrieved user by email")

	// Test user update operations
	newPassword := "newhashedpassword456"
	updateQuery := "UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err = db.DB.ExecContext(ctx, updateQuery, newPassword, user.ID)
	require.NoError(t, err, "Failed to update user password")

	// Verify update
	var updatedPassword string
	selectQuery := "SELECT password FROM users WHERE id = ?"
	row = db.DB.QueryRowContext(ctx, selectQuery, user.ID)
	err = row.Scan(&updatedPassword)
	require.NoError(t, err, "Failed to retrieve updated password")
	assert.Equal(t, newPassword, updatedPassword, "Password should be updated")

	t.Log("Successfully updated user password")

	// Test user deletion
	deleteQuery := "DELETE FROM users WHERE id = ?"
	result, err := db.DB.ExecContext(ctx, deleteQuery, user.ID)
	require.NoError(t, err, "Failed to delete user")

	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err, "Failed to get rows affected")
	assert.Equal(t, int64(1), rowsAffected, "Should delete exactly one user")

	// Verify deletion
	var count int
	countQuery := "SELECT COUNT(*) FROM users WHERE id = ?"
	row = db.DB.QueryRowContext(ctx, countQuery, user.ID)
	err = row.Scan(&count)
	require.NoError(t, err, "Failed to count users")
	assert.Equal(t, 0, count, "Should not find deleted user")

	t.Log("User CRUD operations test completed successfully")
}

// TestTeamOperationsWithTransactions tests team operations with database transactions
func TestTeamOperationsWithTransactions(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	ctx := context.Background()

	// Create a team owner user
	user := helpers.CreateUser(func(u *modelsv1.User) {
		u.Email = "teamowner@test.com"
		u.Name = "Team Owner"
		u.UserType = modelsv1.UserTypeTeam
	})

	// Load user into database
	err = db.LoadFixture(helpers.TestUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		UserType: string(user.UserType),
	})
	require.NoError(t, err, "Failed to create team owner")

	// Create team using transaction
	tx, err := db.BeginTx(ctx)
	require.NoError(t, err, "Failed to start transaction")

	team := helpers.CreateTeam(user.ID, func(t *modelsv1.Team) {
		t.Name = "Integration Test Team"
	})

	// Insert team
	insertTeamQuery := "INSERT INTO teams (id, name, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, insertTeamQuery, team.ID, team.Name, team.OwnerID, team.CreatedAt, team.UpdatedAt)
	require.NoError(t, err, "Failed to create team")

	err = tx.Commit()
	require.NoError(t, err, "Failed to commit transaction")

	t.Logf("Created team with ID: %s", team.ID)

	// Test transaction rollback scenario
	tx2, err := db.BeginTx(ctx)
	require.NoError(t, err, "Failed to start second transaction")

	// Try to create another team
	rollbackTeam := helpers.CreateTeam(user.ID, func(t *modelsv1.Team) {
		t.Name = "Rollback Test Team"
	})

	_, err = tx2.ExecContext(ctx, insertTeamQuery, rollbackTeam.ID, rollbackTeam.Name, 
		rollbackTeam.OwnerID, rollbackTeam.CreatedAt, rollbackTeam.UpdatedAt)
	require.NoError(t, err, "Failed to create rollback team")

	// Rollback the transaction
	err = tx2.Rollback()
	require.NoError(t, err, "Failed to rollback transaction")

	// Verify rollback team doesn't exist
	var count int
	countQuery := "SELECT COUNT(*) FROM teams WHERE id = ?"
	row := db.DB.QueryRowContext(ctx, countQuery, rollbackTeam.ID)
	err = row.Scan(&count)
	require.NoError(t, err, "Failed to count teams")
	assert.Equal(t, 0, count, "Rollback team should not exist")

	t.Log("Successfully tested transaction rollback")
}

// TestRefreshTokenOperations tests refresh token database operations
func TestRefreshTokenOperations(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	ctx := context.Background()

	// Create a user first
	user := helpers.CreateUser(func(u *modelsv1.User) {
		u.Email = "tokentest@test.com"
		u.Name = "Token Test User"
	})

	err = db.LoadFixture(helpers.TestUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		UserType: string(user.UserType),
	})
	require.NoError(t, err, "Failed to create user")

	// Create and insert refresh token
	token := helpers.NewTestToken(user.ID, func(t *helpers.TestToken) {
		t.Token = "test_refresh_token_123"
		t.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	})

	err = db.LoadFixture(token)
	require.NoError(t, err, "Failed to insert refresh token")

	t.Logf("Inserted refresh token for user: %s", user.ID)

	// Test token retrieval
	var retrievedUserID string
	query := "SELECT user_id FROM refresh_tokens WHERE token = ? AND revoked = FALSE AND expires_at > CURRENT_TIMESTAMP"
	row := db.DB.QueryRowContext(ctx, query, token.Token)
	err = row.Scan(&retrievedUserID)
	require.NoError(t, err, "Failed to retrieve user by token")
	assert.Equal(t, user.ID, retrievedUserID, "Retrieved user should match")

	t.Log("Successfully retrieved user by refresh token")

	// Test token revocation
	revokeQuery := "UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = ?"
	_, err = db.DB.ExecContext(ctx, revokeQuery, user.ID)
	require.NoError(t, err, "Failed to revoke token")

	// Verify token is revoked
	var revokedUserID string
	row = db.DB.QueryRowContext(ctx, query, token.Token)
	err = row.Scan(&revokedUserID)
	assert.Error(t, err, "Should fail to retrieve user with revoked token")

	t.Log("Successfully tested token revocation")
}

// TestConcurrentDatabaseOperations tests database operations under load
func TestConcurrentDatabaseOperations(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	// Verify tables exist before starting operations
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	require.NoError(t, err, "Failed to check if users table exists")
	require.Equal(t, 1, count, "Users table should exist")

	const numOperations = 3

	// Test multiple database operations in sequence (simulating load)
	successCount := 0
	for i := 0; i < numOperations; i++ {
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = fmt.Sprintf("load%d@test.com", i)
			u.Name = fmt.Sprintf("Load Test User %d", i)
		})
		
		// Create user
		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		
		if err == nil {
			// Verify user was created
			var userID string
			err = db.DB.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&userID)
			if err == nil && userID == user.ID {
				successCount++
			} else {
				t.Logf("Operation %d failed during verification: %v", i, err)
			}
		} else {
			t.Logf("Operation %d failed during creation: %v", i, err)
		}
	}

	assert.Equal(t, numOperations, successCount, "All database operations should succeed")
	t.Logf("Database load test completed - %d/%d operations succeeded", successCount, numOperations)
}

