-- +goose Up
-- +goose StatementBegin

-- 1. Refactor POSTS table for multi-channel support
ALTER TABLE posts ADD COLUMN IF NOT EXISTS channel_ids JSONB DEFAULT '[]'::jsonb;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS providers TEXT[] DEFAULT '{}';
ALTER TABLE posts ADD COLUMN IF NOT EXISTS provider_post_ids JSONB DEFAULT '{}'::jsonb;

-- Migrate existing data: Move current 'provider' into new 'providers' array
UPDATE posts SET 
    channel_ids = jsonb_build_array(channel_id),
    providers = ARRAY[provider],
    provider_post_ids = jsonb_build_object(provider, provider_post_id)
WHERE (channel_ids = '[]'::jsonb OR channel_ids IS NULL) AND provider IS NOT NULL;

-- Make old columns nullable (to support multi-channel posts where single columns don't apply)
ALTER TABLE posts ALTER COLUMN channel_id DROP NOT NULL;
ALTER TABLE posts ALTER COLUMN provider DROP NOT NULL;
ALTER TABLE posts ALTER COLUMN provider_post_id DROP NOT NULL;
-- Drop unique constraint that assumes single channel
ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_channel_id_provider_post_id_key;

-- 2. Update existing analytics tables with new Meta metrics
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS views INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS accounts_engaged INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS total_interactions INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS bio_link_clicks INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS phone_call_clicks INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS text_message_clicks INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS get_directions_clicks INTEGER;

-- Update instagram_post_analytics with new fields
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS views INTEGER;
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS total_interactions INTEGER;
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS engagement_rate NUMERIC(5,4);
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS plays INTEGER;

-- 3. Generic Analytics Snapshots (for Channels AND Posts)
CREATE TABLE analytics_snapshots (
    id BIGSERIAL PRIMARY KEY,
    entity_type TEXT NOT NULL,     -- 'channel', 'post'
    entity_id BIGINT NOT NULL,     -- channel_id or post_id
    period_type TEXT NOT NULL,     -- 'daily', 'weekly', 'monthly'
    
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    
    metrics JSONB NOT NULL DEFAULT '{}'::jsonb, -- Flexible storage for all metrics
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE (entity_type, entity_id, period_type, start_date)
);

CREATE INDEX idx_analytics_snapshots_entity ON analytics_snapshots(entity_type, entity_id);
CREATE INDEX idx_analytics_snapshots_period ON analytics_snapshots(period_type, start_date);

-- +goose StatementEnd
