-- schema.sql for sqlc generation, DO NOT use with atlas; use golang-migrate instead.

CREATE TABLE IF NOT EXISTS metadata (
    id VARCHAR(255) NOT NULL,
    request_kind VARCHAR(255) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::JSONB,
    PRIMARY KEY (id, request_kind)
);

CREATE TABLE IF NOT EXISTS users (
    email VARCHAR(255) PRIMARY KEY NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS users_metadata_through (
    email VARCHAR(255) NOT NULL REFERENCES users(email) ON DELETE CASCADE,
    id VARCHAR(255) NOT NULL REFERENCES metadata(id) ON DELETE CASCADE,
    request_kind VARCHAR(255) NOT NULL REFERENCES metadata(request_kind) ON DELETE CASCADE,
    PRIMARY KEY (email, id, request_kind)
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
    views BIGINT NOT NULL
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
    views BIGINT NOT NULL
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

-- twitch clip views
CREATE TABLE IF NOT EXISTS twitch_clip_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views BIGINT NOT NULL
);

-- twitch video views
CREATE TABLE IF NOT EXISTS twitch_video_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views BIGINT NOT NULL
);

-- twitch stream views
CREATE TABLE IF NOT EXISTS twitch_stream_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views BIGINT NOT NULL
);

-- twitch user past dec avg views
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_avg_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    avg_views REAL NOT NULL
);

-- twitch user past dec med views
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_med_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    med_views REAL NOT NULL
);

-- twitch user past dec std views
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_std_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    std_views REAL NOT NULL
);

-- twitch user past dec avg duration
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_avg_duration (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    avg_duration REAL NOT NULL
);

-- twitch user past dec med duration
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_med_duration (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    med_duration REAL NOT NULL
);

-- twitch user past dec std duration
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_std_duration (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    std_duration REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS reddit_user_subscriptions (
    name VARCHAR(255) PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS reddit_subreddit_subscriptions (
    name VARCHAR(255) PRIMARY KEY NOT NULL
);
