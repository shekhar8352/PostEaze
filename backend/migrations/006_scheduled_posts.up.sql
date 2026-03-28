-- Outbound scheduled content (distinct from synced `posts` rows)

CREATE TABLE IF NOT EXISTS scheduled_posts (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_ids BIGINT[] NOT NULL DEFAULT '{}',
    platforms TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    scheduled_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'submitting', 'scheduled', 'failed', 'cancelled', 'published')),
    post_type TEXT NOT NULL
        CHECK (post_type IN ('image', 'video', 'carousel', 'reel', 'story')),
    caption TEXT,
    media JSONB NOT NULL DEFAULT '{}'::jsonb,
    provider_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scheduled_posts_owner ON scheduled_posts (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_scheduled_at ON scheduled_posts (scheduled_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_channel_ids ON scheduled_posts USING GIN (channel_ids);
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_status ON scheduled_posts (status);
