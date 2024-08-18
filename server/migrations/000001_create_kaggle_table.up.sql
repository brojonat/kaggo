BEGIN;
CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS vector;
-- kaggle notebook votes
CREATE TABLE IF NOT EXISTS kaggle_notebook_votes (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
SELECT create_hypertable('kaggle_notebook_votes', 'ts', if_not_exists => TRUE);
CREATE INDEX kaggle_notebook_votes_id ON kaggle_notebook_votes (id, ts);

-- kaggle notebook downloads
CREATE TABLE IF NOT EXISTS kaggle_notebook_downloads (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
SELECT create_hypertable('kaggle_notebook_downloads', 'ts', if_not_exists => TRUE);
CREATE INDEX kaggle_notebook_votes_downloads ON kaggle_notebook_downloads (id, ts);

-- kaggle dataset votes
CREATE TABLE IF NOT EXISTS kaggle_dataset_votes (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
SELECT create_hypertable('kaggle_dataset_votes', 'ts', if_not_exists => TRUE);
CREATE INDEX kaggle_dataset_votes_id ON kaggle_dataset_votes (id, ts);

-- kaggle dataset views
CREATE TABLE IF NOT EXISTS kaggle_dataset_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
SELECT create_hypertable('kaggle_dataset_views', 'ts', if_not_exists => TRUE);
CREATE INDEX kaggle_dataset_views_id ON kaggle_dataset_views (id, ts);

-- kaggle dataset downloads
CREATE TABLE IF NOT EXISTS kaggle_dataset_downloads (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
SELECT create_hypertable('kaggle_dataset_downloads', 'ts', if_not_exists => TRUE);
CREATE INDEX kaggle_dataset_downloads_id ON kaggle_dataset_downloads (id, ts);
COMMIT;