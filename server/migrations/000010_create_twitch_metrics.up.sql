BEGIN;

-- twitch clip views
CREATE TABLE IF NOT EXISTS twitch_clip_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views BIGINT NOT NULL
);
SELECT create_hypertable('twitch_clip_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_clip_views_id ON twitch_clip_views (id, ts);

-- twitch video views
CREATE TABLE IF NOT EXISTS twitch_video_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views BIGINT NOT NULL
);
SELECT create_hypertable('twitch_video_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_video_views_id ON twitch_video_views (id, ts);

-- twitch stream views
CREATE TABLE IF NOT EXISTS twitch_stream_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    views BIGINT NOT NULL
);
SELECT create_hypertable('twitch_stream_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_stream_views_id ON twitch_stream_views (id, ts);

-- twitch user past dec avg views
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_avg_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    avg_views REAL NOT NULL
);
SELECT create_hypertable('twitch_user_past_dec_avg_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_user_past_dec_avg_views_id ON twitch_user_past_dec_avg_views (id, ts);

-- twitch user past dec med views
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_med_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    med_views REAL NOT NULL
);
SELECT create_hypertable('twitch_user_past_dec_med_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_user_past_dec_med_views_id ON twitch_user_past_dec_med_views (id, ts);

-- twitch user past dec std views
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_std_views (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    std_views REAL NOT NULL
);
SELECT create_hypertable('twitch_user_past_dec_std_views', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_user_past_dec_std_views_id ON twitch_user_past_dec_std_views (id, ts);

-- twitch user past dec avg duration
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_avg_duration (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    avg_duration REAL NOT NULL
);
SELECT create_hypertable('twitch_user_past_dec_avg_duration', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_user_past_dec_avg_duration_id ON twitch_user_past_dec_avg_duration (id, ts);

-- twitch user past dec med duration
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_med_duration (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    med_duration REAL NOT NULL
);
SELECT create_hypertable('twitch_user_past_dec_med_duration', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_user_past_dec_med_duration_id ON twitch_user_past_dec_med_duration (id, ts);

-- twitch user past dec std duration
CREATE TABLE IF NOT EXISTS twitch_user_past_dec_std_duration (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    std_duration REAL NOT NULL
);
SELECT create_hypertable('twitch_user_past_dec_std_duration', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS twitch_user_past_dec_std_duration_id ON twitch_user_past_dec_std_duration (id, ts);

COMMIT;