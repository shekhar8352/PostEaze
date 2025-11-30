package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shekhar8352/PostEaze/services/email_service"
)

func main() {
	// Load .env file
	// adjusted to look for .env in project root relative to this script if run from here,
	// or just rely on user running it correctly.
	// For now keeping original logic but adding a comment or trying to load from specific path might be better.
	// Let's keep it simple and just copy.
	if err := godotenv.Load("../../.env"); err != nil {
		// Fallback to default if running from root
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Error loading .env file")
		}
	}

	// Debug: Check loaded environment variables
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fmt.Printf("DEBUG: SMTP_EMAIL='%s'\n", smtpEmail)
	if len(smtpPass) > 4 {
		fmt.Printf("DEBUG: SMTP_PASSWORD='...%s' (len=%d)\n", smtpPass[len(smtpPass)-4:], len(smtpPass))
	} else {
		fmt.Printf("DEBUG: SMTP_PASSWORD='%s' (len=%d)\n", smtpPass, len(smtpPass))
	}

	emailService := email_service.NewGmailEmailService()

	// Test Notification Email
	to := "abc@gmail.com" // Sending to self for testing
	err := emailService.SendNotificationEmail(to, "This is a test notification from PostEaze backend.")
	if err != nil {
		log.Fatalf("Failed to send notification email: %v", err)
	}
	fmt.Println("Notification email sent successfully!")

	// Test Team Invite Email
	inviteLink := "https://posteaze.com/join/team/123"
	err = emailService.SendTeamInviteEmail(to, inviteLink)
	if err != nil {
		log.Fatalf("Failed to send team invite email: %v", err)
	}
	fmt.Println("Team invite email sent successfully!")
}
