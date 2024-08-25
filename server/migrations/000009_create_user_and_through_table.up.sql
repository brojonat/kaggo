BEGIN;

CREATE TABLE IF NOT EXISTS users (
    email VARCHAR(255) PRIMARY KEY NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE IF NOT EXISTS users_metadata_through (
    email VARCHAR(255),
    id VARCHAR(255),
    request_kind VARCHAR(255),
    PRIMARY KEY (email, id, request_kind)
);
ALTER TABLE users_metadata_through
ADD FOREIGN KEY (email) REFERENCES users ON DELETE CASCADE;
ALTER TABLE users_metadata_through
ADD FOREIGN KEY (id, request_kind) REFERENCES metadata ON DELETE CASCADE;


COMMIT;