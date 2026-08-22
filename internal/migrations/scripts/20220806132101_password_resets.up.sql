CREATE TABLE IF NOT EXISTS password_resets
(
    id            VARCHAR(100) PRIMARY KEY,
    credential_id VARCHAR(100)      NOT NULL,
    ott           VARCHAR(255) NOT NULL,
    created_at    timestamptz    NOT NULL,
    updated_at    timestamptz    NOT NULL
);
