package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
)

// TestEndToEndUserJourney tests a complete user journey from signup to logout
func TestEndToEndUserJourney(t *testing.T) {
	// Setup simple test environment
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()

	t.Run("Individual User Complete Journey", func(t *testing.T) {
		// Step 1: Create and setup user
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "journey@test.com"
			u.Name = "Journey Test User"
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to setup user")

		t.Logf("Step 1 completed: User created with ID: %s", user.ID)

		// Step 2: User login
		loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = user.Email
			r.Password = "testpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, 200, response.Code, "Login should succeed")

		t.Log("Step 2 completed: User logged in successfully")

		// Step 3: Test authenticated request
		response = http.AuthRequest("GET", "/api/v1/auth/refresh", user.ID, nil)
		assert.Equal(t, 200, response.Code, "Authenticated request should succeed")

		t.Log("Step 3 completed: Authenticated request successful")

		// Step 4: Test logout
		response = http.Request("POST", "/api/v1/auth/logout", nil)
		assert.Equal(t, 200, response.Code, "Logout should succeed")

		t.Log("Step 4 completed: User logged out successfully")

		t.Log("Individual user complete journey test completed successfully")
	})

	t.Run("Team User Complete Journey", func(t *testing.T) {
		// Step 1: Create team owner
		teamOwner := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "teamjourney@test.com"
			u.Name = "Team Journey Owner"
			u.UserType = modelsv1.UserTypeTeam
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       teamOwner.ID,
			Name:     teamOwner.Name,
			Email:    teamOwner.Email,
			Password: teamOwner.Password,
			UserType: string(teamOwner.UserType),
		})
		require.NoError(t, err, "Failed to setup team owner")

		// Step 2: Create team
		team := helpers.CreateTeam(teamOwner.ID, func(t *modelsv1.Team) {
			t.Name = "Journey Test Team"
		})

		err = db.LoadFixture(helpers.TestTeam{
			ID:      team.ID,
			Name:    team.Name,
			OwnerID: team.OwnerID,
		})
		require.NoError(t, err, "Failed to setup team")

		t.Logf("Steps 1-2 completed: Team owner and team created")

		// Step 3: Team owner login
		loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = teamOwner.Email
			r.Password = "testpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, 200, response.Code, "Team owner login should succeed")

		t.Log("Step 3 completed: Team owner logged in successfully")

		// Step 4: Test team operations (simulated)
		response = http.AuthRequestWithRole("GET", "/api/v1/auth/refresh", teamOwner.ID, "admin", nil)
		assert.Equal(t, 200, response.Code, "Team admin request should succeed")

		t.Log("Step 4 completed: Team operations successful")

		t.Log("Team user complete journey test completed successfully")
	})
}

// TestEndToEndErrorRecovery tests error recovery scenarios
func TestEndToEndErrorRecovery(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()

	t.Run("Authentication Error Recovery", func(t *testing.T) {
		// Test invalid login
		invalidLogin := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = "nonexistent@test.com"
			r.Password = "wrongpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", invalidLogin)
		assert.Equal(t, 401, response.Code, "Invalid login should return 401")

		// Test recovery with valid login after creating user
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "recovery@test.com"
			u.Name = "Recovery User"
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to setup recovery user")

		validLogin := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = user.Email
			r.Password = "testpassword"
		})

		response = http.Request("POST", "/api/v1/auth/login", validLogin)
		assert.Equal(t, 200, response.Code, "Valid login should succeed after error")

		t.Log("Authentication error recovery test completed successfully")
	})

	t.Run("Database Error Recovery", func(t *testing.T) {
		// Test duplicate user creation (should fail)
		user1 := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "duplicate@test.com"
			u.Name = "First User"
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user1.ID,
			Name:     user1.Name,
			Email:    user1.Email,
			Password: user1.Password,
			UserType: string(user1.UserType),
		})
		require.NoError(t, err, "Failed to create first user")

		// Try to create second user with same email (should fail)
		user2 := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "duplicate@test.com" // Same email
			u.Name = "Second User"
		})

		err = db.LoadFixture(helpers.TestUser{
			ID:       user2.ID,
			Name:     user2.Name,
			Email:    user2.Email,
			Password: user2.Password,
			UserType: string(user2.UserType),
		})
		assert.Error(t, err, "Should fail to create user with duplicate email")

		// Test recovery with unique email
		user3 := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "unique@test.com"
			u.Name = "Unique User"
		})

		err = db.LoadFixture(helpers.TestUser{
			ID:       user3.ID,
			Name:     user3.Name,
			Email:    user3.Email,
			Password: user3.Password,
			UserType: string(user3.UserType),
		})
		require.NoError(t, err, "Should succeed to create user with unique email")

		t.Log("Database error recovery test completed successfully")
	})
}

