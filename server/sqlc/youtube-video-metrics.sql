-- name: InsertYouTubeVideoViews :exec
INSERT INTO youtube_video_views (id, title, ts, views)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @views);

-- name: InsertYouTubeVideoComments :exec
INSERT INTO youtube_video_comments (id, title, ts, comments)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @comments);

-- name: InsertYouTubeVideoLikes :exec
INSERT INTO youtube_video_likes (id, title, ts, likes)
VALUES (@id, @title, NOW()::TIMESTAMPTZ, @likes);

-- name: GetYouTubeVideoMetricsByID :many
SELECT *, 'youtube.video.views' AS "metric"
FROM youtube_video_views y
WHERE y.id = @id
UNION ALL
SELECT *, 'youtube.video.comments' AS "metric"
FROM youtube_video_views y
WHERE y.id = @id
UNION ALL
SELECT *, 'youtube.video.likes' AS "metric"
FROM youtube_video_likes y
WHERE y.id = @id;

-- name: GetYouTubeVideoMetricsByTitle :many
SELECT *, 'youtube.video.views' AS "metric"
FROM youtube_video_views y
WHERE y.title = @title
UNION ALL
SELECT *, 'youtube.video.comments' AS "metric"
FROM youtube_video_views y
WHERE y.title = @title
UNION ALL
SELECT *, 'youtube.video.likes' AS "metric"
FROM youtube_video_likes y
WHERE y.title = @title;
