-- Align posts with application: channel_ids[], owner_id, providers[], provider_post_ids JSONB.
-- Skips if legacy column channel_id is already absent (fresh DBs using updated 001).

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'posts'
          AND column_name = 'channel_id'
    ) THEN
        RAISE NOTICE 'posts.channel_id not found; assuming posts already uses multi-provider schema';
    ELSE
    ALTER TABLE posts ADD COLUMN IF NOT EXISTS channel_ids BIGINT[] DEFAULT '{}';
    ALTER TABLE posts ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id) ON DELETE SET NULL;
    ALTER TABLE posts ADD COLUMN IF NOT EXISTS providers TEXT[] DEFAULT ARRAY[]::TEXT[];
    ALTER TABLE posts ADD COLUMN IF NOT EXISTS provider_post_ids JSONB DEFAULT '{}'::jsonb;

    UPDATE posts p
    SET
        channel_ids = ARRAY[p.channel_id],
        owner_id = c.owner_user_id,
        providers = ARRAY[p.provider],
        provider_post_ids = jsonb_build_object(p.provider, p.provider_post_id)
    FROM channels c
    WHERE c.id = p.channel_id;

    ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_channel_id_provider_post_id_key;

    DROP INDEX IF EXISTS idx_posts_channel;
    DROP INDEX IF EXISTS idx_posts_provider_post;

    ALTER TABLE posts DROP COLUMN channel_id;
    ALTER TABLE posts DROP COLUMN provider;
    ALTER TABLE posts DROP COLUMN provider_post_id;

    ALTER TABLE posts ALTER COLUMN channel_ids SET NOT NULL;
    ALTER TABLE posts ALTER COLUMN channel_ids SET DEFAULT '{}';
    ALTER TABLE posts ALTER COLUMN providers SET NOT NULL;
    ALTER TABLE posts ALTER COLUMN providers SET DEFAULT ARRAY[]::TEXT[];
    ALTER TABLE posts ALTER COLUMN provider_post_ids SET NOT NULL;
    ALTER TABLE posts ALTER COLUMN provider_post_ids SET DEFAULT '{}'::jsonb;

    CREATE INDEX IF NOT EXISTS idx_posts_channel_ids ON posts USING GIN (channel_ids);
    CREATE INDEX IF NOT EXISTS idx_posts_providers ON posts USING GIN (providers);
    END IF;
END $$;
