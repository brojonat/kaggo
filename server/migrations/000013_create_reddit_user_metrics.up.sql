BEGIN;

-- reddit user awardee karma
CREATE TABLE IF NOT EXISTS reddit_user_awardee_karma (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    karma INTEGER NOT NULL
);
SELECT create_hypertable('reddit_user_awardee_karma', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_user_awardee_karma_id ON reddit_user_awardee_karma (id, ts);

-- reddit user awarder karma
CREATE TABLE IF NOT EXISTS reddit_user_awarder_karma (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    karma INTEGER NOT NULL
);
SELECT create_hypertable('reddit_user_awarder_karma', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_user_awarder_karma_id ON reddit_user_awarder_karma (id, ts);

-- reddit user comment karma
CREATE TABLE IF NOT EXISTS reddit_user_comment_karma (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    karma INTEGER NOT NULL
);
SELECT create_hypertable('reddit_user_comment_karma', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_user_comment_karma_id ON reddit_user_comment_karma (id, ts);

-- reddit user link karma
CREATE TABLE IF NOT EXISTS reddit_user_link_karma (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    karma INTEGER NOT NULL
);
SELECT create_hypertable('reddit_user_link_karma', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_user_link_karma_id ON reddit_user_link_karma (id, ts);

-- reddit user total karma
CREATE TABLE IF NOT EXISTS reddit_user_total_karma (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    karma INTEGER NOT NULL
);
SELECT create_hypertable('reddit_user_total_karma', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_user_total_karma_id ON reddit_user_total_karma (id, ts);

COMMIT;