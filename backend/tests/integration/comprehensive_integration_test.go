package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
)

// TestComprehensiveIntegrationWorkflow tests complete integration workflows
func TestComprehensiveIntegrationWorkflow(t *testing.T) {
	// Setup simple test environment
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()

	// Test complete user workflow
	t.Run("Complete User Workflow", func(t *testing.T) {
		// Create test user
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "workflow@test.com"
			u.Name = "Workflow Test User"
		})

		// Load user into database
		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to load user fixture")

		// Test authentication workflow
		loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = user.Email
			r.Password = "testpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, 200, response.Code, "Login should succeed")

		t.Log("Complete user workflow test completed successfully")
	})

	// Test team workflow
	t.Run("Team Workflow", func(t *testing.T) {
		// Create team owner
		teamOwner := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "teamowner@workflow.com"
			u.Name = "Team Owner"
			u.UserType = modelsv1.UserTypeTeam
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       teamOwner.ID,
			Name:     teamOwner.Name,
			Email:    teamOwner.Email,
			Password: teamOwner.Password,
			UserType: string(teamOwner.UserType),
		})
		require.NoError(t, err, "Failed to load team owner")

		// Create team
		team := helpers.CreateTeam(teamOwner.ID, func(t *modelsv1.Team) {
			t.Name = "Workflow Test Team"
		})

		err = db.LoadFixture(helpers.TestTeam{
			ID:      team.ID,
			Name:    team.Name,
			OwnerID: team.OwnerID,
		})
		require.NoError(t, err, "Failed to load team")

		// Test team authentication
		loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = teamOwner.Email
			r.Password = "testpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, 200, response.Code, "Team owner login should succeed")

		t.Log("Team workflow test completed successfully")
	})

	// Test concurrent operations
	t.Run("Concurrent Operations", func(t *testing.T) {
		const numConcurrentOps = 3
		results := make(chan bool, numConcurrentOps)

		// Create test users for concurrent operations
		for i := 0; i < numConcurrentOps; i++ {
			user := helpers.CreateUser(func(u *modelsv1.User) {
				u.Email = fmt.Sprintf("concurrent%d@workflow.com", i)
				u.Name = fmt.Sprintf("Concurrent User %d", i)
			})

			err := db.LoadFixture(helpers.TestUser{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
				UserType: string(user.UserType),
			})
			require.NoError(t, err, "Failed to load concurrent user %d", i)
		}

		// Test concurrent login operations
		for i := 0; i < numConcurrentOps; i++ {
			go func(index int) {
				loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
					r.Email = fmt.Sprintf("concurrent%d@workflow.com", index)
					r.Password = "testpassword"
				})

				response := http.Request("POST", "/api/v1/auth/login", loginReq)
				results <- response.Code == 200
			}(i)
		}

		// Collect results
		successCount := 0
		for i := 0; i < numConcurrentOps; i++ {
			if <-results {
				successCount++
			}
		}

		assert.Equal(t, numConcurrentOps, successCount, "All concurrent operations should succeed")
		t.Logf("Concurrent operations test completed - %d/%d succeeded", successCount, numConcurrentOps)
	})

	// Test error scenarios
	t.Run("Error Scenarios", func(t *testing.T) {
		// Test invalid login
		invalidLogin := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = "nonexistent@test.com"
			r.Password = "wrongpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", invalidLogin)
		assert.Equal(t, 401, response.Code, "Invalid login should return 401")

		// Test invalid endpoints
		response = http.Request("GET", "/api/v1/nonexistent", nil)
		assert.Equal(t, 404, response.Code, "Nonexistent endpoint should return 404")

		t.Log("Error scenarios test completed successfully")
	})

	t.Log("Comprehensive integration workflow test completed successfully")
}

// TestDatabaseIntegrationWorkflow tests database operations across different components
func TestDatabaseIntegrationWorkflow(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	ctx := context.Background()

	// Test user and team creation workflow
	t.Run("User and Team Creation Workflow", func(t *testing.T) {
		// Create user
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "dbworkflow@test.com"
			u.Name = "DB Workflow User"
			u.UserType = modelsv1.UserTypeTeam
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to create user")

		// Create team
		team := helpers.CreateTeam(user.ID, func(t *modelsv1.Team) {
			t.Name = "DB Workflow Team"
		})

		err = db.LoadFixture(helpers.TestTeam{
			ID:      team.ID,
			Name:    team.Name,
			OwnerID: team.OwnerID,
		})
		require.NoError(t, err, "Failed to create team")

		// Create refresh token
		token := helpers.NewTestToken(user.ID, func(t *helpers.TestToken) {
			t.Token = "db_workflow_token"
		})

		err = db.LoadFixture(token)
		require.NoError(t, err, "Failed to create token")

		// Verify all data exists
		var userCount, teamCount, tokenCount int

		// Check user
		row := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", user.ID)
		err = row.Scan(&userCount)
		require.NoError(t, err, "Failed to count users")
		assert.Equal(t, 1, userCount, "User should exist")

		// Check team
		row = db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM teams WHERE id = ?", team.ID)
		err = row.Scan(&teamCount)
		require.NoError(t, err, "Failed to count teams")
		assert.Equal(t, 1, teamCount, "Team should exist")

		// Check token
		row = db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM refresh_tokens WHERE id = ?", token.ID)
		err = row.Scan(&tokenCount)
		require.NoError(t, err, "Failed to count tokens")
		assert.Equal(t, 1, tokenCount, "Token should exist")

		t.Log("Database integration workflow test completed successfully")
	})

	// Test transaction workflow
	t.Run("Transaction Workflow", func(t *testing.T) {
		// Test successful transaction
		tx, err := db.BeginTx(ctx)
		require.NoError(t, err, "Failed to start transaction")

		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "transaction@test.com"
			u.Name = "Transaction User"
		})

		insertQuery := "INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err = tx.ExecContext(ctx, insertQuery, user.ID, user.Name, user.Email, 
			user.Password, string(user.UserType), user.CreatedAt, user.UpdatedAt)
		require.NoError(t, err, "Failed to insert user in transaction")

		err = tx.Commit()
		require.NoError(t, err, "Failed to commit transaction")

		// Verify user exists
		var count int
		row := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", user.ID)
		err = row.Scan(&count)
		require.NoError(t, err, "Failed to count users")
		assert.Equal(t, 1, count, "User should exist after commit")

		// Test rollback transaction
		tx2, err := db.BeginTx(ctx)
		require.NoError(t, err, "Failed to start rollback transaction")

		rollbackUser := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "rollback@test.com"
			u.Name = "Rollback User"
		})

		_, err = tx2.ExecContext(ctx, insertQuery, rollbackUser.ID, rollbackUser.Name, rollbackUser.Email, 
			rollbackUser.Password, string(rollbackUser.UserType), rollbackUser.CreatedAt, rollbackUser.UpdatedAt)
		require.NoError(t, err, "Failed to insert rollback user")

		err = tx2.Rollback()
		require.NoError(t, err, "Failed to rollback transaction")

		// Verify rollback user doesn't exist
		row = db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", rollbackUser.ID)
		err = row.Scan(&count)
		require.NoError(t, err, "Failed to count rollback users")
		assert.Equal(t, 0, count, "Rollback user should not exist")

		t.Log("Transaction workflow test completed successfully")
	})
}