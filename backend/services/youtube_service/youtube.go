package youtube_service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/google"
	"github.com/shekhar8352/PostEaze/provider/youtube"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

type YouTubeService struct {
	OAuth    *google.OAuthProvider
	YTClient *youtube.Client
}

func NewYouTubeService() *YouTubeService {
	return &YouTubeService{
		OAuth:    google.NewOAuthProvider(),
		YTClient: youtube.NewClient(),
	}
}

type ChannelResponse struct {
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

func (s *YouTubeService) CreateChannel(ctx context.Context, code string, ownerUserID uuid.UUID, teamID *uuid.UUID) (*ChannelResponse, error) {
	redirectURI := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
	if redirectURI == "" {
		return nil, fmt.Errorf("GOOGLE_OAUTH_REDIRECT_URI not set")
	}

	tokens, err := s.OAuth.ExchangeCode(code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	ch, err := s.YTClient.GetMyChannel(tokens.AccessToken)
	if err != nil {
		return nil, err
	}

	encAccess, err := encryption.Encrypt(tokens.AccessToken)
	if err != nil {
		return nil, err
	}
	var encRefresh []byte
	if tokens.RefreshToken != "" {
		encRefresh, err = encryption.Encrypt(tokens.RefreshToken)
		if err != nil {
			return nil, err
		}
	}

	var expiresAt *time.Time
	if tokens.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
		expiresAt = &t
	}

	displayName := ch.Snippet.Title
	metadata, _ := json.Marshal(map[string]any{
		"youtube_channel_id": ch.ID,
		"custom_url":         ch.Snippet.CustomURL,
	})

	db := database.GetDB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	channel := &entities.Channel{
		OwnerUserID:       ownerUserID,
		TeamID:            teamID,
		Provider:          "youtube",
		ProviderChannelID: ch.ID,
		DisplayName:       &displayName,
		IsActive:          true,
		Metadata:          metadata,
	}
	if err := repositories.CreateChannel(ctx, tx, channel); err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}

	tokenType := "Bearer"
	scopes := tokens.Scope
	channelToken := &entities.ChannelToken{
		ChannelID:    channel.ID,
		AccessToken:  encAccess,
		RefreshToken: encRefresh,
		TokenType:    &tokenType,
		Scopes:       &scopes,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}
	if err := repositories.CreateChannelToken(ctx, tx, channelToken); err != nil {
		return nil, fmt.Errorf("create channel token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &ChannelResponse{
		ChannelID:   channel.ID,
		ChannelName: displayName,
	}, nil
}

func (s *YouTubeService) GetValidAccessToken(ctx context.Context, channelID int64) (string, error) {
	tok, err := repositories.GetLatestTokenByChannelID(ctx, channelID)
	if err != nil {
		return "", err
	}
	if tok == nil {
		return "", fmt.Errorf("no token for channel %d", channelID)
	}

	plain, err := encryption.Decrypt(tok.AccessToken)
	if err != nil {
		return "", err
	}

	if tok.ExpiresAt != nil && time.Now().Add(2*time.Minute).Before(*tok.ExpiresAt) {
		return plain, nil
	}

	if tok.RefreshToken == nil || len(tok.RefreshToken) == 0 {
		return plain, nil
	}

	refreshPlain, err := encryption.Decrypt(tok.RefreshToken)
	if err != nil {
		return "", err
	}

	refreshed, err := s.OAuth.RefreshToken(refreshPlain)
	if err != nil {
		return "", err
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
	_ = repositories.UpdateChannelTokenAccess(ctx, tok.ID, encAccess, expiresAt)
	return refreshed.AccessToken, nil
}
