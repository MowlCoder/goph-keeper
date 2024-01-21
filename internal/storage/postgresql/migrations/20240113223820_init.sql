-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS log_pass (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    login VARCHAR(512) NOT NULL,
    password VARCHAR(255) NOT NULL,
    source VARCHAR(512) DEFAULT '',
    version INT DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    number VARCHAR(128) NOT NULL,
    expired_at VARCHAR(20) NOT NULL,
    cvv VARCHAR(128) NOT NULL,
    meta TEXT DEFAULT '',
    version INT DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS log_pass;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS cards;
-- +goose StatementEnd
