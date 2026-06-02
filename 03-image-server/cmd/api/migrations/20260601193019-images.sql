
-- +migrate Up
CREATE TABLE images (
    image_id SERIAL PRIMARY KEY,
    task_id TEXT NOT NULL,
    url TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE
);

-- +migrate Down
DROP TABLE images;
