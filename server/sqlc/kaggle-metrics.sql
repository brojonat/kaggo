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

















-- Kaggle Notebook Bucketed Metrics


-- name: GetKaggleNotebookMetricsByIDsBucket15Min :many
SELECT *, 'kaggle.notebook.votes' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(votes::REAL) AS "value"
	FROM kaggle_notebook_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetKaggleNotebookMetricsByIDsBucket1Hr :many
SELECT *, 'kaggle.notebook.votes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(votes::REAL) AS value
	FROM kaggle_notebook_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetKaggleNotebookMetricsByIDsBucket8Hr :many
SELECT *, 'kaggle.notebook.votes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(votes::REAL) AS value
	FROM kaggle_notebook_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetKaggleNotebookMetricsByIDsBucket1Day :many
SELECT *, 'kaggle.notebook.votes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(votes::REAL) AS value
	FROM kaggle_notebook_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;


-- Kaggle Dataset Bucketed Metrics


-- name: GetKaggleDatasetMetricsByIDsBucket15Min :many
SELECT *, 'kaggle.dataset.votes' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(votes::REAL) AS "value"
	FROM kaggle_dataset_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM kaggle_dataset_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.downloads' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '15 minutes', ts) AS "bucket",
	    MAX(downloads::REAL) AS "value"
	FROM kaggle_dataset_downloads
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetKaggleDatasetMetricsByIDsBucket1Hr :many
SELECT *, 'kaggle.dataset.votes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(votes::REAL) AS value
	FROM kaggle_dataset_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM kaggle_dataset_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.downloads' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(downloads::REAL) AS value
	FROM kaggle_dataset_downloads
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetKaggleDatasetMetricsByIDsBucket8Hr :many
SELECT *, 'kaggle.dataset.votes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS bucket,
	    MAX(votes::REAL) AS value
	FROM kaggle_dataset_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.views' AS "metric"
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '8 hours', ts) AS "bucket",
	    MAX(views::REAL) AS "value"
	FROM kaggle_dataset_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab
WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.downloads' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(downloads::REAL) AS value
	FROM kaggle_dataset_downloads
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;

-- name: GetKaggleDatasetMetricsByIDsBucket1Day :many
SELECT *, 'kaggle.dataset.votes' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(votes::REAL) AS value
	FROM kaggle_dataset_votes
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.views' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 day', ts) AS bucket,
	    MAX(views::REAL) AS value
	FROM kaggle_dataset_views
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ
UNION ALL
SELECT *, 'kaggle.dataset.downloads' AS metric
FROM (
	SELECT
		id,
	    time_bucket(INTERVAL '1 hour', ts) AS bucket,
	    MAX(downloads::REAL) AS value
	FROM kaggle_dataset_downloads
	GROUP BY id, bucket
	ORDER BY id, bucket
) AS tab WHERE
    tab.id = ANY(@ids::VARCHAR[]) AND
    tab.bucket >= @ts_start::TIMESTAMPTZ AND
    tab.bucket <= @ts_end::TIMESTAMPTZ;