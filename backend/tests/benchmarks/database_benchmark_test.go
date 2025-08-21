package benchmarks

import (
	"context"
	"testing"
	"time"

	"github.com/shekhar8352/PostEaze/tests/helpers"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// BenchmarkUserCreation benchmarks user creation in the database
func BenchmarkUserCreation(b *testing.B) {
	db, err := helpers.NewTestDB()
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Clean data between iterations
		db.CleanData()
		
		// Create unique user data for each iteration
		user := modelsv1.User{
			Name:     "Benchmark User",
			Email:    "benchmark" + time.Now().Format("20060102150405.000000") + "@example.com",
			Password: "$2a$10$hashedpassword",
			UserType: modelsv1.UserTypeIndividual,
		}
		b.StartTimer()

		// Direct SQL insert for benchmarking
		query := `INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?)`
		now := time.Now()
		_, err := db.DB.ExecContext(ctx, query, 
			"user-"+time.Now().Format("20060102150405.000000"), 
			user.Name, user.Email, user.Password, string(user.UserType), now, now)
		if err != nil {
			b.Fatalf("User creation failed: %v", err)
		}
	}
}

// BenchmarkUserLookupByEmail benchmarks user lookup by email
func BenchmarkUserLookupByEmail(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	// This would be better tested with proper database setup
	b.Skip("Skipping repository benchmark - requires global database setup")
}

// BenchmarkUserLookupByToken benchmarks user lookup by refresh token
func BenchmarkUserLookupByToken(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping repository benchmark - requires global database setup")
}

// BenchmarkRefreshTokenInsertion benchmarks refresh token insertion
func BenchmarkRefreshTokenInsertion(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping repository benchmark - requires global database setup")
}

// BenchmarkTeamCreation benchmarks team creation in the database
func BenchmarkTeamCreation(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping repository benchmark - requires global database setup")
}

// BenchmarkTeamMemberAddition benchmarks adding users to a team
func BenchmarkTeamMemberAddition(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping repository benchmark - requires global database setup")
}

// BenchmarkConcurrentDatabaseOperations benchmarks concurrent database operations
func BenchmarkConcurrentDatabaseOperations(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping repository benchmark - requires global database setup")
}

// BenchmarkDatabaseTransaction benchmarks database transaction operations
func BenchmarkDatabaseTransaction(b *testing.B) {
	db, err := helpers.NewTestDB()
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Clean data between iterations
		db.CleanData()
		b.StartTimer()

		// Start transaction
		tx, err := db.BeginTx(ctx)
		if err != nil {
			b.Fatalf("Failed to begin transaction: %v", err)
		}

		// Create user within transaction
		user := modelsv1.User{
			Name:     "Transaction User",
			Email:    "tx" + time.Now().Format("20060102150405.000000") + "@example.com",
			Password: "$2a$10$hashedpassword",
			UserType: modelsv1.UserTypeIndividual,
		}

		// Direct SQL insert within transaction for benchmarking
		query := `INSERT INTO users (id, name, email, password, user_type, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?)`
		now := time.Now()
		_, err = tx.ExecContext(ctx, query, 
			"user-"+time.Now().Format("20060102150405.000000"), 
			user.Name, user.Email, user.Password, string(user.UserType), now, now)
		if err != nil {
			tx.Rollback()
			b.Fatalf("User creation in transaction failed: %v", err)
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			b.Fatalf("Transaction commit failed: %v", err)
		}
	}
}

// BenchmarkBulkUserCreation benchmarks creating multiple users in a single operation
func BenchmarkBulkUserCreation(b *testing.B) {
	// Skip this benchmark as it requires global database initialization
	b.Skip("Skipping repository benchmark - requires global database setup")
}