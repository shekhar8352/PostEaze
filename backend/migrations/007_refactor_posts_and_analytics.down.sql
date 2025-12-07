-- +goose Down
-- +goose StatementBegin

-- 1. Drop Generic Analytics Snapshots table
DROP TABLE IF EXISTS analytics_snapshots;

-- 2. Revert Analytics Tables Columns
ALTER TABLE instagram_post_analytics 
    DROP COLUMN IF EXISTS views, 
    DROP COLUMN IF EXISTS total_interactions, 
    DROP COLUMN IF EXISTS engagement_rate, 
    DROP COLUMN IF EXISTS plays;

ALTER TABLE instagram_profile_analytics 
    DROP COLUMN IF EXISTS views, 
    DROP COLUMN IF EXISTS accounts_engaged, 
    DROP COLUMN IF EXISTS total_interactions, 
    DROP COLUMN IF EXISTS bio_link_clicks, 
    DROP COLUMN IF EXISTS phone_call_clicks, 
    DROP COLUMN IF EXISTS text_message_clicks, 
    DROP COLUMN IF EXISTS get_directions_clicks;

-- 3. Revert Posts Table Changes
-- Note: Reverting data from JSONB/Arrays back to flat columns is lossy if multiple channels were added.
-- We will attempt to restore the *first* channel/provider if available, but this is a destructive rollback for multi-channel data.

UPDATE posts SET 
    channel_id = (channel_ids->0)::BIGINT,
    provider = providers[1],
    provider_post_id = provider_post_ids->>(providers[1])
WHERE channel_id IS NULL AND jsonb_array_length(channel_ids) > 0;

-- Drop the new columns
ALTER TABLE posts 
    DROP COLUMN IF EXISTS channel_ids, 
    DROP COLUMN IF EXISTS providers, 
    DROP COLUMN IF EXISTS provider_post_ids;

-- Re-enable constraints (this might fail if there are duplicates now, but it's part of the down migration intent)
-- ALTER TABLE posts ADD CONSTRAINT posts_channel_id_provider_post_id_key UNIQUE (channel_id, provider_post_id);

-- +goose StatementEnd
