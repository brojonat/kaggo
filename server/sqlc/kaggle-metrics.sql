-- name: InsertKaggleNotebookVotes :exec
INSERT INTO kaggle_notebook_votes (slug, ts, val)
VALUES (@slug, NOW()::TIMESTAMPTZ, @votes);

-- name: InsertKaggleNotebookDownloads :exec
INSERT INTO kaggle_notebook_downloads (slug, ts, val)
VALUES (@slug, NOW()::TIMESTAMPTZ, @downloads);

-- name: InsertKaggleDatasetVotes :exec
INSERT INTO kaggle_dataset_votes (slug, ts, val)
VALUES (@slug, NOW()::TIMESTAMPTZ, @votes);

-- name: InsertKaggleDatasetDownloads :exec
INSERT INTO kaggle_dataset_downloads (slug, ts, val)
VALUES (@slug, NOW()::TIMESTAMPTZ, @downloads);

-- name: GetKaggleNotebookMetrics :many
SELECT *, 'knv' AS "metric"
FROM kaggle_notebook_votes knv
WHERE knv.slug = @slug
UNION ALL
SELECT *, 'knd' AS "metric"
FROM kaggle_notebook_downloads knd
WHERE knd.slug = @slug;

-- name: GetKaggleDatasetMetrics :many
SELECT *, 'kdv' AS "metric"
FROM kaggle_dataset_votes kdv
WHERE kdv.slug = @slug
UNION ALL
SELECT *, 'kdd' AS "metric"
FROM kaggle_dataset_downloads kdd
WHERE kdd.slug = @slug;;
