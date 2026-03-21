-- Align schema with application upserts (safe if already applied)
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS views INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS accounts_engaged INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS total_interactions INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS bio_link_clicks INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS phone_call_clicks INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS text_message_clicks INTEGER;
ALTER TABLE instagram_profile_analytics ADD COLUMN IF NOT EXISTS get_directions_clicks INTEGER;

ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS views INTEGER;
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS total_interactions INTEGER;
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS engagement_rate NUMERIC(10,6);
ALTER TABLE instagram_post_analytics ADD COLUMN IF NOT EXISTS plays INTEGER;
