CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id         VARCHAR(100) PRIMARY KEY,
    user_id    VARCHAR(100) NOT NULL,
    token      VARCHAR(100) NOT NULL,
    expiry     INTEGER      NOT NULL,
    created_at timestamptz    NOT NULL,
    updated_at timestamptz    NOT NULL
);
