-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS aliases
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias VARCHAR(255) NOT NULL UNIQUE,
    uuid TEXT NOT NULL UNIQUE,
    time_ttl TIMESTAMP NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS aliases;
-- +goose StatementEnd
