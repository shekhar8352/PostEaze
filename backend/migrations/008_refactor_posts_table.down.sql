-- +goose Down
-- +goose StatementBegin

-- 1. Restore provider column and renaming
ALTER TABLE posts ADD COLUMN IF NOT EXISTS provider VARCHAR(50);
-- Restore channel_id
ALTER TABLE posts ADD COLUMN IF NOT EXISTS channel_id BIGINT;

UPDATE posts SET provider = 'instagram';

ALTER TABLE posts RENAME COLUMN instagram_post_id TO provider_post_id;

-- 2. Drop new columns
ALTER TABLE posts DROP COLUMN IF EXISTS channel_ids;
ALTER TABLE posts DROP COLUMN IF EXISTS owner_id;

-- 3. Restore constraint might conflict if duplicates exist, skipping strictly for safety unless critical
-- ALTER TABLE posts ADD CONSTRAINT posts_channel_id_provider_post_id_key UNIQUE (channel_id, provider_post_id);

-- +goose StatementEnd
