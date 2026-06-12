package google_drive_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/google"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

type TokenService struct {
	OAuth *google.OAuthProvider
}

func NewTokenService() *TokenService {
	return &TokenService{OAuth: google.NewOAuthProvider()}
}

func (s *TokenService) GetValidAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	ui, err := repositories.GetUserIntegration(ctx, userID, entities.UserIntegrationProviderGoogleDrive)
	if err != nil {
		return "", err
	}
	if ui == nil {
		return "", fmt.Errorf("google drive is not connected")
	}

	plain, err := encryption.Decrypt(ui.AccessToken)
	if err != nil {
		return "", fmt.Errorf("decrypt access token: %w", err)
	}

	if ui.ExpiresAt != nil && time.Now().Add(2*time.Minute).Before(*ui.ExpiresAt) {
		return plain, nil
	}

	if ui.RefreshToken == nil || len(ui.RefreshToken) == 0 {
		if ui.ExpiresAt != nil && time.Now().After(*ui.ExpiresAt) {
			return "", fmt.Errorf("google drive token expired; reconnect integration")
		}
		return plain, nil
	}

	refreshPlain, err := encryption.Decrypt(ui.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("decrypt refresh token: %w", err)
	}

	refreshed, err := s.OAuth.RefreshToken(refreshPlain)
	if err != nil {
		return "", fmt.Errorf("refresh token: %w", err)
	}

	encAccess, err := encryption.Encrypt(refreshed.AccessToken)
	if err != nil {
		return "", err
	}
	var expiresAt *time.Time
	if refreshed.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second)
		expiresAt = &t
	}
	if err := repositories.UpdateUserIntegrationTokens(ctx, ui.ID, encAccess, expiresAt); err != nil {
		return "", err
	}
	return refreshed.AccessToken, nil
}
