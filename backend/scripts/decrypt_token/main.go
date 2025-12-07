package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Initialize encryption
	if err := encryption.Init(); err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	// Build the database URL from environment variables
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgURL := os.Getenv("POSTGRES_URL")
	pgDB := os.Getenv("POSTGRES_DB")

	if pgUser == "" || pgPass == "" || pgURL == "" || pgDB == "" {
		log.Fatal("Missing database environment variables (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_URL, POSTGRES_DB)")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUser, pgPass, pgURL, pgDB)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Query the channel token from the channel_tokens table joined with channels
	channelID := int64(7)
	var accessTokenEncrypted []byte
	var displayName sql.NullString

	err = db.QueryRow(`
		SELECT c.display_name, ct.access_token 
		FROM channels c 
		JOIN channel_tokens ct ON c.id = ct.channel_id 
		WHERE c.id = $1 AND ct.revoked = false
		ORDER BY ct.created_at DESC
		LIMIT 1`,
		channelID).Scan(&displayName, &accessTokenEncrypted)
	if err != nil {
		log.Fatalf("Failed to query channel: %v", err)
	}

	// Decrypt the access token
	accessToken, err := encryption.Decrypt(accessTokenEncrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt access token: %v", err)
	}

	name := "<no display name>"
	if displayName.Valid {
		name = displayName.String
	}

	fmt.Printf("Channel ID: %d\n", channelID)
	fmt.Printf("Channel Name: %s\n", name)
	fmt.Printf("Access Token: %s\n", accessToken)
}
