-- schema.sql for sqlc generation, DO NOT use with atlas; use golang-migrate instead.

CREATE TABLE metadata (
    id VARCHAR(255) NOT NULL,
    request_kind VARCHAR(255) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::JSONB,
    PRIMARY KEY (id, request_kind)
);

-- internal metric for testing
CREATE TABLE IF NOT EXISTS internal_random (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);

-- youtube video views
CREATE TABLE IF NOT EXISTS youtube_video_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views INTEGER NOT NULL
);

-- youtube video likes
CREATE TABLE IF NOT EXISTS youtube_video_likes (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    likes INTEGER NOT NULL
);

-- youtube video comments
CREATE TABLE IF NOT EXISTS youtube_video_comments (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    comments INTEGER NOT NULL
);

-- youtube channel views
CREATE TABLE IF NOT EXISTS youtube_channel_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views INTEGER NOT NULL
);

-- youtube channel subscribers
CREATE TABLE IF NOT EXISTS youtube_channel_subscribers (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    subscribers INTEGER NOT NULL
);

-- youtube channel videos
CREATE TABLE IF NOT EXISTS youtube_channel_videos (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    videos INTEGER NOT NULL
);

-- kaggle notebook votes
CREATE TABLE IF NOT EXISTS kaggle_notebook_votes (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    votes INTEGER NOT NULL
);

-- kaggle dataset votes
CREATE TABLE IF NOT EXISTS kaggle_dataset_votes (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    votes INTEGER NOT NULL
);

-- kaggle dataset views
CREATE TABLE IF NOT EXISTS kaggle_dataset_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views INTEGER NOT NULL
);

-- kaggle dataset downloads
CREATE TABLE IF NOT EXISTS kaggle_dataset_downloads (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    downloads INTEGER NOT NULL
);

-- reddit post score
CREATE TABLE IF NOT EXISTS reddit_post_score (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    score INTEGER NOT NULL
);

-- reddit post ratio
CREATE TABLE IF NOT EXISTS reddit_post_ratio (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    ratio REAL NOT NULL
);

-- reddit comment score
CREATE TABLE IF NOT EXISTS reddit_comment_score (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    score INTEGER NOT NULL
);

-- reddit comment controversiality
CREATE TABLE IF NOT EXISTS reddit_comment_controversiality (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    controversiality REAL NOT NULL
);

-- reddit subreddit subscribers
CREATE TABLE IF NOT EXISTS reddit_subreddit_subscribers (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    subscribers INTEGER NOT NULL
);

-- reddit subreddit active user counts
CREATE TABLE IF NOT EXISTS reddit_subreddit_active_user_count (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    active_user_count INTEGER NOT NULL
);
