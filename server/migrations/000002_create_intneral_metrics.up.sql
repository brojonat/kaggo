BEGIN;
-- internal metric
CREATE TABLE IF NOT EXISTS internal_random (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    val INTEGER NOT NULL
);
SELECT create_hypertable('internal_random', 'ts', if_not_exists => TRUE);
CREATE INDEX internal_random_value_id ON internal_random (id, ts);
COMMIT;