// TestEndToEndConcurrentOperations tests concurrent end-to-end operations
func TestEndToEndConcurrentOperations(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()
	const numConcurrentUsers = 3

	// Create test users for concurrent operations
	for i := 0; i < numConcurrentUsers; i++ {
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = fmt.Sprintf("concurrent%d@e2e.com", i)
			u.Name = fmt.Sprintf("Concurrent E2E User %d", i)
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to create concurrent user %d", i)
	}

	t.Run("Concurrent User Journeys", func(t *testing.T) {
		results := make(chan bool, numConcurrentUsers)

		// Launch concurrent user journeys
		for i := 0; i < numConcurrentUsers; i++ {
			go func(index int) {
				success := performConcurrentUserJourney(http, index)
				results <- success
			}(i)
		}

		// Collect results
		successCount := 0
		for i := 0; i < numConcurrentUsers; i++ {
			if <-results {
				successCount++
			}
		}

		assert.Equal(t, numConcurrentUsers, successCount, "All concurrent user journeys should succeed")
		t.Logf("Concurrent end-to-end operations completed - %d/%d succeeded", successCount, numConcurrentUsers)
	})
}

// performConcurrentUserJourney performs a complete user journey for concurrent testing
func performConcurrentUserJourney(http *helpers.HTTPHelper, index int) bool {
	// Login
	loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
		r.Email = fmt.Sprintf("concurrent%d@e2e.com", index)
		r.Password = "testpassword"
	})

	response := http.Request("POST", "/api/v1/auth/login", loginReq)
	if response.Code != 200 {
		return false
	}

	// Authenticated request
	userID := fmt.Sprintf("concurrent-user-%d", index)
	response = http.AuthRequest("GET", "/api/v1/auth/refresh", userID, nil)
	if response.Code != 200 {
		return false
	}

	// Logout
	response = http.Request("POST", "/api/v1/auth/logout", nil)
	if response.Code != 200 {
		return false
	}

	return true
}

// TestEndToEndPerformance tests end-to-end performance scenarios
func TestEndToEndPerformance(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()

	// Create multiple users for performance testing
	const numUsers = 10
	for i := 0; i < numUsers; i++ {
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = fmt.Sprintf("perf%d@test.com", i)
			u.Name = fmt.Sprintf("Performance User %d", i)
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to create performance user %d", i)
	}

	t.Run("Sequential Login Performance", func(t *testing.T) {
		successCount := 0

		for i := 0; i < numUsers; i++ {
			loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
				r.Email = fmt.Sprintf("perf%d@test.com", i)
				r.Password = "testpassword"
			})

			response := http.Request("POST", "/api/v1/auth/login", loginReq)
			if response.Code == 200 {
				successCount++
			}
		}

		assert.Equal(t, numUsers, successCount, "All sequential logins should succeed")
		t.Logf("Sequential login performance test completed - %d/%d succeeded", successCount, numUsers)
	})

	t.Log("End-to-end performance test completed successfully")
}