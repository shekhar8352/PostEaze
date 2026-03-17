-- +goose Up
-- +goose StatementBegin

/******************************************************************************************
 * TABLE: channels
 *
 * Represents a connected social media account (Facebook Page, Instagram Account,
 * YouTube Channel, LinkedIn Page, X/Twitter account, Threads, etc.).
 *
 * Key points:
 * - A channel ALWAYS belongs to a user (owner_user_id).
 * - A channel MAY optionally be assigned to a team (team_id).
 * - A channel belongs to exactly one team OR no team (team_id NULL).
 * - provider + provider_channel_id uniquely identifies the external account.
 * - metadata JSONB supports future extensibility.
 ******************************************************************************************/

CREATE TABLE channels (
    id BIGSERIAL PRIMARY KEY,

    -- Owner of the channel: the user who connected it
    owner_user_id UUID NOT NULL,

    -- Channel may optionally belong to ONE team. NULL means unassigned.
    team_id UUID,

    -- Provider/platform info (facebook, instagram, youtube, linkedin, x, threads, etc.)
    provider TEXT NOT NULL,

    -- Unique identifier returned by provider (e.g., page_id, ig_business_id, yt_channel_id)
    provider_channel_id TEXT NOT NULL,

    -- Human-readable info
    display_name TEXT,
    username TEXT,
    avatar_url TEXT,

    -- Status fields
    is_active BOOLEAN DEFAULT TRUE,
    error_status TEXT, -- used to surface connection/token errors in UI

    -- For provider-specific extra info (category, permissions, business flags...)
    metadata JSONB DEFAULT '{}'::jsonb,

    -- Timestamps
    connected_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- Prevent duplicate channels for same provider + channel ID
    UNIQUE (provider, provider_channel_id)
);

-- Helpful indexes for queries
CREATE INDEX idx_channels_owner_user_id ON channels (owner_user_id);
CREATE INDEX idx_channels_team_id ON channels (team_id);
CREATE INDEX idx_channels_provider ON channels (provider);



/******************************************************************************************
 * TABLE: channel_tokens
 *
 * Secure storage for provider access tokens and refresh tokens.
 * 
 * - Stored as encrypted BYTEA values.
 * - A channel may have multiple tokens over time (token rotation).
 * - Always fetch the latest non-revoked token when publishing.
 * - expires_at null = token does not expire (rare but possible on Meta page tokens).
 ******************************************************************************************/

CREATE TABLE channel_tokens (
    id BIGSERIAL PRIMARY KEY,

    -- FK to channels table; cascade delete ensures tokens are removed on channel deletion
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,

    -- Encrypted tokens
    access_token BYTEA NOT NULL,
    refresh_token BYTEA,

    token_type TEXT,
    scopes TEXT,
    expires_at TIMESTAMPTZ,

    revoked BOOLEAN DEFAULT FALSE,

    issued_at TIMESTAMPTZ DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for efficient lookups & token expiration management
CREATE INDEX idx_channel_tokens_channel_id ON channel_tokens (channel_id);
CREATE INDEX idx_channel_tokens_expires_at ON channel_tokens (expires_at);

-- +goose StatementEnd
