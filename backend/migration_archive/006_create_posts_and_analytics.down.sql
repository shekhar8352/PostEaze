-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS youtube_channel_analytics;
DROP TABLE IF EXISTS youtube_video_analytics;

DROP TABLE IF EXISTS facebook_post_analytics;
DROP TABLE IF EXISTS facebook_page_analytics;

DROP TABLE IF EXISTS instagram_profile_analytics;
DROP TABLE IF EXISTS instagram_story_analytics;
DROP TABLE IF EXISTS instagram_post_analytics;

DROP TABLE IF EXISTS posts;

-- +goose StatementEnd
