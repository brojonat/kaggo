CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS vector;


SELECT *, 'knv' AS "metric" FROM kaggle_notebook_votes knv 
UNION ALL
SELECT *, 'knd' AS "metric" FROM kaggle_notebook_downloads knd ;

SELECT
    knv.slug AS "knv_slug",
    knv.ts AS "knv_ts",
    knv.val AS "knv_val",
    knd.slug AS "knd_slug",
    knd.ts AS "knd_ts",
    knd.val AS "knd_val"
FROM kaggle_notebook_votes knv
FULL JOIN kaggle_notebook_downloads knd ON knv.ts = knd.ts;


SELECT * FROM internal_random ir ;
SELECT * FROM youtube_video_views yvv ;

SELECT *, 'youtube.video.views' AS "metric" FROM youtube_video_views yvv
UNION ALL
SELECT *, 'youtube.video.comments' AS "metric" FROM youtube_video_comments yvc
UNION ALL
SELECT *, 'youtube.video.likes' AS "metric" FROM youtube_video_likes yvl;

SELECT *, 'reddit.post.score' AS "metric" FROM reddit_post_score rps
-- UNION ALL
-- SELECT *, 'reddit.post.ratio' AS "metric" FROM reddit_post_ratio rpr
ORDER BY ts DESC;

SELECT *, 'reddit.comment.score' AS "metric" FROM reddit_comment_score rcs
UNION ALL
SELECT *, 'reddit.comment.controversiality' AS "metric" FROM reddit_comment_controversiality rcr
ORDER BY ts ASC;


CREATE TABLE IF NOT EXISTS reddit_comment_controversiality (
    id VARCHAR(255) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    controversiality REAL NOT NULL
);
SELECT create_hypertable('reddit_comment_controversiality', 'ts', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS reddit_comment_controversiality_id ON reddit_comment_controversiality (id, ts);

