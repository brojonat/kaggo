-- name: InsertInternalRandom :exec
INSERT INTO internal_random (id, ts, val)
VALUES (@id, NOW()::TIMESTAMPTZ, @value);

-- name: GetInternalMetrics :many
SELECT *, 'internal_random' AS "metric"
FROM internal_random ir
WHERE ir.id = @id;