-- +goose Up
-- +goose StatementBegin

/*******************************************************************************************
 * POSTS TABLE (GENERAL)
 * Stores posts from ALL providers:
 *  - Insta native posts
 *  - FB native posts
 *  - YT native videos
 *  - Posts published via Posteaze
 *******************************************************************************************/
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,

    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,                    -- instagram, facebook, youtube, ...
    provider_post_id TEXT NOT NULL,            -- ig_media_id, fb_post_id, youtube_video_id

    source TEXT NOT NULL DEFAULT 'native',     -- 'native' | 'posteaze'
    post_type TEXT,                            -- image, video, reel, story, youtube_video, fb_post...

    caption TEXT,
    media JSONB,                               -- array of media (urls, types)
    
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (channel_id, provider_post_id)
);

CREATE INDEX idx_posts_channel ON posts(channel_id);
CREATE INDEX idx_posts_provider_post ON posts(provider, provider_post_id);

/*******************************************************************************************
 * INSTAGRAM ANALYTICS TABLES
 *******************************************************************************************/

-- 1️⃣ Instagram Post Analytics
CREATE TABLE instagram_post_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    impressions INTEGER,
    reach INTEGER,
    likes INTEGER,
    comments INTEGER,
    saves INTEGER,
    shares INTEGER,
    video_views INTEGER,
    profile_visits INTEGER,
    follows INTEGER,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (post_id, date)
);

CREATE INDEX idx_ig_post_analytics_channel_date ON instagram_post_analytics(channel_id, date);


-- 2️⃣ Instagram Story Analytics
CREATE TABLE instagram_story_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    impressions INTEGER,
    reach INTEGER,
    exits INTEGER,
    forwards INTEGER,
    replies INTEGER,
    taps_forward INTEGER,
    taps_backward INTEGER,
    taps_exit INTEGER,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (post_id, date)
);

CREATE INDEX idx_ig_story_analytics_channel_date ON instagram_story_analytics(channel_id, date);


-- 3️⃣ Instagram Profile Analytics (daily)
CREATE TABLE instagram_profile_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    follower_count INTEGER,
    impressions INTEGER,
    profile_views INTEGER,
    reach INTEGER,
    website_clicks INTEGER,
    email_clicks INTEGER,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (channel_id, date)
);

CREATE INDEX idx_ig_profile_analytics_channel_date ON instagram_profile_analytics(channel_id, date);


/*******************************************************************************************
 * FACEBOOK ANALYTICS TABLES
 *******************************************************************************************/

-- 1️⃣ Facebook Page Analytics
CREATE TABLE facebook_page_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    followers INTEGER,
    likes INTEGER,
    reach INTEGER,
    impressions INTEGER,
    page_engaged_users INTEGER,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (channel_id, date)
);

CREATE INDEX idx_fb_page_analytics_channel_date ON facebook_page_analytics(channel_id, date);


-- 2️⃣ Facebook Post Analytics
CREATE TABLE facebook_post_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    impressions INTEGER,
    reach INTEGER,
    likes INTEGER,
    comments INTEGER,
    shares INTEGER,
    photo_views INTEGER,
    video_views INTEGER,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (post_id, date)
);

CREATE INDEX idx_fb_post_analytics_channel_date ON facebook_post_analytics(channel_id, date);


/*******************************************************************************************
 * YOUTUBE ANALYTICS TABLES
 *******************************************************************************************/

-- 1️⃣ YouTube Video Analytics
CREATE TABLE youtube_video_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    views INTEGER,
    likes INTEGER,
    comments INTEGER,
    shares INTEGER,
    impressions INTEGER,
    impressions_click_through_rate NUMERIC,
    average_view_duration_seconds INTEGER,
    average_view_percentage NUMERIC,
    watch_time_seconds BIGINT,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (post_id, date)
);

CREATE INDEX idx_yt_video_analytics_channel_date ON youtube_video_analytics(channel_id, date);


-- 2️⃣ YouTube Channel Analytics
CREATE TABLE youtube_channel_analytics (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,

    date DATE NOT NULL,

    subscribers INTEGER,
    views INTEGER,
    watch_time_seconds BIGINT,
    impressions INTEGER,
    estimated_revenue NUMERIC,

    raw JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (channel_id, date)
);

CREATE INDEX idx_yt_channel_analytics_channel_date ON youtube_channel_analytics(channel_id, date);


-- +goose StatementEnd
