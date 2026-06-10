-- Google Drive integration + storage provider on media versions

CREATE TABLE IF NOT EXISTS user_integrations (
    id                      BIGSERIAL    PRIMARY KEY,
    user_id                 UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider                TEXT         NOT NULL CHECK (provider IN ('google_drive')),
    provider_account_email  TEXT,
    access_token            BYTEA        NOT NULL,
    refresh_token           BYTEA,
    scopes                  TEXT,
    expires_at              TIMESTAMPTZ,
    revoked                 BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, provider)
);

CREATE INDEX IF NOT EXISTS idx_user_integrations_user ON user_integrations (user_id);

ALTER TABLE media_assets
    ADD COLUMN IF NOT EXISTS drive_file_id TEXT;

CREATE INDEX IF NOT EXISTS idx_media_assets_drive_file ON media_assets (drive_file_id)
    WHERE drive_file_id IS NOT NULL;

ALTER TABLE media_versions
    ADD COLUMN IF NOT EXISTS storage_provider TEXT NOT NULL DEFAULT 'blob'
        CHECK (storage_provider IN ('blob', 'google_drive')),
    ADD COLUMN IF NOT EXISTS drive_file_id TEXT,
    ADD COLUMN IF NOT EXISTS drive_revision_id TEXT;

CREATE INDEX IF NOT EXISTS idx_media_versions_storage ON media_versions (storage_provider);
