CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS vector;


SELECT slug AS "id", ts AS "tss", votes AS "value", 'kaggle.notebook' AS "metric"
FROM kaggle_notebook_votes knv WHERE slug != 'foo/bar-baz'
UNION ALL
SELECT *, 'kaggle.dataset.views' AS "metric" 
FROM kaggle_dataset_views kdv WHERE slug != 'foo/bar-baz'
UNION ALL
SELECT *, 'kaggle.dataset.downloads' AS "metric" 
FROM kaggle_dataset_downloads kdd WHERE slug != 'foo/bar-baz'
UNION ALL
SELECT *, 'kaggle.dataset.votes' AS "metric" 
FROM kaggle_dataset_votes kdv2 WHERE slug != 'foo/bar-baz';

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

SELECT id,
	first(title, ts) AS "Title",
   time_bucket(INTERVAL '15 min', ts) AS bucket,
   MAX(views) AS value,
   'foo' AS "foo"
FROM youtube_video_views yvv 
GROUP BY id, bucket;

SELECT * FROM metadata m ;

INSERT INTO metadata (id, request_kind, data)
VALUES ('foo/bar-baz', 'youtube.video', '{}'::JSONB)
ON CONFLICT ON CONSTRAINT metadata_pkey DO UPDATE
SET DATA = EXCLUDED.data;


