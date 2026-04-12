ALTER TABLE media_assets DROP CONSTRAINT IF EXISTS fk_media_assets_current_version;
DROP TABLE IF EXISTS media_versions;
DROP TABLE IF EXISTS media_assets;
