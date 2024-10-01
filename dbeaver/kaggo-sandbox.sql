

-- get
SELECT * FROM youtube_channel_subscriptions ycs ;
-- insert
INSERT INTO youtube_channel_subscriptions(id)
VALUES ('UCKrdjiuS66yXOdEZ_cOD_TA');
-- delete
DELETE FROM youtube_channel_subscriptions WHERE url = ANY('{"foo"}'::VARCHAR[]);


SELECT ycs.id, "data"->>'title'
FROM youtube_channel_subscriptions ycs
LEFT JOIN metadata m ON ycs.id = m.id;


SELECT m.id AS username, children.post_id AS post_id, children."data" AS data
FROM metadata m 
LEFT JOIN (
	SELECT 
		id AS post_id,
		m2."data" ->> 'parent_user_name' AS parent_user_name,
		m2."data" AS "data"
	FROM metadata m2 
	WHERE request_kind = 'reddit.post'
) children ON m.id = children.parent_user_name
WHERE request_kind = 'reddit.user';

SELECT * FROM metadata m 
WHERE request_kind = 'reddit.post' AND id='1d6cphj';


-- find all nsfw tagged entities
SELECT * FROM metadata m
WHERE (m."data" ->> 'tags')::JSONB ? 'NSFW';

-- find all nsfw tagged entities
SELECT m."data" FROM metadata m
WHERE m."data" ->>'parent_user_name' LIKE 'Smart%';

SELECT * FROM reddit_post_score
WHERE id = 'foo';

SELECT * FROM metadata m WHERE request_kind = 'reddit.user';

SELECT * FROM reddit_subreddit_subscriptions rss ;

DELETE FROM reddit_user_subscriptions ;




