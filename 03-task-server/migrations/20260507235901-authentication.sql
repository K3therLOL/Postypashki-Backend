
-- +migrate Up
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE sessions;
DROP TABLE users;
