CREATE TABLE IF NOT EXISTS sessions
(
    id         VARCHAR(100) PRIMARY KEY,
    user_id    VARCHAR(100) NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    token      VARCHAR(255) NOT NULL,
    created_at timestamp    NOT NULL,
    updated_at timestamp    NOT NULL
);
