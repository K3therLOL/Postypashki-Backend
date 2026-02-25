-- +migrate Up
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL
);
-- +migrate Down
DROP TABLE users;
