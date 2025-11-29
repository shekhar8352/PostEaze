package entities

// Channel SQL queries
const (
	CreateChannel = `
		INSERT INTO channels (
			owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, metadata
		) VALUES (
			:owner_user_id, :team_id, :provider, :provider_channel_id,
			:display_name, :username, :avatar_url, :is_active, :metadata
		)
		RETURNING id, owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, error_status,
			metadata, connected_at, created_at, updated_at
	`

	GetChannelByID = `
		SELECT id, owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, error_status,
			metadata, connected_at, created_at, updated_at
		FROM channels
		WHERE id = :id
	`

	GetChannelsByOwner = `
		SELECT id, owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, error_status,
			metadata, connected_at, created_at, updated_at
		FROM channels
		WHERE owner_user_id = :owner_user_id
		ORDER BY created_at DESC
	`

	GetChannelsByTeam = `
		SELECT id, owner_user_id, team_id, provider, provider_channel_id,
			display_name, username, avatar_url, is_active, error_status,
			metadata, connected_at, created_at, updated_at
		FROM channels
		WHERE team_id = :team_id
		ORDER BY created_at DESC
	`

	UpdateChannelStatus = `
		UPDATE channels
		SET is_active = :is_active, error_status = :error_status, updated_at = NOW()
		WHERE id = :id
	`
)

// ChannelToken SQL queries
const (
	CreateChannelToken = `
		INSERT INTO channel_tokens (
			channel_id, access_token, refresh_token, token_type,
			scopes, expires_at, revoked
		) VALUES (
			:channel_id, :access_token, :refresh_token, :token_type,
			:scopes, :expires_at, :revoked
		)
		RETURNING id, channel_id, access_token, refresh_token, token_type,
			scopes, expires_at, revoked, issued_at, last_used_at, created_at
	`

	GetLatestTokenByChannelID = `
		SELECT id, channel_id, access_token, refresh_token, token_type,
			scopes, expires_at, revoked, issued_at, last_used_at, created_at
		FROM channel_tokens
		WHERE channel_id = :channel_id AND revoked = false
		ORDER BY created_at DESC
		LIMIT 1
	`

	RevokeTokensByChannelID = `
		UPDATE channel_tokens
		SET revoked = true
		WHERE channel_id = :channel_id AND revoked = false
	`

	UpdateTokenLastUsed = `
		UPDATE channel_tokens
		SET last_used_at = NOW()
		WHERE id = :id
	`
)
