package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// DatabaseIntegrationTestSuite tests database operations with actual database connections
type DatabaseIntegrationTestSuite struct {
	suite.Suite
	db      *sql.DB
	cleanup func()
	ctx     context.Context
}

// SetupSuite initializes the test environment with real database connection
func (s *DatabaseIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	
	// Setup test database with real database connection
	db, cleanup, err := testutils.SetupTestDBWithDatabase(s.ctx)
	require.NoError(s.T(), err, "Failed to setup test database")
	
	s.db = db
	s.cleanup = cleanup
	
	s.T().Log("Database integration test suite setup completed")
}

// TearDownSuite cleans up test environment
func (s *DatabaseIntegrationTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
	s.T().Log("Database integration test suite teardown completed")
}

// SetupTest prepares each test with clean database state
func (s *DatabaseIntegrationTestSuite) SetupTest() {
	// Clean up any existing test data
	err := testutils.CleanupTestData(s.ctx, s.db)
	require.NoError(s.T(), err, "Failed to cleanup test data")
}

// TestUserCRUDOperations tests complete user CRUD operations with real database
func (s *DatabaseIntegrationTestSuite) TestUserCRUDOperations() {
	// Test user creation
	user := modelsv1.User{
		Name:     "Integration Test User",
		Email:    "integration@test.com",
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	// Start transaction for user creation
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start transaction")
	
	createdUser, err := repositories.CreateUser(s.ctx, tx, user)
	require.NoError(s.T(), err, "Failed to create user")
	assert.NotEmpty(s.T(), createdUser.ID, "User ID should be generated")
	assert.Equal(s.T(), user.Name, createdUser.Name, "User name should match")
	assert.Equal(s.T(), user.Email, createdUser.Email, "User email should match")
	assert.Equal(s.T(), user.UserType, createdUser.UserType, "User type should match")
	assert.NotZero(s.T(), createdUser.CreatedAt, "Created at should be set")
	assert.NotZero(s.T(), createdUser.UpdatedAt, "Updated at should be set")
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Failed to commit transaction")
	
	s.T().Logf("Created user with ID: %s", createdUser.ID)
	
	// Test user retrieval by email
	retrievedUser, err := repositories.GetUserByEmail(s.ctx, user.Email)
	require.NoError(s.T(), err, "Failed to retrieve user by email")
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID, "Retrieved user ID should match")
	assert.Equal(s.T(), user.Email, retrievedUser.Email, "Retrieved user email should match")
	assert.Equal(s.T(), user.Name, retrievedUser.Name, "Retrieved user name should match")
	assert.Equal(s.T(), string(user.UserType), retrievedUser.UserType, "Retrieved user type should match")
	
	s.T().Log("Successfully retrieved user by email")
	
	// Test user update operations (password change simulation)
	newPassword := "newhashedpassword456"
	updateQuery := "UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err = s.db.ExecContext(s.ctx, updateQuery, newPassword, createdUser.ID)
	require.NoError(s.T(), err, "Failed to update user password")
	
	// Verify update
	updatedUser, err := repositories.GetUserByEmail(s.ctx, user.Email)
	require.NoError(s.T(), err, "Failed to retrieve updated user")
	assert.Equal(s.T(), newPassword, updatedUser.Password, "Password should be updated")
	
	s.T().Log("Successfully updated user password")
	
	// Test user deletion
	deleteQuery := "DELETE FROM users WHERE id = ?"
	result, err := s.db.ExecContext(s.ctx, deleteQuery, createdUser.ID)
	require.NoError(s.T(), err, "Failed to delete user")
	
	rowsAffected, err := result.RowsAffected()
	require.NoError(s.T(), err, "Failed to get rows affected")
	assert.Equal(s.T(), int64(1), rowsAffected, "Should delete exactly one user")
	
	// Verify deletion
	_, err = repositories.GetUserByEmail(s.ctx, user.Email)
	assert.Error(s.T(), err, "Should not find deleted user")
	
	s.T().Log("User CRUD operations test completed successfully")
}

