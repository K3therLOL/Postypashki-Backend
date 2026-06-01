
-- +migrate Up
CREATE TABLE images (
    image_id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    url TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE
);

-- +migrate Down
DROP TABLE images;
