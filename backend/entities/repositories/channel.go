package repositories

import (
	"context"
	"database/sql"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

func CreateChannel(ctx context.Context, tx *sql.Tx, channel *entities.Channel) error {
	query := `
		INSERT INTO channels (
			owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, connected_at, created_at, updated_at
	`
	return tx.QueryRowContext(ctx, query,
		channel.OwnerUserID,
		channel.TeamID,
		channel.Provider,
		channel.ProviderChannelID,
		channel.DisplayName,
		channel.Username,
		channel.AvatarURL,
		channel.IsActive,
		channel.Metadata,
	).Scan(&channel.ID, &channel.ConnectedAt, &channel.CreatedAt, &channel.UpdatedAt)
}

func GetChannelByID(ctx context.Context, channelID int64) (*entities.Channel, error) {
	query := `
		SELECT id, owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, error_status,
			metadata, connected_at, created_at, updated_at
		FROM channels
		WHERE id = $1
	`
	channel := &entities.Channel{}
	err := database.GetDB().QueryRowContext(ctx, query, channelID).Scan(
		&channel.ID,
		&channel.OwnerUserID,
		&channel.TeamID,
		&channel.Provider,
		&channel.ProviderChannelID,
		&channel.DisplayName,
		&channel.Username,
		&channel.AvatarURL,
		&channel.IsActive,
		&channel.ErrorStatus,
		&channel.Metadata,
		&channel.ConnectedAt,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func CreateChannelToken(ctx context.Context, tx *sql.Tx, token *entities.ChannelToken) error {
	query := `
		INSERT INTO channel_tokens (
			channel_id, access_token, refresh_token, token_type,
			scopes, expires_at, revoked
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, issued_at, created_at
	`
	return tx.QueryRowContext(ctx, query,
		token.ChannelID,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.Scopes,
		token.ExpiresAt,
		token.Revoked,
	).Scan(&token.ID, &token.IssuedAt, &token.CreatedAt)
}

func GetLatestTokenByChannelID(ctx context.Context, channelID int64) (*entities.ChannelToken, error) {
	query := `
		SELECT id, channel_id, access_token, refresh_token, token_type,
			scopes, expires_at, revoked, issued_at, last_used_at, created_at
		FROM channel_tokens
		WHERE channel_id = $1 AND revoked = false
		ORDER BY created_at DESC
		LIMIT 1
	`
	token := &entities.ChannelToken{}
	err := database.GetDB().QueryRowContext(ctx, query, channelID).Scan(
		&token.ID,
		&token.ChannelID,
		&token.AccessToken,
		&token.RefreshToken,
		&token.TokenType,
		&token.Scopes,
		&token.ExpiresAt,
		&token.Revoked,
		&token.IssuedAt,
		&token.LastUsedAt,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GetChannelsByUserID retrieves all channels for a user, optionally filtered by provider
func GetChannelsByUserID(ctx context.Context, userID string, provider string) ([]entities.Channel, error) {
	db := database.GetDB()
	var rows *sql.Rows
	var err error

	if provider != "" {
		// Filter by provider
		rows, err = db.QueryContext(ctx, entities.GetChannelsByUserIDAndProviderQuery(), userID, provider)
	} else {
		// Get all channels for user
		rows, err = db.QueryContext(ctx, entities.GetChannelsByUserIDQuery(), userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []entities.Channel
	for rows.Next() {
		var channel entities.Channel
		err := rows.Scan(
			&channel.ID,
			&channel.OwnerUserID,
			&channel.TeamID,
			&channel.Provider,
			&channel.ProviderChannelID,
			&channel.DisplayName,
			&channel.Username,
			&channel.AvatarURL,
			&channel.IsActive,
			&channel.ErrorStatus,
			&channel.Metadata,
			&channel.ConnectedAt,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}
