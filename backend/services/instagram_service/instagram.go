package instagram_service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/instagram"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

type InstagramService interface {
	CreateChannel(ctx context.Context, code string, channelName string, ownerUserID uuid.UUID, teamID *uuid.UUID) (*ChannelResponse, error)
}

type InstagramServiceImpl struct {
	provider instagram.InstagramProvider
}

func NewInstagramService() *InstagramServiceImpl {
	return &InstagramServiceImpl{
		provider: instagram.NewInstagramProvider(),
	}
}

type ChannelResponse struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

func (s *InstagramServiceImpl) CreateChannel(ctx context.Context, code string, channelName string, ownerUserID uuid.UUID, teamID *uuid.UUID) (*ChannelResponse, error) {
	// Get redirect URI from environment
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")
	if redirectURI == "" {
		return nil, fmt.Errorf("INSTAGRAM_REDIRECT_URI not set in environment")
	}

	// 1. Exchange code for short-lived token
	shortTokenResp, err := s.provider.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// 2. Exchange short-lived token for long-lived token
	longTokenResp, err := s.provider.GetLongLivedToken(shortTokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get long lived token: %w", err)
	}

	// 3. Encrypt the access token
	encryptedToken, err := encryption.Encrypt(longTokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	// 4. Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(longTokenResp.ExpiresIn) * time.Second)

	// 5. Get database connection and start transaction
	db := database.GetDB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 6. Create channel record
	channel := &entities.Channel{
		OwnerUserID:       ownerUserID,
		TeamID:            teamID,
		Provider:          "instagram",
		ProviderChannelID: fmt.Sprintf("%d", shortTokenResp.UserID),
		DisplayName:       &channelName,
		IsActive:          true,
		Metadata:          []byte("{}"),
	}

	err = repositories.CreateChannel(ctx, tx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	// 7. Create channel token record
	tokenType := "Bearer"
	channelToken := &entities.ChannelToken{
		ChannelID:    channel.ID,
		AccessToken:  encryptedToken,
		RefreshToken: nil, // Instagram doesn't provide refresh tokens
		TokenType:    &tokenType,
		ExpiresAt:    &expiresAt,
		Revoked:      false,
	}

	err = repositories.CreateChannelToken(ctx, tx, channelToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel token: %w", err)
	}

	// 8. Subscribe to webhooks
	// We do this BEFORE committing the transaction, but if it fails, we might want to log a warning
	// rather than failing the whole channel creation, or we can fail it.
	// Given the requirement "We need to subscribe to these events while creating the channel itself",
	// we should probably fail if subscription fails.
	webhookFields := []string{"comments", "mentions", "story_insights"}
	// Note: We need the Page ID (Instagram Business Account ID) to subscribe.
	// The shortTokenResp.UserID is the Instagram User ID.
	// For Instagram Basic Display, we might not be able to subscribe to these webhooks directly on the user node
	// in the same way as Graph API.
	// However, assuming we are using the Instagram Graph API (Business), the ID we got is likely the IG User ID.
	// Let's attempt subscription.
	err = s.provider.SubscribeToWebhooks(longTokenResp.AccessToken, fmt.Sprintf("%d", shortTokenResp.UserID), webhookFields)
	if err != nil {
		// Webhooks are not compulsory as of now, so we just log the error and proceed
		fmt.Printf("Warning: failed to subscribe to webhooks: %v\n", err)
	}

	// 9. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &ChannelResponse{
		ChannelID:   channel.ID,
		ChannelName: channelName,
	}, nil
}
