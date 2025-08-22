package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tests/helpers"
)

// TestCompleteAuthenticationFlow tests the complete authentication workflow
func TestCompleteAuthenticationFlow(t *testing.T) {
	// Setup simple test environment
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()

	// Test individual user signup workflow
	t.Run("Individual User Signup and Login", func(t *testing.T) {
		// Create test user data
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "individual@test.com"
			u.Name = "Individual User"
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

		// Test login request
		loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = user.Email
			r.Password = "testpassword" // Raw password for login
		})

		response := http.Request("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, 200, response.Code, "Login should succeed")

		// Parse response
		var result map[string]interface{}
		err = helpers.GetResponseJSON(response, &result)
		require.NoError(t, err, "Should parse login response")

		assert.Equal(t, "success", result["status"])
		t.Logf("Individual user login test completed successfully")
	})

	// Test team user workflow
	t.Run("Team User Signup and Login", func(t *testing.T) {
		// Create test team user
		teamUser := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = "teamowner@test.com"
			u.Name = "Team Owner"
			u.UserType = modelsv1.UserTypeTeam
		})

		// Load user into database
		err := db.LoadFixture(helpers.TestUser{
			ID:       teamUser.ID,
			Name:     teamUser.Name,
			Email:    teamUser.Email,
			Password: teamUser.Password,
			UserType: string(teamUser.UserType),
		})
		require.NoError(t, err, "Failed to load team user fixture")

		// Create associated team
		team := helpers.CreateTeam(teamUser.ID, func(t *modelsv1.Team) {
			t.Name = "Test Team"
		})

		err = db.LoadFixture(helpers.TestTeam{
			ID:      team.ID,
			Name:    team.Name,
			OwnerID: team.OwnerID,
		})
		require.NoError(t, err, "Failed to load team fixture")

		// Test login
		loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = teamUser.Email
			r.Password = "testpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, 200, response.Code, "Team user login should succeed")

		t.Logf("Team user login test completed successfully")
	})

	// Test authentication error scenarios
	t.Run("Authentication Error Scenarios", func(t *testing.T) {
		// Test invalid login
		invalidLogin := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = "nonexistent@test.com"
			r.Password = "wrongpassword"
		})

		response := http.Request("POST", "/api/v1/auth/login", invalidLogin)
		assert.Equal(t, 401, response.Code, "Invalid login should return 401")

		// Test missing credentials
		emptyLogin := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
			r.Email = ""
			r.Password = ""
		})

		response = http.Request("POST", "/api/v1/auth/login", emptyLogin)
		assert.Equal(t, 400, response.Code, "Empty credentials should return 400")

		t.Logf("Authentication error scenarios tested successfully")
	})
}

// TestConcurrentAuthenticationRequests tests concurrent authentication requests
func TestConcurrentAuthenticationRequests(t *testing.T) {
	db, err := helpers.NewTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Cleanup()

	http := helpers.NewHTTPTest()
	const numConcurrentRequests = 3

	// Create test users for concurrent login
	for i := 0; i < numConcurrentRequests; i++ {
		user := helpers.CreateUser(func(u *modelsv1.User) {
			u.Email = fmt.Sprintf("concurrent%d@test.com", i)
			u.Name = fmt.Sprintf("Concurrent User %d", i)
		})

		err := db.LoadFixture(helpers.TestUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserType: string(user.UserType),
		})
		require.NoError(t, err, "Failed to load user fixture %d", i)
	}

	// Test concurrent login requests
	results := make(chan bool, numConcurrentRequests)

	for i := 0; i < numConcurrentRequests; i++ {
		go func(index int) {
			loginReq := helpers.NewLoginRequest(func(r *helpers.LoginRequest) {
				r.Email = fmt.Sprintf("concurrent%d@test.com", index)
				r.Password = "testpassword"
			})

			response := http.Request("POST", "/api/v1/auth/login", loginReq)
			results <- response.Code == 200
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < numConcurrentRequests; i++ {
		if <-results {
			successCount++
		}
	}

	assert.Equal(t, numConcurrentRequests, successCount, "All concurrent requests should succeed")
	t.Logf("Concurrent authentication test completed - %d/%d requests succeeded", successCount, numConcurrentRequests)
}