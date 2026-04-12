-- Media workspace: assets with step-by-step versioning

CREATE TABLE IF NOT EXISTS media_assets (
    id                 BIGSERIAL    PRIMARY KEY,
    owner_user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id            BIGINT       REFERENCES teams(id) ON DELETE SET NULL,
    title              TEXT         NOT NULL DEFAULT '',
    asset_type         TEXT         NOT NULL CHECK (asset_type IN ('photo', 'video')),
    status             TEXT         NOT NULL DEFAULT 'draft'
                                    CHECK (status IN ('draft', 'ready', 'published')),
    current_version_id BIGINT,      -- FK added after media_versions exists
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_media_assets_owner     ON media_assets (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_media_assets_status    ON media_assets (status);
CREATE INDEX IF NOT EXISTS idx_media_assets_team      ON media_assets (team_id);

CREATE TABLE IF NOT EXISTS media_versions (
    id              BIGSERIAL    PRIMARY KEY,
    media_asset_id  BIGINT       NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,
    version_number  INT          NOT NULL,
    label           TEXT         NOT NULL DEFAULT 'raw',
    blob_url        TEXT         NOT NULL,
    blob_path_key   TEXT         NOT NULL,
    file_name       TEXT         NOT NULL DEFAULT '',
    content_type    TEXT         NOT NULL DEFAULT '',
    file_size       BIGINT       NOT NULL DEFAULT 0,
    metadata        JSONB        NOT NULL DEFAULT '{}'::jsonb,
    notes           TEXT         NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_media_versions_asset ON media_versions (media_asset_id, version_number);

-- Now add the FK from media_assets -> media_versions
ALTER TABLE media_assets
    ADD CONSTRAINT fk_media_assets_current_version
    FOREIGN KEY (current_version_id) REFERENCES media_versions(id) ON DELETE SET NULL;
