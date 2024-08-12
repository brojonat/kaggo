-- name: InsertYouTubeVideoViews :exec
INSERT INTO youtube_video_views (id, slug, ts, views)
VALUES (@id, @slug, NOW()::TIMESTAMPTZ, @views);

-- name: InsertYouTubeVideoComments :exec
INSERT INTO youtube_video_comments (id, slug, ts, comments)
VALUES (@id, @slug, NOW()::TIMESTAMPTZ, @comments);

-- name: InsertYouTubeVideoLikes :exec
INSERT INTO youtube_video_likes (id, slug, ts, likes)
VALUES (@id, @slug, NOW()::TIMESTAMPTZ, @likes);

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

-- name: GetYouTubeVideoMetricsBySlug :many
SELECT *, 'youtube.video.views' AS "metric"
FROM youtube_video_views y
WHERE y.slug = @slug
UNION ALL
SELECT *, 'youtube.video.comments' AS "metric"
FROM youtube_video_views y
WHERE y.slug = @slug
UNION ALL
SELECT *, 'youtube.video.likes' AS "metric"
FROM youtube_video_likes y
WHERE y.slug = @slug;
