-- schema.sql for sqlc generation, DO NOT use with atlas; use golang-migrate instead.

-- youtube video views
CREATE TABLE IF NOT EXISTS youtube_video_views (
    id VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views INTEGER NOT NULL
);
SELECT create_hypertable('youtube_video_views', 'ts', if_not_exists => TRUE);

-- youtube video likes
CREATE TABLE IF NOT EXISTS youtube_video_likes (
    id VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    likes INTEGER NOT NULL
);
SELECT create_hypertable('youtube_video_likes', 'ts', if_not_exists => TRUE);

-- youtube video comments
CREATE TABLE IF NOT EXISTS youtube_video_comments (
    id VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    comments INTEGER NOT NULL
);
SELECT create_hypertable('youtube_video_comments', 'ts', if_not_exists => TRUE);

-- internal metric for testing
CREATE TABLE IF NOT EXISTS internal_random (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);

-- kaggle notebook votes
CREATE TABLE IF NOT EXISTS kaggle_notebook_votes (
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);

-- kaggle notebook downloads
CREATE TABLE IF NOT EXISTS kaggle_notebook_downloads (
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);

-- kaggle dataset votes
CREATE TABLE IF NOT EXISTS kaggle_dataset_votes (
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);

-- kaggle dataset downloads
CREATE TABLE IF NOT EXISTS kaggle_dataset_downloads (
    slug VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
