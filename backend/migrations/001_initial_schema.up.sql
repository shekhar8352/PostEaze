-- Combined initial schema for PostEaze
-- This migration replaces the earlier 001–008 sequence for a fresh dev database.
-- It defines all core tables as they should exist now.

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ============================================================================
-- USERS & AUTH
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    platforms TEXT[] DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_platforms ON users USING GIN(platforms);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_refresh_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================================
-- TEAMS & MEMBERS
-- ============================================================================

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    owner_id UUID NOT NULL,
    description TEXT,
    avatar_url TEXT,
    visibility TEXT CHECK (visibility IN ('private', 'public')) DEFAULT 'private',
    status TEXT CHECK (status IN ('active', 'archived', 'deleted')) DEFAULT 'active',
    settings JSONB DEFAULT '{}'::jsonb,
    owner_role_override TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_teams_owner_id ON teams(owner_id);
CREATE INDEX IF NOT EXISTS idx_teams_name_trgm ON teams USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_active_teams ON teams(status) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS team_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL,
    status TEXT CHECK (status IN ('active','invited','pending','removed')) DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    invited_by UUID REFERENCES users(id),
    permissions JSONB DEFAULT '{}'::jsonb,
    is_primary BOOLEAN DEFAULT false,
    last_active_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (team_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_team_members_team_id ON team_members(team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id);
CREATE INDEX IF NOT EXISTS idx_team_members_team_user ON team_members(team_id, user_id);
CREATE INDEX IF NOT EXISTS idx_active_team_members ON team_members(team_id) WHERE status = 'active';

-- ============================================================================
-- CHANNELS & TOKENS
-- ============================================================================

CREATE TABLE IF NOT EXISTS channels (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id UUID NOT NULL,
    team_id UUID,
    provider TEXT NOT NULL,
    provider_channel_id TEXT NOT NULL,
    display_name TEXT,
    username TEXT,
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    error_status TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    connected_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (provider, provider_channel_id)
);

CREATE INDEX IF NOT EXISTS idx_channels_owner_user_id ON channels (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_channels_team_id ON channels (team_id);
CREATE INDEX IF NOT EXISTS idx_channels_provider ON channels (provider);

CREATE TABLE IF NOT EXISTS channel_tokens (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
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

CREATE INDEX IF NOT EXISTS idx_channel_tokens_channel_id ON channel_tokens (channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_tokens_expires_at ON channel_tokens (expires_at);

-- ============================================================================
-- POSTS & ANALYTICS
-- ============================================================================

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_post_id TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'native',
    post_type TEXT,
    caption TEXT,
    media JSONB,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (channel_id, provider_post_id)
);

CREATE INDEX IF NOT EXISTS idx_posts_channel ON posts(channel_id);
CREATE INDEX IF NOT EXISTS idx_posts_provider_post ON posts(provider, provider_post_id);

CREATE TABLE IF NOT EXISTS instagram_post_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    impressions INTEGER,
    reach INTEGER,
    likes INTEGER,
    comments INTEGER,
    saves INTEGER,
    shares INTEGER,
    video_views INTEGER,
    profile_visits INTEGER,
    follows INTEGER,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (post_id, date)
);

CREATE INDEX IF NOT EXISTS idx_ig_post_analytics_channel_date ON instagram_post_analytics(channel_id, date);

CREATE TABLE IF NOT EXISTS instagram_story_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    impressions INTEGER,
    reach INTEGER,
    exits INTEGER,
    forwards INTEGER,
    replies INTEGER,
    taps_forward INTEGER,
    taps_backward INTEGER,
    taps_exit INTEGER,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (post_id, date)
);

CREATE INDEX IF NOT EXISTS idx_ig_story_analytics_channel_date ON instagram_story_analytics(channel_id, date);

CREATE TABLE IF NOT EXISTS instagram_profile_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    follower_count INTEGER,
    impressions INTEGER,
    profile_views INTEGER,
    reach INTEGER,
    website_clicks INTEGER,
    email_clicks INTEGER,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (channel_id, date)
);

CREATE INDEX IF NOT EXISTS idx_ig_profile_analytics_channel_date ON instagram_profile_analytics(channel_id, date);

CREATE TABLE IF NOT EXISTS facebook_page_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    followers INTEGER,
    likes INTEGER,
    reach INTEGER,
    impressions INTEGER,
    page_engaged_users INTEGER,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (channel_id, date)
);

CREATE INDEX IF NOT EXISTS idx_fb_page_analytics_channel_date ON facebook_page_analytics(channel_id, date);

CREATE TABLE IF NOT EXISTS facebook_post_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    impressions INTEGER,
    reach INTEGER,
    likes INTEGER,
    comments INTEGER,
    shares INTEGER,
    photo_views INTEGER,
    video_views INTEGER,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (post_id, date)
);

CREATE INDEX IF NOT EXISTS idx_fb_post_analytics_channel_date ON facebook_post_analytics(channel_id, date);

CREATE TABLE IF NOT EXISTS youtube_video_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    views INTEGER,
    likes INTEGER,
    comments INTEGER,
    shares INTEGER,
    impressions INTEGER,
    impressions_click_through_rate NUMERIC,
    average_view_duration_seconds INTEGER,
    average_view_percentage NUMERIC,
    watch_time_seconds BIGINT,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (post_id, date)
);

CREATE INDEX IF NOT EXISTS idx_yt_video_analytics_channel_date ON youtube_video_analytics(channel_id, date);

CREATE TABLE IF NOT EXISTS youtube_channel_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    subscribers INTEGER,
    views INTEGER,
    watch_time_seconds BIGINT,
    impressions INTEGER,
    estimated_revenue NUMERIC,
    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (channel_id, date)
);

CREATE INDEX IF NOT EXISTS idx_yt_channel_analytics_channel_date ON youtube_channel_analytics(channel_id, date);

