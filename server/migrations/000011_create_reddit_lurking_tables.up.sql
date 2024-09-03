BEGIN;

CREATE TABLE IF NOT EXISTS reddit_user_subscriptions (
    name VARCHAR(255) PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS reddit_subreddit_subscriptions (
    name VARCHAR(255) PRIMARY KEY NOT NULL
);

COMMIT;