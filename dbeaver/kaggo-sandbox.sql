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
	
SELECT * FROM metadata;

SELECT * FROM reddit_subreddit_subscribers rss ;

SELECT * FROM youtube_video_views ycs ;

SELECT * FROM users;

SELECT * FROM metadata m WHERE request_kind LIKE 'twitch%';

SELECT m.data->>'id' AS foo
FROM metadata m
WHERE id = 'lj99vxv' AND request_kind = 'reddit.comment';

SELECT * FROM youtube_channel_views;

SELECT * FROM users_metadata_through umt ;

SELECT * FROM twitch_stream_views ;

SELECT * FROM metadata m WHERE request_kind like 'twitch.stream';


SELECT u.email, u.data AS "user_metadata", m.id, m.request_kind, m.data AS "metric_metadata"
FROM users u
INNER JOIN users_metadata_through umt ON u.email = umt.email
INNER JOIN metadata m ON umt.id = m.id AND umt.request_kind = m.request_kind
WHERE u.email = 'brojonat@gmail.com';

-- get
SELECT * FROM reddit_user_subscriptions;
-- insert
INSERT INTO reddit_user_subscriptions (name)
VALUES ('smartastic');
-- delete
DELETE  FROM reddit_user_subscriptions WHERE name = 'miaipanema';

SELECT id, DATA->>'owner' AS owner, DATA->>'link' AS link, data
FROM metadata 
WHERE request_kind = 'reddit.post';


-- get
SELECT * FROM reddit_subreddit_subscriptions;
-- insert
INSERT INTO reddit_subreddit_subscriptions (name) 
VALUES ('orangecounty');
-- delete
DELETE FROM reddit_subreddit_subscriptions WHERE name = ANY('{"golang","orangecounty"}'::VARCHAR[]);




