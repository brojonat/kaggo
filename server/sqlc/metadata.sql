-- name: InsertMetadata :exec
INSERT INTO metadata (id, request_kind, data)
VALUES (@id, @request_kind, @data)
ON CONFLICT ON CONSTRAINT metadata_pkey DO UPDATE
SET data = EXCLUDED.data;

-- name: GetMetadataByIDs :many
SELECT id, request_kind, data
FROM metadata
WHERE id = ANY(@ids::VARCHAR[]);
