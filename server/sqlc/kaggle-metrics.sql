-- name: InsertKaggleNotebookVotes :exec
INSERT INTO kaggle_notebook_votes (slug, ts, votes)
VALUES (@slug, NOW()::TIMESTAMPTZ, @votes);

-- name: InsertKaggleDatasetVotes :exec
INSERT INTO kaggle_dataset_votes (slug, ts, votes)
VALUES (@slug, NOW()::TIMESTAMPTZ, @votes);

-- name: InsertKaggleDatasetViews :exec
INSERT INTO kaggle_dataset_views (slug, ts, views)
VALUES (@slug, NOW()::TIMESTAMPTZ, @views);

-- name: InsertKaggleDatasetDownloads :exec
INSERT INTO kaggle_dataset_downloads (slug, ts, downloads)
VALUES (@slug, NOW()::TIMESTAMPTZ, @downloads);

-- name: GetKaggleNotebookMetrics :many
SELECT *, 'knv' AS "metric"
FROM kaggle_notebook_votes knv
WHERE knv.slug = @slug;

-- name: GetKaggleDatasetMetrics :many
SELECT *, 'kaggle.dataset.votes' AS "metric"
FROM kaggle_dataset_votes kdv
WHERE kdv.slug = @slug
UNION ALL
SELECT *, 'kaggle.dataset.views' AS "metric"
FROM kaggle_dataset_views kdvw
WHERE kdvw.slug = @slug
UNION ALL
SELECT *, 'kaggle.dataset.downloads' AS "metric"
FROM kaggle_dataset_downloads kdd
WHERE kdd.slug = @slug;

