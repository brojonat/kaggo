


-- get
SELECT * FROM reddit_user_subscriptions;
-- insert
INSERT INTO reddit_user_subscriptions (name)
VALUES ('Francis_Star');
-- delete
DELETE  FROM reddit_user_subscriptions WHERE name = 'miaipanema';


-- get
SELECT * FROM reddit_subreddit_subscriptions;
-- insert
INSERT INTO reddit_subreddit_subscriptions (name) 
VALUES ('orangecounty');
-- delete
DELETE FROM reddit_subreddit_subscriptions WHERE name = ANY('{"golang","orangecounty"}'::VARCHAR[]);


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



SELECT * FROM metadata m 
WHERE request_kind = 'reddit.subreddit';

-- find all nsfw tagged entities
SELECT * FROM metadata m 
WHERE (m."data" ->> 'tags')::JSONB ? 'NSFW';

-- find all nsfw tagged entities
SELECT m."data" ->>'owner' FROM metadata m
WHERE request_kind = 'reddit.post';

SELECT * FROM reddit_post_score
WHERE id = '1fafkp7';




