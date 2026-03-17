-- Down migration for 001_initial_schema.up.sql
-- Drops all objects created in the corresponding up migration.
-- NOTE: This is intended for dev environments. It will DROP ALL DATA.

-- Drop analytics tables (depend on posts/channels)
DROP TABLE IF EXISTS youtube_channel_analytics CASCADE;
DROP TABLE IF EXISTS youtube_video_analytics CASCADE;
DROP TABLE IF EXISTS facebook_post_analytics CASCADE;
DROP TABLE IF EXISTS facebook_page_analytics CASCADE;
DROP TABLE IF EXISTS instagram_profile_analytics CASCADE;
DROP TABLE IF EXISTS instagram_story_analytics CASCADE;
DROP TABLE IF EXISTS instagram_post_analytics CASCADE;

-- Drop posts
DROP TABLE IF EXISTS posts CASCADE;

-- Drop channel tokens and channels
DROP TABLE IF EXISTS channel_tokens CASCADE;
DROP TABLE IF EXISTS channels CASCADE;

-- Drop team membership and teams
DROP TABLE IF EXISTS team_members CASCADE;
DROP TABLE IF EXISTS teams CASCADE;

-- Drop auth tables
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Optionally drop extensions (safe if other DB objects don't rely on them)
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS "pgcrypto";

