ALTER TABLE instagram_post_analytics DROP COLUMN IF EXISTS plays;
ALTER TABLE instagram_post_analytics DROP COLUMN IF EXISTS engagement_rate;
ALTER TABLE instagram_post_analytics DROP COLUMN IF EXISTS total_interactions;
ALTER TABLE instagram_post_analytics DROP COLUMN IF EXISTS views;

ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS get_directions_clicks;
ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS text_message_clicks;
ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS phone_call_clicks;
ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS bio_link_clicks;
ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS total_interactions;
ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS accounts_engaged;
ALTER TABLE instagram_profile_analytics DROP COLUMN IF EXISTS views;
