BEGIN;

CREATE TABLE IF NOT EXISTS metadata (
    id VARCHAR(255) NOT NULL,
    request_kind VARCHAR(255) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::JSONB,
    PRIMARY KEY (id, request_kind)
);

COMMIT;