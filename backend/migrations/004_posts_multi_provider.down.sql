-- Best-effort rollback for 004 (single-channel posts only if each row has one channel_id).

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'posts'
          AND column_name = 'channel_id'
    ) THEN
        RAISE NOTICE 'posts.channel_id already exists; skipping 004 down';
    ELSE
    ALTER TABLE posts
        ADD COLUMN channel_id BIGINT REFERENCES channels(id) ON DELETE CASCADE,
        ADD COLUMN provider TEXT,
        ADD COLUMN provider_post_id TEXT;

    UPDATE posts
    SET
        channel_id = channel_ids[1],
        provider = providers[1],
        provider_post_id = COALESCE(
            provider_post_ids->>providers[1],
            (SELECT value FROM jsonb_each_text(provider_post_ids) LIMIT 1)
        )
    WHERE cardinality(channel_ids) >= 1
      AND cardinality(providers) >= 1;

    ALTER TABLE posts ALTER COLUMN channel_id SET NOT NULL;
    ALTER TABLE posts ALTER COLUMN provider SET NOT NULL;
    ALTER TABLE posts ALTER COLUMN provider_post_id SET NOT NULL;

    DROP INDEX IF EXISTS idx_posts_channel_ids;
    DROP INDEX IF EXISTS idx_posts_providers;

    ALTER TABLE posts DROP COLUMN IF EXISTS channel_ids;
    ALTER TABLE posts DROP COLUMN IF EXISTS owner_id;
    ALTER TABLE posts DROP COLUMN IF EXISTS providers;
    ALTER TABLE posts DROP COLUMN IF EXISTS provider_post_ids;

    CREATE UNIQUE INDEX IF NOT EXISTS posts_channel_id_provider_post_id_key
        ON posts (channel_id, provider_post_id);
    CREATE INDEX IF NOT EXISTS idx_posts_channel ON posts (channel_id);
    CREATE INDEX IF NOT EXISTS idx_posts_provider_post ON posts (provider, provider_post_id);
    END IF;
END $$;
