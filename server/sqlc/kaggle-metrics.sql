-- name: InsertKaggleNotebookVotes :exec
INSERT INTO kaggle_notebook_votes (id, ts, votes)
VALUES (@id, NOW()::TIMESTAMPTZ, @votes);

-- name: InsertKaggleDatasetVotes :exec
INSERT INTO kaggle_dataset_votes (id, ts, votes)
VALUES (@id, NOW()::TIMESTAMPTZ, @votes);

-- name: InsertKaggleDatasetViews :exec
INSERT INTO kaggle_dataset_views (id, ts, views)
VALUES (@id, NOW()::TIMESTAMPTZ, @views);

-- name: InsertKaggleDatasetDownloads :exec
INSERT INTO kaggle_dataset_downloads (id, ts, downloads)
VALUES (@id, NOW()::TIMESTAMPTZ, @downloads);

-- name: GetKaggleNotebookMetrics :many
SELECT
    k.id AS "id",
    k.ts AS "ts",
    k.votes::REAL AS "value",
    'kaggle.notebook.votes' AS "metric"
FROM kaggle_notebook_votes k
WHERE
    k.id = ANY(@ids::VARCHAR[]) AND
    k.ts >= @ts_start AND
    k.ts <= @ts_end;

-- name: GetKaggleDatasetMetrics :many
SELECT
    id AS "id",
    ts AS "ts",
    votes::REAL AS "value",
    'kaggle.dataset.votes' AS "metric"
FROM kaggle_dataset_votes AS k
WHERE
    k.id = ANY(@ids::VARCHAR[]) AND
    k.ts >= @ts_start AND
    k.ts <= @ts_end
UNION ALL
SELECT
    id AS "id",
    ts AS "ts",
    views::REAL AS "value",
    'kaggle.dataset.views' AS "metric"
FROM kaggle_dataset_views AS k
WHERE
    k.id = ANY(@ids::VARCHAR[]) AND
    k.ts >= @ts_start AND
    k.ts <= @ts_end
UNION ALL
SELECT
    id AS "id",
    ts AS "ts",
    downloads::REAL AS "value",
    'kaggle.dataset.downloads' AS "metric"
FROM kaggle_dataset_downloads AS k
WHERE
    k.id = ANY(@ids::VARCHAR[]) AND
    k.ts >= @ts_start AND
    k.ts <= @ts_end;

