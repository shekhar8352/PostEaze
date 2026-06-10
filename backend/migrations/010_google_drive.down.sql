ALTER TABLE media_versions
    DROP COLUMN IF EXISTS drive_revision_id,
    DROP COLUMN IF EXISTS drive_file_id,
    DROP COLUMN IF EXISTS storage_provider;

ALTER TABLE media_assets
    DROP COLUMN IF EXISTS drive_file_id;

DROP TABLE IF EXISTS user_integrations;
