-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS threads (
    id BIGINT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS threads;