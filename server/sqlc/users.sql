-- name: InsertUser :exec
INSERT INTO users (email, data)
VALUES (@email, @data)
ON CONFLICT ON CONSTRAINT users_pkey DO UPDATE
SET data = EXCLUDED.data;

-- name: GrantMetricToUser :exec
INSERT INTO users_metadata_through (email, id, request_kind)
VALUES (@email, @id, @request_kind);

-- name: GrantMetricGroupToUser :exec
INSERT INTO users_metadata_through (email, id, request_kind)
SELECT u.email, m.id, m.request_kind
FROM users u
CROSS JOIN metadata m
WHERE u.email = @email AND m.request_kind = @request_kind;

-- name: RemoveMetricGroupFromUser :exec
DELETE FROM users_metadata_through
WHERE email = @email AND request_kind = @request_kind;

-- name: RemoveMetricFromUser :exec
DELETE FROM users_metadata_through
WHERE email = @email AND id = @id AND request_kind = @request_kind;

-- name: GetUsers :many
SELECT *
FROM users
WHERE email = ANY(@emails::VARCHAR[]);

-- name: DeleteUsers :exec
DELETE FROM users
WHERE email = ANY(@emails::VARCHAR[]);

-- name: GetUserMetrics :many
SELECT u.email, u.data AS "user_metadata", m.id, m.request_kind, m.data AS "metric_metadata"
FROM users u
INNER JOIN users_metadata_through umt ON u.email = umt.email
INNER JOIN metadata m ON umt.id = m.id AND umt.request_kind = m.request_kind
WHERE u.email = @email;
