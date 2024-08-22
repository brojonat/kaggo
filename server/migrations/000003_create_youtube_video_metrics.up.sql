BEGIN;

-- youtube video views
CREATE TABLE IF NOT EXISTS youtube_video_views (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views INTEGER NOT NULL
);
SELECT create_hypertable('youtube_video_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS youtube_video_views_id ON youtube_video_views (id, ts);
CREATE INDEX IF NOT EXISTS youtube_video_views_title ON youtube_video_views (title, ts);

-- youtube video likes
CREATE TABLE IF NOT EXISTS youtube_video_likes (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    likes INTEGER NOT NULL
);
SELECT create_hypertable('youtube_video_likes', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS youtube_video_likes_id ON youtube_video_likes (id, ts);

-- youtube video comments
CREATE TABLE IF NOT EXISTS youtube_video_comments (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    comments INTEGER NOT NULL
);
SELECT create_hypertable('youtube_video_comments', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS youtube_video_comments_id ON youtube_video_comments (id, ts);

COMMIT;
