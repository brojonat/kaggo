BEGIN;

-- reddit post score
CREATE TABLE IF NOT EXISTS reddit_post_score (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    score INTEGER NOT NULL
);
SELECT create_hypertable('reddit_post_score', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_post_score_id ON reddit_post_score (id, ts);
CREATE INDEX IF NOT EXISTS reddit_post_score_title ON reddit_post_score (title, ts);

-- reddit post upvote ratio
CREATE TABLE IF NOT EXISTS reddit_post_ratio (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    ratio REAL NOT NULL
);
SELECT create_hypertable('reddit_post_ratio', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_post_ratio_id ON reddit_post_ratio (id, ts);
CREATE INDEX IF NOT EXISTS reddit_post_ratio_title ON reddit_post_ratio (title, ts);

-- reddit comment score
CREATE TABLE IF NOT EXISTS reddit_comment_score (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    score INTEGER NOT NULL
);
SELECT create_hypertable('reddit_comment_score', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_comment_score_id ON reddit_comment_score (id, ts);

-- reddit comment controversiality
CREATE TABLE IF NOT EXISTS reddit_comment_controversiality (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    controversiality REAL NOT NULL
);
SELECT create_hypertable('reddit_comment_controversiality', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_comment_controversiality_id ON reddit_comment_controversiality (id, ts);

COMMIT;