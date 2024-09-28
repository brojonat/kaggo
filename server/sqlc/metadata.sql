-- name: InsertMetadata :exec
INSERT INTO metadata (id, request_kind, data)
VALUES (@id, @request_kind, @data)
ON CONFLICT ON CONSTRAINT metadata_pkey DO UPDATE
SET data = EXCLUDED.data;

-- name: GetMetadataByIDs :many
SELECT id, request_kind, data
FROM metadata
WHERE id = ANY(@ids::VARCHAR[]);

-- name: GetMetadatum :one
SELECT id, request_kind, data
FROM metadata
WHERE request_kind = @request_kind AND LOWER(id) = LOWER(@id);


-- name: GetChildrenMetadataByID :many
SELECT
    m.id AS username,
    m.request_kind AS "parent_request_kind",
    children.id AS child_id,
    children.request_kind AS "child_request_kind",
    children."data" AS data
FROM metadata m
LEFT JOIN (
	SELECT
		id AS id,
        m2.request_kind AS request_kind,
		m2."data" AS "data",
		m2."data" ->> 'parent_user_name' AS parent_id
	FROM metadata m2
	WHERE m2.request_kind = @child_request_kind
) children ON m.id = children.parent_id
WHERE m.id = @id AND m.request_kind = @parent_request_kind;
