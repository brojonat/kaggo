BEGIN;

-- reddit subreddit subscribers
CREATE TABLE IF NOT EXISTS reddit_subreddit_subscribers (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    subscribers INTEGER NOT NULL
);
SELECT create_hypertable('reddit_subreddit_subscribers', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_subreddit_subscribers_id ON reddit_subreddit_subscribers (id, ts);

-- reddit subreddit active user counts
CREATE TABLE IF NOT EXISTS reddit_subreddit_active_user_count (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    active_user_count INTEGER NOT NULL
);
SELECT create_hypertable('reddit_subreddit_active_user_count', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_subreddit_active_user_count_id ON reddit_subreddit_active_user_count (id, ts);

COMMIT;