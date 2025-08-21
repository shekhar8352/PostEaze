package benchmarks

import (
	"testing"
	"time"

	"github.com/shekhar8352/PostEaze/utils"
)

// BenchmarkCompleteUserSignupWorkflow benchmarks the complete user signup workflow
func BenchmarkCompleteUserSignupWorkflow(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping workflow benchmark - requires global database setup")
}

// BenchmarkCompleteUserLoginWorkflow benchmarks the complete user login workflow
func BenchmarkCompleteUserLoginWorkflow(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping workflow benchmark - requires global database setup")
}

// BenchmarkUserAuthenticationFlow benchmarks user authentication validation
func BenchmarkUserAuthenticationFlow(b *testing.B) {
	// Generate test tokens
	userID := "test-user-123"
	role := "individual"
	
	accessToken, err := utils.GenerateAccessToken(userID, role)
	if err != nil {
		b.Fatalf("Failed to generate test token: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Step 1: Parse and validate token
		claims, err := utils.ParseToken(accessToken, false)
		if err != nil {
			b.Fatalf("Token parsing failed: %v", err)
		}

		// Step 2: Extract user ID from token
		extractedUserID, err := utils.GetUserIDFromToken(accessToken)
		if err != nil {
			b.Fatalf("User ID extraction failed: %v", err)
		}

		// Verify extracted data
		if claims.UserID != userID || extractedUserID != userID {
			b.Fatalf("Token validation failed: expected %s, got %s", userID, claims.UserID)
		}
	}
}

// BenchmarkTeamCreationWorkflow benchmarks the complete team creation workflow
func BenchmarkTeamCreationWorkflow(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping workflow benchmark - requires global database setup")
}

// BenchmarkUserSessionManagement benchmarks user session management operations
func BenchmarkUserSessionManagement(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping workflow benchmark - requires global database setup")
}

// BenchmarkUserDataValidation benchmarks user data validation operations
func BenchmarkUserDataValidation(b *testing.B) {
	// Test data for validation
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"firstname+lastname@company.org",
	}
	
	passwords := []string{
		"password123",
		"complexPassword!@#",
		"simplepass",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Cycle through test data
		email := validEmails[i%len(validEmails)]
		password := passwords[i%len(passwords)]

		// Step 1: Hash password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			b.Fatalf("Password hashing failed: %v", err)
		}

		// Step 2: Validate password hash
		if !utils.CheckPasswordHash(password, hashedPassword) {
			b.Fatalf("Password validation failed")
		}

		// Step 3: Basic email validation (simple check)
		if email == "" || len(email) < 5 {
			b.Fatalf("Email validation failed")
		}

		// Step 4: User type validation
		userType := "individual"
		if userType != "individual" && userType != "team" {
			b.Fatalf("User type validation failed")
		}
	}
}

// BenchmarkConcurrentUserOperations benchmarks concurrent user operations
func BenchmarkConcurrentUserOperations(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Generate and validate a token (no database required)
			userID := "concurrent-user-" + time.Now().Format("20060102150405.000000")
			token, err := utils.GenerateAccessToken(userID, "individual")
			if err == nil {
				_, _ = utils.ParseToken(token, false)
			}
		}
	})
}