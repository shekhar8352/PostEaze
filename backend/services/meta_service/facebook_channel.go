package meta_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/meta"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

// CreateFacebookPageChannel connects a Facebook Page using OAuth code and persists channel + page access token.
func (s *MetaServiceImpl) CreateFacebookPageChannel(
	ctx context.Context,
	code string,
	redirectURI string,
	pageID string,
	channelName string,
	ownerUserID uuid.UUID,
	teamID *uuid.UUID,
) (channelID int64, displayName string, err error) {
	tokenResp, err := s.provider.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		return 0, "", fmt.Errorf("exchange code: %w", err)
	}
	longLived, err := s.provider.GetLongLivedToken(tokenResp.AccessToken)
	if err != nil {
		return 0, "", fmt.Errorf("long-lived token: %w", err)
	}

	pages, err := s.provider.GetPages(longLived.AccessToken)
	if err != nil {
		return 0, "", fmt.Errorf("list pages: %w", err)
	}

	var selected *meta.Page
	for i := range pages {
		if pages[i].ID == pageID {
			selected = &pages[i]
			break
		}
	}
	if selected == nil {
		return 0, "", fmt.Errorf("page %s not found for this account or missing permissions", pageID)
	}
	if selected.AccessToken == "" {
		return 0, "", fmt.Errorf("page access token missing; ensure pages_show_list and Page tasks are granted")
	}

	return s.persistFacebookPageChannel(ctx, selected, channelName, ownerUserID, teamID)
}

// CreateFacebookPageChannelFromPageAccessToken verifies a page access token with the Graph API and persists the channel.
// Used after POST /meta/callback returns pages (OAuth code already consumed).
func (s *MetaServiceImpl) CreateFacebookPageChannelFromPageAccessToken(
	ctx context.Context,
	pageID string,
	pageAccessToken string,
	channelName string,
	ownerUserID uuid.UUID,
	teamID *uuid.UUID,
) (channelID int64, displayName string, err error) {
	selected, err := s.provider.FetchPageByAccessToken(pageID, pageAccessToken)
	if err != nil {
		return 0, "", fmt.Errorf("verify page token: %w", err)
	}
	if selected.AccessToken == "" {
		return 0, "", fmt.Errorf("page access token missing")
	}

	return s.persistFacebookPageChannel(ctx, selected, channelName, ownerUserID, teamID)
}

func (s *MetaServiceImpl) persistFacebookPageChannel(
	ctx context.Context,
	selected *meta.Page,
	channelName string,
	ownerUserID uuid.UUID,
	teamID *uuid.UUID,
) (channelID int64, displayName string, err error) {
	encryptedToken, err := encryption.Encrypt(selected.AccessToken)
	if err != nil {
		return 0, "", fmt.Errorf("encrypt token: %w", err)
	}

	metaJSON, err := json.Marshal(map[string]interface{}{
		"page_id":    selected.ID,
		"name":       selected.Name,
		"category":   selected.Category,
		"page_tasks": selected.Tasks,
	})
	if err != nil {
		return 0, "", fmt.Errorf("metadata: %w", err)
	}

	display := channelName
	if display == "" {
		display = selected.Name
	}
	if display == "" {
		display = "Page " + selected.ID
	}

	db := database.GetDB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback()

	ch := &entities.Channel{
		OwnerUserID:       ownerUserID,
		TeamID:            teamID,
		Provider:          "facebook",
		ProviderChannelID: selected.ID,
		DisplayName:       &display,
		Username:          nil,
		AvatarURL:         nil,
		IsActive:          true,
		Metadata:          metaJSON,
	}

	if err := repositories.CreateChannel(ctx, tx, ch); err != nil {
		return 0, "", fmt.Errorf("create channel: %w", err)
	}

	tokenType := "Bearer"
	// Page tokens from long-lived users are effectively long-lived; no fixed expiry from this response.
	var expiresAt *time.Time
	chTok := &entities.ChannelToken{
		ChannelID:   ch.ID,
		AccessToken: encryptedToken,
		TokenType:   &tokenType,
		ExpiresAt:   expiresAt,
		Revoked:     false,
	}
	if err := repositories.CreateChannelToken(ctx, tx, chTok); err != nil {
		return 0, "", fmt.Errorf("create channel token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, "", err
	}

	return ch.ID, display, nil
}
