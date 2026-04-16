-- Provider-agnostic stored comments (Instagram first; Facebook-ready).

CREATE TABLE IF NOT EXISTS social_comments (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    post_id BIGINT NULL REFERENCES posts (id) ON DELETE SET NULL,
    provider_content_id TEXT NOT NULL,
    provider_comment_id TEXT NOT NULL,
    parent_provider_comment_id TEXT NULL,
    text TEXT NULL,
    author_provider_user_id TEXT NULL,
    author_username TEXT NULL,
    raw_payload JSONB NULL,
    comment_created_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_social_comments_channel_provider_comment UNIQUE (channel_id, provider, provider_comment_id),
    CONSTRAINT chk_social_comments_provider CHECK (provider IN ('instagram', 'facebook'))
);

CREATE INDEX IF NOT EXISTS idx_social_comments_channel_post
    ON social_comments (channel_id, post_id);

CREATE INDEX IF NOT EXISTS idx_social_comments_channel_content
    ON social_comments (channel_id, provider, provider_content_id);

CREATE INDEX IF NOT EXISTS idx_social_comments_channel_created
    ON social_comments (channel_id, created_at DESC);