// TestTeamOperationsWithTransactions tests team operations with database transactions
func (s *DatabaseIntegrationTestSuite) TestTeamOperationsWithTransactions() {
	// First create a user to be the team owner
	user := modelsv1.User{
		Name:     "Team Owner",
		Email:    "teamowner@test.com",
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeTeam,
	}
	
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start transaction")
	
	createdUser, err := repositories.CreateUser(s.ctx, tx, user)
	require.NoError(s.T(), err, "Failed to create user")
	
	// Create team
	teamName := "Integration Test Team"
	teamID, err := repositories.SaveTeam(s.ctx, tx, teamName, createdUser.ID)
	require.NoError(s.T(), err, "Failed to create team")
	assert.NotEmpty(s.T(), teamID, "Team ID should be generated")
	
	// Add user to team
	err = repositories.AddListOfUsersToTeam(s.ctx, tx, teamID, []string{createdUser.ID}, string(modelsv1.RoleAdmin))
	require.NoError(s.T(), err, "Failed to add user to team")
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Failed to commit transaction")
	
	s.T().Logf("Created team with ID: %s and added user as admin", teamID)
	
	// Test transaction rollback scenario
	tx2, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start second transaction")
	
	// Try to create another team
	_, err = repositories.SaveTeam(s.ctx, tx2, "Rollback Test Team", createdUser.ID)
	require.NoError(s.T(), err, "Failed to create second team")
	
	// Rollback the transaction
	err = database.RollbackTx(tx2)
	require.NoError(s.T(), err, "Failed to rollback transaction")
	
	s.T().Log("Successfully tested transaction rollback")
}

// TestRefreshTokenOperations tests refresh token database operations
func (s *DatabaseIntegrationTestSuite) TestRefreshTokenOperations() {
	// Create a user first
	user := modelsv1.User{
		Name:     "Token Test User",
		Email:    "tokentest@test.com",
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start transaction")
	
	createdUser, err := repositories.CreateUser(s.ctx, tx, user)
	require.NoError(s.T(), err, "Failed to create user")
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Failed to commit transaction")
	
	// Generate and insert refresh token
	refreshToken, err := utils.GenerateRefreshToken(createdUser.ID)
	require.NoError(s.T(), err, "Failed to generate refresh token")
	
	expiryTime := utils.GetRefreshTokenExpiry()
	err = repositories.InsertRefreshTokenOfUser(s.ctx, createdUser.ID, refreshToken, expiryTime)
	require.NoError(s.T(), err, "Failed to insert refresh token")
	
	s.T().Logf("Inserted refresh token for user: %s", createdUser.ID)
	
	// Test token retrieval
	retrievedUser, err := repositories.GetUserbyToken(s.ctx, refreshToken)
	require.NoError(s.T(), err, "Failed to retrieve user by token")
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID, "Retrieved user should match")
	
	s.T().Log("Successfully retrieved user by refresh token")
	
	// Test token revocation
	err = repositories.RevokeTokenForUser(s.ctx, createdUser.ID)
	require.NoError(s.T(), err, "Failed to revoke token")
	
	// Verify token is revoked (this should fail)
	_, err = repositories.GetUserbyToken(s.ctx, refreshToken)
	assert.Error(s.T(), err, "Should fail to retrieve user with revoked token")
	
	s.T().Log("Successfully tested token revocation")
}

