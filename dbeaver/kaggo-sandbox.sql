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

SELECT v.id, bucket, v.value, m.DATA ->> 'link' AS "title" FROM (
	SELECT yvv.id,
	   time_bucket(INTERVAL '15 min', ts) AS bucket,
	   MAX(views) AS value
	FROM youtube_video_views yvv 
	GROUP BY yvv.id, bucket
) v
LEFT JOIN metadata m ON v.id = m.id 
WHERE m.request_kind = 'youtube.video';
	
SELECT * FROM metadata m WHERE request_kind LIKE 'youtube.channel';

SELECT * FROM reddit_subreddit_subscribers rss ;

SELECT * FROM youtube_channel_subscribers ycs ;



