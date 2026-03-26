-- Lifetime audience breakdowns from Instagram Graph (audience_city, audience_country, etc.), one row per channel per UTC day.
CREATE TABLE IF NOT EXISTS instagram_audience_snapshots (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    raw JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (channel_id, snapshot_date)
);

CREATE INDEX IF NOT EXISTS idx_ig_audience_snapshots_channel_date ON instagram_audience_snapshots(channel_id, snapshot_date DESC);