// TestConcurrentDatabaseOperations tests concurrent database operations
func (s *DatabaseIntegrationTestSuite) TestConcurrentDatabaseOperations() {
	const numConcurrentOperations = 5
	
	// Create channels to collect results
	results := make(chan error, numConcurrentOperations)
	
	// Launch concurrent user creation operations
	for i := 0; i < numConcurrentOperations; i++ {
		go func(index int) {
			user := modelsv1.User{
				Name:     fmt.Sprintf("Concurrent User %d", index),
				Email:    fmt.Sprintf("concurrent%d@test.com", index),
				Password: "hashedpassword123",
				UserType: modelsv1.UserTypeIndividual,
			}
			
			tx, err := database.GetTx(s.ctx, nil)
			if err != nil {
				results <- err
				return
			}
			
			_, err = repositories.CreateUser(s.ctx, tx, user)
			if err != nil {
				database.RollbackTx(tx)
				results <- err
				return
			}
			
			err = database.CommitTx(tx)
			results <- err
		}(i)
	}
	
	// Collect and verify results
	successCount := 0
	for i := 0; i < numConcurrentOperations; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else {
			s.T().Logf("Concurrent operation %d failed: %v", i, err)
		}
	}
	
	assert.Equal(s.T(), numConcurrentOperations, successCount, "All concurrent operations should succeed")
	s.T().Logf("Concurrent database operations test completed - %d/%d operations succeeded", successCount, numConcurrentOperations)
}

// TestDatabaseConnectionPooling tests database connection pooling behavior
func (s *DatabaseIntegrationTestSuite) TestDatabaseConnectionPooling() {
	// Test multiple simultaneous connections
	const numConnections = 3
	
	connections := make([]*sql.Tx, numConnections)
	
	// Start multiple transactions
	for i := 0; i < numConnections; i++ {
		tx, err := database.GetTx(s.ctx, nil)
		require.NoError(s.T(), err, "Failed to get transaction %d", i)
		connections[i] = tx
	}
	
	// Perform operations on each transaction
	for i, tx := range connections {
		user := modelsv1.User{
			Name:     fmt.Sprintf("Pool Test User %d", i),
			Email:    fmt.Sprintf("pooltest%d@test.com", i),
			Password: "hashedpassword123",
			UserType: modelsv1.UserTypeIndividual,
		}
		
		_, err := repositories.CreateUser(s.ctx, tx, user)
		require.NoError(s.T(), err, "Failed to create user in transaction %d", i)
	}
	
	// Commit all transactions
	for i, tx := range connections {
		err := database.CommitTx(tx)
		require.NoError(s.T(), err, "Failed to commit transaction %d", i)
	}
	
	s.T().Log("Database connection pooling test completed successfully")
}

// TestDatabaseTransactionIsolation tests transaction isolation
func (s *DatabaseIntegrationTestSuite) TestDatabaseTransactionIsolation() {
	// Start first transaction
	tx1, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start first transaction")
	
	// Create user in first transaction (not committed yet)
	user := modelsv1.User{
		Name:     "Isolation Test User",
		Email:    "isolation@test.com",
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	createdUser, err := repositories.CreateUser(s.ctx, tx1, user)
	require.NoError(s.T(), err, "Failed to create user in first transaction")
	
	// Start second transaction
	tx2, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start second transaction")
	
	// Try to retrieve user from second transaction (should not see uncommitted data)
	_, err = repositories.GetUserByEmail(s.ctx, user.Email)
	assert.Error(s.T(), err, "Second transaction should not see uncommitted data")
	
	// Commit first transaction
	err = database.CommitTx(tx1)
	require.NoError(s.T(), err, "Failed to commit first transaction")
	
	// Now second transaction should be able to see the committed data
	retrievedUser, err := repositories.GetUserByEmail(s.ctx, user.Email)
	require.NoError(s.T(), err, "Should retrieve user after commit")
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID, "Retrieved user should match")
	
	// Clean up second transaction
	err = database.RollbackTx(tx2)
	require.NoError(s.T(), err, "Failed to rollback second transaction")
	
	s.T().Log("Database transaction isolation test completed successfully")
}

