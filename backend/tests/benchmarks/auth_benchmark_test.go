package benchmarks

import (
	"testing"

	"github.com/shekhar8352/PostEaze/utils"
)

// BenchmarkJWTGeneration benchmarks JWT access token generation
func BenchmarkJWTGeneration(b *testing.B) {
	userID := "test-user-123"
	role := "individual"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := utils.GenerateAccessToken(userID, role)
		if err != nil {
			b.Fatalf("Failed to generate access token: %v", err)
		}
	}
}

// BenchmarkJWTRefreshTokenGeneration benchmarks JWT refresh token generation
func BenchmarkJWTRefreshTokenGeneration(b *testing.B) {
	userID := "test-user-123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := utils.GenerateRefreshToken(userID)
		if err != nil {
			b.Fatalf("Failed to generate refresh token: %v", err)
		}
	}
}

// BenchmarkJWTValidation benchmarks JWT token validation
func BenchmarkJWTValidation(b *testing.B) {
	// Generate a token to validate
	userID := "test-user-123"
	role := "individual"
	token, err := utils.GenerateAccessToken(userID, role)
	if err != nil {
		b.Fatalf("Failed to generate test token: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := utils.ParseToken(token, false)
		if err != nil {
			b.Fatalf("Failed to parse token: %v", err)
		}
	}
}

// BenchmarkPasswordHashing benchmarks password hashing
func BenchmarkPasswordHashing(b *testing.B) {
	password := "testpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := utils.HashPassword(password)
		if err != nil {
			b.Fatalf("Failed to hash password: %v", err)
		}
	}
}

// BenchmarkPasswordValidation benchmarks password validation
func BenchmarkPasswordValidation(b *testing.B) {
	password := "testpassword123"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password for benchmark: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valid := utils.CheckPasswordHash(password, hashedPassword)
		if !valid {
			b.Fatalf("Password validation failed")
		}
	}
}

// BenchmarkSignupFlow benchmarks the complete signup business logic
func BenchmarkSignupFlow(b *testing.B) {
	// Skip this benchmark as it requires full database initialization
	// This would be better tested in integration benchmarks
	b.Skip("Skipping business logic benchmark - requires full database setup")
}

// BenchmarkLoginFlow benchmarks the complete login business logic
func BenchmarkLoginFlow(b *testing.B) {
	// Skip this benchmark as it requires full database initialization
	// This would be better tested in integration benchmarks
	b.Skip("Skipping business logic benchmark - requires full database setup")
}

// BenchmarkRefreshTokenFlow benchmarks the refresh token business logic
func BenchmarkRefreshTokenFlow(b *testing.B) {
	// Skip this benchmark as it requires full database initialization
	// This would be better tested in integration benchmarks
	b.Skip("Skipping business logic benchmark - requires full database setup")
}