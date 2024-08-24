BEGIN;

-- youtube channel views
CREATE TABLE IF NOT EXISTS youtube_channel_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views INTEGER NOT NULL
);
SELECT create_hypertable('youtube_channel_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS youtube_channel_views_id ON youtube_channel_views (id, ts);

-- youtube channel subscribers
CREATE TABLE IF NOT EXISTS youtube_channel_subscribers (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    subscribers INTEGER NOT NULL
);
SELECT create_hypertable('youtube_channel_subscribers', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS youtube_channel_subscribers_id ON youtube_channel_subscribers (id, ts);

-- youtube channel videos
CREATE TABLE IF NOT EXISTS youtube_channel_videos (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    videos INTEGER NOT NULL
);
SELECT create_hypertable('youtube_channel_videos', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS youtube_channel_videos_id ON youtube_channel_videos (id, ts);

COMMIT;