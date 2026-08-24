-- migrate:up
CREATE TABLE IF NOT EXISTS album(
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    total_cards INTEGER NOT NULL,
    created_at INTEGER NOT NULL
);

-- migrate:down
DROP TABLE IF NOT EXISTS album;