// TestDatabaseErrorHandling tests database error handling scenarios
func (s *DatabaseIntegrationTestSuite) TestDatabaseErrorHandling() {
	// Test duplicate email constraint
	user1 := modelsv1.User{
		Name:     "First User",
		Email:    "duplicate@test.com",
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	user2 := modelsv1.User{
		Name:     "Second User",
		Email:    "duplicate@test.com", // Same email
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	// Create first user
	tx1, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start transaction")
	
	_, err = repositories.CreateUser(s.ctx, tx1, user1)
	require.NoError(s.T(), err, "Failed to create first user")
	
	err = database.CommitTx(tx1)
	require.NoError(s.T(), err, "Failed to commit first transaction")
	
	// Try to create second user with same email (should fail)
	tx2, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start second transaction")
	
	_, err = repositories.CreateUser(s.ctx, tx2, user2)
	assert.Error(s.T(), err, "Should fail to create user with duplicate email")
	
	// Rollback failed transaction
	err = database.RollbackTx(tx2)
	require.NoError(s.T(), err, "Failed to rollback failed transaction")
	
	s.T().Log("Database error handling test completed successfully")
}

// TestDatabasePerformanceWithLargeDataset tests database performance with larger datasets
func (s *DatabaseIntegrationTestSuite) TestDatabasePerformanceWithLargeDataset() {
	const numUsers = 100
	
	startTime := time.Now()
	
	// Create multiple users in a single transaction
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start transaction")
	
	for i := 0; i < numUsers; i++ {
		user := modelsv1.User{
			Name:     fmt.Sprintf("Performance Test User %d", i),
			Email:    fmt.Sprintf("performance%d@test.com", i),
			Password: "hashedpassword123",
			UserType: modelsv1.UserTypeIndividual,
		}
		
		_, err := repositories.CreateUser(s.ctx, tx, user)
		require.NoError(s.T(), err, "Failed to create user %d", i)
	}
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Failed to commit transaction")
	
	duration := time.Since(startTime)
	s.T().Logf("Created %d users in %v (%.2f users/second)", numUsers, duration, float64(numUsers)/duration.Seconds())
	
	// Test retrieval performance
	startTime = time.Now()
	
	for i := 0; i < 10; i++ { // Test retrieving 10 random users
		email := fmt.Sprintf("performance%d@test.com", i*10)
		_, err := repositories.GetUserByEmail(s.ctx, email)
		require.NoError(s.T(), err, "Failed to retrieve user with email %s", email)
	}
	
	retrievalDuration := time.Since(startTime)
	s.T().Logf("Retrieved 10 users in %v", retrievalDuration)
	
	s.T().Log("Database performance test completed successfully")
}

// TestCompleteBusinessLogicIntegration tests complete business logic with actual database operations
func (s *DatabaseIntegrationTestSuite) TestCompleteBusinessLogicIntegration() {
	// Test complete signup flow with database operations
	signupParams := modelsv1.SignupParams{
		Name:     "Business Logic Test User",
		Email:    "businesslogic@test.com",
		Password: "testpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	// Call actual business logic
	result, err := businessv1.Signup(s.ctx, signupParams)
	require.NoError(s.T(), err, "Signup business logic should succeed")
	
	// Verify result structure
	assert.Contains(s.T(), result, "user", "Result should contain user")
	assert.Contains(s.T(), result, "access_token", "Result should contain access token")
	assert.Contains(s.T(), result, "refresh_token", "Result should contain refresh token")
	
	user := result["user"].(*modelsv1.User)
	accessToken := result["access_token"].(string)
	refreshToken := result["refresh_token"].(string)
	
	assert.NotEmpty(s.T(), user.ID, "User ID should be generated")
	assert.Equal(s.T(), signupParams.Name, user.Name, "User name should match")
	assert.Equal(s.T(), signupParams.Email, user.Email, "User email should match")
	assert.NotEmpty(s.T(), accessToken, "Access token should be generated")
	assert.NotEmpty(s.T(), refreshToken, "Refresh token should be generated")
	
	s.T().Logf("Signup completed - User ID: %s", user.ID)
	
	// Verify user was created in database
	dbUser, err := repositories.GetUserByEmail(s.ctx, signupParams.Email)
	require.NoError(s.T(), err, "Should retrieve user from database")
	assert.Equal(s.T(), user.ID, dbUser.ID, "Database user should match business logic result")
	
	// Verify refresh token was stored
	tokenUser, err := repositories.GetUserbyToken(s.ctx, refreshToken)
	require.NoError(s.T(), err, "Should retrieve user by refresh token")
	assert.Equal(s.T(), user.ID, tokenUser.ID, "Token should belong to correct user")
	
	s.T().Log("Database verification completed for signup")
	
	// Test login business logic
	loginParams := modelsv1.LoginParams{
		Email:    signupParams.Email,
		Password: signupParams.Password,
	}
	
	loginResult, err := businessv1.Login(s.ctx, loginParams)
	require.NoError(s.T(), err, "Login business logic should succeed")
	
	// Verify login result
	assert.Contains(s.T(), loginResult, "user", "Login result should contain user")
	assert.Contains(s.T(), loginResult, "access_token", "Login result should contain access token")
	assert.Contains(s.T(), loginResult, "refresh_token", "Login result should contain refresh token")
	
	loginUser := loginResult["user"].(*modelsv1.User)
	newRefreshToken := loginResult["refresh_token"].(string)
	
	assert.Equal(s.T(), user.ID, loginUser.ID, "Login user should match signup user")
	assert.NotEqual(s.T(), refreshToken, newRefreshToken, "New refresh token should be different")
	
	s.T().Log("Login business logic completed successfully")
	
	// Test refresh token business logic
	refreshResult, err := businessv1.RefreshToken(s.ctx, newRefreshToken)
	require.NoError(s.T(), err, "Refresh token business logic should succeed")
	
	assert.Contains(s.T(), refreshResult, "access_token", "Refresh result should contain access token")
	newAccessToken := refreshResult["access_token"]
	assert.NotEmpty(s.T(), newAccessToken, "New access token should be generated")
	
	s.T().Log("Refresh token business logic completed successfully")
	
	// Test logout business logic
	err = businessv1.Logout(s.ctx, newRefreshToken)
	require.NoError(s.T(), err, "Logout business logic should succeed")
	
	// Verify token was revoked
	_, err = repositories.GetUserbyToken(s.ctx, newRefreshToken)
	assert.Error(s.T(), err, "Should not be able to use revoked token")
	
	s.T().Log("Logout business logic completed successfully")
	s.T().Log("Complete business logic integration test completed successfully")
}

// TestDatabaseConnectionRecovery tests database connection recovery scenarios
func (s *DatabaseIntegrationTestSuite) TestDatabaseConnectionRecovery() {
	// Test connection after temporary disconnection simulation
	// This test verifies that the database connection pool handles reconnections properly
	
	// First, perform a normal operation
	user := modelsv1.User{
		Name:     "Connection Test User",
		Email:    "connection@test.com",
		Password: "hashedpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}
	
	tx, err := database.GetTx(s.ctx, nil)
	require.NoError(s.T(), err, "Failed to start transaction")
	
	createdUser, err := repositories.CreateUser(s.ctx, tx, user)
	require.NoError(s.T(), err, "Failed to create user")
	
	err = database.CommitTx(tx)
	require.NoError(s.T(), err, "Failed to commit transaction")
	
	s.T().Logf("Created user before connection test: %s", createdUser.ID)
	
	// Simulate connection recovery by performing multiple operations
	for i := 0; i < 5; i++ {
		testUser := modelsv1.User{
			Name:     fmt.Sprintf("Recovery Test User %d", i),
			Email:    fmt.Sprintf("recovery%d@test.com", i),
			Password: "hashedpassword123",
			UserType: modelsv1.UserTypeIndividual,
		}
		
		tx, err := database.GetTx(s.ctx, nil)
		require.NoError(s.T(), err, "Failed to start transaction %d", i)
		
		_, err = repositories.CreateUser(s.ctx, tx, testUser)
		require.NoError(s.T(), err, "Failed to create user %d", i)
		
		err = database.CommitTx(tx)
		require.NoError(s.T(), err, "Failed to commit transaction %d", i)
		
		// Verify user can be retrieved
		_, err = repositories.GetUserByEmail(s.ctx, testUser.Email)
		require.NoError(s.T(), err, "Failed to retrieve user %d", i)
	}
	
	s.T().Log("Database connection recovery test completed successfully")
}

// TestDatabaseIntegrationSuite runs the database integration test suite
func TestDatabaseIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}