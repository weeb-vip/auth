CREATE TABLE IF NOT EXISTS credentials
(
    id         VARCHAR(100) PRIMARY KEY,
    username   varchar(255) NOT NULL,
    type       varchar(255) NOT NULL,
    value      varchar(255) NOT NULL,
    created_at timestamp    NOT NULL,
    updated_at timestamp    NOT NULL
);
