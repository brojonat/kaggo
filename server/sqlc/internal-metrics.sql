-- name: InsertInternalRandom :exec
INSERT INTO internal_random (id, ts, val)
VALUES (@id, NOW()::TIMESTAMPTZ, @value);

-- name: GetInternalMetricsByIDs :many
SELECT
    id AS "id",
    ts AS "ts",
    val AS "value",
    'internal.random' AS "metric"
FROM internal_random AS i
WHERE
    i.id = ANY(@ids::VARCHAR[]) AND
    i.ts >= @ts_start AND
    i.ts <= @ts_end;
