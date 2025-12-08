-- +goose Up
-- +goose StatementBegin

-- Refactor generic posts table to support multi-channel and explicit owner
ALTER TABLE posts ADD COLUMN IF NOT EXISTS channel_ids BIGINT[] DEFAULT '{}';
ALTER TABLE posts ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id);

-- Clean up generic provider columns if they were added by previous partial migrations
ALTER TABLE posts DROP COLUMN IF EXISTS provider_post_ids;

-- Rename provider_post_id to instagram_post_id
ALTER TABLE posts RENAME COLUMN provider_post_id TO instagram_post_id;

-- Populate new channel_ids from existing channel_id
UPDATE posts SET channel_ids = ARRAY[channel_id] WHERE channel_ids = '{}' AND channel_id IS NOT NULL;

-- Populate owner_id from the channel's owner
UPDATE posts p
SET owner_id = c.owner_user_id
FROM channels c
WHERE p.channel_id = c.id AND p.owner_id IS NULL;

-- Drop deprecated columns
ALTER TABLE posts DROP COLUMN IF EXISTS provider;
ALTER TABLE posts DROP COLUMN IF EXISTS channel_id;

-- +goose StatementEnd
