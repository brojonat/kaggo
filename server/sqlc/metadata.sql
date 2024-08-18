-- name: InsertMetadata :exec
INSERT INTO metadata (id, metric_kind, data)
VALUES (@id, @metric_kind, @data);

-- name: GetMetadataByIDs :many
SELECT id, metric_kind, data
FROM metadata
WHERE id = ANY(@ids::VARCHAR[]);
