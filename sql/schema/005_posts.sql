-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text NOT NULL,
    url text NOT NULL UNIQUE,
    description text,
    published_at timestamp,
    feed_id uuid NOT NULL REFERENCES feeds ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;

