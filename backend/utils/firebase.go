package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	client *auth.Client
}

type FirebaseUser struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var firebaseService *FirebaseService

// InitializeFirebase initializes the Firebase Admin SDK
func InitializeFirebase(ctx context.Context) error {
	// Create service account key from environment variables
	serviceAccount := map[string]string{
		"type":                        "service_account",
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        fmt.Sprintf("https://www.googleapis.com/robot/v1/metadata/x509/%s", os.Getenv("FIREBASE_CLIENT_EMAIL")),
	}

	serviceAccountJSON, err := json.Marshal(serviceAccount)
	if err != nil {
		return fmt.Errorf("failed to marshal service account: %w", err)
	}

	opt := option.WithCredentialsJSON(serviceAccountJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase Auth client: %w", err)
	}

	firebaseService = &FirebaseService{client: client}
	return nil
}

// GetFirebaseService returns the initialized Firebase service
func GetFirebaseService() *FirebaseService {
	return firebaseService
}

// ValidateToken validates a Firebase ID token and returns user information
func (fs *FirebaseService) ValidateToken(ctx context.Context, idToken string) (*FirebaseUser, error) {
	if fs.client == nil {
		return nil, fmt.Errorf("Firebase client not initialized")
	}

	token, err := fs.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Firebase token: %w", err)
	}

	// Get user record for additional information
	userRecord, err := fs.client.GetUser(ctx, token.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user record: %w", err)
	}

	firebaseUser := &FirebaseUser{
		UID:   token.UID,
		Name:  userRecord.DisplayName,
		Email: userRecord.Email,
	}

	// If display name is empty, try to construct from email or use UID
	if firebaseUser.Name == "" {
		if firebaseUser.Email != "" {
			// Extract name from email (part before @)
			if atIndex := len(firebaseUser.Email); atIndex > 0 {
				for i, char := range firebaseUser.Email {
					if char == '@' {
						atIndex = i
						break
					}
				}
				firebaseUser.Name = firebaseUser.Email[:atIndex]
			}
		} else {
			firebaseUser.Name = "User " + token.UID[:8] // Use first 8 chars of UID
		}
	}

	return firebaseUser, nil
}